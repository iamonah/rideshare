package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/iamonah/rideshare/services/trip-service/internal/domain/trip"
	"github.com/iamonah/rideshare/shared/contracts/events"
	eventcontracts "github.com/iamonah/rideshare/shared/contracts/events"
	"github.com/iamonah/rideshare/shared/messaging"
	"github.com/iamonah/rideshare/shared/types"
)

type DriverConsumer struct {
	rabbitmq *messaging.RabbitMQClient
	tb       trip.ExtTripBusiness
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQClient, tb *trip.TripBusiness) *DriverConsumer {
	return &DriverConsumer{
		rabbitmq: rabbitmq,
		tb:       tb,
	}
}

func (c *DriverConsumer) Listen(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("consumer is required")
	}

	if c.rabbitmq == nil {
		return fmt.Errorf("rabbitmq client is required")
	}
	if c.tb == nil {
		return fmt.Errorf("trip business is required")
	}

	// This queue is the driver-matching inbox. It receives freshly created trips and
	// retry events after a driver declines so the service can choose the next driver.
	err := c.rabbitmq.Consume(ctx, messaging.DriverTripResponseQueue, func(ctx context.Context, msg messaging.Message) error {
		log.Printf("trip-service received message on queue %s with routing key %s", messaging.DriverTripResponseQueue, msg.RoutingKey)

		var envelope messaging.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			log.Printf("failed to decode AMQP envelope for routing key %s: %s", msg.RoutingKey, string(msg.Body))
			return fmt.Errorf("decode amqp message envelope: %w", err)
		}

		var event eventcontracts.DriverTripResponseData
		if err := json.Unmarshal(envelope.Data, &event); err != nil {
			return fmt.Errorf("decode trip created event: %w", err)
		}
		log.Printf("hello world: %+v", event) //why is this not printing?
		switch msg.RoutingKey {
		case messaging.DriverEventTripAccepted:
			return c.handleTripAccepted(ctx, event)
		case messaging.DriverEventTripDeclined:
			return c.handleTripDeclined(ctx, event.TripID, event.RiderID)
		default:
			log.Printf("ignoring unsupported driver matching event: %s", msg.RoutingKey)
			return nil
		}
	})

	if err != nil {
		return fmt.Errorf("consume driver trip events queue: %w", err)
	}

	return nil
}

func (c *DriverConsumer) handleTripDeclined(ctx context.Context, tripId, userId string) error {
	// When a driver declines, we should try to find another driver
	currentTrip, err := c.tb.GetTripByID(ctx, tripId)
	if err != nil {
		return err
	}

	if currentTrip == nil {
		payload, err := json.Marshal(map[string]string{
			"tripId":  tripId,
			"message": "Trip has expired. Please book another trip.",
		})
		if err != nil {
			return err
		}

		if err = c.rabbitmq.Publish(ctx, messaging.TripEventNotFound, messaging.AmqpMessage{
			OwnerID: userId,
			Data:    payload,
		}); err != nil {
			return err
		}
		return nil
	}

	//nil driver snapshot indicates a declined trip and signals the driver matching workflow to find another driver.
	retryEvent, err := marshalTripEvent(currentTrip, nil)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.Publish(ctx, messaging.DriverEventDriverNotInterested,
		messaging.AmqpMessage{
			OwnerID: userId,
			Data:    retryEvent,
		},
	); err != nil {
		return err
	}

	return nil
}

func (c *DriverConsumer) handleTripAccepted(ctx context.Context, driver eventcontracts.DriverTripResponseData) error {
	currentTrip, err := c.tb.GetTripByID(ctx, driver.TripID)
	if err != nil {
		return err
	}

	if currentTrip == nil {
		payload, err := json.Marshal(map[string]string{
			"tripId":  driver.TripID,
			"message": "Trip has expired. Please book another trip.",
		})
		if err != nil {
			return err
		}

		if err = c.rabbitmq.Publish(ctx, messaging.TripEventNotFound, messaging.AmqpMessage{
			OwnerID: driver.RiderID,
			Data:    payload,
		}); err != nil {
			return err
		}
		return nil
	}

	if err := c.tb.UpdateTrip(ctx, driver.TripID, "accepted", assignedDriverSnapshot(driver)); err != nil {
		log.Printf("Failed to update the trip: %v", err)
		return err
	}

	currentTrip.Status = "accepted"
	assignedEvent, err := marshalTripEvent(currentTrip, assignedEventDriverSnapshot(driver))
	if err != nil {
		return err
	}

	log.Printf("currenttrip: %+v\n", currentTrip)
	// Notify the rider that a driver has been assigned
	if err := c.rabbitmq.Publish(ctx, messaging.DriverEventDriverAssigned, messaging.AmqpMessage{
		OwnerID: currentTrip.UserID,
		Data:    assignedEvent,
	}); err != nil {
		return err
	}

	marshalledPayload, err := json.Marshal(events.PaymentCreateSessionCommand{
		TripID:   driver.TripID,
		UserID:   currentTrip.UserID,
		DriverID: driver.Driver.Id,
		Amount:   currentTrip.RideFare.TotalPriceInCents,
		Currency: "USD",
	})
	if err != nil {
		return err
	}

	if err := c.rabbitmq.Publish(ctx, messaging.PaymentCmdCreateSession, messaging.AmqpMessage{
		OwnerID: currentTrip.UserID,
		Data:    marshalledPayload,
	}); err != nil {
		return err
	}

	return nil
}

func assignedEventDriverSnapshot(driver eventcontracts.DriverTripResponseData) *eventcontracts.AssignedDriverSnapshot {
	return &eventcontracts.AssignedDriverSnapshot{
		ID:             driver.Driver.Id,
		Name:           driver.Driver.Name,
		ProfilePicture: driver.Driver.ProfilePicture,
		CarPlate:       driver.Driver.CarPlate,
		PackageSlug:    driver.Driver.PackageSlug,
	}
}

func assignedDriverSnapshot(driver eventcontracts.DriverTripResponseData) *trip.AssignedDriverSnapshot {
	return &trip.AssignedDriverSnapshot{
		ID:        driver.Driver.Id,
		FirstName: driver.Driver.Name,
	}
}

func marshalTripEvent(t *trip.Trip, driver *eventcontracts.AssignedDriverSnapshot) ([]byte, error) {
	if t == nil || t.RideFare == nil || t.RideFare.Route == nil {
		return nil, fmt.Errorf("trip event payload requires trip, fare, and route")
	}

	geometry := make([]types.Coordinate, 0, len(t.RideFare.Route.Geometry.Coordinates))
	for _, pair := range t.RideFare.Route.Geometry.Coordinates {
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid route coordinate pair")
		}
		geometry = append(geometry, types.Coordinate{
			Longitude: pair[0],
			Latitude:  pair[1],
		})
	}
	if len(geometry) == 0 {
		return nil, fmt.Errorf("trip route geometry is required")
	}

	payload := eventcontracts.TripCreatedEvent{
		TripID: t.ID.Hex(),
		UserID: t.UserID,
		Status: t.Status,
		Fare: eventcontracts.TripFareSnapshot{
			ID:          t.RideFare.ID.Hex(),
			PackageSlug: t.RideFare.PackageSlug.String(),
		},
		Pickup:          geometry[0],
		Dropoff:         geometry[len(geometry)-1],
		DistanceMeters:  t.RideFare.Route.Distance,
		DurationSeconds: t.RideFare.Route.Duration,
		Route: eventcontracts.TripRouteSnapshot{
			Geometry: geometry,
		},
		Driver: driver,
	}

	return json.Marshal(payload)
}
