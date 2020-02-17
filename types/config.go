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

// ConfigMessage for configuring
type ConfigMessage struct {
	Position
	Config json.RawMessage
}
