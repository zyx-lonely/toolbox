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

// FixRegistryItems 修复注册表问题
func (a *App) FixRegistryItems(items []optimize.RegistryScanResult) []optimize.RegistryFixResult {
	return optimize.FixRegistryItems(items)
}

// ============================================================
//  环境变量管理
// ============================================================

// GetEnvironmentVariables 获取环境变量
func (a *App) GetEnvironmentVariables(scope string) []optimize.EnvVar {
	return optimize.GetEnvironmentVariables(scope)
}

// SetEnvironmentVariable 设置环境变量
func (a *App) SetEnvironmentVariable(name, value string, scope string) error {
	return optimize.SetEnvironmentVariable(name, value, scope)
}

// DeleteEnvironmentVariable 删除环境变量
func (a *App) DeleteEnvironmentVariable(name string, scope string) error {
	return optimize.DeleteEnvironmentVariable(name, scope)
}

// GetPathVariable 获取 PATH 环境变量
func (a *App) GetPathVariable(scope string) []string {
	return optimize.GetPathVariable(scope)
}

// AddToPathVariable 向 PATH 添加路径
func (a *App) AddToPathVariable(newPath string, scope string) error {
	return optimize.AddToPathVariable(newPath, scope)
}

// RemoveFromPathVariable 从 PATH 移除路径
func (a *App) RemoveFromPathVariable(pathToRemove string, scope string) error {
	return optimize.RemoveFromPathVariable(pathToRemove, scope)
}

// ============================================================
//  磁盘健康监控
// ============================================================

// GetDiskHealthInfo 获取磁盘健康信息
func (a *App) GetDiskHealthInfo() []optimize.DiskHealthInfo {
	return optimize.GetDiskHealthInfo()
}

// GetDiskSpaceInfo 获取磁盘空间信息
func (a *App) GetDiskSpaceInfo() []map[string]interface{} {
	return optimize.GetDiskSpaceInfo()
}

// CheckDiskErrors 检查磁盘错误
func (a *App) CheckDiskErrors(driveLetter string) []string {
	return optimize.CheckDiskErrors(driveLetter)
}

// OptimizeDrive 优化磁盘
func (a *App) OptimizeDrive(driveLetter string) error {
	return optimize.OptimizeDrive(driveLetter)
}

// ============================================================
//  服务依赖分析
// ============================================================

// GetServiceDependencyGraph 获取服务依赖关系图
func (a *App) GetServiceDependencyGraph() (*optimize.ServiceDependencyGraph, error) {
	return optimize.GetServiceDependencyGraph()
}

// AnalyzeServiceImpact 分析禁用服务的影响
func (a *App) AnalyzeServiceImpact(serviceName string) ([]string, error) {
	return optimize.AnalyzeServiceImpact(serviceName)
}

// CanDisableService 检查是否可以安全禁用服务
func (a *App) CanDisableService(serviceName string) (bool, []string, error) {
	return optimize.CanDisableService(serviceName)
}

// ============================================================
//  文件关联管理
// ============================================================

// GetFileAssociations 获取文件关联
func (a *App) GetFileAssociations() []optimize.FileAssociation {
	return optimize.GetFileAssociations()
}

// SetFileAssociation 设置文件关联
func (a *App) SetFileAssociation(extension, progID, command string) error {
	return optimize.SetFileAssociation(extension, progID, command)
}

// RemoveFileAssociation 删除文件关联
func (a *App) RemoveFileAssociation(extension string) error {
	return optimize.RemoveFileAssociation(extension)
}

// GetOpenWithList 获取"打开方式"列表
func (a *App) GetOpenWithList(extension string) []string {
	return optimize.GetOpenWithList(extension)
}

// ============================================================
//  开机自启管理
// ============================================================

// GetAllStartupPrograms 获取所有开机自启程序
func (a *App) GetAllStartupPrograms() ([]optimize.StartupProgram, optimize.StartupGroup) {
	return optimize.GetAllStartupPrograms()
}

// ToggleStartupProgram 启用/禁用开机自启程序
func (a *App) ToggleStartupProgram(name, location string, enable bool) error {
	return optimize.ToggleStartupProgram(name, location, enable)
}

// DeleteStartupProgram 删除开机自启程序
func (a *App) DeleteStartupProgram(name, location string) error {
	return optimize.DeleteStartupProgram(name, location)
}

// SearchStartupPrograms 搜索启动项
func (a *App) SearchStartupPrograms(keyword string) []optimize.StartupProgram {
	return optimize.SearchStartupPrograms(keyword)
}

// ============================================================
//  自动清理调度
// ============================================================

// GetCleanupTasks 获取所有清理任务
func (a *App) GetCleanupTasks() []optimize.CleanupTask {
	return optimize.GetAllTasks()
}

// AddCleanupTask 添加清理任务
func (a *App) AddCleanupTask(task optimize.CleanupTask) error {
	return optimize.AddTask(task)
}

// UpdateCleanupTask 更新清理任务
func (a *App) UpdateCleanupTask(task optimize.CleanupTask) error {
	return optimize.UpdateTask(task)
}

// DeleteCleanupTask 删除清理任务
func (a *App) DeleteCleanupTask(id string) error {
	return optimize.DeleteTask(id)
}

// ToggleCleanupTask 启用/禁用清理任务
func (a *App) ToggleCleanupTask(id string, enabled bool) error {
	return optimize.ToggleTask(id, enabled)
}

// RunCleanupTask 立即执行清理任务
func (a *App) RunCleanupTask(id string) (*optimize.CleanupLog, error) {
	return optimize.RunTask(id)
}

// GetCleanupLogs 获取清理日志
func (a *App) GetCleanupLogs(limit int) []optimize.CleanupLog {
	return optimize.GetLogs(limit)
}

// GetCleanupStats 获取清理统计
func (a *App) GetCleanupStats() map[string]interface{} {
	return optimize.GetSchedulerStats()
}
