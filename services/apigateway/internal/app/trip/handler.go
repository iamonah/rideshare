package trip

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	httpcommon "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/common"
	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/errs"
)

type Handler struct {
	upstream PreviewTripUpstream
}

func NewHandler(upstream PreviewTripUpstream) *Handler {
	return &Handler{upstream: upstream}
}

func (h *Handler) HandlePreview(w http.ResponseWriter, r *http.Request) {
	var reqBody PreviewTripInput
	if err := httpcommon.ReadJSON(r, &reqBody); err != nil {
		if writeErr := httpcommon.WriteAPIError(w, errs.New(errs.InvalidArgument, errors.New("failed to parse JSON data"))); writeErr != nil {
			log.Printf("failed to write preview trip invalid JSON response: %v", writeErr)
		}
		return
	}
	defer r.Body.Close()

	if err := errs.Validate(reqBody); err != nil {
		if writeErr := httpcommon.WriteAPIError(w, errs.New(errs.InvalidArgument, err)); writeErr != nil {
			log.Printf("failed to write preview trip validation error response: %v", writeErr)
		}
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	body, err := h.upstream.PreviewTrip(ctx, reqBody)
	if err != nil {
		log.Printf("failed to preview a trip: %v", err)
		if writeErr := httpcommon.WriteUpstreamGRPCError(w, "trip service", err); writeErr != nil {
			log.Printf("failed to write preview trip upstream error response: %v", writeErr)
		}
		return
	}

	payload, err := json.Marshal(body)
	if err != nil {
		if writeErr := httpcommon.WriteAPIError(w, errs.New(errs.Internal, errors.New("internal service error"))); writeErr != nil {
			log.Printf("failed to write preview trip internal error response: %v", writeErr)
		}
		return
	}

	if err := httpcommon.WriteJSON(w, http.StatusOK, contracts.APIResponse{Data: payload}); err != nil {
		log.Printf("failed to write preview trip success response: %v", err)
	}
}
