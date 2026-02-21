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
	"github.com/iamonah/rideshare/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway...")

	mux := mux.NewRouter()

	mux.HandleFunc("/trip/preview", apigateway.HandleTripPreview).Methods("POST")
	mux.HandleFunc("/ws/drivers", apigateway.HandleDriversWebsocket).Methods("GET")
	mux.HandleFunc("/ws/riders", apigateway.HandleRidersWebsocket).Methods("GET")

	log.Println("Listening on", httpAddr)

	muxHandler := apigateway.WithCORS(mux)
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
