package events

import (
	"context"
	"fmt"
	"log"

	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/messaging"
)

type TripEventHandler = messaging.MessageHandler

type TripCreatedHandler interface {
	HandleTripCreated(ctx context.Context, fn TripEventHandler) error
}

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQClient
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQClient) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
	}
}

func HandleTripCreated(ctx context.Context, body []byte) error {
	log.Printf("received trip created event: %s", string(body))
	return nil
}

// my thoughts on possible approaches for handling this
func (c *TripConsumer) HandleTripCreated(ctx context.Context, fn TripEventHandler) error {
	err := c.rabbitmq.Consume(ctx, contracts.TripEventCreated, fn)
	if err != nil {
		return fmt.Errorf("consume tripcreated event: %w", err)
	}
	return nil
}
