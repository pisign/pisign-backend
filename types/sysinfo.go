package types

type SysInfoResponse struct {
	MemUsed         uint64
	MemTotal        uint64
	MemUsedPercent  float64
	DiskUsed        uint64
	DiskTotal       uint64
	DiskFree        uint64
	DiskUsedPercent float64
}

type SysInfoConfig struct {
}
