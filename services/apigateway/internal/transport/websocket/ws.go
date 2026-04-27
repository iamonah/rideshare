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
	"github.com/iamonah/rideshare/shared/contracts"
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

type EventHandler func(client *Client, event contracts.WSMessage) error

type Manager struct {
	clientsList ClientList
	sync.Mutex
	handers map[string]EventHandler
}

func NewManager() *Manager {
	return &Manager{
		clientsList: make(ClientList),
		handers:     make(map[string]EventHandler),
	}
}

func (m *Manager) RegisterHandler(eventType []string, handler EventHandler) {
	for _, et := range eventType {
		m.handers[et] = handler
	}
}

func (m *Manager) routeEventHandler(c *Client, event contracts.WSMessage) error {
	if handler, ok := m.handers[event.Type]; ok {
		if err := handler(c, event); err != nil {
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

func (m *Manager) getClient(id string) (*Client, bool) {
	m.Lock()
	defer m.Unlock()

	client, ok := m.clientsList[id]
	return client, ok
}

func (m *Manager) SendToClient(id string, event contracts.WSMessage) error {
	client, ok := m.getClient(id)
	if !ok {
		return ErrConnectionNotFound
	}

	return client.Send(event)
}

func (m *Manager) removeClient(id string) {
	m.Lock()
	client, ok := m.clientsList[id]
	if ok {
		delete(m.clientsList, id)
	}
	m.Unlock()

	if ok {
		client.Close()
	}
}

type EventRouter func(client *Client, event contracts.WSMessage) error
type ClientCloseHandler func(id string)

func NewClient(conn *websocket.Conn, id string, onEvent EventRouter, onClose ClientCloseHandler) *Client {
	return &Client{
		Conn:    conn,
		ID:      id,
		Egress:  make(chan contracts.WSMessage),
		done:    make(chan struct{}),
		onEvent: onEvent,
		onClose: onClose,
	}
}

type Client struct {
	Conn    *websocket.Conn
	ID      string
	Egress  chan contracts.WSMessage
	done    chan struct{}
	once    sync.Once
	onEvent EventRouter
	onClose ClientCloseHandler
}

func (c *Client) Close() {
	c.once.Do(func() {
		close(c.done)
		_ = c.Conn.Close()
	})
}

func (c *Client) Send(event contracts.WSMessage) error {
	select {
	case <-c.done:
		return ErrConnectionNotFound
	case c.Egress <- event:
		return nil
	}
}

func (c *Client) ReadMessage() {
	defer func() {
		fmt.Printf("closing websocket connetion id: %s", c.ID)
		if c.onClose != nil {
			c.onClose(c.ID)
		}
	}()

	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("error setting read deadline: %v", err)
		return
	}

	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

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

		if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Printf("error refreshing read deadline for client %s: %v", c.ID, err)
			return
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
	ticker := time.NewTicker(pinggInterval)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-c.done:
			return
		case v := <-c.Egress:
			bytes, err := json.Marshal(v)
			if err != nil {
				log.Printf("error marshaling message for client %s: %v", c.ID, err)
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				log.Printf("error writing message to client %s: %v", c.ID, err)
				return
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error sending ping to client %s: %v", c.ID, err)
				return
			}
		}
	}
}
