package tray

import (
	_ "embed"
	"sync"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var iconData []byte

var (
	appContext interface {
		MenuShowApp()
		MenuQuit()
	}
	mu      sync.Mutex
	started bool
)

// SetApp sets the app reference for tray callbacks
func SetApp(a interface{ MenuShowApp(); MenuQuit() }) {
	appContext = a
}

// Run starts the system tray in a goroutine
func Run() {
	mu.Lock()
	if started {
		mu.Unlock()
		return
	}
	started = true
	mu.Unlock()

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("电脑工具箱")
	systray.SetTooltip("电脑工具箱 - 系统监控与优化")

	if len(iconData) > 0 {
		systray.SetIcon(iconData)
	} else {
		// 备用：1x1 透明像素
		systray.SetIcon([]byte{0x89, 0x50, 0x4E, 0x47})
	}

	showItem := systray.AddMenuItem("显示主窗口", "显示电脑工具箱主窗口")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("退出", "退出电脑工具箱")

	go func() {
		for {
			select {
			case <-showItem.ClickedCh:
				if appContext != nil {
					appContext.MenuShowApp()
				}
			case <-quitItem.ClickedCh:
				if appContext != nil {
					appContext.MenuQuit()
				}
				return
			}
		}
	}()
}

func onExit() {
	// 清理资源
}
