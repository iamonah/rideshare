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
	GetTripByID(ctx context.Context, id string) (*Trip, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *AssignedDriverSnapshot) error
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
	PublishTripCreated(ctx context.Context, trip *Trip) error
}

func NewTripBusiness(repo TripRepository, route RouteProvider, eventPublisher TripEventPublisher) *TripBusiness {
	return &TripBusiness{
		repo:           repo,
		route:          route,
		eventPublisher: eventPublisher,
	}
}

func (tb *TripBusiness) GetTripByID(ctx context.Context, id string) (*Trip, error) {
	return tb.repo.GetTripByID(ctx, id)
}

func (tb *TripBusiness) UpdateTrip(ctx context.Context, tripID string, status string, driver *AssignedDriverSnapshot) error {
	return tb.repo.UpdateTrip(ctx, tripID, status, driver)
}
