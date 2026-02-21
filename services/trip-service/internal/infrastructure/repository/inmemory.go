package repository

import (
	"context"

	"github.com/iamonah/rideshare/services/trip-service/internal/business"
)

type inmemRepository struct {
	trips     map[string]*business.TripModel
	rideFares map[string]*business.RideFareModel
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*business.TripModel),
		rideFares: make(map[string]*business.RideFareModel),
	}
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *business.TripModel) (*business.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}
