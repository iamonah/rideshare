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

	"github.com/iamonah/rideshare/services/apigateway"
	"github.com/iamonah/rideshare/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway...")

	mux := http.NewServeMux()

	mux.HandleFunc(http.MethodPost+" /trip/preview", apigateway.HandleTripPreview)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	shutDown := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			shutDown <- err
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("Signal received: %s", (<-signalChan).String())

	select {
	case <-shutDown:
		log.Fatal("HTTP server closed unexpectedly")
	case <-signalChan:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("graceful shutdown failed: %s", err)
			_ = server.Close() // Force close if graceful shutdown fails
		}
	}
}
