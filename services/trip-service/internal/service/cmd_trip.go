package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ExTripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*OsrmApiResponse, error)
}

func (s *tripService) CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error) {
	t := &TripModel{
		ID:       bson.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}

	createdTrip, err := s.repo.CreateTrip(ctx, t)
	if err != nil {
		return nil, errs.Wrap(errs.Internal, "failed to create trip", err)
	}

	return createdTrip, nil
}

func (s *tripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errs.Wrap(errs.Internal, "failed to build route request", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errs.Wrap(errs.Unavailable, "route provider unavailable", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.New(errs.Unavailable, "route provider unavailable")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Wrap(errs.Internal, "failed to read route response", err)
	}

	var routeResp OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, errs.Wrap(errs.Internal, "failed to parse route response", err)
	}

	return &routeResp, nil
}
