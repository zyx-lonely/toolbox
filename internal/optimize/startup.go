package optimize

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
	Delay       int    `json:"delay"`  // 延迟秒数，0 表示无延迟
}

// disabledStartupStore 记录被禁用的启动项（用于重新启用）
// 存储格式: 名称 -> {命令, 位置}
var (
	disabledStartupStore = make(map[string]struct {
		command  string
		location string
	})
	disabledStartupMu sync.Mutex
)

// GetStartupItems 获取所有启动项
func GetStartupItems() []StartupItem {
	var items []StartupItem

	// 加载延迟配置
	config := loadStartupConfig()

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

	// 填充延迟值
	for i := range items {
		if delay, ok := config.Delays[items[i].Name]; ok {
			items[i].Delay = delay
		}
	}

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
		root    registry.Key
		subKey  string
	}{
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`},
		{registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\Run`},
	}

	for _, loc := range locations {
		k, err := registry.OpenKey(loc.root, loc.subKey, registry.READ|registry.WRITE)
		if err != nil {
			continue
		}

		// 读取原始值用于恢复
		val, valType, err := k.GetStringValue(name)
		if err == nil && valType == registry.SZ {
			// 保存到恢复存储
			disabledStartupMu.Lock()
			disabledStartupStore[name] = struct {
				command  string
				location string
			}{command: val, location: loc.subKey}
			disabledStartupMu.Unlock()
		}

		if err := k.DeleteValue(name); err == nil {
			k.Close()
			return nil
		}
		k.Close()
	}

	return fmt.Errorf("启动项 %s 未找到", name)
}

// enableStartupItem 重新启用启动项（从备份恢复）
func enableStartupItem(name string) error {
	// 查找被禁用的启动项信息
	disabledStartupMu.Lock()
	info, ok := disabledStartupStore[name]
	disabledStartupMu.Unlock()
	if !ok {
		return fmt.Errorf("未找到启动项 %s 的备份记录，无法重新启用", name)
	}

	// 根据位置恢复
	var root registry.Key
	var subKey string

	switch info.location {
	case "HKCU-Run", "HKCU-RunOnce":
		root = registry.CURRENT_USER
		subKey = info.location
	case "HKLM-Run", "HKLM-RunOnce", "HKLM-WOW-Run":
		root = registry.LOCAL_MACHINE
		subKey = info.location
	default:
		return fmt.Errorf("不支持的启动项位置: %s", info.location)
	}

	k, err := registry.OpenKey(root, subKey, registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表键失败: %v", err)
	}
	defer k.Close()

	if err := k.SetStringValue(name, info.command); err != nil {
		return fmt.Errorf("恢复启动项失败: %v", err)
	}

	// 从备份存储中移除
	disabledStartupMu.Lock()
	delete(disabledStartupStore, name)
	disabledStartupMu.Unlock()
	return nil
}

func estimateImpact(command string) string {
	// 简单启发式判断启动影响
	lowImpact := []string{"update", "scheduler", "notifier", "tray"}
	cmd := strings.ToLower(command)
	for _, low := range lowImpact {
		if strings.Contains(cmd, low) {
			return "low"
		}
	}
	return "medium"
}

// GetStartupDelay 获取启动项延迟（秒）
func GetStartupDelay(name string) int {
	config := loadStartupConfig()
	if config.Delays != nil {
		if delay, ok := config.Delays[name]; ok {
			return delay
		}
	}
	return 0
}

// SetStartupDelay 设置启动项延迟（秒，0 表示无延迟）
func SetStartupDelay(name string, delay int) error {
	config := loadStartupConfig()
	if config.Delays == nil {
		config.Delays = make(map[string]int)
	}
	config.Delays[name] = delay
	return saveStartupConfig(config)
}

// StartupConfig 启动项配置
type StartupConfig struct {
	Delays map[string]int `json:"delays"` // 启动项名称 -> 延迟秒数
}

func loadStartupConfig() StartupConfig {
	configPath := getStartupConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return StartupConfig{Delays: make(map[string]int)}
	}
	var config StartupConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return StartupConfig{Delays: make(map[string]int)}
	}
	if config.Delays == nil {
		config.Delays = make(map[string]int)
	}
	return config
}

func saveStartupConfig(config StartupConfig) error {
	configPath := getStartupConfigPath()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	// 确保目录存在
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	return nil
}

func getStartupConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("APPDATA")
	}
	return filepath.Join(configDir, "pc-toolbox", "startup_config.json")
}
