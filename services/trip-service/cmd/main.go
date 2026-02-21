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

	"github.com/iamonah/rideshare/services/trip-service/internal/business"
	"github.com/iamonah/rideshare/services/trip-service/internal/infrastructure/repository"
)

func main() {
	log.Println("--- Trip Service Initializing... ---")
	inmemRepo := repository.NewInmemRepository()
	svc := business.NewService(inmemRepo)
	app := business.NewHttpHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodPost+" /preview", app.HandleTripPreview)
	server := http.Server{
		Addr:    ":8083",
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
