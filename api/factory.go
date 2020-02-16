// Package api spins up new api instances to connect to client widgets based on name
package api

import (
	"fmt"

	"github.com/pisign/pisign-backend/types"

	"github.com/pisign/pisign-backend/api/clock"
	"github.com/pisign/pisign-backend/api/weather"
)

// InvalidAPIError error for missing API
type InvalidAPIError struct {
	APIName string
}

// NewAPI returns a new instance of a specific API based on the name
func NewAPI(name string, configChan chan types.ConfigMessage, pool types.Pool) (types.API, error) {
	switch name {
	case "weather":
		return weather.NewAPI(configChan, pool), nil
	case "clock":
		return clock.NewAPI(configChan, pool), nil
	default:
		return nil, InvalidAPIError{name}
	}
}

func (e InvalidAPIError) Error() string {
	return fmt.Sprintf("Invalid API: %s", e.APIName)
}
