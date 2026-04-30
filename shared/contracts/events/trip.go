package events

import "github.com/iamonah/rideshare/shared/types"

type TripFareSnapshot struct {
	ID          string `json:"id"`
	PackageSlug string `json:"packageSlug"`
}

type TripRouteSnapshot struct {
	Geometry []types.Coordinate `json:"geometry"`
}

type AssignedDriverSnapshot struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	ProfilePicture string            `json:"profilePicture"`
	CarPlate       string            `json:"carPlate"`
	PackageSlug    string            `json:"packageSlug"`
	GeoHash        string            `json:"geoHash"`
	Location       *types.Coordinate `json:"location,omitempty"`
}

// TripCreatedEvent is published by trip-service when a rider creates a trip.
type TripCreatedEvent struct {
	TripID          string                  `json:"tripId"`
	UserID          string                  `json:"userId"`
	Status          string                  `json:"status"`
	Fare            TripFareSnapshot        `json:"fare"`
	Pickup          types.Coordinate        `json:"pickup"`
	Dropoff         types.Coordinate        `json:"dropoff"`
	DistanceMeters  float64                 `json:"distanceMeters"`
	DurationSeconds float64                 `json:"durationSeconds"`
	Route           TripRouteSnapshot       `json:"route"`
	Driver          *AssignedDriverSnapshot `json:"driver,omitempty"`
}
