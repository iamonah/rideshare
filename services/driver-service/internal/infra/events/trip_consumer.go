package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/iamonah/rideshare/shared/contracts"
	eventcontracts "github.com/iamonah/rideshare/shared/contracts/events"
	"github.com/iamonah/rideshare/shared/messaging"
)

type TripCreatedHandler interface {
	HandleTripCreated(ctx context.Context, fn func(context.Context, *eventcontracts.TripCreatedEvent) error) error
}

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQClient
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQClient) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
	}
}

func (c *TripConsumer) HandleTripCreated(ctx context.Context, fn func(context.Context, *eventcontracts.TripCreatedEvent) error) error {
	if fn == nil {
		return fmt.Errorf("trip created handler is required")
	}

	err := c.rabbitmq.Consume(ctx, contracts.FindAvailableDriversQueue, func(ctx context.Context, body []byte) error {
		var envelope contracts.AmqpMessage
		if err := json.Unmarshal(body, &envelope); err != nil {
			return fmt.Errorf("decode amqp message envelope: %w", err)
		}

		var event eventcontracts.TripCreatedEvent
		if err := json.Unmarshal(envelope.Data, &event); err != nil {
			return fmt.Errorf("decode trip created event: %w", err)
		}

		log.Printf("received trip created event for trip ID: %s", event.TripID)
		return fn(ctx, &event)
	})
	if err != nil {
		return fmt.Errorf("consume trip created event: %w", err)
	}

	return nil
}
