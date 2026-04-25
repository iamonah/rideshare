package websockettransport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/messaging"
)

func (h *Handler) ReceiveDriverEvents(client *Client, event contracts.WSMessage) error {
	var data messaging.AmqpMessage
	err := json.Unmarshal(event.Data, &data)
	if err != nil {
		return fmt.Errorf("ReceiveDriverEvents: failed to unmarshal message data: %w", err)
	}
	switch event.Type {
	case messaging.DriverCmdLocation:
		// Handle driver location update in the future
		return nil
	case messaging.DriverCmdTripAccept, messaging.DriverCmdTripDecline:
		// Forward the message to RabbitMQ
		if err := h.rabbitmq.Publish(context.Background(), messaging.DriverCommandsExchange, event.Type, data); err != nil {
			log.Printf("Error publishing message to RabbitMQ: %v", err)
		}
	default:
		log.Printf("Unknown message type: %s", event.Type)
	}
	return nil
}
