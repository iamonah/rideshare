package trip

import (
	"context"

	"github.com/iamonah/rideshare/shared/errs"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *TripBusiness) CreateTrip(ctx context.Context, fare *RideFare) (*Trip, error) {
	t := &Trip{
		ID:       bson.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}

	createdTrip, err := s.repo.CreateTrip(ctx, t)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to create trip")
	}

	return createdTrip, nil
}
