package main

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsWin "github.com/wailsapp/wails/v2/pkg/options/windows"
	"golang.org/x/sys/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// AppVersion 应用版本号
const AppVersion = "1.2.0"

// BuildDate 构建日期
const BuildDate = "20260604"

func main() {
	// 单实例互斥锁 —— 防止多开
	const mutexName = "Global\\PCToolbox-{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}"
	mu, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if mu == 0 {
		fmt.Printf("创建互斥锁失败: %v\n", err)
		return
	}
	defer windows.CloseHandle(mu)

	// 检查是否已有实例
	if err == windows.ERROR_ALREADY_EXISTS {
		fmt.Println("电脑工具箱已在运行，无法重复启动")
		return
	}

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
			BackdropType:         wailsWin.Mica,
		},
		BackgroundColour: &options.RGBA{R: 245, G: 245, B: 245, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
