package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/iamonah/rideshare/shared/types"
)

type HttpHandler struct {
	Service TripService
}

func NewHttpHandler(service TripService) *HttpHandler {
	return &HttpHandler{
		Service: service,
	}
}

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (s *HttpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	fmt.Println("Received trip preview request:", reqBody)

	fare := &RideFareModel{
		UserID: reqBody.UserID,
	}

	ctx := r.Context()

	t, err := s.Service.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Created trip:", t)

	writeJSON(w, http.StatusOK, t)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
