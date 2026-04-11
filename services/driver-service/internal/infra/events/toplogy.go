package events

import (
	"github.com/iamonah/rideshare/shared/contracts"
	sharedmessaging "github.com/iamonah/rideshare/shared/messaging"
)

const (
	TripCreatedQueue = "driver.trip.created.queue"
)

func Topology() sharedmessaging.Topology {
	return sharedmessaging.Topology{
		Exchanges: []sharedmessaging.ExchangeSpec{
			{
				Name:    contracts.TripEventsExchange,
				Kind:    sharedmessaging.TopicExchangeKind,
				Durable: true,
			},
		},
		Queues: []sharedmessaging.QueueSpec{
			sharedmessaging.DurableQueueWithDLX(TripCreatedQueue),
		},
		Bindings: []sharedmessaging.BindingSpec{
			{
				Queue:      TripCreatedQueue,
				Exchange:   contracts.TripEventsExchange,
				RoutingKey: contracts.TripEventCreated,
			},
		},
	}
}

func BootstrapDriverTopology(client *sharedmessaging.RabbitMQClient) error {
	return client.BootstrapTopology(Topology())
}
