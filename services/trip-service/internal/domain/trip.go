package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TripModel struct {
	ID       bson.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   string         `bson:"user_id" json:"user_id"`
	Status   string         `bson:"status" json:"status"`
	RideFare *RideFareModel `bson:"ride_fare" json:"ride_fare"`
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
}
