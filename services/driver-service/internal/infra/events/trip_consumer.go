package events

import (
	"context"
	"log"

	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/messaging"
)

type TripCreatedHandler interface {
	HandleTripCreated(ctx context.Context, body []byte) error
}

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQClient
	handler  TripCreatedHandler
	queue    string
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQClient, handler TripCreatedHandler) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
		handler:  handler,
		queue:    contracts.TripEventCreated,
	}
}

func (c *TripConsumer) Start(ctx context.Context) error {
	err := c.rabbitmq.Consume(ctx, c.queue, func(ctx context.Context, body []byte) error {
		log.Printf("driver recieved message: %v", body)
		return nil
	})

	if err != nil {
		log.Printf("failed to start trip consumer: %v", err)
		return err
	}

	return nil
}
