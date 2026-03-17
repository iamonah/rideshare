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
		StartLocation: &trippb.Coordinate{
			Latitude:  *req.Pickup.Latitude,
			Longitude: *req.Pickup.Longitude,
		},
		EndLocation: &trippb.Coordinate{
			Latitude:  *req.Destination.Latitude,
			Longitude: *req.Destination.Longitude,
		},
	}
}
