package main

import (
	"pc-toolbox/internal/screenshot"
)

// ============================================================
//  截屏工具
// ============================================================

// CaptureAllScreens 截取所有屏幕
func (a *App) CaptureAllScreens() screenshot.CaptureResult {
	return screenshot.CaptureAllScreens()
}

// CaptureScreen 截取指定屏幕
func (a *App) CaptureScreen(displayIndex int) screenshot.CaptureResult {
	return screenshot.CaptureScreen(displayIndex)
}

// OpenInExplorer 在资源管理器中打开文件
func (a *App) OpenInExplorer(path string) error {
	return screenshot.OpenInExplorer(path)
}
