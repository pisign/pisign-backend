package types

import (
	"encoding/json"
)

// Widget interface
type Widget interface {
	json.Unmarshaler
	Send(interface{})
	Read()
	Close() chan bool
}
