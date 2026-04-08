package rabbitmq

import (
	"fmt"
	"net/url"

	"github.com/streadway/amqp"
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

	return &RabbitMQClient{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *RabbitMQClient) Close() {
	c.channel.Close()
	c.conn.Close()
}

func (c *RabbitMQClient) Publish(exchange, routingKey string, body []byte) error {
	return c.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (c *RabbitMQClient) Consume(queue string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}