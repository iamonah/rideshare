package driverservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/iamonah/rideshare/services/driver-service/internal/infra/events"
	eventcontracts "github.com/iamonah/rideshare/shared/contracts/events"
	"github.com/iamonah/rideshare/shared/errs"
	driverpb "github.com/iamonah/rideshare/shared/proto/pb/driverpb"
	"github.com/iamonah/rideshare/shared/util"
	"github.com/mmcloughlin/geohash"

	mathrand "math/rand/v2"
)

type driverInMap struct {
	Driver *driverpb.Driver
	// Index int
	// TODO: route
}

type Service struct {
	mu           sync.RWMutex
	drivers      []*driverInMap
	tripConsumer events.TripCreatedHandler
}

func NewService(tripConsumer *events.TripConsumer) *Service {
	return &Service{
		drivers:      make([]*driverInMap, 0),
		tripConsumer: tripConsumer,
	}
}

func (s *Service) RegisterEventConsumers(ctx context.Context) error {
	if err := s.tripConsumer.HandleTripCreated(ctx, s.HandleTripCreated); err != nil {
		return fmt.Errorf("register trip created consumer: %w", err)
	}

	return nil
}

type PackageSlug string

var packageSlugs = map[string]PackageSlug{
	"van":    "van",
	"suv":    "suv",
	"sedan":  "sedan",
	"luxury": "luxury",
}

func parsePackageSlug(s string) (PackageSlug, bool) {
	slug, ok := packageSlugs[strings.ToLower(strings.TrimSpace(s))]
	return slug, ok
}

type Driver struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	ProfilePicture string      `json:"profilePicture"`
	CarPlate       string      `json:"carPlate"`
	PackageSlug    PackageSlug `json:"packageSlug"`
}

type TripRequest struct {
	TripID          string
	RiderID         string
	Status          string
	PackageSlug     string
	Pickup          Coordinate
	Dropoff         Coordinate
	Route           []Coordinate
	DistanceMeters  float64
	DurationSeconds float64
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

func (s *Service) RegisterDriver(req *driverpb.RegisterDriverRequest) (*driverpb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	driverID := req.GetDriverId()

	packageSlug, ok := parsePackageSlug(req.GetPackageSlug())
	if !ok {
		return nil, errs.New(errs.InvalidArgument, errors.New("unsupported package slug"))
	}

	randomRoute := PredefinedRoutes[mathrand.IntN(len(PredefinedRoutes))]
	latitude := randomRoute[0][0]
	longitude := randomRoute[0][1]

	driver := &driverpb.Driver{
		Id:             driverID,
		Name:           "Lando Norris",
		ProfilePicture: util.GetRandomAvatar(mathrand.IntN(10)),
		CarPlate:       GenerateRandomPlate(),
		GeoHash:        geohash.EncodeWithPrecision(latitude, longitude, 7),
		PackageSlug:    string(packageSlug),
		Location: &driverpb.Location{
			Latitude:  latitude,
			Longitude: longitude,
		},
	}

	s.drivers = append(s.drivers, &driverInMap{Driver: driver})
	return driver, nil
}

func (s *Service) HandleTripCreated(ctx context.Context, event *eventcontracts.TripCreatedEvent) error {
	_ = ctx

	request := tripCreatedEventToRequest(event)
	log.Printf("driver service consumed trip request: %+v", request)
	return nil
}

func tripCreatedEventToRequest(event *eventcontracts.TripCreatedEvent) TripRequest {
	if event == nil {
		return TripRequest{}
	}

	return TripRequest{
		TripID:      event.TripID,
		RiderID:     event.UserID,
		Status:      event.Status,
		PackageSlug: event.Fare.PackageSlug,
		Pickup: Coordinate{
			Latitude:  event.Pickup.Latitude,
			Longitude: event.Pickup.Longitude,
		},
		Dropoff: Coordinate{
			Latitude:  event.Dropoff.Latitude,
			Longitude: event.Dropoff.Longitude,
		},
		Route:           eventRouteToLocalCoordinates(event.Route.Geometry),
		DistanceMeters:  event.DistanceMeters,
		DurationSeconds: event.DurationSeconds,
	}
}

func eventRouteToLocalCoordinates(route []eventcontracts.Coordinate) []Coordinate {
	coordinates := make([]Coordinate, 0, len(route))
	for _, coordinate := range route {
		coordinates = append(coordinates, Coordinate{
			Latitude:  coordinate.Latitude,
			Longitude: coordinate.Longitude,
		})
	}

	return coordinates
}

func (s *Service) UnregisterDriver(req *driverpb.RegisterDriverRequest) (*driverpb.Driver, error) {
	driverID := req.GetDriverId()
	packageSlug := req.GetPackageSlug()

	s.mu.Lock()
	defer s.mu.Unlock()

	for index, driverEntry := range s.drivers {
		driver := driverEntry.GetDriver()
		if driver == nil {
			continue
		}
		if driver.GetId() != driverID || driver.GetPackageSlug() != packageSlug {
			continue
		}

		s.drivers = append(s.drivers[:index], s.drivers[index+1:]...)
		return driver, nil
	}

	return nil, errs.New(errs.NotFound, errors.New("driver not found"))
}

func (d *driverInMap) GetDriver() *driverpb.Driver {
	if d == nil {
		return nil
	}

	return d.Driver
}
