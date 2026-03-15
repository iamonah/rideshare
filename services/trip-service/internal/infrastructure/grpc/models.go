package grpc_Handler

import (
	"fmt"

	"github.com/iamonah/rideshare/services/trip-service/internal/service"
	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/proto/pb/trippb"
)

func toProto(os service.OsrmApiResponse) (*trippb.Route, error) {
	if len(os.Routes) == 0 {
		return nil, errs.New(errs.NotFound, "no route found")
	}

	route := os.Routes[0]
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
		return nil, errs.Wrap(
			errs.Unavailable,
			"route service temporarily unavailable",
			fmt.Errorf("osrm returned route geometry with no coordinates"),
		)
	}
	grpcCoords := make([]*trippb.Coordinate, 0, len(coords))
	for _, pair := range coords {
		if len(pair) != 2 {
			return nil, errs.Wrap(
				errs.Unavailable,
				"service temporarily unavailable",
				fmt.Errorf("osrm returned invalid geometry pair length: %d", len(pair)),
			)
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
