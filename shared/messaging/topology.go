package messaging

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rm *RabbitMQClient) declareExchange(name, kind string) error {
	if err := rm.channel.ExchangeDeclare(name, kind, true, false, false, false, nil); err != nil {
		return fmt.Errorf("exchange %q: %w", name, err)
	}

	return nil
}

func (rm *RabbitMQClient) declareQueue(name string, args amqp.Table) error {
	if _, err := rm.channel.QueueDeclare(name, true, false, false, false, args); err != nil {
		return fmt.Errorf("queue %q: %w", name, err)
	}

	return nil
}

func (rm *RabbitMQClient) bindQueue(exchange, queue string, routingKeys []string) error {
	for _, routingKey := range routingKeys {
		if err := rm.channel.QueueBind(queue, routingKey, exchange, false, nil); err != nil {
			return fmt.Errorf("bind queue %q to exchange %q with routing key %q: %w", queue, exchange, routingKey, err)
		}
	}

	return nil
}

func (rm *RabbitMQClient) declareQueueAndBind(queue string, routingKeys []string, args amqp.Table) error {
	if err := rm.declareQueue(queue, args); err != nil {
		return err
	}

	if len(routingKeys) == 0 {
		return nil
	}

	if err := rm.bindQueue(RideShareExchange, queue, routingKeys); err != nil {
		return err
	}

	return nil
}

func (rm *RabbitMQClient) setupTopicExchange(name string) error {
	if err := rm.declareExchange(name, TopicExchangeKind); err != nil {
		return fmt.Errorf("setup topic exchange %q: %w", name, err)
	}

	return nil
}

func (rm *RabbitMQClient) setupQueueBindings(exchange, queue string, routingKeys []string, args amqp.Table) error {
	if err := rm.declareQueueAndBind(queue, routingKeys, args); err != nil {
		return fmt.Errorf("setup queue %q on exchange %q: %w", queue, exchange, err)
	}

	return nil
}
