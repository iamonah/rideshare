package events

import (
	"context"

	"github.com/iamonah/rideshare/shared/contracts"
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

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, data string) error {
	body := []byte(data)

	return p.rabbitmq.Publish(ctx, contracts.TripEventsExchange, contracts.TripEventCreated, body)
}
