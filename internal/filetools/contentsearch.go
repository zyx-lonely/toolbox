package filetools

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// SearchResult 文件内容搜索结果
type SearchResult struct {
	Path     string `json:"path"`
	FileName string `json:"fileName"`
	Line     int    `json:"line"`
	Content  string `json:"content"`
	FileSize int64  `json:"fileSize"`
}

// SearchFileContent 在目录中搜索文件内容
func SearchFileContent(rootDir string, keyword string, fileTypes string) []SearchResult {
	var results []SearchResult
	var mu sync.Mutex

	exts := parseExtensions(fileTypes)

	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil { return nil }
		if info.IsDir() {
			name := info.Name()
			if name == "Windows" || name == "System32" || name == "$Recycle.Bin" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// 文件类型过滤
		if len(exts) > 0 {
			ext := strings.ToLower(filepath.Ext(path))
			if !containsExt(exts, ext) {
				return nil
			}
		}

		// 检查文件大小（跳过过大文件）
		if info.Size() > 10*1024*1024 {
			return nil
		}

		// 读取文件前100KB
		data, err := os.ReadFile(path)
		if err != nil { return nil }

		// 只检查文本
		text := string(data)
		if len(text) > 100000 {
			text = text[:100000]
		}

		if strings.Contains(strings.ToLower(text), strings.ToLower(keyword)) {
			mu.Lock()
			results = append(results, SearchResult{
				Path: path, FileName: info.Name(),
				FileSize: info.Size(),
			})
			if len(results) >= 50 {
				mu.Unlock()
				return filepath.SkipDir
			}
			mu.Unlock()
		}

		return nil
	})

	return results
}

func parseExtensions(s string) []string {
	if s == "" { return nil }
	parts := strings.Split(s, ",")
	var exts []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			if !strings.HasPrefix(p, ".") {
				p = "." + p
			}
			exts = append(exts, strings.ToLower(p))
		}
	}
	return exts
}

func containsExt(exts []string, ext string) bool {
	for _, e := range exts {
		if e == ext { return true }
	}
	return false
}
