package messaging

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rm *RabbitMQClient) declareExchange(name, kind string) error {
	err := rm.channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare exchange %q: %w", name, err)
	}

	return nil
}

func (rm *RabbitMQClient) declareQueue(name string, args amqp.Table) error {
	_, err := rm.channel.QueueDeclare(name, true, false, false, false, args)
	if err != nil {
		return fmt.Errorf("declare queue %q: %w", name, err)
	}

	return nil
}

func (rm *RabbitMQClient) bindQueue(exchange, queue string, routingKeys []string) error {
	for _, routingKey := range routingKeys {
		err := rm.channel.QueueBind(queue, routingKey, exchange, false, nil)
		if err != nil {
			return fmt.Errorf("bind queue %q to %q with %q: %w", queue, exchange, routingKey, err)
		}
	}

	return nil
}

func (rm *RabbitMQClient) declareQueueAndBind(exchange, queue string, routingKeys []string, args amqp.Table) error {
	if err := rm.declareQueue(queue, args); err != nil {
		return err
	}

	if len(routingKeys) == 0 {
		return nil
	}

	if err := rm.bindQueue(exchange, queue, routingKeys); err != nil {
		return err
	}

	return nil
}
