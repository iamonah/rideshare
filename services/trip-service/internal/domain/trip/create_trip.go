package trip

import (
	"context"
	"fmt"

	"github.com/iamonah/rideshare/shared/errs"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *TripBusiness) CreateTrip(ctx context.Context, userID string, rideFareID bson.ObjectID) (*Trip, error) {
	tripFare, err := s.repo.GetRideFareByID(ctx, rideFareID.String())
	if err != nil {
		return nil, errs.New(errs.NotFound, fmt.Errorf("ride fare not found: %s", rideFareID.Hex()))
	}

	if tripFare.UserID != userID {
		return nil, errs.New(errs.PermissionDenied, fmt.Errorf("user %s is not authorized to create trip with fare %s", userID, rideFareID.Hex()))
	}
	
	t := &Trip{
		ID:       bson.NewObjectID(),
		UserID:   userID,
		Status:   "pending",
		RideFare: tripFare,
		Driver: &AssignedDriverSnapshot{},
	}

	createdTrip, err := s.repo.CreateTrip(ctx, t)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to create trip")
	}

	//publish trip created event
	err = s.eventPublisher.PublishTripCreated(ctx, fmt.Sprintf("Trip created with ID: %s for user: %s", createdTrip.ID.Hex(), userID))
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to publish trip created event")
	}

	return createdTrip, nil
}
