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
			{
				Name:    contracts.DriverCommandsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
			{
				Name:    contracts.DriverEventsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
			{
				Name:    contracts.PaymentEventsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
			{
				Name:    contracts.PaymentCommandsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
		},
		Queues: []QueueSpec{
			{
				Name:       contracts.DriverTripEventsQueue,
				Durable:    true,
				AutoDelete: false,
			},
			{
				Name:       contracts.DriverTripRequestsQueue,
				Durable:    true,
				AutoDelete: false,
			},
			{
				Name:       contracts.TripDriverEventsQueue,
				Durable:    true,
				AutoDelete: false,
			},
		},
		Bindings: []BindingSpec{
			{
				Queue:      contracts.DriverTripEventsQueue,
				Exchange:   contracts.TripEventsExchange,
				RoutingKey: contracts.TripEventCreated,
			},
			{
				Queue:      contracts.DriverTripEventsQueue,
				Exchange:   contracts.DriverEventsExchange,
				RoutingKey: contracts.DriverEventDriverNotInterested,
			},
			{
				Queue:      contracts.DriverTripRequestsQueue,
				Exchange:   contracts.DriverCommandsExchange,
				RoutingKey: contracts.DriverCmdTripRequest,
			},
			{
				Queue:      contracts.TripDriverEventsQueue,
				Exchange:   contracts.DriverEventsExchange,
				RoutingKey: contracts.DriverEventNoDriversFound,
			},
			{
				Queue:      contracts.TripDriverEventsQueue,
				Exchange:   contracts.DriverEventsExchange,
				RoutingKey: contracts.DriverEventDriverAssigned,
			},
		},
	}
}
