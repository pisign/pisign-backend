// Package api spins up new api instances to connect to client widgets based on name
package api

import (
	"fmt"

	"github.com/pisign/pisign-backend/api/text"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/types"

	"github.com/pisign/pisign-backend/api/clock"
	"github.com/pisign/pisign-backend/api/weather"
	"github.com/pisign/pisign-backend/api/twitter"
)

// InvalidAPIError error for missing API
type InvalidAPIError struct {
	APIName string
}

// NewAPI returns a new instance of a specific API based on the name
func NewAPI(name string, sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) (types.API, error) {
	switch name {
	case "weather":
		return weather.NewAPI(sockets, pool, id), nil
	case "clock":
		return clock.NewAPI(sockets, pool, id), nil
	case "text":
		return text.NewAPI(sockets, pool, id), nil
	case "twitter":
		return twitter.NewAPI(sockets, pool, id), nil
	default:
		return nil, InvalidAPIError{name}
	}
}

func (e InvalidAPIError) Error() string {
	return fmt.Sprintf("Invalid API: %s", e.APIName)
}
