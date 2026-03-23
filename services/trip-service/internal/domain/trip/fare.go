package trip

import "go.mongodb.org/mongo-driver/v2/bson"

type RideFare struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID            string        `bson:"user_id" json:"user_id"`
	PackageSlug       string        `bson:"package_slug" json:"package_slug"` // ex: van, luxury, sedan
	TotalPriceInCents float64       `bson:"total_price_in_cents" json:"total_price_in_cents"`
}
