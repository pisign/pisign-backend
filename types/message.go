package types

type MessageStatus string

const (
	StatusSuccess MessageStatus = "success"
	StatusFailure MessageStatus = "failure"
)

// BaseMessage is the base message we are sending to frontend
// All structs being send to the frontend should inherit from this
type BaseMessage struct {
	Status MessageStatus
	Error  string      `json:",omitempty"`
	Data   interface{} `json:",omitempty"`
}
