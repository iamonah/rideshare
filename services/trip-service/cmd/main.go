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

	"github.com/iamonah/rideshare/services/trip-service/internal/domain"
	"github.com/iamonah/rideshare/services/trip-service/internal/infrastructure/repository"
)

func main() {
	log.Println("--- Trip Service Initializing... ---")
	inmemRepo := repository.NewInmemRepository()
	svc := domain.NewService(inmemRepo)
	app := domain.NewHttpHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodPost+" /preview", app.HandleTripPreview)
	server := http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	httpError := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			httpError <- err
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("Signal received: %s", (<-signalChan).String())

	select {
	case <-httpError:
		log.Fatal("HTTP server closed unexpectedly")
	case <-signalChan:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			_ = server.Close() // Force close if graceful shutdown fails
			log.Fatal("graceful shutdown failed: %w", err)
		}
	}
}
