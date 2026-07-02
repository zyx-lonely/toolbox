package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"pc-toolbox/internal/common"
	"pc-toolbox/internal/logger"
	"pc-toolbox/internal/tray"
)

// App struct
type App struct {
	ctx       context.Context
	config    *common.AppConfig
	isQuiting bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg, err := common.LoadConfig(common.GetDefaultConfigPath())
	if err != nil || cfg == nil {
		logger.Warn("加载配置失败，使用默认配置: %v", err)
		cfg = &common.AppConfig{Theme: "light", Language: "zh-CN"}
	}
	a := &App{config: cfg}
	tray.SetApp(a)
	return a
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go tray.Run()

	// 启动配置热重载监听
	watcher := common.NewConfigWatcher(func(cfg *common.AppConfig) {
		a.config = cfg
		// 通知前端配置已变更
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "config:updated", cfg)
		}
	})
	watcher.Start()
}

// MenuShowApp 从系统托盘显示窗口
func (a *App) MenuShowApp() {
	if a.ctx != nil {
		runtime.WindowShow(a.ctx)
	}
}

// MenuQuit 退出应用
func (a *App) MenuQuit() {
	a.isQuiting = true
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

// beforeClose is called when the app is about to close.
func (a *App) beforeClose(ctx context.Context) bool {
	if a.isQuiting {
		return false // 允许退出
	}
	// 关闭时最小化到托盘而不是退出
	runtime.WindowHide(ctx)
	return true // 阻止关闭，改为隐藏
}

// OpenExternalURL 在系统浏览器中打开 URL
func (a *App) OpenExternalURL(url string) error {
	return common.OpenURL(url)
}

// OpenDirectory 在资源管理器中打开目录
func (a *App) OpenDirectory(path string) error {
	return common.OpenURL(path)
}

// GetAppVersion 获取应用版本号
func (a *App) GetAppVersion() string {
	return AppVersion
}

// GetBuildDate 获取构建日期
func (a *App) GetBuildDate() string {
	return BuildDate
}

// GetConfig 获取当前配置
func (a *App) GetConfig() *common.AppConfig {
	return a.config
}

// SaveConfig 保存配置
func (a *App) SaveConfig(cfg *common.AppConfig) error {
	a.config = cfg
	return common.SaveConfig(common.GetDefaultConfigPath(), cfg)
}

// SelectDirectory 选择目录对话框
func (a *App) SelectDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择目录",
	})
	if err != nil {
		return ""
	}
	return dir
}

// SelectFile 选择文件对话框
func (a *App) SelectFile() string {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
	})
	if err != nil {
		return ""
	}
	return file
}
