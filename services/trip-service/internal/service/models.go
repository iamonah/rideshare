package service

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RideFareModel struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID            string        `bson:"user_id" json:"user_id"`
	PackageSlug       string        `bson:"package_slug" json:"package_slug"` // ex: van, luxury, sedan
	TotalPriceInCents float64       `bson:"total_price_in_cents" json:"total_price_in_cents"`
}
type OsrmApiResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}
