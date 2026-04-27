package websockettransport

import (
	"context"
	"fmt"

	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/messaging"
)

func (s *Server) ReceiveDriverEvents(client *Client, event contracts.WSMessage) error {
	switch event.Type {
	case messaging.DriverCmdLocation:
		// Handle driver location update in the future
		return nil
	case messaging.DriverCmdTripAccept, messaging.DriverCmdTripDecline:
		envelope := messaging.AmqpMessage{
			OwnerID: client.ID,
			Data:    event.Data,
		}
		return s.broker.Publish(context.Background(), event.Type, envelope)
	default:
		return fmt.Errorf("unknown driver event type: %s", event.Type)
	}
}

func (s *Server) ReceiveRiderEvents(client *Client, event contracts.WSMessage) error {
	switch event.Type {
	// Handle rider events in the future
	}
	return nil
}
