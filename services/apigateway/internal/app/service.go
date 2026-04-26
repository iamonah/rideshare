package app

import (
	"context"
)

type TripUpstream interface {
	PreviewTrip(ctx context.Context, input PreviewTripInput) (*PreviewTripOutput, error)
	CreateTrip(ctx context.Context, input CreateTripInput) (*CreateTripOutput, error)
	Close() error
}
