package trip

import (
	"context"
	"errors"

	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/types"
)

func (s *TripBusiness) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*Route, error) {
	if s.route == nil {
		return nil, errs.New(errs.Internal, errors.New("route provider is not configured"))
	}

	return s.route.GetRoute(ctx, pickup, destination)
}

// total_fare = base_fare + distance_charge + time_charge + platform_fee
// driver_payout = total_fare - platform_commission
func (s *TripBusiness) EstimatePackagesWithRoutes(ctx context.Context, route Route) ([]*RideFare, error) {
	return nil, nil
}

func (s *TripBusiness) GenerateTripFare(ctx context.Context, fare []*RideFare, userID string) ([]*RideFare, error) {
	return nil, nil
}
