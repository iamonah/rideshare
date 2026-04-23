package websockettransport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iamonah/rideshare/services/apigateway/internal/infra/client"
	httpcommon "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/common"
	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/proto/pb/driverpb"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	driverClient *client.DriverClient
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ClientList map[string]*Client

type Client struct {
	Conn    *websocket.Conn
	ID      string
	Message Message
	Egress  chan []byte
}

func NewClient(conn *websocket.Conn, id string) *Client {
	return &Client{
		Conn:   conn,
		ID:     id,
		Egress: make(chan []byte),
	}
}

func NewHandler(dc *client.DriverClient) *Handler {
	return &Handler{
		driverClient: dc,
	}
}

func (c *Client) ReadMessage() {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Websocket closed for client %s: %v", c.ID, err)
			}
			return
		}

		var req Message
		if err := json.Unmarshal(data, &req); err != nil {
			log.Printf("Failed to unmarshal message from client %s: %v", c.ID, err)
			continue
		}

		//handle message based on req.Type
	}
}

func (c *Client) WriteMessage() {}

func (h *Handler) HandleDriversWebsocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("Missing userID in query parameters")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		conn.Close()
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

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: data,
	}

	// a defer function to unregister the driver when the connection is closed
	defer func() {
		h.driverClient.UnregisterDriver(r.Context(), &driverpb.RegisterDriverRequest{
			DriverId:    userID,
			PackageSlug: packageSlug,
		})
	}()
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Failed to send registration message to driver %s: %v", userID, err)
		conn.Close()
		return
	}

	client := NewClient(conn, userID)
	go client.ReadMessage()
	go client.WriteMessage()
}

func (h *Handler) HandleRidersWebsocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("Missing userID in query parameters")
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		conn.Close()
		return
	}

	client := NewClient(conn, userID)
	go client.ReadMessage()
	go client.WriteMessage()
}
