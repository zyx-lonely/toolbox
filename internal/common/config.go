package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"pc-toolbox/internal/logger"
)

// AppConfig 应用配置
type AppConfig struct {
	Theme    string `json:"theme"`
	Language string `json:"language"`
}

// GetDefaultConfigPath 获取默认配置文件路径（%APPDATA%/pc-toolbox/）
func GetDefaultConfigPath() string {
	appData, err := os.UserConfigDir()
	if err != nil {
		// 回退到可执行文件目录
		execPath, _ := os.Executable()
		return filepath.Join(filepath.Dir(execPath), "config.json")
	}
	cfgDir := filepath.Join(appData, "pc-toolbox")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		// 目录创建失败，使用当前目录作为回退
		logger.Warn("创建配置目录失败: %v", err)
		return "config.json"
	}
	return filepath.Join(cfgDir, "config.json")
}

var defaultConfig = AppConfig{
	Theme:    "light",
	Language: "zh-CN",
}

// LoadConfig 加载配置
func LoadConfig(path string) (*AppConfig, error) {
	if path == "" {
		path = GetDefaultConfigPath()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return &defaultConfig, nil // 返回默认配置
	}
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &defaultConfig, nil
	}
	return &cfg, nil
}

// SaveConfig 保存配置
func SaveConfig(path string, cfg *AppConfig) error {
	if path == "" {
		path = GetDefaultConfigPath()
	}
	if cfg == nil {
		cfg = &defaultConfig
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ConfigWatcher 配置热重载监听器
type ConfigWatcher struct {
	mu       sync.Mutex
	path     string
	modTime  time.Time
	callback func(*AppConfig)
	stopChan chan struct{}
}

// NewConfigWatcher 创建配置监听器
func NewConfigWatcher(callback func(*AppConfig)) *ConfigWatcher {
	return &ConfigWatcher{
		path:     GetDefaultConfigPath(),
		callback: callback,
		stopChan: make(chan struct{}),
	}
}

// Start 启动监听（每 2 秒检查文件变化）
func (w *ConfigWatcher) Start() {
	// 记录初始修改时间
	if info, err := os.Stat(w.path); err == nil {
		w.modTime = info.ModTime()
	}

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				w.check()
			case <-w.stopChan:
				return
			}
		}
	}()
}

// Stop 停止监听
func (w *ConfigWatcher) Stop() {
	close(w.stopChan)
}

func (w *ConfigWatcher) check() {
	info, err := os.Stat(w.path)
	if err != nil {
		return
	}

	w.mu.Lock()
	changed := info.ModTime().After(w.modTime)
	if changed {
		w.modTime = info.ModTime()
	}
	w.mu.Unlock()

	if changed {
		cfg, err := LoadConfig(w.path)
		if err != nil {
			logger.Warn("热重载配置失败: %v", err)
			return
		}
		logger.Info("配置已热重载")
		if w.callback != nil {
			w.callback(cfg)
		}
	}
}
