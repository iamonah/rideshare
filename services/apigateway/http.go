package apigateway

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/iamonah/rideshare/services/apigateway/grpc_client"
	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/errs"
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
		writeAPIError(w, errs.New(errs.InvalidArgument, errors.New("failed to parse JSON data")))
		return
	}

	defer r.Body.Close()

	if err := errs.Validate(reqBody); err != nil {
		writeAPIError(w, errs.New(errs.InvalidArgument, err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	//grpc_client call
	grpcResp, err := h.tripclient.Client.PreviewTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("failed to preview a trip: %v", err)
		writeUpstreamGRPCError(w, "trip service", err)
		return
	}

	body, err := protojson.Marshal(grpcResp)
	if err != nil {
		writeAPIError(w, errs.New(errs.Internal, errors.New("internal service error")))
		return
	}

	gatewayResp := contracts.APIResponse{
		Data: body,
	}

	writeJSON(w, http.StatusOK, gatewayResp)
}
