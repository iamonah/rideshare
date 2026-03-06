package service

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TripModel struct {
	ID       bson.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   string         `bson:"user_id" json:"user_id"`
	Status   string         `bson:"status" json:"status"`
	RideFare *RideFareModel `bson:"ride_fare" json:"ride_fare"`
}
