package main

import (
	"pc-toolbox/internal/filetools"
)

var folderMonitor *filetools.FolderMonitor

// StartFolderMonitor 开始监控文件夹
func (a *App) StartFolderMonitor(path string) error {
	if folderMonitor != nil {
		folderMonitor.Stop()
	}
	folderMonitor = filetools.NewFolderMonitor(path, 1000)
	return folderMonitor.Start()
}

// StopFolderMonitor 停止监控
func (a *App) StopFolderMonitor() {
	if folderMonitor != nil {
		folderMonitor.Stop()
		folderMonitor = nil
	}
}

// GetFolderMonitorEvents 获取监控事件
func (a *App) GetFolderMonitorEvents() []filetools.FolderChangeEvent {
	if folderMonitor == nil {
		return nil
	}
	return folderMonitor.GetEvents()
}

// ClearFolderMonitorEvents 清空监控事件
func (a *App) ClearFolderMonitorEvents() {
	if folderMonitor != nil {
		folderMonitor.ClearEvents()
	}
}

// IsFolderMonitoring 是否正在监控
func (a *App) IsFolderMonitoring() bool {
	if folderMonitor == nil {
		return false
	}
	return folderMonitor.IsRunning()
}
