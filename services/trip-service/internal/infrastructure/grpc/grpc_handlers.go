package grpc_Handler

import (
	"context"
	"errors"

	"github.com/iamonah/rideshare/services/trip-service/internal/service"
	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/errs/grpcerrs"
	"github.com/iamonah/rideshare/shared/proto/pb/trippb"
	"github.com/iamonah/rideshare/shared/types"
	"google.golang.org/grpc"
)

type TripService struct {
	trips service.ExTripService
	trippb.UnimplementedTripServiceServer
}

type previewTripInput struct {
	UserID      string            `json:"user_id" validate:"required"`
	Pickup      *types.Coordinate `json:"pickup" validate:"required"`
	Destination *types.Coordinate `json:"destination" validate:"required"`
}

func NewTripServer(server *grpc.Server, trips service.ExTripService) *TripService {
	svc := &TripService{
		trips: trips,
	}
	trippb.RegisterTripServiceServer(server, svc)
	return svc
}

func (s *TripService) PreviewTrip(ctx context.Context, req *trippb.PreviewTripRequest) (
	*trippb.PreviewTripResponse, error) {
	if req == nil {
		return nil, grpcerrs.ToStatus(errs.New(errs.InvalidArgument, errors.New("request is required")))
	}

	input := previewTripInput{
		UserID: req.GetUserId(),
	}

	if req.GetPickup() != nil {
		input.Pickup = &types.Coordinate{
			Latitude:  req.GetPickup().GetLatitude(),
			Longitude: req.GetPickup().GetLongitude(),
		}
	}

	if req.GetDestination() != nil {
		input.Destination = &types.Coordinate{
			Latitude:  req.GetDestination().GetLatitude(),
			Longitude: req.GetDestination().GetLongitude(),
		}
	}

	if err := errs.Validate(input); err != nil {
		return nil, grpcerrs.ToStatus(errs.New(errs.InvalidArgument, err))
	}

	route, err := s.trips.GetRoute(ctx, input.Pickup, input.Destination)
	if err != nil {
		return nil, grpcerrs.ToStatus(err)
	}
	protoRoute, err := toProto(*route)
	if err != nil {
		return nil, grpcerrs.ToStatus(err)
	}

	return &trippb.PreviewTripResponse{
		Route:     protoRoute,
		RideFares: []*trippb.RideFare{},
	}, nil
}
