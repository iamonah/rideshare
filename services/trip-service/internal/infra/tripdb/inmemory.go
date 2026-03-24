package tripdb

import (
	"context"
	"fmt"

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

func (r *inmemRepository) CreateTrip(ctx context.Context, t *tripdomain.Trip) (*tripdomain.Trip, error) {
	r.trips[t.ID.Hex()] = t
	return t, nil
}

func (r *inmemRepository) GetTripByID(ctx context.Context, id string) (*tripdomain.Trip, error) {
	t, ok := r.trips[id]
	if !ok {
		return nil, nil
	}

	return t, nil
}

func (r *inmemRepository) SaveRideFare(ctx context.Context, fare *tripdomain.RideFare) error {
	r.rideFares[fare.ID.Hex()] = fare
	return nil
}

func (r *inmemRepository) GetRideFareByID(ctx context.Context, id string) (*tripdomain.RideFare, error) {
	fare, ok := r.rideFares[id]
	if !ok {
		return nil, fmt.Errorf("ride fare not found: %s", id)
	}

	return fare, nil
}
