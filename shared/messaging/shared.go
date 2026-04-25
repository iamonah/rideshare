package messaging

func sharedTopologySetup() Topology {
	return Topology{
		Exchanges: []ExchangeSpec{
			{
				Name:    TripEventsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
			{
				Name:    DriverCommandsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
			{
				Name:    DriverEventsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
			{
				Name:    PaymentEventsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
			{
				Name:    PaymentCommandsExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
		},
		Queues: []QueueSpec{
			{
				Name:       DriverTripEventsQueue,
				Durable:    true,
				AutoDelete: false,
			},
			{
				Name:       DriverTripRequestsQueue,
				Durable:    true,
				AutoDelete: false,
			},
			{
				Name:       TripDriverEventsQueue,
				Durable:    true,
				AutoDelete: false,
			},
			{
				Name:       TripDriverCommandsQueue,
				Durable:    true,
				AutoDelete: false,
			},
		},
		Bindings: []BindingSpec{
			{
				Queue:      DriverTripEventsQueue,
				Exchange:   TripEventsExchange,
				RoutingKey: TripEventCreated,
			},
			{
				Queue:      DriverTripEventsQueue,
				Exchange:   DriverEventsExchange,
				RoutingKey: DriverEventDriverNotInterested,
			},
			{
				Queue:      DriverTripRequestsQueue,
				Exchange:   DriverCommandsExchange,
				RoutingKey: DriverCmdTripRequest,
			},
			{
				Queue:      TripDriverEventsQueue,
				Exchange:   DriverEventsExchange,
				RoutingKey: DriverEventNoDriversFound,
			},
			{
				Queue:      TripDriverEventsQueue,
				Exchange:   DriverEventsExchange,
				RoutingKey: DriverEventDriverAssigned,
			},
			{
				Queue:      TripDriverCommandsQueue,
				Exchange:   DriverCommandsExchange,
				RoutingKey: DriverCmdTripAccept,
			},
			{
				Queue:      TripDriverCommandsQueue,
				Exchange:   DriverCommandsExchange,
				RoutingKey: DriverCmdTripDecline,
			},
		},
	}
}
