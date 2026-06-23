package main

import (
	"pc-toolbox/internal/common"
	"pc-toolbox/internal/system"
)

// ============================================================
//  系统信息模块
// ============================================================

// GetSystemInfo 获取完整系统信息
func (a *App) GetSystemInfo() (*system.SystemInfo, error) {
	return system.GetSystemInfo()
}

// GetHardwareMonitor 获取实时硬件监控数据
func (a *App) GetHardwareMonitor() (*system.HardwareMonitor, error) {
	return system.GetHardwareMonitor()
}

// GetShortcutKeys 获取系统快捷键列表
func (a *App) GetShortcutKeys() []system.ShortcutKey {
	return system.GetShortcutKeys()
}

// GetPowerPlans 获取电源方案列表
func (a *App) GetPowerPlans() []system.PowerPlan {
	return system.GetPowerPlans()
}

// SetPowerPlan 切换电源方案
func (a *App) SetPowerPlan(guid string) error {
	return system.SetPowerPlan(guid)
}

// GetActivationInfo 获取激活信息
func (a *App) GetActivationInfo() system.ActivationInfo {
	return system.GetActivationInfo()
}

// GetActivationTools 获取激活工具列表
func (a *App) GetActivationTools() []system.ActivationTool {
	return system.GetActivationTools()
}

// GetKMSMethods 获取 KMS 激活方法
func (a *App) GetKMSMethods() []system.KMSMethod {
	return system.GetKMSMethods()
}

// GetTemperatures 获取硬件温度
func (a *App) GetTemperatures() ([]system.TemperatureInfo, error) {
	return system.GetTemperatures()
}

// GetBrowserExtensions 获取浏览器扩展列表
func (a *App) GetBrowserExtensions() []system.BrowserExtension {
	return system.GetBrowserExtensions()
}

// DisableBrowserExtension 禁用浏览器扩展
func (a *App) DisableBrowserExtension(browser string, extID string) common.APIResponse {
	return system.DisableBrowserExtension(browser, extID)
}

// EnableBrowserExtension 启用浏览器扩展
func (a *App) EnableBrowserExtension(browser string, extID string) common.APIResponse {
	return system.EnableBrowserExtension(browser, extID)
}

// RemoveBrowserExtension 删除浏览器扩展
func (a *App) RemoveBrowserExtension(browser string, extID string) common.APIResponse {
	return system.RemoveBrowserExtension(browser, extID)
}


