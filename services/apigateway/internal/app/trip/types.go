package trip

import "github.com/iamonah/rideshare/shared/types"

type PreviewTripInput struct {
	UserID      string           `json:"userId" validate:"required"`
	Pickup      types.Coordinate `json:"pickup" validate:"required"`
	Destination types.Coordinate `json:"destination" validate:"required"`
}

type PreviewTripOutput struct {
	TripID    string            `json:"tripId,omitempty"`
	Route     Route             `json:"route"`
	RideFares []PreviewRideFare `json:"rideFares"`
}

type Route struct {
	Distance float64     `json:"distance"`
	Duration float64     `json:"duration"`
	Geometry []*Geometry `json:"geometry"`
}

type Geometry struct {
	Coordinates []*types.Coordinate `json:"coordinates"`
}

type PreviewRideFare struct {
	ID                string  `json:"id"`
	UserID            string  `json:"userId"`
	PackageSlug       string  `json:"packageSlug"`
	TotalPriceInCents float64 `json:"totalPriceInCents"`
}
