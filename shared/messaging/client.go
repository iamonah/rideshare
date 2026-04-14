package messaging

import (
	"context"
	"fmt"
	"log"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type RabbitMqConfig struct {
	Username string
	Password string
	Host     string
	Vhost    string
	Port     int16
}

func NewRabbitMQClient(config RabbitMqConfig) (*RabbitMQClient, error) {
	cleanPassword := url.QueryEscape(config.Password)
	address := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", config.Username, cleanPassword,
		config.Host, config.Port, config.Vhost)

	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	rc := RabbitMQClient{
		conn:    conn,
		channel: ch,
	}

	err = rc.BootstrapTopology(SetupTopology())
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("bootstrapSharedTopology: %w", err)
	}

	return &rc, nil
}

func (rm *RabbitMQClient) Close() {
	if rm != nil {
		if rm.channel != nil {
			_ = rm.channel.Close()
		}
		if rm.conn != nil {
			_ = rm.conn.Close()
		}
	}
}

func (rm *RabbitMQClient) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	return rm.channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

type MessageHandler func(ctx context.Context, body []byte) error

func (rm *RabbitMQClient) Consume(ctx context.Context, queue string, fn MessageHandler) error {
	delivery, err := rm.channel.ConsumeWithContext(
		ctx,
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return err
	}

	err = rm.channel.Qos(1, 0, false) // prefetch count of 1 for fair dispatch
	if err != nil {
		return fmt.Errorf("Qos: %w", err)
	}

	go func() {
		for msg := range delivery {
			if err := fn(ctx, msg.Body); err != nil {
				log.Printf("failed to handle trip created event: %v", err)
				if msg.Redelivered {
					msg.Nack(false, false) // or false
					continue
				}
				msg.Nack(false, true)
				continue
			}
			if err := msg.Ack(false); err != nil {
				log.Printf("ack failed, channel likely closed: %v", err)
			}
		}

	}()

	return nil
}
