package business

import (
	"context"
	"fmt"

	"github.com/iamonah/rideshare/shared/pb/trip"
	"github.com/iamonah/rideshare/shared/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TripService struct {
	trips ExTripBusiness
	trip.UnimplementedTripServiceServer
}

func NewTripService(trips ExTripBusiness) *TripService {
	return &TripService{
		trips: trips,
	}
}

func (s *TripService) PreviewTrip(ctx context.Context, req *trip.PreviewTripRequest) (
	*trip.PreviewTripResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if req.GetStartLocation() == nil || req.GetEndLocation() == nil {
		return nil, status.Error(codes.InvalidArgument, "start_location and end_location are required")
	}

	route, err := s.trips.GetRoute(ctx, &types.Coordinate{
		Latitude:  req.StartLocation.Latitude,
		Longitude: req.StartLocation.Longitude,
	}, &types.Coordinate{
		Latitude:  req.EndLocation.Latitude,
		Longitude: req.EndLocation.Longitude,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}
	if len(route.Routes) == 0 {
		return nil, status.Error(codes.NotFound, "no route found")
	}

	first := route.Routes[0]
	geometry, err := mapOSRMGeometry(first.Geometry.Coordinates)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map route geometry: %v", err)
	}

	return &trip.PreviewTripResponse{
		Route: &trip.Route{
			Distance: first.Distance,
			Duration: first.Duration,
			Geometry: geometry,
		},
	}, nil
}

func mapOSRMGeometry(coords [][]float64) ([]*trip.Geometry, error) {
	geoCoords := make([]*trip.Coordinate, 0, len(coords))
	for _, c := range coords {
		if len(c) < 2 {
			return nil, fmt.Errorf("invalid coordinate length: %d", len(c))
		}
		geoCoords = append(geoCoords, &trip.Coordinate{
			Latitude:  c[1],
			Longitude: c[0],
		})
	}

	return []*trip.Geometry{
		{Coordinates: geoCoords},
	}, nil
}
