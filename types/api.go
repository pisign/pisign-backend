package types

import (
	"encoding/json"

	"github.com/pisign/pisign-backend/utils"
)

// BaseAPI base for all APIs
type BaseAPI struct {
	Position
	Name       string
	ConfigChan chan *json.RawMessage `json:"-"`
	Pool       Pool                  `json:"-"`
}

// Init Initialization
func (b *BaseAPI) Init(name string, configChan chan *json.RawMessage, pool Pool) {
	b.Name = name
	b.ConfigChan = configChan
	b.Pool = pool
}

// GetName returns the name (or type) of the api
func (b *BaseAPI) GetName() string {
	return b.Name
}

// ConfigurePosition configures position
func (b *BaseAPI) ConfigurePosition(body *json.RawMessage) {
	err := utils.ParseJSON(*body, &b.Position)
	if err != nil {
		panic("OH NO")
	}
}
