// Package api spins up new api instances to connect to client widgets based on name
package api

import (
	"fmt"

	"github.com/pisign/pisign-backend/api/clock"
	"github.com/pisign/pisign-backend/api/slideshow"
	"github.com/pisign/pisign-backend/api/text"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/types"

	"github.com/pisign/pisign-backend/api/weather"
)

// InvalidAPIError error for missing API
type InvalidAPIError struct {
	APIName string
}

func factory(name string, sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID, create bool) (types.API, error) {
	switch name {
	case "weather":
		if create {
			return weather.NewAPI(sockets, pool, id), nil
		} else {
			return new(weather.API), nil
		}
	case "clock":
		if create {
			return clock.NewAPI(sockets, pool, id), nil
		} else {
			return new(clock.API), nil
		}
	case "text":
		if create {
			return text.NewAPI(sockets, pool, id), nil
		} else {
			return new(text.API), nil
		}
	case "slideshow":
		if create {
			return slideshow.NewAPI(sockets, pool, id), nil
		} else {
			return new(slideshow.API), nil
		}
	default:
		return nil, InvalidAPIError{name}
	}
}

// NewAPI returns a new instance of a specific API based on the name
func NewAPI(name string, sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) (types.API, error) {
	return factory(name, sockets, pool, id, true)
}

func ValidateAPI(name string) error {
	_, err := factory(name, nil, nil, uuid.New(), false)
	return err
}

func (e InvalidAPIError) Error() string {
	return fmt.Sprintf("Invalid API: %s", e.APIName)
}
