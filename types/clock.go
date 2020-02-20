package types

// ClockResponse main format for data coming out of clock api
type ClockResponse struct {
	Time int64
}

// ClockConfig configuration arguments for clock api
type ClockConfig struct {
	Location string
}
