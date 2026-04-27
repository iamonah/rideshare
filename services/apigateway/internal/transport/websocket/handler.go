package websockettransport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/iamonah/rideshare/services/apigateway/internal/infra/client"
	httpcommon "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/common"
	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/messaging"
	"github.com/iamonah/rideshare/shared/proto/pb/driverpb"
)

type Server struct {
	manager *Manager
	broker  *messaging.RabbitMQClient
	drivers *client.DriverClient
}

func NewServer(dc *client.DriverClient, rmq *messaging.RabbitMQClient) *Server {
	return &Server{
		manager: NewManager(),
		drivers: dc,
		broker:  rmq,
	}
}

func (s *Server) HandleDriversWebsocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("Missing userID in query parameters")
		http.Error(w, "missing userID", http.StatusBadRequest)
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		http.Error(w, "failed to upgrade websocket connection", http.StatusBadRequest)
		return
	}

	data, err := s.drivers.RegisterDriver(r.Context(), &driverpb.RegisterDriverRequest{
		DriverId:    userID,
		PackageSlug: packageSlug,
	})

	if err != nil {
		httpcommon.WriteUpstreamGRPCError(w, "driver_service", err)
		conn.Close()
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal registration message for driver %s: %v", userID, err)
		conn.Close()
		return
	}

	msg := contracts.WSMessage{
		Type: messaging.DriverCmdRegister,
		Data: payload,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Failed to send registration message to driver %s: %v", userID, err)
		conn.Close()
	}

	fmt.Println("Start drvier websocket for driverID:", userID)
	client := NewClient(conn, userID, s.manager.routeEventHandler, s.manager.removeClient)
	s.manager.addClient(userID, client)
	go client.ReadMessage()
	go client.WriteMessage()
}

func (s *Server) HandleRidersWebsocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("Missing userID in query parameters")
		http.Error(w, "missing userID", http.StatusBadRequest)
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		http.Error(w, "failed to upgrade websocket connection", http.StatusBadRequest)
		return
	}

	fmt.Println("Start rider websocket for riderID:", userID)
	client := NewClient(conn, userID, s.manager.routeEventHandler, s.manager.removeClient)
	s.manager.addClient(userID, client)
	go client.ReadMessage()
	go client.WriteMessage()
}

// This queue carries trip request commands that are forwarded to the connected driver websocket.
func (s *Server) SendToDriver(ctx context.Context) error {
	return s.broker.Consume(ctx, messaging.DriverCmdTripRequestQueue, func(ctx context.Context, msg messaging.Message) error {
		var envelope messaging.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			return fmt.Errorf("listenDriverTripRequestsQueue: decode amqp message envelope: %w", err)
		}

		event := contracts.WSMessage{
			Type: msg.RoutingKey,
			Data: envelope.Data,
		}

		log.Println("send to driver websocket with event:", event.Type)
		if err := s.manager.SendToClient(envelope.OwnerID, event); err != nil {
			return fmt.Errorf("listenDriverTripRequestsQueue: send to driver %s: %w", envelope.OwnerID, err)
		}

		return nil
	})
}

// Rider notifications are split across multiple queues following the tutorial naming.
// We consume each queue and forward its payload to the rider websocket identified by ownerId.
func (s *Server) SendToRider(ctx context.Context) error {
	for _, queue := range []string{
		messaging.NotifyRiderNoDriversFoundQueue,
		messaging.NotifyDriverAssignQueue,
		messaging.NotifyPaymentSessionCreatedQueue,
		messaging.NotifyPaymentSuccessQueue,
	} {
		if err := s.consumeRiderQueue(ctx, queue); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) consumeRiderQueue(ctx context.Context, queue string) error {
	return s.broker.Consume(ctx, queue, func(ctx context.Context, msg messaging.Message) error {
		var envelope messaging.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			return fmt.Errorf("listenRiderEventsQueue: decode amqp message envelope: %w", err)
		}

		event := contracts.WSMessage{
			Type: msg.RoutingKey,
			Data: envelope.Data,
		}
		log.Println("send to rider websocket with event:", event.Type)
		if err := s.manager.SendToClient(envelope.OwnerID, event); err != nil {
			return fmt.Errorf("listenRiderEventsQueue: send to rider %s: %w", envelope.OwnerID, err)
		}

		return nil
	})
}
