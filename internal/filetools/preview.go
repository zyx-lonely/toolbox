package filetools

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FilePreviewResult struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	IsText      bool   `json:"isText"`
	IsImage     bool   `json:"isImage"`
	ContentType string `json:"contentType"`
	TextContent string `json:"textContent,omitempty"`
	ImageBase64 string `json:"imageBase64,omitempty"`
	Error       string `json:"error,omitempty"`
	Truncated   bool   `json:"truncated,omitempty"`
	TotalLines  int    `json:"totalLines,omitempty"`
}

var textExts = map[string]bool{
	".txt": true, ".md": true, ".log": true, ".csv": true,
	".json": true, ".xml": true, ".yaml": true, ".yml": true,
	".toml": true, ".ini": true, ".conf": true, ".cfg": true,
	".go": true, ".py": true, ".js": true, ".ts": true,
	".java": true, ".c": true, ".cpp": true, ".h": true,
	".cs": true, ".rb": true, ".php": true, ".swift": true,
	".rs": true, ".lua": true, ".sh": true, ".bat": true,
	".ps1": true, ".html": true, ".css": true, ".scss": true,
	".sql": true, ".vue": true, ".jsx": true, ".tsx": true,
	".env": true, ".gitignore": true, ".dockerignore": true,
}

var imageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".bmp": true, ".webp": true, ".ico": true, ".svg": true,
}

var binaryExts = map[string]bool{
	".exe": true, ".dll": true, ".so": true, ".dylib": true,
	".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
	".mp3": true, ".mp4": true, ".avi": true, ".mkv": true, ".flv": true,
	".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
	".ppt": true, ".pptx": true, ".ttf": true, ".otf": true, ".woff": true,
	".db": true, ".sqlite": true, ".bak": true,
}

func PreviewFile(path string) FilePreviewResult {
	name := filepath.Base(path)
	info, err := os.Stat(path)
	if err != nil {
		return FilePreviewResult{Path: path, Name: name, Error: err.Error()}
	}

	ext := strings.ToLower(filepath.Ext(name))

	result := FilePreviewResult{
		Path: path,
		Name: name,
		Size: info.Size(),
	}

	if imageExts[ext] {
		result.IsImage = true
		result.ContentType = "image"
		data, err := os.ReadFile(path)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		if ext == ".svg" {
			result.ImageBase64 = "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString(data)
		} else {
			result.ImageBase64 = "data:" + getMime(ext) + ";base64," + base64.StdEncoding.EncodeToString(data)
		}
		return result
	}

	if textExts[ext] || (!isBinaryExt(ext) && info.Size() < 2*1024*1024) {
		result.IsText = true
		result.ContentType = "text"
		data, err := os.ReadFile(path)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		content := string(data)
		lines := strings.Split(content, "\n")
		result.TotalLines = len(lines)
		const maxLines = 5000
		if len(lines) > maxLines {
			content = strings.Join(lines[:maxLines], "\n") + fmt.Sprintf("\n\n... (共 %d 行，仅显示前 %d 行)", len(lines), maxLines)
			result.Truncated = true
		}
		result.TextContent = content
		return result
	}

	result.ContentType = "binary"
	result.Error = "不支持预览此类型文件（二进制文件）"
	return result
}

func isBinaryExt(ext string) bool {
	return binaryExts[ext]
}

func getMime(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
