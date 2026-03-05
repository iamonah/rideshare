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

	"github.com/gorilla/mux"
	"github.com/iamonah/rideshare/services/apigateway"
	"github.com/iamonah/rideshare/services/apigateway/grpc_client"
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

	mux := mux.NewRouter()

	if tripServiceAddr == "" {
		log.Fatal("TRIP_SERVICE_GRPC_URL is required")
	}

	log.Printf("Trip gRPC client initialized for %s (non-blocking dial)", tripServiceAddr)
	tripClient, err := grpc_client.NewTripClient(tripServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer tripClient.Close()

	apigatewayHandler := apigateway.NewHandlerApiGateway(tripClient)
	mux.HandleFunc("/trip/preview", apigatewayHandler.HandleTripPreview).Methods("POST")
	mux.HandleFunc("/ws/drivers", apigatewayHandler.HandleDriversWebsocket).Methods("GET")
	mux.HandleFunc("/ws/riders", apigatewayHandler.HandleRidersWebsocket).Methods("GET")

	log.Println("Listening on", httpAddr)

	muxHandler := apigatewayHandler.WithCORS(mux)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: muxHandler,
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
