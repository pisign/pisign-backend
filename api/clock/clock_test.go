package clock

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNewAPI(t *testing.T) {
	api := NewAPI()
	apiCompare := new(API)
	apiCompare.APIName = "clock"
	apiCompare.Location = "Local"

	if !reflect.DeepEqual(api, apiCompare) {
		t.Errorf("clock.NewAPI: invalid api creation")
	}
}

func TestLoc(t *testing.T) {
	api := NewAPI()
	locLocal := time.Local

	// Test local
	loc := api.loc()
	if loc != locLocal {
		t.Errorf("clock.API.loc(): local returns wrong location")
	}

	// Test unknown defaulting to local
	api.Location = "unknown"
	loc = api.loc()
	if loc != locLocal {
		t.Errorf("clock.API.loc(): undefined location does not return 'Local' time")
	}

	api.Location = "America/Chicago"
	loc = api.loc()
	if loc.String() != "America/Chicago" {
		t.Errorf("clock.API.loc(): 'America/Chicago' returned wrong location")
	}

}

func TestConfigure(t *testing.T) {
	messages := []struct {
		body json.RawMessage
		api  *API
	}{
		{
			[]byte(`{"Location": "America/Chicago"}`),
			NewAPI(),
		},
		{
			nil,
			NewAPI(),
		},
		{
			[]byte(`{"Location": "Atlantis"}`),
			NewAPI(),
		},
	}
	messages[0].api.Location = "America/Chicago"
	messages[1].api.Location = "Local"
	messages[2].api.Location = "Local"

	for _, message := range messages {
		api := NewAPI()
		api.Configure(&message.body)
		if !reflect.DeepEqual(api, message.api) {
			t.Errorf("clock.API.Configure(): Configuration failed with message: %s", message.body)
		}
	}
}
