package trip

import (
	"context"
)

type PreviewTripUpstream interface {
	PreviewTrip(ctx context.Context, input PreviewTripInput) (*PreviewTripOutput, error)
}

type Service struct {
	upstream PreviewTripUpstream
}

func NewService(upstream PreviewTripUpstream) *Service {
	return &Service{upstream: upstream}
}

func (s *Service) PreviewTrip(ctx context.Context, input PreviewTripInput) (*PreviewTripOutput, error) {
	return s.upstream.PreviewTrip(ctx, input)
}
