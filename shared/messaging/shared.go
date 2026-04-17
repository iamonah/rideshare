package messaging

import "github.com/iamonah/rideshare/shared/contracts"

func sharedTopologySetup() Topology {
	return Topology{
		Exchanges: []ExchangeSpec{
			{
				Name:    contracts.TripEventsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
		},
		Queues: []QueueSpec{
			{
				Name:       contracts.FindAvailableDriversQueue, //driver-service consumption
				Durable:    true,
				AutoDelete: false,
			},
		},
		Bindings: []BindingSpec{
			{
				Queue:      contracts.FindAvailableDriversQueue,
				Exchange:   contracts.TripEventsExchange,
				RoutingKey: contracts.TripEventCreated,
			},
			{
				Queue:      contracts.FindAvailableDriversQueue,
				Exchange:   contracts.TripEventsExchange,
				RoutingKey: contracts.TripEventDriverNotInterested,
			},
		},
	}
}
