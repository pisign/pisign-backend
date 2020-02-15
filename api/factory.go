package api

import (
	"fmt"

	"github.com/pisign/pisign-backend/api/clock"
	"github.com/pisign/pisign-backend/api/weather"
	"github.com/pisign/pisign-backend/types"
)

// InvalidAPIError error for missing API
type InvalidAPIError struct {
	APIName string
}

// NewAPI returns a new instance of a specific API based on the name
func NewAPI(name string) (types.API, error) {
	switch name {
	case "weather":
		return weather.NewAPI(), nil
	case "clock":
		return clock.NewAPI(), nil
	default:
		return nil, InvalidAPIError{name}
	}
}

func (e InvalidAPIError) Error() string {
	return fmt.Sprintf("Invalid API: %s", e.APIName)
}
