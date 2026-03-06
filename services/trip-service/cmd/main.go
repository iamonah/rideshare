package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_Handler "github.com/iamonah/rideshare/services/trip-service/internal/infrastructure/grpc"
	tripdb "github.com/iamonah/rideshare/services/trip-service/internal/infrastructure/tripdb"
	tripservice "github.com/iamonah/rideshare/services/trip-service/internal/service"
	"github.com/iamonah/rideshare/shared/env"
	"github.com/iamonah/rideshare/shared/proto/pb/trip"
	"google.golang.org/grpc"
)

var (
	grpcAdd         = env.GetString("HTTP_ADDR", ":8081")
	tripServiceAddr = env.GetString("TRIP_SERVICE_GRPC_URL", "")
)

func main() {
	log.Println("--- Trip Service Initializing... ---")
	inmemRepo := tripdb.NewInmemRepository()
	svc := tripservice.NewService(inmemRepo)

	tripServer := grpc_Handler.NewTripServer(svc)

	listener, err := net.Listen("tcp", "9093")
	if err != nil {
		log.Fatalf("netListen %v", err)
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(nil))

	trip.RegisterTripServiceServer(grpcServer, tripServer)

	shutDown := make(chan error, 1)
	go func() {
		log.Printf("Trip gRPC server listening on %s", "9093")
		if err := grpcServer.Serve(listener); err != nil {
			shutDown <- err
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
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
