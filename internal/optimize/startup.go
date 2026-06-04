package optimize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// StartupItem 启动项
type StartupItem struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	Location    string `json:"location"` // "HKCU-Run", "HKLM-Run", "StartupFolder"
	Publisher   string `json:"publisher"`
	Enabled     bool   `json:"enabled"`
	Impact      string `json:"impact"` // "high", "medium", "low"
}

// GetStartupItems 获取所有启动项
func GetStartupItems() []StartupItem {
	var items []StartupItem

	// 1. HKCU 注册表启动项
	items = append(items, getRegistryStartupItems(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`, "HKCU-Run")...)

	// 2. HKCU RunOnce
	items = append(items, getRegistryStartupItems(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\RunOnce`, "HKCU-RunOnce")...)

	// 3. HKLM 注册表启动项
	items = append(items, getRegistryStartupItems(registry.LOCAL_MACHINE,
		`Software\Microsoft\Windows\CurrentVersion\Run`, "HKLM-Run")...)

	// 4. HKLM RunOnce
	items = append(items, getRegistryStartupItems(registry.LOCAL_MACHINE,
		`Software\Microsoft\Windows\CurrentVersion\RunOnce`, "HKLM-RunOnce")...)

	// 5. HKLM Wow6432Node (32位程序在64位系统)
	items = append(items, getRegistryStartupItems(registry.LOCAL_MACHINE,
		`Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Run`, "HKLM-WOW-Run")...)

	// 6. 启动文件夹
	items = append(items, getStartupFolderItems()...)

	return items
}

func getRegistryStartupItems(root registry.Key, subKey string, location string) []StartupItem {
	var items []StartupItem

	k, err := registry.OpenKey(root, subKey, registry.READ)
	if err != nil {
		return items
	}
	defer k.Close()

	names, err := k.ReadValueNames(100)
	if err != nil {
		return items
	}

	for _, name := range names {
		val, valType, err := k.GetStringValue(name)
		if err != nil || valType != registry.SZ {

			continue
		}
		items = append(items, StartupItem{
			Name:     name,
			Command:  val,
			Location: location,
			Enabled:  true,
			Impact:   estimateImpact(val),
		})
	}

	return items
}

func getStartupFolderItems() []StartupItem {
	var items []StartupItem

	// 获取当前用户的启动文件夹
	startupDir := filepath.Join(os.Getenv("APPDATA"),
		"Microsoft", "Windows", "Start Menu", "Programs", "Startup")

	entries, err := os.ReadDir(startupDir)
	if err != nil {
		return items
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		items = append(items, StartupItem{
			Name:     entry.Name(),
			Command:  filepath.Join(startupDir, entry.Name()),
			Location: "StartupFolder",
			Enabled:  true,
			Impact:   estimateImpact(entry.Name()),
		})
		_ = info
	}

	return items
}

// ToggleStartupItem 启用/禁用启动项
func ToggleStartupItem(name string, enable bool) error {
	if enable {
		return enableStartupItem(name)
	}
	return disableStartupItem(name)
}

func disableStartupItem(name string) error {
	// 尝试在 4 个位置查找并删除
	locations := []struct {
		root  registry.Key
		subKey string
	}{
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`},
		{registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\Run`},
	}

	for _, loc := range locations {
		k, err := registry.OpenKey(loc.root, loc.subKey, registry.WRITE)
		if err != nil {
			continue
		}
		if err := k.DeleteValue(name); err == nil {
			k.Close()
			return nil
		}
		k.Close()
	}

	return fmt.Errorf("启动项 %s 未找到", name)
}

func enableStartupItem(name string) error {
	// 重新启用较为复杂，实际实现需要保存已禁用的值
	return fmt.Errorf("暂不支持重新启用，请手动操作")
}

func estimateImpact(command string) string {
	// 简单启发式判断启动影响
	lowImpact := []string{"update", "scheduler", "notifier", "tray"}
	cmd := command
	for _, low := range lowImpact {
		if contains(cmd, low) {
			return "low"
		}
	}
	return "medium"
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(s) > 0 && len(substr) > 0 &&
		strings.Contains(s, substr)
}
