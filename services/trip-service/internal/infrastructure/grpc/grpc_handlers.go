package grpc_Handler

import (
	"context"

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
		return nil, grpcerrs.ToStatus(errs.New(errs.InvalidArgument, "request is required"))
	}

	fieldErrs := errs.NewFieldErrors()
	if req.GetUserId() == "" {
		fieldErrs.AddMessage("user_id", "is required")
	}
	if req.GetStartLocation() == nil {
		fieldErrs.AddMessage("start_location", "is required")
	}
	if req.GetEndLocation() == nil {
		fieldErrs.AddMessage("end_location", "is required")
	}

	pickup := &types.Coordinate{
		Latitude:  req.GetStartLocation().GetLatitude(),
		Longitude: req.GetStartLocation().GetLongitude(),
	}

	destination := &types.Coordinate{
		Latitude:  req.GetEndLocation().GetLatitude(),
		Longitude: req.GetEndLocation().GetLongitude(),
	}

	validateRouteCoordinates(fieldErrs, pickup, destination)
	if err := fieldErrs.ToError(); err != nil {
		return nil, grpcerrs.ToStatus(err)
	}
	route, err := s.trips.GetRoute(ctx, pickup, destination)
	if err != nil {
		//log error
		return nil, grpcerrs.ToStatus(err)
	}
	protoRoute, err := toProto(*route)
	if err != nil {
		return nil, grpcerrs.ToStatus(err)
	}

	return &trippb.PreviewTripResponse{
		Route: protoRoute,
	}, nil
}

func validateRouteCoordinates(fieldErrs errs.FieldErrors, pickup, destination *types.Coordinate) error {
	validateCoordinate(fieldErrs, "start_location", pickup)
	validateCoordinate(fieldErrs, "end_location", destination)
	if len(fieldErrs) > 0 {
		return errs.Validation(fieldErrs)
	}

	return nil
}

func validateCoordinate(fieldErrs errs.FieldErrors, name string, coord *types.Coordinate) {
	if coord.Latitude < -90 || coord.Latitude > 90 {
		fieldErrs.AddMessage(name+".latitude", "must be between -90 and 90")
	}
	if coord.Longitude < -180 || coord.Longitude > 180 {
		fieldErrs.AddMessage(name+".longitude", "must be between -180 and 180")
	}
}
