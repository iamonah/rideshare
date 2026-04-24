package trip

import "context"

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *Trip) (*Trip, error)
	SaveRideFare(ctx context.Context, f *RideFare) error
	GetRideFareByID(ctx context.Context, id string) (*RideFare, error)
	GetTripByID(ctx context.Context, id string) (*Trip, error)
	// UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error
}
