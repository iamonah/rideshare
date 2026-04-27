package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamonah/rideshare/services/apigateway/internal/infra/client"
	httptransport "github.com/iamonah/rideshare/services/apigateway/internal/transport/http"
	websockettransport "github.com/iamonah/rideshare/services/apigateway/internal/transport/websocket"
	"github.com/iamonah/rideshare/shared/env"
	"github.com/iamonah/rideshare/shared/messaging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr          = env.GetString("HTTP_ADDR", ":8081")
	tripServiceAddr   = env.GetString("TRIP_SERVICE_GRPC_URL", "")
	driverServiceAddr = env.GetString("DRIVER_SERVICE_GRPC_URL", "")

	rabbitUsername = env.GetString("RABBITMQ_DEFAULT_USER", "")
	rabbitPassword = env.GetString("RABBITMQ_DEFAULT_PASS", "")
	rabbitHost     = env.GetString("RABBITMQ_HOST", "")
	rabbitVhost    = env.GetString("RABBITMQ_VHOST", "")
	rabbitPort     = env.GetInt("RABBITMQ_PORT", 5672)
)

func main() {
	log.Println("Starting API Gateway...")

	if tripServiceAddr == "" {
		log.Fatal("TRIP_SERVICE_GRPC_URL is required")
	}

	if driverServiceAddr == "" {
		log.Fatal("DRIVER_SERVICE_GRPC_URL is required")
	}

	log.Printf("Trip gRPC client initialized for %s (non-blocking dial)", tripServiceAddr)
	tripClient, err := client.NewClient(tripServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer tripClient.Close()

	log.Printf("Driver gRPC client initialized for %s (non-blocking dial)", driverServiceAddr)
	driverClient, err := client.NewDriverClient(driverServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer driverClient.Close()

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
	websocketServer := websockettransport.NewServer(driverClient, rabbitClient)
	err = websocketServer.ListenDriverTripRequestsQueue(context.Background())
	if err != nil {
		log.Fatalf("failed to listen on driver trip requests queue: %v", err)
	}
	err = websocketServer.ListenRiderEventsQueue(context.Background())
	if err != nil {
		log.Fatalf("failed to listen on rider events queue: %v", err)
	}

	log.Println("Listening on", httpAddr)

	server := &http.Server{
		Addr: httpAddr,
		Handler: httptransport.NewRouter(httptransport.Dependencies{
			Handlers: tripClient,
			// Drivers:    driverClient,
			Websockets: websocketServer,
		}),
	}

	shutDown := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			shutDown <- err
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-shutDown:
		log.Fatalf("HTTP server closed unexpectedly: %v", err)
	case signal := <-signalChan:
		fmt.Printf("Signal received: %s", signal.String())
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("graceful shutdown failed: %s", err)
			_ = server.Close() // Force close if graceful shutdown fails
		}
	}
}
