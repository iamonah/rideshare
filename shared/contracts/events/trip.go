package events

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type TripFareSnapshot struct {
	ID          string `json:"id"`
	PackageSlug string `json:"packageSlug"`
}

type TripRouteSnapshot struct {
	Geometry []Coordinate `json:"geometry"`
}

type AssignedDriverSnapshot struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
	CarPlate       string `json:"carPlate"`
	PackageSlug    string `json:"packageSlug"`
}

// TripCreatedEvent is published by trip-service when a rider creates a trip.
type TripCreatedEvent struct {
	TripID          string                  `json:"tripId"`
	UserID          string                  `json:"userId"`
	Status          string                  `json:"status"`
	Fare            TripFareSnapshot        `json:"fare"`
	Pickup          Coordinate              `json:"pickup"`
	Dropoff         Coordinate              `json:"dropoff"`
	DistanceMeters  float64                 `json:"distanceMeters"`
	DurationSeconds float64                 `json:"durationSeconds"`
	Route           TripRouteSnapshot       `json:"route"`
	Driver          *AssignedDriverSnapshot `json:"driver,omitempty"`
}
