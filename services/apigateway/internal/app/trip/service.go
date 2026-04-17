package trip

import (
	"context"
)

type Upstream interface {
	PreviewTrip(ctx context.Context, input PreviewTripInput) (*PreviewTripOutput, error)
	CreateTrip(ctx context.Context, input CreateTripInput) (*CreateTripOutput, error)
}
