package types

// MessageStatus indicates the broad type of message sent back to the client
type MessageStatus string

// Enumerates different possible message statuses
const (
	StatusSuccess MessageStatus = "success"
	StatusFailure MessageStatus = "failure"
)

// BaseMessage is the base message we are sending to frontend
// All structs being sent to the frontend should inherit from this
type BaseMessage struct {
	Position
	Status MessageStatus
	Error  string      `json:",omitempty"`
	Data   interface{} `json:",omitempty"`
}
