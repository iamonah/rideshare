package grpc_Handler

import (
	"errors"

	tripdomain "github.com/iamonah/rideshare/services/trip-service/internal/domain/trip"
	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/proto/pb/trippb"
)

func toProto(routeResp tripdomain.Route) (*trippb.Route, error) {
	if len(routeResp.Routes) == 0 {
		return nil, errs.New(errs.NotFound, errors.New("no route found"))
	}

	route := routeResp.Routes[0]
	geometry, err := mapOSRMGeometry(route.Geometry.Coordinates)
	if err != nil {
		return nil, err
	}

	return &trippb.Route{
		Distance: route.Distance,
		Duration: route.Duration,
		Geometry: geometry,
	}, nil
}

func mapOSRMGeometry(coords [][]float64) ([]*trippb.Geometry, error) {
	if len(coords) == 0 {
		return nil, errs.New(errs.Unavailable, errors.New("route service temporarily unavailable"))
	}
	grpcCoords := make([]*trippb.Coordinate, 0, len(coords))
	for _, pair := range coords {
		if len(pair) != 2 {
			return nil, errs.New(errs.Unavailable, errors.New("service temporarily unavailable"))
		}
		grpcCoords = append(grpcCoords, &trippb.Coordinate{
			Latitude:  pair[1],
			Longitude: pair[0],
		})
	}

	return []*trippb.Geometry{
		{Coordinates: grpcCoords},
	}, nil
}
