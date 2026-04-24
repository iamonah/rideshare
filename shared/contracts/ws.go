package contracts

import "encoding/json"

type Type string

// WSMessage is the message structure for the WebSocket.
type WSMessage struct {
	Type Type            `json:"type"`
	Data json.RawMessage `json:"data"`
}

type WSDriverMessage struct {
	Type Type            `json:"type"`
	Data json.RawMessage `json:"data"`
}
