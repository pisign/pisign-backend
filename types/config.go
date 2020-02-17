package types

import "encoding/json"

// Position information for api widgets
type Position struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	W int    `json:"w"`
	H int    `json:"h"`
	I string `json:"i"`
}

// ClientAction is fun
type ClientAction string

// Types of actions
const (
	Configure ClientAction = "CONFIGURE"
	Delete    ClientAction = "DELETE"
)

// ClientMessage for configuring
type ClientMessage struct {
	Action   ClientAction
	Position Position
	Config   json.RawMessage
}
