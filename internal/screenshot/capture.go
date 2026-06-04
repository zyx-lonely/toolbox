package screenshot

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
	"time"

	kb "github.com/kbinani/screenshot"
)

// CaptureResult 截屏结果
type CaptureResult struct {
	Path     string `json:"path"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int64  `json:"size"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

// CaptureAllScreens 截取所有屏幕
func CaptureAllScreens() CaptureResult {
	n := kb.NumActiveDisplays()
	if n <= 0 {
		return CaptureResult{Success: false, Error: "未检测到显示器"}
	}

	// 合并所有显示器为一个图像
	var bounds image.Rectangle
	for i := 0; i < n; i++ {
		b := kb.GetDisplayBounds(i)
		if i == 0 {
			bounds = b
		} else {
			bounds = bounds.Union(b)
		}
	}

	img, err := kb.CaptureRect(bounds)
	if err != nil {
		return CaptureResult{Success: false, Error: err.Error()}
	}

	return saveImage(img, bounds.Dx(), bounds.Dy())
}

// CaptureScreen 截取指定屏幕
func CaptureScreen(displayIndex int) CaptureResult {
	n := kb.NumActiveDisplays()
	if displayIndex < 0 || displayIndex >= n {
		return CaptureResult{Success: false, Error: fmt.Sprintf("显示器索引 %d 无效，共有 %d 个显示器", displayIndex, n)}
	}

	bounds := kb.GetDisplayBounds(displayIndex)
	img, err := kb.CaptureRect(bounds)
	if err != nil {
		return CaptureResult{Success: false, Error: err.Error()}
	}

	return saveImage(img, bounds.Dx(), bounds.Dy())
}

func saveImage(img *image.RGBA, w, h int) CaptureResult {
	saveDir := filepath.Join(os.TempDir(), "pc-toolbox-screenshots")
	os.MkdirAll(saveDir, 0755)

	filename := fmt.Sprintf("screenshot_%s.png", time.Now().Format("20060102_150405"))
	savePath := filepath.Join(saveDir, filename)

	f, err := os.Create(savePath)
	if err != nil {
		return CaptureResult{Success: false, Error: err.Error()}
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return CaptureResult{Success: false, Error: err.Error()}
	}

	info, _ := f.Stat()

	return CaptureResult{
		Path:    savePath,
		Width:   w,
		Height:  h,
		Size:    info.Size(),
		Success: true,
	}
}

// OpenInExplorer 在资源管理器中打开
func OpenInExplorer(path string) error {
	cmd := exec.Command("explorer", "/select,", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}
