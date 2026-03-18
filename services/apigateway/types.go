package apigateway

import (
	"github.com/iamonah/rideshare/shared/proto/pb/trippb"
	"github.com/iamonah/rideshare/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userID" validate:"required"`
	Pickup      types.Coordinate `json:"pickup" validate:"required"`
	Destination types.Coordinate `json:"destination" validate:"required"`
}

func (req *previewTripRequest) toProto() *trippb.PreviewTripRequest {
	return &trippb.PreviewTripRequest{
		UserId: req.UserID,
		Pickup: &trippb.Coordinate{
			Latitude:  req.Pickup.Latitude,
			Longitude: req.Pickup.Longitude,
		},
		Destination: &trippb.Coordinate{
			Latitude:  req.Destination.Latitude,
			Longitude: req.Destination.Longitude,
		},
	}
}

type startTripRequest struct {
	RideFareID string `json:"rideFareID" validate:"required"`
	UserID     string `json:"userID" validate:"required"`
}

func (c *startTripRequest) toProto() *trippb.CreateTripRequest {
	return &trippb.CreateTripRequest{
		RideFareId: c.RideFareID,
		UserId:     c.UserID,
	}
}
