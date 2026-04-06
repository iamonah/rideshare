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
	PackageSlug       PackageSlug   `bson:"package_slug" json:"package_slug"` // ex: van, luxury, sedan
	TotalPriceInCents float64       `bson:"total_price_in_cents" json:"total_price_in_cents"`
	Route             *RouteSummary `bson:"route" json:"route"` //driver visualization
}

func (r *RideFare) ToProto() *trippb.RideFare {
	if r == nil {
		return nil
	}

	return &trippb.RideFare{
		Id:                r.ID.Hex(),
		UserId:            r.UserID,
		PackageSlug:       r.PackageSlug.String(),
		TotalPriceInCents: r.TotalPriceInCents,
		Route:             routeSummaryToProto(r.Route),
	}
}

func routeSummaryToProto(route *RouteSummary) *trippb.Route {
	if route == nil {
		return nil
	}

	coordinates := make([]*trippb.Coordinate, 0, len(route.Geometry.Coordinates))
	for _, pair := range route.Geometry.Coordinates {
		if len(pair) != 2 {
			continue
		}

		coordinates = append(coordinates, &trippb.Coordinate{
			Latitude:  pair[1],
			Longitude: pair[0],
		})
	}

	return &trippb.Route{
		Distance: route.Distance,
		Duration: route.Duration,
		Geometry: []*trippb.Geometry{
			{Coordinates: coordinates},
		},
	}
}

func ToRideFaresProto(fares []*RideFare) []*trippb.RideFare {
	var protoFares []*trippb.RideFare
	for _, f := range fares {
		protoFares = append(protoFares, f.ToProto())
	}
	return protoFares
}
