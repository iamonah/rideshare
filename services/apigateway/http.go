package apigateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/iamonah/rideshare/services/apigateway/grpc_client"
	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/proto/pb/trip"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type HandlerApiGateway struct {
	tripclient *grpc_client.TripClient
}

func NewHandlerApiGateway(tripclient *grpc_client.TripClient) *HandlerApiGateway {
	return &HandlerApiGateway{tripclient: tripclient}
}

func (h *HandlerApiGateway) HandleTripPreview(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	//grpc_client call
	grpcResp, err := h.tripclient.Client.PreviewTrip(ctx, &trip.PreviewTripRequest{
		UserId: reqBody.UserID,
		StartLocation: &trip.Coordinate{
			Latitude:  reqBody.Pickup.Latitude,
			Longitude: reqBody.Pickup.Longitude,
		},
		EndLocation: &trip.Coordinate{
			Latitude:  reqBody.Destination.Latitude,
			Longitude: reqBody.Destination.Longitude,
		},
	})
	if err != nil {
		writeTripError(w, err)
		return
	}

	body, err := protojson.Marshal(grpcResp)
	if err != nil {
		http.Error(w, "failed to encode trip response", http.StatusInternalServerError)
		return
	}

	gatewayResp := contracts.APIResponse{
		Data: body,
	}

	writeJSON(w, http.StatusOK, gatewayResp)
}

func writeTripError(w http.ResponseWriter, err error) {
	httpStatus := http.StatusBadGateway
	apiErr := &contracts.APIError{
		Code:    "TRIP_UPSTREAM_ERROR",
		Message: "failed to call trip service",
	}

	if errors.Is(err, context.DeadlineExceeded) {
		httpStatus = http.StatusGatewayTimeout
		apiErr.Code = "TRIP_TIMEOUT"
		apiErr.Message = "trip service request timed out"
		writeJSON(w, httpStatus, contracts.APIResponse{Error: apiErr})
		return
	}

	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
			apiErr.Code = "TRIP_INVALID_ARGUMENT"
		case codes.NotFound:
			httpStatus = http.StatusNotFound
			apiErr.Code = "TRIP_NOT_FOUND"
		case codes.DeadlineExceeded:
			httpStatus = http.StatusGatewayTimeout
			apiErr.Code = "TRIP_TIMEOUT"
		case codes.Unavailable:
			httpStatus = http.StatusServiceUnavailable
			apiErr.Code = "TRIP_UNAVAILABLE"
		}
		if st.Message() != "" {
			apiErr.Message = st.Message()
		}
	}

	writeJSON(w, httpStatus, contracts.APIResponse{Error: apiErr})
}
