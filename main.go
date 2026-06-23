package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsWin "github.com/wailsapp/wails/v2/pkg/options/windows"
	"golang.org/x/sys/windows"

	"pc-toolbox/internal/logger"
)

//go:embed all:frontend/dist
var assets embed.FS

// AppVersion 应用版本号
const AppVersion = "1.0.1"

// BuildDate 构建日期
const BuildDate = "20260623"

func main() {
	// 初始化日志系统
	logDir := getLogDir()
	if err := logger.Init(logDir, logger.LevelInfo); err != nil {
		// 日志初始化失败，使用 fmt 输出到控制台（此时 logger 不可用）
		fmt.Printf("初始化日志系统失败: %v\n", err)
	}
	defer logger.GetLogger().Close()

	logger.Info("应用启动，版本: %s", AppVersion)

	// 单实例互斥锁 —— 防止多开
	const mutexName = "Global\\PCToolbox-{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}"
	mu, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if err != nil {
		if err == windows.ERROR_ALREADY_EXISTS {
			logger.Warn("检测到应用已在运行，退出")
			return
		}
		logger.Error("创建互斥锁失败: %v", err)
		return
	}
	defer windows.CloseHandle(mu)

	logger.Info("初始化应用...")
	app := NewApp()

	err = wails.Run(&options.App{
		Title:     "电脑工具箱 v" + AppVersion,
		Width:     1100,
		Height:    720,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Windows: &wailsWin.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		BackgroundColour: &options.RGBA{R: 245, G: 245, B: 245, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		logger.Error("应用启动失败: %v", err)
	}
}

// getLogDir 获取日志目录
func getLogDir() string {
	// 优先使用可执行文件所在目录
	execPath, err := os.Executable()
	if err == nil {
		return filepath.Join(filepath.Dir(execPath), "logs")
	}

	// 回退到用户目录
	homeDir, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(homeDir, ".pc-toolbox", "logs")
	}

	return "logs"
}
