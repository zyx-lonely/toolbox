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

// AppVersion 搴旂敤鐗堟湰鍙?const AppVersion = "1.0.0"

// BuildDate 鏋勫缓鏃ユ湡
const BuildDate = "20260623"

func main() {
	// 鍒濆鍖栨棩蹇楃郴缁?	logDir := getLogDir()
	if err := logger.Init(logDir, logger.LevelInfo); err != nil {
		// 鏃ュ織鍒濆鍖栧け璐ワ紝浣跨敤 fmt 杈撳嚭鍒版帶鍒跺彴锛堟鏃?logger 涓嶅彲鐢級
		fmt.Printf("鍒濆鍖栨棩蹇楃郴缁熷け璐? %v\n", err)
	}
	defer logger.GetLogger().Close()

	logger.Info("搴旂敤鍚姩锛岀増鏈? %s", AppVersion)

	// 鍗曞疄渚嬩簰鏂ラ攣 鈥斺€?闃叉澶氬紑
	const mutexName = "Global\\PCToolbox-{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}"
	mu, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if err != nil {
		if err == windows.ERROR_ALREADY_EXISTS {
			logger.Warn("妫€娴嬪埌搴旂敤宸插湪杩愯锛岄€€鍑?)
			return
		}
		logger.Error("鍒涘缓浜掓枼閿佸け璐? %v", err)
		return
	}
	defer windows.CloseHandle(mu)

	logger.Info("鍒濆鍖栧簲鐢?..")
	app := NewApp()

	err = wails.Run(&options.App{
		Title:     "鐢佃剳宸ュ叿绠?v" + AppVersion,
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
		logger.Error("搴旂敤鍚姩澶辫触: %v", err)
	}
}

// getLogDir 鑾峰彇鏃ュ織鐩綍
func getLogDir() string {
	// 浼樺厛浣跨敤鍙墽琛屾枃浠舵墍鍦ㄧ洰褰?	execPath, err := os.Executable()
	if err == nil {
		return filepath.Join(filepath.Dir(execPath), "logs")
	}

	// 鍥為€€鍒扮敤鎴风洰褰?	homeDir, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(homeDir, ".pc-toolbox", "logs")
	}

	return "logs"
}
