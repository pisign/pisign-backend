package types

// BaseAPI base for all APIs
type BaseAPI struct {
	Position
	Name       string             `json:"name"`
	ConfigChan chan ConfigMessage `json:"-"`
	Pool       Pool               `json:"-"`
}

// Init Initialization
func (b *BaseAPI) Init(name string, configChan chan ConfigMessage, pool Pool) {
	b.Name = name
	b.ConfigChan = configChan
	b.Pool = pool
}

// GetName returns the name (or type) of the api
func (b *BaseAPI) GetName() string {
	return b.Name
}

// ConfigurePosition configures position
func (b *BaseAPI) ConfigurePosition(pos Position) {
	b.Position = pos
}
