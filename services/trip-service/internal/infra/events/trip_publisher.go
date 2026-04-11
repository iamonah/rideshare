package events

import (
	"context"

	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/messaging"
)

type TripEventPublisher struct {
	rabbitmq   *messaging.RabbitMQClient
	exchange   string
	routingKey string
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMQClient) *TripEventPublisher {
	return &TripEventPublisher{
		rabbitmq:   rabbitmq,
		exchange:   contracts.TripEventsExchange,
		routingKey: contracts.TripEventCreated,
	}
}

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, data string) error {
	body := []byte(data)

	return p.rabbitmq.Publish(ctx, p.exchange, p.routingKey, body)
}
