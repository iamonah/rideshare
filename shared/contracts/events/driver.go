package events

import (
	"github.com/iamonah/rideshare/shared/types"
)

type DriverTripResponseData struct {
	TripID  string `json:"tripID"`
	RiderID string `json:"riderID"`
	Driver  struct {
		Id             string            `json:"id"`
		Name           string            `json:"name"`
		ProfilePicture string            `json:"profilePicture"`
		CarPlate       string            `json:"carPlate"`
		GeoHash        string            `json:"geoHash"`
		GeoHashAlt     string            `json:"geohash"`
		PackageSlug    string            `json:"packageSlug"`
		Location       *types.Coordinate `json:"location,omitempty"`
	} `json:"driver"`
}
