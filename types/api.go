package types

import (
	"github.com/google/uuid"
)

// BaseAPI base for all APIs
type BaseAPI struct {
	Position
	Name       string
	ConfigChan chan ClientMessage `json:"-"`
	Pool       Pool               `json:"-"`
	UUID       uuid.UUID
}

// Init Initialization
func (b *BaseAPI) Init(name string, configChan chan ClientMessage, pool Pool) {
	b.Name = name
	b.ConfigChan = configChan
	b.Pool = pool
	b.UUID = uuid.New()
}

// GetName returns the name (or type) of the api
func (b *BaseAPI) GetName() string {
	return b.Name
}

func (b *BaseAPI) GetUUID() uuid.UUID {
	return b.UUID
}

// ConfigurePosition configures position
func (b *BaseAPI) ConfigurePosition(pos Position) {
	b.Position = pos
}
