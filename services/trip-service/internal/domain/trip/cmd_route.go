package trip

import (
	"context"
	"errors"
	"fmt"

	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *TripBusiness) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*Route, error) {
	if s.route == nil {
		return nil, errs.New(errs.Internal, errors.New("route provider is not configured"))
	}

	return s.route.GetRoute(ctx, pickup, destination)
}


func (s *TripBusiness) EstimatePackagesWithRoutes(ctx context.Context, route Route) []*RideFare {
	baseFares := getBaseFares()
	estimatedFares := make([]*RideFare, len(baseFares))

	for i, f := range baseFares {
		estimatedFares[i] = estimateFareRoute(f, &route)
	}

	return estimatedFares
}

func (s *TripBusiness) GenerateTripFares(ctx context.Context, rideFares []*RideFare, userID string) ([]*RideFare, error) {
	fares := make([]*RideFare, len(rideFares))

	for i, f := range rideFares {
		id := bson.NewObjectID()

		fare := &RideFare{
			UserID:            userID,
			ID:                id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
		}

		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, errs.New(errs.Internal, fmt.Errorf("failed to save trip fare: %w", err))
		}

		fares[i] = fare
	}

	return fares, nil
}

func estimateFareRoute(f *RideFare, route *Route) *RideFare {
	pricingCfg := DefaultPricingConfig()
	carPackagePrice := f.TotalPriceInCents

	distanceKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	// total_fare = base_fare + distance_charge + time_charge + platform_fee
	// driver_payout = total_fare - platform_commission

	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricingPerMinute
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &RideFare{
		TotalPriceInCents: totalPrice,
		PackageSlug:       f.PackageSlug,
	}
}

func getBaseFares() []*RideFare {
	return []*RideFare{
		{
			PackageSlug:       PackageSlugSUV.String(),
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       PackageSlugSedan.String(),
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       PackageSlugVan.String(),
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       PackageSlugLuxury.String(),
			TotalPriceInCents: 1000,
		},
	}
}
