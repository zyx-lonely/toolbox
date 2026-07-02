package daily

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// ScreenRecordConfig 录屏配置
type ScreenRecordConfig struct {
	OutputDir string `json:"outputDir"`
	Duration  int    `json:"duration"` // 秒，0 表示手动停止
	FPS       int    `json:"fps"`
	Audio     bool   `json:"audio"`
}

// ScreenRecordResult 录屏结果
type ScreenRecordResult struct {
	FilePath string `json:"filePath"`
	Size     int64  `json:"size"`
	Duration string `json:"duration"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

// 录屏进程
var (
	recordingProcess *exec.Cmd
	isRecording      bool
)

// StartScreenRecording 开始录屏
func StartScreenRecording(config ScreenRecordConfig) (*ScreenRecordResult, error) {
	if isRecording {
		return nil, fmt.Errorf("正在录屏中")
	}

	if config.FPS <= 0 {
		config.FPS = 30
	}
	if config.Duration <= 0 {
		config.Duration = 60 // 默认60秒
	}
	if config.OutputDir == "" {
		config.OutputDir = filepath.Join(os.Getenv("USERPROFILE"), "Videos")
	}

	// 确保输出目录存在
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("创建输出目录失败: %w", err)
	}

	outputFile := filepath.Join(config.OutputDir,
		fmt.Sprintf("录屏_%s.mp4", time.Now().Format("20060102_150405")))

	// 检查是否有 ffmpeg
	ffmpegPath := findFFmpeg()
	if ffmpegPath == "" {
		// 使用 Windows 内置的 screenclip 或 PowerToys
		return startWithPowerShell(config, outputFile)
	}

	// 使用 FFmpeg 录屏
	args := []string{
		"-f", "gdigrab",
		"-framerate", fmt.Sprintf("%d", config.FPS),
		"-i", "desktop",
		"-t", fmt.Sprintf("%d", config.Duration),
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-pix_fmt", "yuv420p",
		outputFile,
	}

	if config.Audio {
		args = append(args, "-f", "dshow", "-i", "audio=Microphone", "-c:a", "aac")
	}

	cmd := exec.Command(ffmpegPath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	isRecording = true
	recordingProcess = cmd

	if err := cmd.Start(); err != nil {
		isRecording = false
		return nil, fmt.Errorf("启动录屏失败: %w", err)
	}

	// 等待录屏完成
	go func() {
		cmd.Wait()
		isRecording = false
		recordingProcess = nil
	}()

	return &ScreenRecordResult{
		FilePath: outputFile,
		Success:  true,
	}, nil
}

// StopScreenRecording 停止录屏
func StopScreenRecording() *ScreenRecordResult {
	if !isRecording || recordingProcess == nil {
		return &ScreenRecordResult{Success: false, Error: "没有正在录屏"}
	}

	// 发送终止信号
	recordingProcess.Process.Kill()
	recordingProcess.Wait()
	isRecording = false
	recordingProcess = nil

	return &ScreenRecordResult{Success: true}
}

// IsRecording 是否正在录屏
func IsRecording() bool {
	return isRecording
}

// startWithPowerShell 使用 PowerShell 录屏
func startWithPowerShell(config ScreenRecordConfig, outputFile string) (*ScreenRecordResult, error) {
	// 使用 PowerShell 的 Windows API 录屏
	psScript := fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing

$screen = [System.Windows.Forms.Screen]::PrimaryScreen
$bounds = $screen.Bounds

$bitmap = New-Object System.Drawing.Bitmap($bounds.Width, $bounds.Height)
$graphics = [System.Drawing.Graphics]::FromImage($bitmap)

$duration = %d
$fps = %d
$frames = $duration * $fps
$output = "%s"

for ($i = 0; $i -lt $frames; $i++) {
    $graphics.CopyFromScreen($bounds.Location, [System.Drawing.Point]::Empty, $bounds.Size)
    Start-Sleep -Milliseconds (1000 / $fps)
}

$bitmap.Save($output)
`, config.Duration, config.FPS, outputFile)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	isRecording = true
	recordingProcess = cmd

	if err := cmd.Start(); err != nil {
		isRecording = false
		return nil, fmt.Errorf("启动录屏失败: %w", err)
	}

	go func() {
		cmd.Wait()
		isRecording = false
		recordingProcess = nil
	}()

	return &ScreenRecordResult{
		FilePath: outputFile,
		Success:  true,
	}, nil
}

func findFFmpeg() string {
	// 检查常见位置
	paths := []string{
		"ffmpeg",
		`C:\ffmpeg\bin\ffmpeg.exe`,
		filepath.Join(os.Getenv("USERPROFILE"), "ffmpeg", "bin", "ffmpeg.exe"),
	}

	for _, p := range paths {
		if _, err := exec.LookPath(p); err == nil {
			return p
		}
	}

	// 尝试直接运行
	cmd := exec.Command("where", "ffmpeg")
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) > 0 {
			return strings.TrimSpace(lines[0])
		}
	}

	return ""
}
