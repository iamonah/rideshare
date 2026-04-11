package events

import (
	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/messaging"
)

func Topology() messaging.Topology {
	return messaging.Topology{
		Exchanges: []messaging.ExchangeSpec{
			{
				Name:    contracts.TripEventsExchange,
				Kind:    messaging.TopicExchangeKind,
				Durable: true,
			},
		},
	}
}

func BootstrapTripTopology(client *messaging.RabbitMQClient) error {
	return client.BootstrapTopology(Topology())
}
