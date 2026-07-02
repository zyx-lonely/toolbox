package tray

import (
	_ "embed"
	"sync"

	"github.com/getlantern/systray"
	"pc-toolbox/internal/common"
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
	mu.Lock()
	defer mu.Unlock()
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

	common.GoSafe(func() {
		systray.Run(onReady, onExit)
	})
}

func onReady() {
	systray.SetTitle("电脑工具箱")
	systray.SetTooltip("电脑工具箱 - 系统监控与优化")

	if len(iconData) > 0 {
		systray.SetIcon(iconData)
	} else {
		// 备用：生成内置 16x16 蓝色图标
		systray.SetIcon(generateFallbackIcon())
	}

	showItem := systray.AddMenuItem("显示主窗口", "显示电脑工具箱主窗口")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("退出", "退出电脑工具箱")

	go func() {
		for {
			select {
			case <-showItem.ClickedCh:
				mu.Lock()
				ctx := appContext
				mu.Unlock()
				if ctx != nil {
					ctx.MenuShowApp()
				}
			case <-quitItem.ClickedCh:
				mu.Lock()
				ctx := appContext
				mu.Unlock()
				if ctx != nil {
					ctx.MenuQuit()
				}
				return
			}
		}
	}()
}

func onExit() {
	// 清理资源
}

// generateFallbackIcon 生成一个简单的 16x16 ICO 格式图标（蓝色方块）
func generateFallbackIcon() []byte {
	// ICO 文件头
	header := []byte{
		0x00, 0x00,             // 保留
		0x01, 0x00,             // 类型: 1 = ICO
		0x01, 0x00,             // 数量: 1 张图片
	}
	// 图片目录 (16x16, 24bit, no palette)
	dir := []byte{
		0x10,                   // 宽度: 16
		0x10,                   // 高度: 16
		0x00,                   // 调色板数: 0
		0x00,                   // 保留
		0x01, 0x00,             // 颜色平面: 1
		0x18, 0x00,             // 位深度: 24
		0x00, 0x00, 0x00, 0x00, // 图像大小 (稍后填充)
		0x16, 0x00, 0x00, 0x00, // 偏移量: 22 (6 + 16)
	}

	// 16x16 蓝色背景像素数据 (BGR 格式) + AND 掩码
	pixelSize := 16 * 16 * 3 // 768 bytes
	maskSize := 16 * 4       // 64 bytes (16 rows * 4 bytes per row)
	imageSize := pixelSize + maskSize

	// 填充目录中的图像大小
	dir[8] = byte(imageSize)
	dir[9] = byte(imageSize >> 8)
	dir[10] = byte(imageSize >> 16)
	dir[11] = byte(imageSize >> 24)

	// 位图信息头 (BITMAPINFOHEADER, 40 bytes)
	infoHeader := []byte{
		0x28, 0x00, 0x00, 0x00, // 头大小: 40
		0x10, 0x00, 0x00, 0x00, // 宽度: 16
		0x20, 0x00, 0x00, 0x00, // 高度: 32 (XOR + AND)
		0x01, 0x00,             // 平面数: 1
		0x18, 0x00,             // 位深度: 24
		0x00, 0x00, 0x00, 0x00, // 压缩: 无
	}

	// 填充信息头中的图像大小
	infoHeader[12] = byte(imageSize)
	infoHeader[13] = byte(imageSize >> 8)
	infoHeader[14] = byte(imageSize >> 16)
	infoHeader[15] = byte(imageSize >> 24)

	// 像素数据: 蓝色方块 (BGR: 0x4A, 0x90, 0xD6)
	pixels := make([]byte, pixelSize)
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			offset := (y*16 + x) * 3
			pixels[offset] = 0x4A     // B
			pixels[offset+1] = 0x90 // G
			pixels[offset+2] = 0xD6 // R
		}
	}

	// AND 掩码 (全 0 = 不透明)
	mask := make([]byte, maskSize)

	// 组装 ICO 文件
	ico := make([]byte, 0, 6+16+40+imageSize)
	ico = append(ico, header...)
	ico = append(ico, dir...)
	ico = append(ico, infoHeader...)
	ico = append(ico, pixels...)
	ico = append(ico, mask...)

	return ico
}
