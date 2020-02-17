package types

// DataObject holds the data from the external API
type DataObject interface {
	// Build builds the data object
	Update(interface{}) error

	// Transform turns the data object into a front-end parsable object
	Transform() interface{}
}

// API is the entrance point of all apis to connect to a client
type API interface {
	// Configure settings from raw json message
	Configure(message ConfigMessage) error

	// Main loop that faciliates interaction between outside world and the client widet
	Run(w Socket)

	// Data gets the data to send
	Data() interface{}

	GetName() string
}

// Socket interface, needed to avoid circular dependency with Socket package
type Socket interface {
	// Read information from the client
	Read()

	// Send information to the client
	Send(interface{})

	// Close the socket
	Close() chan bool

	// SendErrorMessage sends error message
	SendErrorMessage(string)
}

// Pool pool
type Pool interface {
	Start()
	Register(API)
	Unregister(API)
	Save()
}
