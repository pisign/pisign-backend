package types

// ClockResponse main format for data coming out of clock api
type ClockResponse struct {
	Time string
}

// ClockConfig configuration arguments for clock api
type ClockConfig struct {
	Location string
}
