package types

// ClockOut main format for data coming out of clock api
type ClockOut struct {
	Time string
}

// ClockConfig configuration arguments for clock api
type ClockConfig struct {
	Location string
}
