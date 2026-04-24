package websockettransport

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	httpcommon "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/common"
	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/proto/pb/driverpb"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
	pongWait              = 10 * time.Second
	pinggInterval         = (9 * pongWait) / 10
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Todo: move this to redis or in-memory store to manage connections across multiple instances of the API Gateway
type ClientList map[string]*Client

type EventHandler func(client *Client, data json.RawMessage) error

type Manager struct {
	clientsList ClientList
	sync.Mutex
	handers map[contracts.Type]EventHandler
}

func NewManager() *Manager {
	return &Manager{
		clientsList: make(ClientList),
		handers:     make(map[contracts.Type]EventHandler),
	}
}

func (m *Manager) RegisterHandler(eventType contracts.Type, handler EventHandler) {
	m.handers[eventType] = handler
}

func (m *Manager) routeEventHandler(c *Client, event contracts.WSMessage) error {
	if handler, ok := m.handers[event.Type]; ok {
		if err := handler(c, event.Data); err != nil {
			log.Printf("error handling event: %v", err)
			return err
		}
		return nil
	} else {
		return fmt.Errorf("no handler for event type: %s", event.Type)
	}
}
func (m *Manager) addClient(id string, client *Client) {
	m.Lock()
	defer m.Unlock()

	if existing, ok := m.clientsList[id]; ok {
		_ = existing.Conn.Close()
	}

	m.clientsList[id] = client
}

func (m *Manager) removeClient(id string) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clientsList[id]; !ok {
		return
	}

	m.clientsList[id].Conn.Close()
	delete(m.clientsList, id)
}

type EventRouter func(client *Client, event contracts.WSMessage) error
type ClientCloseHandler func(id string)

func NewClient(conn *websocket.Conn, id string, onEvent EventRouter, onClose ClientCloseHandler) *Client {
	return &Client{
		Conn:    conn,
		ID:      id,
		Egress:  make(chan []byte),
		onEvent: onEvent,
		onClose: onClose,
	}
}

type Client struct {
	Conn    *websocket.Conn
	ID      string
	Egress  chan []byte
	onEvent EventRouter
	onClose ClientCloseHandler
}

func (c *Client) ReadMessage() {
	defer func() {
		close(c.Egress)
		if c.onClose != nil {
			c.onClose(c.ID) //removes the client and closes the connection
		}
	}()

	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting read deadline: %v", err)
		return
	}

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Websocket closed for client %s: %v", c.ID, err)
			}
			return
		}

		var req contracts.WSMessage
		if err := json.Unmarshal(data, &req); err != nil {
			log.Printf("Failed to unmarshal message from client %s: %v", c.ID, err)
			continue
		}

		if c.onEvent == nil {
			log.Printf("no event router configured for client %s", c.ID)
			return
		}

		if err := c.onEvent(c, req); err != nil {
			log.Printf("error routing event: %v", err)
			return
		}
	}
}

func (c *Client) WriteMessage() {

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
		Type: contracts.Type("register_driver"),
		Data: payload,
	}

	// Todo: unregister the driver when the connection is closed
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Failed to send registration message to driver %s: %v", userID, err)
		conn.Close()
	}

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

	client := NewClient(conn, userID, h.Manager.routeEventHandler, h.Manager.removeClient)
	h.Manager.addClient(userID, client)
	go client.ReadMessage()
	go client.WriteMessage()
}
