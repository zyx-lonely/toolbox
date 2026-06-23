package common

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// GetTotalSize 递归计算目录总大小
func GetTotalSize(path string) (uint64, error) {
	var total uint64
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过无法访问的路径
		}
		if !info.IsDir() {
			total += uint64(info.Size())
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("计算目录大小失败: %v", err)
	}
	return total, nil
}

// FormatBytes 格式化字节数为人类可读格式
func FormatBytes[T uint64 | int64](bytes T) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	val := float64(bytes)
	i := 0
	for val >= 1024 && i < len(units)-1 {
		val /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%.0f %s", val, units[i])
	}
	return fmt.Sprintf("%.2f %s", val, units[i])
}

// IsDir 判断路径是否为目录
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// SafeFileName 移除文件名中的非法字符
func SafeFileName(name string) string {
	invalid := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, c := range invalid {
		result = strings.ReplaceAll(result, c, "_")
	}
	return result
}

// FileExists 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// OpenURL 在系统默认浏览器中打开 URL
func OpenURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL 不能为空")
	}
	c := exec.Command("cmd", "/c", "start", "", url)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := c.Start(); err != nil {
		return fmt.Errorf("打开 URL 失败: %v", err)
	}
	return nil
}

// ReadFileAsBase64 读取文件并返回 Base64 编码的数据 URI
func ReadFileAsBase64(path string) string {
	if path == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	mime := getMimeType(path)
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data)
}

func getMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
