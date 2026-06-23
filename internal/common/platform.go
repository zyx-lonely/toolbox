package common

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Platform 平台抽象接口
type Platform interface {
	GetOSName() string
	GetArchitecture() string
	IsAdmin() bool
	GetTempDir() string
	GetUserHomeDir() string
}

// DefaultPlatform 默认平台实现
type DefaultPlatform struct{}

func (p *DefaultPlatform) GetOSName() string {
	return runtime.GOOS
}

func (p *DefaultPlatform) GetArchitecture() string {
	return runtime.GOARCH
}

func (p *DefaultPlatform) IsAdmin() bool {
	if runtime.GOOS != "windows" {
		return os.Getuid() == 0
	}
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	return err == nil
}

func (p *DefaultPlatform) GetTempDir() string {
	return os.TempDir()
}

func (p *DefaultPlatform) GetUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

// Contains 检查字符串切片是否包含指定字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
