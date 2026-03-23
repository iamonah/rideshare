package trip

import "context"

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *Trip) (*Trip, error)
}
