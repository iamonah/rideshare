package apigateway

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/iamonah/rideshare/shared/contracts"
)

func HandleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	//type expected to be send by front end for service forwarding
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// validation
	if reqBody.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	// Call trip service
	payload, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, "failed to encode request", http.StatusBadRequest)
		return
	}

	tripServiceURL := os.Getenv("TRIP_SERVICE_URL")
	if tripServiceURL == "" {
		http.Error(w, "trip service not configured", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		tripServiceURL+"/preview",
		bytes.NewReader(payload),
	)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "failed to call trip service", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		http.Error(w, "trip service error", http.StatusBadGateway)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read trip service response", http.StatusBadGateway)
		return
	}

	gatewayResp := contracts.APIResponse{
		Data: json.RawMessage(body),
	}

	writeJSON(w, http.StatusOK, gatewayResp)
}
