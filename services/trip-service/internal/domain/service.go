package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type service struct {
	repo TripRepository
}

func NewService(repo TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error) {
	t := &TripModel{
		ID:       bson.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}

	return s.repo.CreateTrip(ctx, t)
}
