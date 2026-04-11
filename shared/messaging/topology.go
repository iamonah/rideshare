package messaging

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeSpec struct {
	Name    string
	Kind    string
	Durable bool
}

type QueueSpec struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Args       amqp.Table
}

type BindingSpec struct {
	Queue      string
	RoutingKey string
	Exchange   string
}

type Topology struct {
	Exchanges []ExchangeSpec
	Queues    []QueueSpec
	Bindings  []BindingSpec
}

func (rm *RabbitMQClient) declareExchange(spec ExchangeSpec) error {
	err := rm.channel.ExchangeDeclare(spec.Name, spec.Kind, spec.Durable, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declare exchange %q: %w", spec.Name, err)
	}

	return nil
}

func (rm *RabbitMQClient) declareQueue(spec QueueSpec) error {
	_, err := rm.channel.QueueDeclare(spec.Name, spec.Durable, spec.AutoDelete, false, false, spec.Args)
	if err != nil {
		return fmt.Errorf("declare queue %q: %w", spec.Name, err)
	}

	return nil
}

func (rm *RabbitMQClient) bindQueue(spec BindingSpec) error {
	err := rm.channel.QueueBind(spec.Queue, spec.RoutingKey, spec.Exchange, false, nil)
	if err != nil {
		return fmt.Errorf("bind queue %q to %q with %q: %w", spec.Queue, spec.Exchange, spec.RoutingKey, err)
	}

	return nil
}

func (rm *RabbitMQClient) BootstrapTopology(topology Topology) error {
	for _, exchange := range topology.Exchanges {
		if err := rm.declareExchange(exchange); err != nil {
			return err
		}
	}

	for _, queue := range topology.Queues {
		if err := rm.declareQueue(queue); err != nil {
			return err
		}
	}

	for _, binding := range topology.Bindings {
		if err := rm.bindQueue(binding); err != nil {
			return err
		}
	}

	return nil
}
