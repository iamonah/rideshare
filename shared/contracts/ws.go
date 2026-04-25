package contracts

import "encoding/json"

// WSMessage is the message structure for the WebSocket.
type WSMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type WSDriverMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
