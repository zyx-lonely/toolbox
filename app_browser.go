package main

import (
	"pc-toolbox/internal/browser"
)

// ============================================================
//  浏览器数据管理
// ============================================================

// DetectBrowsers 检测浏览器
func (a *App) DetectBrowsers() []browser.Browser {
	return browser.DetectBrowsers()
}

// ScanBrowserData 扫描浏览器可清理数据
func (a *App) ScanBrowserData() []browser.CleanItem {
	return browser.ScanBrowserData()
}

// CleanBrowserData 清理浏览器数据
func (a *App) CleanBrowserData(items []browser.CleanItem) []browser.CleanResult {
	return browser.CleanBrowserData(items)
}
