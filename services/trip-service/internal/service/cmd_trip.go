package service

import (
	"context"
	"encoding/json"
	"errors"
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
		return nil, errs.Newf(errs.Internal, err, "failed to create trip")
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
		return nil, errs.Newf(errs.Internal, err, "failed to build route request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errs.Newf(errs.Unavailable, err, "route provider unavailable") //how will this error message look like when sent to the client
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to read route response")
	}

	if resp.StatusCode != http.StatusOK {
		var routeErr OsrmErrorResponse
		if err := json.Unmarshal(body, &routeErr); err != nil {
			return nil, errs.Newf(errs.Internal, err, "json.Unmarshal route provider error response (status=%d)", resp.StatusCode)
		}

		return nil, classifyOSRMError(resp.StatusCode, &routeErr)
	}

	var routeResp OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to parse route response")
	}

	return &routeResp, nil
}

func classifyOSRMError(statusCode int, routeErr *OsrmErrorResponse) error {
	code := ""
	if routeErr != nil {
		code = routeErr.Code
	}

	switch code {
	case "NoRoute":
		return errs.New(errs.NotFound, errors.New("no route found between pickup and destination"))
	case "NoSegment":
		return errs.New(errs.InvalidArgument, errors.New("pickup or destination is not on a routable road"))
	case "InvalidValue", "InvalidQuery", "InvalidOptions", "InvalidUrl":
		return errs.New(errs.InvalidArgument, errors.New("route request is invalid"))
	}

	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		return errs.New(errs.InvalidArgument, errors.New("route request is invalid"))
	}

	return errs.Newf(
		errs.Unavailable,
		fmt.Errorf("route provider returned status %d with code %q", statusCode, code),
		"route provider unavailable",
	)
}
