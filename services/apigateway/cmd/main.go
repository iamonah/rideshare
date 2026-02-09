package main

import (
	"log"
	"net/http"

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

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
