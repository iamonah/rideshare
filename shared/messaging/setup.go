package messaging

import amqp "github.com/rabbitmq/amqp091-go"

const (
	//exchange kind
	TopicExchangeKind  = "topic"
	FanoutExchangeKind = "fanout"
	DirectExchangeKind = "direct"

	//setup deadletter exchange and queue
	DeadLetterExchange   = "rideshare.dlx"
	DeadLetterQueue      = "rideshare.dlx.queue"
	DeadLetterBindingKey = "#"
)

func DeadLetterTopology() Topology {
	return Topology{
		Exchanges: []ExchangeSpec{
			{
				Name:    DeadLetterExchange,
				Kind:    TopicExchangeKind,
				Durable: true,
			},
		},
		Queues: []QueueSpec{
			{
				Name:       DeadLetterQueue,
				Durable:    true,
				AutoDelete: false,
				Args: amqp.Table{
					// Keep this explicit so redeclares match existing broker state.
					"x-dead-letter-exchange": "",
					// Keep this explicit so redeclares match existing broker state.
					"x-dead-letter-routing-key": "",
					// Expire messages from the DLQ after 60 seconds instead of re-dead-lettering them.
					"x-message-ttl": int32(60000),
				},
			},
		},
		Bindings: []BindingSpec{
			{
				Queue:      DeadLetterQueue,
				Exchange:   DeadLetterExchange,
				RoutingKey: DeadLetterBindingKey,
			},
		},
	}
}
