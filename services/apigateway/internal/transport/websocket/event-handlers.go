package websockettransport

import (
	"context"
	"fmt"

	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/messaging"
)

func (h *Handler) ReceiveDriverEvents(client *Client, event contracts.WSMessage) error {
	switch event.Type {
	case messaging.DriverCmdLocation:
		// Handle driver location update in the future
		return nil
	case messaging.DriverCmdTripAccept, messaging.DriverCmdTripDecline:
		envelope := messaging.AmqpMessage{
			OwnerID: client.ID,
			Data:    event.Data,
		}
		return h.rabbitmq.Publish(context.Background(), messaging.DriverCommandsExchange, event.Type, envelope)
	default:
		return fmt.Errorf("unknown driver event type: %s", event.Type)
	}
}
