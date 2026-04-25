package messaging

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

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

func (rm *RabbitMQClient) setupDeadLetterInfrastructure() error {
	if err := rm.setupTopicExchange(DeadLetterExchange); err != nil {
		return err
	}

	if err := rm.setupQueueBindings(DeadLetterExchange, DeadLetterQueue, []string{
		DeadLetterBindingKey,
	}, amqp.Table{
		// Keep this explicit so redeclares match existing broker state.
		"x-dead-letter-exchange": "",
		// Keep this explicit so redeclares match existing broker state.
		"x-dead-letter-routing-key": "",
		// Expire messages from the DLQ after 60 seconds instead of re-dead-lettering them.
		"x-message-ttl": int32(60000),
	}); err != nil {
		return err
	}

	return nil
}
