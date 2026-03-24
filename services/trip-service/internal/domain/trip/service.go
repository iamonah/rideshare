package trip

import (
	"context"

	"github.com/iamonah/rideshare/shared/types"
)

type PreviewTripResult struct {
	Route *Route
	Fares []*RideFare
}

type ExtTripBusiness interface {
	PreviewTrip(ctx context.Context, pickup, destination *types.Coordinate, userID string) (*PreviewTripResult, error)
	CreateTrip(ctx context.Context, fare *RideFare) (*Trip, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*Route, error)
	EstimatePackagesWithRoutes(ctx context.Context, route Route) []*RideFare
	GenerateTripFares(ctx context.Context, fare []*RideFare, userID string) ([]*RideFare, error)
}

type RouteProvider interface {
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*Route, error)
}

type TripBusiness struct {
	repo  TripRepository
	route RouteProvider
}

func NewTripBusiness(repo TripRepository, route RouteProvider) *TripBusiness {
	return &TripBusiness{
		repo:  repo,
		route: route,
	}
}
