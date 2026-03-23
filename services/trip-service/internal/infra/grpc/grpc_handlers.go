package grpc_Handler

import (
	"context"
	"errors"

	tripdomain "github.com/iamonah/rideshare/services/trip-service/internal/domain/trip"
	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/errs/grpcerrs"
	"github.com/iamonah/rideshare/shared/proto/pb/trippb"
	"github.com/iamonah/rideshare/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc"
)

type TripService struct {
	trips tripdomain.ExtTripBusiness
	trippb.UnimplementedTripServiceServer
}

type previewTripInput struct {
	UserID      string            `json:"user_id" validate:"required"`
	Pickup      *types.Coordinate `json:"pickup" validate:"required"`
	Destination *types.Coordinate `json:"destination" validate:"required"`
}

func NewTripServer(server *grpc.Server, trips tripdomain.ExtTripBusiness) *TripService {
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

type createTripInput struct {
	RideFareID string `json:"ride_fare_id" validate:"required"`
	UserID     string `json:"user_id" validate:"required"`
}

func (s *TripService) CreateTrip(ctx context.Context, req *trippb.CreateTripRequest) (
	*trippb.CreateTripResponse, error) {
	if req == nil {
		return nil, grpcerrs.ToStatus(errs.New(errs.InvalidArgument, errors.New("request is required")))
	}

	input := createTripInput{
		RideFareID: req.GetRideFareId(),
		UserID:     req.GetUserId(),
	}

	if err := errs.Validate(input); err != nil {
		return nil, grpcerrs.ToStatus(errs.New(errs.InvalidArgument, err))
	}

	rideFareID, err := bson.ObjectIDFromHex(input.RideFareID)
	if err != nil {
		return nil, grpcerrs.ToStatus(errs.New(errs.InvalidArgument, errors.New("ride_fare_id must be a valid id")))
	}

	createTrip, err := s.trips.CreateTrip(ctx, &tripdomain.RideFare{
		ID:     rideFareID,
		UserID: input.UserID,
	})
	if err != nil {
		return nil, grpcerrs.ToStatus(err)
	}

	return &trippb.CreateTripResponse{
		TripId: createTrip.ID.Hex(),
	}, nil
}
