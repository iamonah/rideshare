package main

import (
	"log"
	"net/http"

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

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
