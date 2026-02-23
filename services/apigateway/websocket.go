package apigateway

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/util"
)

type PackageSlug string

var PackageSlugs = make(map[string]PackageSlug)

func newPackageSlug(s string) PackageSlug {
	ps := PackageSlug(s)
	PackageSlugs[strings.ToLower(s)] = ps
	return ps
}

var (
	packageSlugVan    = newPackageSlug("van")
	PackageSlugSUV    = newPackageSlug("suv")
	PackageSlugSedan  = newPackageSlug("sedan")
	PackageSlugLuxury = newPackageSlug("luxury")
)

func parsePackageSlug(s string) (PackageSlug, bool) {
	ps, ok := PackageSlugs[strings.ToLower(s)]
	return ps, ok
}

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

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

		//send to the route handler

		c.Egress <- data
	}
}

func (c *Client) WriteMessage() {

}
func NewClient(conn *websocket.Conn, id string) *Client {
	return &Client{
		Conn:   conn,
		ID:     id,
		Egress: make(chan []byte),
	}
}
func (h *HandlerApiGateway)HandleDriversWebsocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("Missing userID in query parameters")
		return
	}

	//
	packageSlug := r.URL.Query().Get("packageSlug")
	slug, ok := parsePackageSlug(packageSlug)
	if !ok {
		log.Printf("Invalid packageSlug '%s' in query parameters", packageSlug)
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		return
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID:             userID,
			Name:           "Victor",
			ProfilePicture: util.GetRandomAvatar(5),
			CarPlate:       "ABC-123",
			PackageSlug:    slug,
		},
	}
	// client := NewClient(conn, userID)
	// go client.ReadMessage()
	// go client.WriteMessage()
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("Failed to send registration message to driver %s: %v", userID, err)
		conn.Close()
		return
	}
}

func (h *HandlerApiGateway)HandleRidersWebsocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("Missing userID in query parameters")
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		return
	}

	client := NewClient(conn, userID)
	go client.ReadMessage()
	go client.WriteMessage()
}

type Driver struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	ProfilePicture string      `json:"profilePicture"`
	CarPlate       string      `json:"carPlate"`
	PackageSlug    PackageSlug `json:"packageSlug"`
}
