package api

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pisign/pisign-backend/api/clock"
	"github.com/pisign/pisign-backend/api/weather"
	"github.com/pisign/pisign-backend/types"
)

func TestNewApiGood(t *testing.T) {
	tables := []struct {
		name string
		api  types.API
	}{
		{"clock", new(clock.API)},
		{"weather", new(weather.API)},
	}

	for _, table := range tables {
		api, err := NewAPI(table.name, nil)
		if err != nil || api == nil {
			t.Errorf("manager.NewAPI failed to create '%s' api instance\n", table.name)
			return
		}
		if reflect.TypeOf(api) != reflect.TypeOf(table.api) {
			t.Errorf("manager.NewAPI created wrong type of API: wanted '%s' but made '%s'", table.name, api.GetName())
			return
		}
	}
}

func TestNewApiBad(t *testing.T) {
	tables := []struct {
		name string
	}{
		{"unknown"},
		{"ClOcK"},
	}

	for _, table := range tables {
		api, err := NewAPI(table.name, nil)
		if api != nil || err == nil {
			t.Errorf("manager.NewAPI falsely created '%s' api instance\n", table.name)
			return
		}
		var e *InvalidAPIError
		if !errors.As(err, &e) {
			t.Errorf("manager.NewAPI threw wrong type of error: %T\n", e)
		}
	}
}
