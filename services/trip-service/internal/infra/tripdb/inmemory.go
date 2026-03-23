package repository

import (
	"context"

	tripdomain "github.com/iamonah/rideshare/services/trip-service/internal/domain/trip"
)

type inmemRepository struct {
	trips     map[string]*tripdomain.Trip
	rideFares map[string]*tripdomain.RideFare
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*tripdomain.Trip),
		rideFares: make(map[string]*tripdomain.RideFare),
	}
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *tripdomain.Trip) (*tripdomain.Trip, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}
