package main

import (
	"pc-toolbox/internal/system"
)

// GetBatteryInfo 获取电池信息
func (a *App) GetBatteryInfo() []system.BatteryInfo {
	info, _ := system.GetBatteryInfo()
	return info
}

// HasBattery 检测是否有电池
func (a *App) HasBattery() bool {
	return system.HasBattery()
}
