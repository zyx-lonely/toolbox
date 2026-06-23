package main

import (
	"pc-toolbox/internal/common"
	"pc-toolbox/internal/optimize"
)

// ============================================================
//  Windows 优化模块
// ============================================================

// ScanCleanupPaths 扫描可清理的路径
func (a *App) ScanCleanupPaths() []optimize.CleanupTarget {
	return optimize.ScanCleanupPaths()
}

// CleanTargets 执行清理
func (a *App) CleanTargets(paths []string) []optimize.CleanResult {
	return optimize.CleanTargets(paths)
}

// GetStartupItems 获取启动项
func (a *App) GetStartupItems() []optimize.StartupItem {
	return optimize.GetStartupItems()
}

// ToggleStartupItem 切换启动项状态
func (a *App) ToggleStartupItem(name string, enable bool) error {
	return optimize.ToggleStartupItem(name, enable)
}

// GetStartupItemDelay 获取启动项延迟（秒）
func (a *App) GetStartupItemDelay(name string) int {
	return optimize.GetStartupDelay(name)
}

// SetStartupItemDelay 设置启动项延迟（秒，0 表示无延迟）
func (a *App) SetStartupItemDelay(name string, delay int) common.APIResponse {
	err := optimize.SetStartupDelay(name, delay)
	if err != nil {
		return common.NewErrorResponseStr(err.Error())
	}
	return common.NewSuccessResponse(nil)
}

// GetServices 获取 Windows 服务列表
func (a *App) GetServices() ([]optimize.ServiceInfo, error) {
	return optimize.GetServices()
}

// ChangeService 修改服务
func (a *App) ChangeService(name string, action string) error {
	return optimize.ChangeService(name, action)
}

// GetOptimizationProfiles 获取优化方案
func (a *App) GetOptimizationProfiles() []optimize.OptimizationProfile {
	return optimize.GetOptimizationProfiles()
}

// ApplyOptimizationProfile 应用优化方案
func (a *App) ApplyOptimizationProfile(profileName string) ([]optimize.ChangeResult, error) {
	return optimize.ApplyOptimizationProfile(profileName)
}

// RestoreService 还原服务配置
func (a *App) RestoreService(name string) error {
	return optimize.RestoreService(name)
}

// GetServiceBackups 获取已备份的服务列表
func (a *App) GetServiceBackups() []optimize.ServiceBackup {
	return optimize.GetServiceBackups()
}

// ClearServiceBackups 清除服务备份
func (a *App) ClearServiceBackups() {
	optimize.ClearServiceBackups()
}

// GetHostsEntries 获取 hosts 文件
func (a *App) GetHostsEntries() ([]optimize.HostsEntry, error) {
	return optimize.GetHostsEntries()
}

// SaveHostsEntries 保存 hosts 文件
func (a *App) SaveHostsEntries(entries []optimize.HostsEntry) error {
	return optimize.SaveHostsEntries(entries)
}

// CreateRestorePoint 创建系统还原点
func (a *App) CreateRestorePoint(description string) error {
	return optimize.CreateRestorePoint(description)
}

// ScanRegistry 扫描注册表
func (a *App) ScanRegistry() []optimize.RegistryScanResult {
	return optimize.ScanRegistry()
}

// RunHealthCheck 执行健康体检
func (a *App) RunHealthCheck() *optimize.HealthReport {
	return optimize.RunHealthCheck()
}

// InstallContextMenu 安装右键菜单
func (a *App) InstallContextMenu() optimize.ContextMenuStatus {
	return optimize.InstallContextMenu()
}

// InstallContextMenuDir 安装文件夹右键菜单
func (a *App) InstallContextMenuDir() optimize.ContextMenuStatus {
	return optimize.InstallContextMenuDir()
}

// UninstallContextMenu 卸载右键菜单
func (a *App) UninstallContextMenu() optimize.ContextMenuStatus {
	return optimize.UninstallContextMenu()
}

// CheckContextMenu 检查右键菜单状态
func (a *App) CheckContextMenu() optimize.ContextMenuStatus {
	return optimize.CheckContextMenu()
}

// GetWindowsUpdates 获取已安装更新
func (a *App) GetWindowsUpdates() []optimize.UpdateInfo {
	return optimize.GetWindowsUpdates()
}

// GetPendingUpdates 获取待安装更新
func (a *App) GetPendingUpdates() []optimize.UpdateInfo {
	return optimize.GetPendingUpdates()
}

// GetRestorePoints 获取还原点列表
func (a *App) GetRestorePoints() ([]optimize.RestorePointInfo, error) {
	return optimize.GetRestorePoints()
}
