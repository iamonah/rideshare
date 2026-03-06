package repository

import (
	"context"

	service "github.com/iamonah/rideshare/services/trip-service/internal/service"
)

type inmemRepository struct {
	trips     map[string]*service.TripModel
	rideFares map[string]*service.RideFareModel
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		trips:     make(map[string]*service.TripModel),
		rideFares: make(map[string]*service.RideFareModel),
	}
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *service.TripModel) (*service.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}
