package driverservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	eventcontracts "github.com/iamonah/rideshare/shared/contracts/events"
	"github.com/iamonah/rideshare/shared/messaging"
)

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQClient
	s        *Service
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQClient, s *Service) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
		s:        s,
	}
}

func (c *TripConsumer) ListenDriverTripEventsQueue(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("trip consumer is required")
	}
	if c.rabbitmq == nil {
		return fmt.Errorf("rabbitmq client is required")
	}
	if c.s == nil {
		return fmt.Errorf("driver service is required")
	}

	// This queue is the driver-matching inbox. It receives freshly created trips and
	// retry events after a driver declines so the service can choose the next driver.
	err := c.rabbitmq.Consume(ctx, messaging.FindAvailableDriversQueue, func(ctx context.Context, msg messaging.Message) error {
		log.Printf("driver-service received message on queue %s with routing key %s", messaging.FindAvailableDriversQueue, msg.RoutingKey)

		var envelope messaging.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			log.Printf("failed to decode AMQP envelope for routing key %s: %s", msg.RoutingKey, string(msg.Body))
			return fmt.Errorf("decode amqp message envelope: %w", err)
		}

		switch msg.RoutingKey {
		case messaging.TripEventCreated, messaging.DriverEventDriverNotInterested:

			var event eventcontracts.TripCreatedEvent
			if err := json.Unmarshal(envelope.Data, &event); err != nil {
				return fmt.Errorf("decode trip created event: %w", err)
			}

			return c.HandleFindAndNotifyDriver(ctx, &event)
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

func (c *TripConsumer) HandleFindAndNotifyDriver(ctx context.Context, event *eventcontracts.TripCreatedEvent) error {
	if event == nil {
		return fmt.Errorf("trip dispatch event is required")
	}

	suitableDrivers := c.s.FindAvailableDrivers(event.Fare.PackageSlug, rejectedDriverID(event))
	if len(suitableDrivers) == 0 {
		log.Printf("found no suitable drivers for trip %s", event.TripID)
		data, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal no drivers found event: %w", err)
		}
		// publish driver-facing event to notify that no drivers were found for the trip request
			if err := c.rabbitmq.Publish(ctx, messaging.DriverEventNoDriversFound, messaging.AmqpMessage{
				OwnerID: event.UserID,
				Data:    data,
			}); err != nil {
			return fmt.Errorf("publish no drivers found: %w", err)
		}

		return nil
	}

	suitableDriver := suitableDrivers[0]
	event.Driver = &eventcontracts.AssignedDriverSnapshot{
		ID:             suitableDriver.GetId(),
		Name:           suitableDriver.GetName(),
		ProfilePicture: suitableDriver.GetProfilePicture(),
		CarPlate:       suitableDriver.GetCarPlate(),
		PackageSlug:    event.Fare.PackageSlug,
	}
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal trip request event: %w", err)
	}


	//driverfacing event to notify driver of trip request
	if err := c.rabbitmq.Publish(ctx, messaging.DriverCmdTripRequest, messaging.AmqpMessage{
		OwnerID: suitableDriver.GetId(),
		Data:    data,
	}); err != nil {
		return fmt.Errorf("publish trip request to driver %s: %w", suitableDriver.GetId(), err)
	}
	return nil
}

func rejectedDriverID(event *eventcontracts.TripCreatedEvent) string {
	if event == nil || event.Driver == nil {
		return ""
	}

	return event.Driver.ID
}


func (c *TripConsumer) HandleTripAccepted(ctx context.Context) error {
	return nil
}
