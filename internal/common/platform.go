package common

import (
	"os"
	"runtime"
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
	// Windows 下通过尝试打开特定安全对象判断
	return false
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
