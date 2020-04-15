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

// ClientAction determines which action should be taken when the client sends a message
type ClientAction string

// Types of actions
// ConfigureAPI sets API configuration data
// ConfigurePosition sets API position data
// Initialize sets both API configuration AND position data
// Delete deletes that specific API instance
// ChangeAPI changes the API type to another (e.g. clock --> weather)
const (
	ConfigureAPI      ClientAction = "ConfigureAPI"
	Initialize        ClientAction = "Init"
	ConfigurePosition ClientAction = "ConfigurePosition"
	Delete            ClientAction = "Delete"
	ChangeAPI         ClientAction = "ChangeAPI"
)

// ClientMessage is the structure of what a client should be sending (in JSON form)
//		Only some fields will be used in any given message, depending on the ClientAction specified
type ClientMessage struct {
	Action   ClientAction
	Position Position
	Config   json.RawMessage
	Name     string
}
