package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tripdomain "github.com/iamonah/rideshare/services/trip-service/internal/domain/trip"
	"github.com/iamonah/rideshare/services/trip-service/internal/infra/external/osrm"
	grpc_Handler "github.com/iamonah/rideshare/services/trip-service/internal/infra/grpc"
	tripdb "github.com/iamonah/rideshare/services/trip-service/internal/infra/tripdb"
	"github.com/iamonah/rideshare/shared/env"
	"github.com/iamonah/rideshare/shared/rabbitmq"
	"google.golang.org/grpc"
)

var (
	grpcAddr       = env.GetString("GRPC_ADDR", ":9093")
	osrmURL        = env.GetString("OSRM_BASE_URL", "")
	osrmTimeout    = env.GetDuration("OSRM_TIMEOUT", 5*time.Second)
	rabbitUsername = env.GetString("RABBITMQ_DEFAULT_USER", "")
	rabbitPassword = env.GetString("RABBITMQ_DEFAULT_PASS", "")
	rabbitHost     = env.GetString("RABBITMQ_HOST", "")
	rabbitVhost    = env.GetString("RABBITMQ_VHOST", "")
	rabbitPort     = env.GetInt("RABBITMQ_PORT", 5672)
)

func main() {
	log.Println("--- Trip Service Initializing... ---")
	inmemRepo := tripdb.NewInmemRepository()
	routeHTTPClient := &http.Client{Timeout: osrmTimeout}
	var routeProvider tripdomain.RouteProvider = osrm.NewClient(routeHTTPClient, osrmURL)

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("netListen %v", err)
	}

	rabbitClient, err := rabbitmq.NewRabbitMQClient(rabbitmq.RabbitMqConfig{
		Username: rabbitUsername,
		Password: rabbitPassword,
		Host:     rabbitHost,
		Vhost:    rabbitVhost,
		Port:     int16(rabbitPort),
	})
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Close()
	log.Println("starting RabbitMQ client...")

	svc := tripdomain.NewTripBusiness(inmemRepo, routeProvider, rabbitClient)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(nil))
	grpc_Handler.NewTripServer(grpcServer, svc)

	shutDown := make(chan error, 1)
	go func() {
		log.Printf("Trip gRPC server listening on %s", grpcAddr)
		if err := grpcServer.Serve(listener); err != nil {
			shutDown <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-shutDown:
		log.Fatalf("failed to serve: %v", err)
	case <-quit:
		log.Println("Shutting down Trip Service...")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			log.Println("Trip Service stopped gracefully")
		case <-ctx.Done():
			grpcServer.Stop()
			log.Fatal("Trip Service shutdown timed out, forcing exit")
		}
	}

}
