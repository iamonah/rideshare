package trip

import (
	"context"
	"errors"
	"fmt"

	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Route struct {
	Routes []RouteSummary
}

type RouteSummary struct {
	Distance float64
	Duration float64
	Geometry RouteGeometry
}

type RouteGeometry struct {
	Coordinates [][]float64
}

func (s *TripBusiness) PreviewTrip(ctx context.Context, pickup, destination *types.Coordinate, userID string) (*PreviewTripResult, error) {
	route, err := s.GetRoute(ctx, pickup, destination)
	if err != nil {
		return nil, err
	}
	if route == nil || len(route.Routes) == 0 {
		return nil, errs.New(errs.NotFound, errors.New("no route found between pickup and destination"))
	}

	estimatedFares := s.EstimatePackagesWithRoutes(ctx, *route)
	fares, err := s.GenerateTripFares(ctx, estimatedFares, userID, &route.Routes[0])
	if err != nil {
		return nil, err
	}

	for _, fare := range fares {
		if err := s.repo.SaveRideFare(ctx, fare); err != nil {
			return nil, errs.New(errs.Internal, fmt.Errorf("failed to save trip fare: %w", err))
		}
	}

	return &PreviewTripResult{
		Route: route,
		Fares: fares,
	}, nil
}

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

func (s *TripBusiness) GenerateTripFares(ctx context.Context, rideFares []*RideFare, userID string, route *RouteSummary) ([]*RideFare, error) {
	_ = ctx

	fares := make([]*RideFare, 0, len(rideFares))

	for _, f := range rideFares {
		id := bson.NewObjectID()

		fare := &RideFare{
			UserID:            userID,
			ID:                id,
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
			Route:             route,
		}

		fares = append(fares, fare)
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
			PackageSlug:       PackageSlugSUV,
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       PackageSlugSedan,
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       PackageSlugVan,
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       PackageSlugLuxury,
			TotalPriceInCents: 1000,
		},
	}
}
