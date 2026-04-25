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

	err := c.rabbitmq.Consume(ctx, messaging.DriverTripEventsQueue, func(ctx context.Context, msg messaging.Message) error {
		var envelope messaging.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
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
		if err := c.rabbitmq.Publish(ctx, messaging.DriverEventsExchange, messaging.DriverEventNoDriversFound, messaging.AmqpMessage{
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

	log.Printf("publishing trip request command for driver: %s", suitableDriver.GetId())
	if err := c.rabbitmq.Publish(ctx, messaging.DriverCommandsExchange, messaging.DriverCmdTripRequest, messaging.AmqpMessage{
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
