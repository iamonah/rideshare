package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/iamonah/rideshare/shared/contracts"
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
		return nil, fmt.Errorf("dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("connection channel: %w", err)
	}

	rc := RabbitMQClient{
		conn:    conn,
		channel: ch,
	}
	rc.logReturnedMessages()

	if err = rc.BootstrapTopology(DeadLetterTopology()); err != nil {
		return nil, fmt.Errorf("bootstrap dead letter infrastructure: %w", err)
	}
	if err = rc.BootstrapTopology(sharedTopologySetup()); err != nil {
		return nil, fmt.Errorf("bootstrap shared broker infrastructure: %w", err)
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

func (rm *RabbitMQClient) logReturnedMessages() {
	returns := rm.channel.NotifyReturn(make(chan amqp.Return, 16))

	go func() {
		for returned := range returns {
			log.Printf(
				"rabbitmq returned unroutable message exchange=%s routing_key=%s reply_code=%d reply_text=%q body=%s",
				returned.Exchange,
				returned.RoutingKey,
				returned.ReplyCode,
				returned.ReplyText,
				string(returned.Body),
			)
		}
	}()
}

func (rm *RabbitMQClient) Publish(ctx context.Context, exchange, routingKey string, msg contracts.AmqpMessage) error {
	log.Printf("publishing messages with routing key: %s", routingKey)
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal amqp message: %w", err)
	}

	return rm.channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		true,       // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

type Message struct {
	Body       []byte
	RoutingKey string
}

type MessageHandler func(ctx context.Context, msg Message) error

func (rm *RabbitMQClient) Consume(ctx context.Context, queue string, fn MessageHandler) error {
	err := rm.channel.Qos(1, 0, false) // prefetch count of 1 for fair dispatch
	if err != nil {
		return fmt.Errorf("Qos: %w", err)
	}

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

	log.Printf("started consuming queue: %s", queue)
	go func() {
		defer log.Printf("stopped consuming queue: %s", queue)

		for delivery := range delivery {
			msg := Message{
				Body:       delivery.Body,
				RoutingKey: delivery.RoutingKey,
			}
			if err := fn(ctx, msg); err != nil {
				log.Printf("failed to handle message from queue %s: %v", queue, err)
				if delivery.Redelivered {
					delivery.Nack(false, false) // or false
					continue
				}
				delivery.Nack(false, true)
				continue
			}
			if err := delivery.Ack(false); err != nil {
				log.Printf("ack failed, channel likely closed: %v", err)
			}
		}

	}()

	return nil
}
