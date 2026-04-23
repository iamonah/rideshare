package driverservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/iamonah/rideshare/shared/contracts"
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

func (c *TripConsumer) ListenConsumer(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("trip consumer is required")
	}
	if c.rabbitmq == nil {
		return fmt.Errorf("rabbitmq client is required")
	}
	if c.s == nil {
		return fmt.Errorf("driver service is required")
	}

	err := c.rabbitmq.Consume(ctx, contracts.DriverTripEventsQueue, func(ctx context.Context, msg messaging.Message) error {
		switch msg.RoutingKey {
		case contracts.TripEventCreated, contracts.DriverEventDriverNotInterested:
		default:
			log.Printf("ignoring unsupported driver matching event: %s", msg.RoutingKey)
			return nil
		}

		var envelope contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			return fmt.Errorf("decode amqp message envelope: %w", err)
		}

		var event eventcontracts.TripCreatedEvent
		if err := json.Unmarshal(envelope.Data, &event); err != nil {
			return fmt.Errorf("decode trip created event: %w", err)
		}

		return c.HandleFindAndNotifyDriver(ctx, &event)
	})

	if err != nil {
		return fmt.Errorf("consume trip created event: %w", err)
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
		if err := c.rabbitmq.Publish(ctx, contracts.DriverEventsExchange, contracts.DriverEventNoDriversFound, contracts.AmqpMessage{
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
	if err := c.rabbitmq.Publish(ctx, contracts.DriverCommandsExchange, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
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
