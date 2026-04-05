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
	if len(route.Geometry.Coordinates) == 0 {
		return nil, errs.New(errs.Unavailable, errors.New("route service temporarily unavailable"))
	}
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

func faresToProto(fares []*tripdomain.RideFare) []*trippb.RideFare {
	protoFares := make([]*trippb.RideFare, len(fares))
	for i, f := range fares {
		protoFares[i] = f.ToProto()
	}

	return protoFares
}

func toProtoTrip(trip *tripdomain.Trip) *trippb.Trip {
	if trip == nil {
		return nil
	}

	protoTrip := &trippb.Trip{
		Id:     trip.ID.Hex(),
		UserId: trip.UserID,
		Status: trip.Status,
		SelectedFare: &trippb.RideFare{
			Id:                trip.RideFare.ID.Hex(),
			UserId:            trip.RideFare.UserID,
			PackageSlug:       trip.RideFare.PackageSlug.String(),
			TotalPriceInCents: trip.RideFare.TotalPriceInCents,
		},
	}

	return protoTrip
}
