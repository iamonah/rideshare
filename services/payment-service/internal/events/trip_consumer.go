package events

import (
	"context"
	"encoding/json"
	"fmt"

	eventcontracts "github.com/iamonah/rideshare/shared/contracts/events"
	"github.com/iamonah/rideshare/shared/messaging"

	"github.com/iamonah/rideshare/services/payment-service/internal/domain"
)

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQClient
	service  domain.Service
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQClient, service domain.Service) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *TripConsumer) Listen(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("trip consumer is required")
	}
	if c.rabbitmq == nil {
		return fmt.Errorf("rabbitmq client is required")
	}
	if c.service == nil {
		return fmt.Errorf("payment service is required")
	}

	return c.rabbitmq.Consume(ctx, messaging.PaymentTripResponseQueue, func(ctx context.Context, msg messaging.Message) error {
		var envelope messaging.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			return fmt.Errorf("decode amqp message envelope: %w", err)
		}

		switch msg.RoutingKey {
		case messaging.PaymentCmdCreateSession:
			var payload eventcontracts.PaymentTripResponseData
			if err := json.Unmarshal(envelope.Data, &payload); err != nil {
				return fmt.Errorf("decode payment create session payload: %w", err)
			}
			return c.handleTripAccepted(ctx, payload)
		}

		return nil
	})
}

func (c *TripConsumer) handleTripAccepted(ctx context.Context, payload eventcontracts.PaymentTripResponseData) error {
	paymentSession, err := c.service.CreatePaymentSession(
		ctx,
		payload.TripID,
		payload.UserID,
		payload.DriverID,
		int64(payload.Amount),
		payload.Currency,
	)
	if err != nil {
		return fmt.Errorf("create payment session: %w", err)
	}

	// Publish payment session created event
	paymentPayload := eventcontracts.PaymentEventSessionCreatedData{
		TripID:    payload.TripID,
		SessionID: paymentSession.StripeSessionID,
		Amount:    float64(paymentSession.Amount) / 100.0, // Convert from cents to dollars
		Currency:  paymentSession.Currency,
	}

	payloadBytes, err := json.Marshal(paymentPayload)
	if err != nil {
		return fmt.Errorf("marshal payment session payload: %w", err)
	}

	if err := c.rabbitmq.Publish(ctx, messaging.PaymentEventSessionCreated,
		messaging.AmqpMessage{
			OwnerID: payload.UserID,
			Data:    payloadBytes,
		},
	); err != nil {
		return fmt.Errorf("publish payment session created event: %w", err)
	}

	return nil
}
