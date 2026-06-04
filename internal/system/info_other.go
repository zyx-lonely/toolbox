//go:build !windows

package system

// GetSystemInfo 获取系统信息（非 Windows 占位）
func GetSystemInfo() (*SystemInfo, error) {
	return &SystemInfo{}, nil
}

// GetHardwareMonitor 获取实时监控数据（非 Windows 占位）
func GetHardwareMonitor() (*HardwareMonitor, error) {
	return &HardwareMonitor{}, nil
}
