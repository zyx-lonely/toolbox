package system

// CPUInfo CPU 信息
type CPUInfo struct {
	Name         string  `json:"name"`
	Cores        int     `json:"cores"`
	LogicalCores int     `json:"logicalCores"`
	BaseClock    string  `json:"baseClock"`
	Usage        float64 `json:"usage"`
	Temperature  float64 `json:"temperature"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Usage     float64 `json:"usage"`
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Label      string  `json:"label"`
	FileSystem string  `json:"fileSystem"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Free       uint64  `json:"free"`
	Usage      float64 `json:"usage"`
}

// NetworkAdapter 网络适配器信息
type NetworkAdapter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	IP          string `json:"ip"`
	MAC         string `json:"mac"`
	IsConnected bool   `json:"isConnected"`
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	Adapters      []NetworkAdapter `json:"adapters"`
	DownloadSpeed uint64           `json:"downloadSpeed"`
	UploadSpeed   uint64           `json:"uploadSpeed"`
}

// OSInfo 操作系统信息
type OSInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	BuildNumber  string `json:"buildNumber"`
	Architecture string `json:"architecture"`
	InstallDate  string `json:"installDate"`
	Uptime       string `json:"uptime"`
}

// SystemInfo 完整系统信息
type SystemInfo struct {
	OS      OSInfo        `json:"os"`
	CPU     CPUInfo       `json:"cpu"`
	Memory  MemoryInfo    `json:"memory"`
	Disks   []DiskInfo    `json:"disks"`
	Network NetworkInfo   `json:"network"`
}

// HardwareMonitor 硬件实时监控数据
type HardwareMonitor struct {
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	MemoryUsed  uint64  `json:"memoryUsed"`
	MemoryTotal uint64  `json:"memoryTotal"`
	DiskIO      uint64  `json:"diskIO"`
	NetDown     uint64  `json:"netDown"`
	NetUp       uint64  `json:"netUp"`
	Uptime      string  `json:"uptime"`
}
