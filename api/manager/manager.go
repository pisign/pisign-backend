package manager

import (
	"fmt"

	"github.com/pisign/pisign-backend/api"
	"github.com/pisign/pisign-backend/api/clock"
	"github.com/pisign/pisign-backend/api/weather"
)

// InvalidAPIError error for missing API
type InvalidAPIError struct {
	APIName string
}

func (e InvalidAPIError) Error() string {
	return fmt.Sprintf("Invalid API: %s", e.APIName)
}

// Connect connects a client to a new API
func Connect(name string) (api.API, error) {
	switch name {
	case "weather":
		return weather.NewAPI(), nil
	case "clock":
		return clock.NewAPI(), nil
	default:
		return nil, InvalidAPIError{name}
	}
}
