package trip

import (
	"context"

	"github.com/iamonah/rideshare/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PreviewTripResult struct {
	Route *Route
	Fares []*RideFare
}

type ExtTripBusiness interface {
	PreviewTrip(ctx context.Context, pickup, destination *types.Coordinate, userID string) (*PreviewTripResult, error)
	CreateTrip(ctx context.Context, userID string, rideFareID bson.ObjectID) (*Trip, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*Route, error)
	EstimatePackagesWithRoutes(ctx context.Context, route Route) []*RideFare
	GenerateTripFares(ctx context.Context, fare []*RideFare, userID string, route *RouteSummary) ([]*RideFare, error)
}

type RouteProvider interface {
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*Route, error)
}

type TripBusiness struct {
	repo           TripRepository
	route          RouteProvider
	eventPublisher TripEventPublisher
}

type TripEventPublisher interface {
	PublishTripCreated(ctx context.Context, data string) error
}

func NewTripBusiness(repo TripRepository, route RouteProvider, eventPublisher TripEventPublisher) *TripBusiness {
	return &TripBusiness{
		repo:           repo,
		route:          route,
		eventPublisher: eventPublisher,
	}
}
