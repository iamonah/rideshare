package trip

import (
	"context"
)

type PreviewTripUpstream interface {
	PreviewTrip(ctx context.Context, input PreviewTripInput) (*PreviewTripOutput, error)
}
