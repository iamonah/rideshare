package trip

import "github.com/iamonah/rideshare/shared/types"

type PreviewTripInput struct {
	UserID      string           `json:"userId" validate:"required"`
	Pickup      types.Coordinate `json:"pickup" validate:"required"`
	Destination types.Coordinate `json:"destination" validate:"required"`
}

type CreateTripInput struct {
	RideFareID string `json:"rideFareId" validate:"required"`
	UserID     string `json:"userId" validate:"required"`
}

type PreviewTripOutput struct {
	TripID    string            `json:"tripId,omitempty"`
	Route     Route             `json:"route"`
	RideFares []PreviewRideFare `json:"rideFares"`
}

type CreateTripOutput struct {
	TripID string `json:"tripId"`
	Trip   *Trip  `json:"trip,omitempty"`
}

type Trip struct {
	ID           string      `json:"id"`
	SelectedFare *RideFare   `json:"selectedFare,omitempty"`
	Route        Route       `json:"route,omitempty"`
	Status       string      `json:"status"`
	UserID       string      `json:"userId"`
	Driver       *TripDriver `json:"driver,omitempty"`
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
	Route             Route   `json:"route"`
}

type RideFare = PreviewRideFare

type TripDriver struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
	CarPlate       string `json:"carPlate"`
}
