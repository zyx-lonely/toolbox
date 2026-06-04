package common

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppConfig 应用配置
type AppConfig struct {
	Theme    string `json:"theme"`
	Language string `json:"language"`
}

// GetDefaultConfigPath 获取默认配置文件路径
func GetDefaultConfigPath() string {
	execPath, _ := os.Executable()
	return filepath.Join(filepath.Dir(execPath), "config.json")
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
