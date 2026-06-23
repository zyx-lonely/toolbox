package system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pc-toolbox/internal/common"
	"pc-toolbox/internal/logger"
)

// BrowserExtension 浏览器扩展信息
type BrowserExtension struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Browser     string `json:"browser"` // chrome, edge, firefox
	Enabled     bool   `json:"enabled"`
	ID          string `json:"id"`
}

// GetBrowserExtensions 获取浏览器扩展列表
func GetBrowserExtensions() []BrowserExtension {
	var extensions []BrowserExtension

	// Chrome 扩展
	chromeExts := getChromeExtensions("chrome")
	extensions = append(extensions, chromeExts...)

	// Edge 扩展 (基于 Chromium，路径类似)
	edgeExts := getChromeExtensions("edge")
	extensions = append(extensions, edgeExts...)

	// Firefox 扩展
	firefoxExts := getFirefoxExtensions()
	extensions = append(extensions, firefoxExts...)

	return extensions
}

// getChromeExtensions 获取 Chrome/Edge 扩展
func getChromeExtensions(browser string) []BrowserExtension {
	var exts []BrowserExtension
	var extPath string

	home := os.Getenv("USERPROFILE")

	switch browser {
	case "chrome":
		extPath = filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Extensions")
	case "edge":
		extPath = filepath.Join(home, "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Extensions")
	default:
		return exts
	}

	entries, err := os.ReadDir(extPath)
	if err != nil {
		return exts
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		extID := entry.Name()
		// 跳过非扩展目录 (扩展ID通常是32个字符)
		if len(extID) != 32 {
			continue
		}

		// 读取扩展版本目录
		extDir := filepath.Join(extPath, extID)
		versions, err := os.ReadDir(extDir)
		if err != nil {
			continue
		}

		for _, ver := range versions {
			if !ver.IsDir() {
				continue
			}
			manifestPath := filepath.Join(extDir, ver.Name(), "manifest.json")
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				continue
			}

			var manifest struct {
				Name        string `json:"name"`
				Version     string `json:"version"`
				Description string `json:"description"`
			}
			if err := json.Unmarshal(data, &manifest); err != nil {
				continue
			}

			exts = append(exts, BrowserExtension{
				Name:        manifest.Name,
				Version:     manifest.Version,
				Description: manifest.Description,
				Path:        filepath.Join(extDir, ver.Name()),
				Browser:     browser,
				Enabled:     true, // Windows 下默认启用
				ID:          extID,
			})
			break // 只取最新版本
		}
	}

	return exts
}

// getFirefoxExtensions 获取 Firefox 扩展
func getFirefoxExtensions() []BrowserExtension {
	var exts []BrowserExtension

	home := os.Getenv("APPDATA")
	extPath := filepath.Join(home, "Mozilla", "Firefox", "Profiles")

	profiles, err := os.ReadDir(extPath)
	if err != nil {
		return exts
	}

	for _, profile := range profiles {
		if !profile.IsDir() || !strings.HasSuffix(profile.Name(), ".default-release") {
			continue
		}

		extDir := filepath.Join(extPath, profile.Name(), "extensions")
		entries, err := os.ReadDir(extDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".xpi") {
				continue
			}

			exts = append(exts, BrowserExtension{
				Name:    strings.TrimSuffix(entry.Name(), ".xpi"),
				Version: "unknown",
				Path:    filepath.Join(extDir, entry.Name()),
				Browser: "firefox",
				Enabled: true,
			})
		}
	}

	return exts
}

// DisableBrowserExtension 禁用浏览器扩展
func DisableBrowserExtension(browser, extID string) common.APIResponse {
	// 通过重命名扩展目录来禁用
	var extPath string
	home := os.Getenv("USERPROFILE")

	switch browser {
	case "chrome":
		extPath = filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Extensions", extID)
	case "edge":
		extPath = filepath.Join(home, "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Extensions", extID)
	default:
		return common.NewErrorResponseStr("不支持的浏览器")
	}

	disabledPath := extPath + ".disabled"
	if err := os.Rename(extPath, disabledPath); err != nil {
		return common.NewErrorResponseStr(fmt.Sprintf("禁用扩展失败: %v", err))
	}

	logger.Info("禁用浏览器扩展: %s (%s)", extID, browser)
	return common.NewSuccessResponse(nil)
}

// EnableBrowserExtension 启用浏览器扩展
func EnableBrowserExtension(browser, extID string) common.APIResponse {
	var disabledPath string
	home := os.Getenv("USERPROFILE")

	switch browser {
	case "chrome":
		disabledPath = filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Extensions", extID+".disabled")
	case "edge":
		disabledPath = filepath.Join(home, "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Extensions", extID+".disabled")
	default:
		return common.NewErrorResponseStr("不支持的浏览器")
	}

	extPath := strings.TrimSuffix(disabledPath, ".disabled")
	if err := os.Rename(disabledPath, extPath); err != nil {
		return common.NewErrorResponseStr(fmt.Sprintf("启用扩展失败: %v", err))
	}

	logger.Info("启用浏览器扩展: %s (%s)", extID, browser)
	return common.NewSuccessResponse(nil)
}

// RemoveBrowserExtension 删除浏览器扩展
func RemoveBrowserExtension(browser, extID string) common.APIResponse {
	var extPath string
	home := os.Getenv("USERPROFILE")

	switch browser {
	case "chrome":
		extPath = filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Extensions", extID)
	case "edge":
		extPath = filepath.Join(home, "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Extensions", extID)
	default:
		return common.NewErrorResponseStr("不支持的浏览器")
	}

	if err := os.RemoveAll(extPath); err != nil {
		return common.NewErrorResponseStr(fmt.Sprintf("删除扩展失败: %v", err))
	}

	logger.Info("删除浏览器扩展: %s (%s)", extID, browser)
	return common.NewSuccessResponse(nil)
}
