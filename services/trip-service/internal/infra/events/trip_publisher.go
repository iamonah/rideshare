package events

import (
	"context"
	"encoding/json"
	"fmt"

	tripdomain "github.com/iamonah/rideshare/services/trip-service/internal/domain/trip"
	"github.com/iamonah/rideshare/shared/contracts"
	eventcontracts "github.com/iamonah/rideshare/shared/contracts/events"
	"github.com/iamonah/rideshare/shared/messaging"
)

type TripEventPublisher struct {
	rabbitmq *messaging.RabbitMQClient
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMQClient) *TripEventPublisher {
	return &TripEventPublisher{
		rabbitmq: rabbitmq,
	}
}

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, trip *tripdomain.Trip) error {
	event, err := toTripCreatedEvent(trip)
	if err != nil {
		return fmt.Errorf("map trip created event: %w", err)
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal trip created event: %w", err)
	}

	return p.rabbitmq.Publish(ctx, contracts.TripEventsExchange, contracts.TripEventCreated, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    jsonData,
	})
}

func toTripCreatedEvent(trip *tripdomain.Trip) (*eventcontracts.TripCreatedEvent, error) {
	if trip == nil {
		return nil, fmt.Errorf("trip is required")
	}
	if trip.RideFare == nil {
		return nil, fmt.Errorf("trip ride fare is required")
	}
	if trip.RideFare.Route == nil {
		return nil, fmt.Errorf("trip route is required")
	}

	route := trip.RideFare.Route
	geometry, err := routeGeometry(route.Geometry.Coordinates)
	if err != nil {
		return nil, err
	}
	pickup := geometry[0]
	dropoff := geometry[len(geometry)-1]

	return &eventcontracts.TripCreatedEvent{
		TripID: trip.ID.Hex(),
		UserID: trip.UserID,
		Status: trip.Status,
		Fare: eventcontracts.TripFareSnapshot{
			ID:          trip.RideFare.ID.Hex(),
			PackageSlug: trip.RideFare.PackageSlug.String(),
		},
		Pickup:          pickup,
		Dropoff:         dropoff,
		DistanceMeters:  route.Distance,
		DurationSeconds: route.Duration,
		Route: eventcontracts.TripRouteSnapshot{
			Geometry: geometry,
		},
	}, nil
}

func routeGeometry(coordinates [][]float64) ([]eventcontracts.Coordinate, error) {
	if len(coordinates) == 0 {
		return nil, fmt.Errorf("trip route geometry is required")
	}

	geometry := make([]eventcontracts.Coordinate, 0, len(coordinates))
	for _, pair := range coordinates {
		coordinate, err := coordinateFromOSRMPair(pair)
		if err != nil {
			return nil, fmt.Errorf("map route coordinate: %w", err)
		}

		geometry = append(geometry, coordinate)
	}

	return geometry, nil
}

func coordinateFromOSRMPair(pair []float64) (eventcontracts.Coordinate, error) {
	if len(pair) != 2 {
		return eventcontracts.Coordinate{}, fmt.Errorf("invalid coordinate pair")
	}

	return eventcontracts.Coordinate{
		Longitude: pair[0],
		Latitude:  pair[1],
	}, nil
}
