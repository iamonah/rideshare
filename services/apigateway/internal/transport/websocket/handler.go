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

type Handler struct {
	Manager      *Manager
	rabbitmq     *messaging.RabbitMQClient
	driverClient *client.DriverClient
}

func NewHandler(dc *client.DriverClient, rmq *messaging.RabbitMQClient) *Handler {
	return &Handler{
		Manager:      NewManager(),
		driverClient: dc,
		rabbitmq:     rmq,
	}
}

func (h *Handler) HandleDriversWebsocket(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.driverClient.RegisterDriver(r.Context(), &driverpb.RegisterDriverRequest{
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
	client := NewClient(conn, userID, h.Manager.routeEventHandler, h.Manager.removeClient)
	h.Manager.addClient(userID, client)
	go client.ReadMessage()
	go client.WriteMessage()
}

func (h *Handler) HandleRidersWebsocket(w http.ResponseWriter, r *http.Request) {
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
	client := NewClient(conn, userID, h.Manager.routeEventHandler, h.Manager.removeClient)
	h.Manager.addClient(userID, client)
	go client.ReadMessage()
	go client.WriteMessage()
}

// This queue carries trip request commands that are forwarded to the connected driver websocket.
func (h *Handler) ListenDriverTripRequestsQueue(ctx context.Context) error {
	return h.rabbitmq.Consume(ctx, messaging.DriverCmdTripRequestQueue, func(ctx context.Context, msg messaging.Message) error {
		var envelope messaging.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			return fmt.Errorf("listenDriverTripRequestsQueue: decode amqp message envelope: %w", err)
		}

		if msg.RoutingKey != messaging.DriverCmdTripRequest {
			return nil
		}

		event := contracts.WSMessage{
			Type: msg.RoutingKey,
			Data: envelope.Data,
		}

		if err := h.Manager.SendToClient(envelope.OwnerID, event); err != nil {
			return fmt.Errorf("listenDriverTripRequestsQueue: send to driver %s: %w", envelope.OwnerID, err)
		}

		return nil
	})
}

// Rider notifications are split across multiple queues following the tutorial naming.
// We consume each queue and forward its payload to the rider websocket identified by ownerId.
func (h *Handler) ListenRiderEventsQueue(ctx context.Context) error {
	for _, queue := range []string{
		messaging.NotifyDriverNoDriversFoundQueue,
		messaging.NotifyDriverAssignQueue,
		messaging.NotifyPaymentSessionCreatedQueue,
		messaging.NotifyPaymentSuccessQueue,
	} {
		if err := h.consumeRiderQueue(ctx, queue); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) consumeRiderQueue(ctx context.Context, queue string) error {
	return h.rabbitmq.Consume(ctx, queue, func(ctx context.Context, msg messaging.Message) error {
		var envelope messaging.AmqpMessage
		if err := json.Unmarshal(msg.Body, &envelope); err != nil {
			return fmt.Errorf("listenRiderEventsQueue: decode amqp message envelope: %w", err)
		}

		event := contracts.WSMessage{
			Type: msg.RoutingKey,
			Data: envelope.Data,
		}

		if err := h.Manager.SendToClient(envelope.OwnerID, event); err != nil {
			return fmt.Errorf("listenRiderEventsQueue: send to rider %s: %w", envelope.OwnerID, err)
		}

		return nil
	})
}
