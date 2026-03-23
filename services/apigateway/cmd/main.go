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

	appdriver "github.com/iamonah/rideshare/services/apigateway/internal/app/driver"
	apppayment "github.com/iamonah/rideshare/services/apigateway/internal/app/payment"
	apptrip "github.com/iamonah/rideshare/services/apigateway/internal/app/trip"
	tripgrpc "github.com/iamonah/rideshare/services/apigateway/internal/infra/tripgrpc"
	httptransport "github.com/iamonah/rideshare/services/apigateway/internal/transport/http"
	websockettransport "github.com/iamonah/rideshare/services/apigateway/internal/transport/websocket"
	"github.com/iamonah/rideshare/shared/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr        = env.GetString("HTTP_ADDR", ":8081")
	tripServiceAddr = env.GetString("TRIP_SERVICE_GRPC_URL", "")
)

func main() {
	log.Println("Starting API Gateway...")

	if tripServiceAddr == "" {
		log.Fatal("TRIP_SERVICE_GRPC_URL is required")
	}

	log.Printf("Trip gRPC client initialized for %s (non-blocking dial)", tripServiceAddr)
	tripClient, err := tripgrpc.NewClient(tripServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer tripClient.Close()

	tripService := apptrip.NewService(tripClient)
	driverService := appdriver.NewService()
	paymentService := apppayment.NewService()
	websocketHandler := websockettransport.NewHandler()

	log.Println("Listening on", httpAddr)

	server := &http.Server{
		Addr: httpAddr,
		Handler: httptransport.NewRouter(httptransport.Dependencies{
			Trips:      tripService,
			Drivers:    driverService,
			Payments:   paymentService,
			Websockets: websocketHandler,
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
