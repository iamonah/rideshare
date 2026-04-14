package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	driverservice "github.com/iamonah/rideshare/services/driver-service/internal"
	"github.com/iamonah/rideshare/services/driver-service/internal/infra/events"

	"github.com/iamonah/rideshare/shared/env"
	"github.com/iamonah/rideshare/shared/messaging"
	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9092"
var (
	rabbitUsername = env.GetString("RABBITMQ_DEFAULT_USER", "")
	rabbitPassword = env.GetString("RABBITMQ_DEFAULT_PASS", "")
	rabbitHost     = env.GetString("RABBITMQ_HOST", "")
	rabbitVhost    = env.GetString("RABBITMQ_VHOST", "")
	rabbitPort     = env.GetInt("RABBITMQ_PORT", 5672)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rabbitClient, err := messaging.NewRabbitMQClient(messaging.RabbitMqConfig{
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

	if err := events.BootstrapDriverTopology(rabbitClient); err != nil {
		log.Fatalf("failed to bootstrap driver RabbitMQ topology: %v", err)
	}

	tripConsumer := events.NewTripConsumer(rabbitClient)
	svc := driverservice.NewService(tripConsumer)
	

	// Starting the gRPC server
	grpcServer := grpcserver.NewServer()
	driverservice.NewGRPCHandler(grpcServer, svc)

	log.Printf("Starting gRPC server Driver service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
