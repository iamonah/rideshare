package trip

import (
	"github.com/iamonah/rideshare/shared/proto/pb/trippb"
	"go.mongodb.org/mongo-driver/v2/bson"
)


type PricingConfig struct {
	PricePerUnitOfDistance float64
	PricingPerMinute       float64
}

func DefaultPricingConfig() *PricingConfig {
	return &PricingConfig{
		PricePerUnitOfDistance: 1.5,
		PricingPerMinute:       0.25,
	}
} 

type RideFare struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID            string        `bson:"user_id" json:"user_id"`
	PackageSlug       string        `bson:"package_slug" json:"package_slug"` // ex: van, luxury, sedan
	TotalPriceInCents float64       `bson:"total_price_in_cents" json:"total_price_in_cents"`
}

func (r *RideFare) ToProto() *trippb.RideFare {
	return &trippb.RideFare{
		Id:                r.ID.Hex(),
		UserId:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(fares []*RideFare) []*trippb.RideFare {
	var protoFares []*trippb.RideFare
	for _, f := range fares {
		protoFares = append(protoFares, f.ToProto())
	}
	return protoFares
}
