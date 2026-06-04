package browser

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"pc-toolbox/internal/common"
)

// Browser 浏览器类型
type Browser struct {
	Name    string `json:"name"`
	ID      string `json:"id"` // chrome, edge, firefox
	Path    string `json:"path"`
	Enabled bool   `json:"enabled"`
}

// CleanItem 可清理的项目
type CleanItem struct {
	BrowserID   string `json:"browserId"`
	BrowserName string `json:"browserName"`
	Category    string `json:"category"` // cache, cookies, history, downloads, passwords, sessions
	Label       string `json:"label"`
	Path        string `json:"path"`
	FileCount   int    `json:"fileCount"`
	TotalSize   uint64 `json:"totalSize"`
	Checked     bool   `json:"checked"`
	Error       string `json:"error,omitempty"`
}

// CleanResult 清理结果
type CleanResult struct {
	BrowserID   string `json:"browserId"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	Success     bool   `json:"success"`
	FileCount   int    `json:"fileCount"`
	FreedBytes  uint64 `json:"freedBytes"`
	Error       string `json:"error,omitempty"`
}

// BrowserCategory 浏览器数据类型定义
type BrowserCategory struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Paths []string `json:"paths"` // 相对 UserData 目录的路径
}

var browsers = []Browser{
	{Name: "Google Chrome", ID: "chrome", Path: "%LOCALAPPDATA%\\Google\\Chrome\\User Data"},
	{Name: "Microsoft Edge", ID: "edge", Path: "%LOCALAPPDATA%\\Microsoft\\Edge\\User Data"},
	{Name: "Firefox", ID: "firefox", Path: "%APPDATA%\\Mozilla\\Firefox\\Profiles"},
}

var categories = map[string]BrowserCategory{
	"cache":    {ID: "cache", Label: "缓存文件", Paths: []string{"Default\\Cache", "Default\\Code Cache", "Default\\Service Worker\\CacheStorage"}},
	"cookies":  {ID: "cookies", Label: "Cookies", Paths: []string{"Default\\Cookies", "Default\\Network\\Cookies"}},
	"history":  {ID: "history", Label: "浏览历史", Paths: []string{"Default\\History", "Default\\Visited Links"}},
	"sessions": {ID: "sessions", Label: "会话数据", Paths: []string{"Default\\Sessions", "Default\\Session Storage"}},
	"downloads": {ID: "downloads", Label: "下载记录", Paths: []string{"Default\\History Provider Cache"}},
}

var firefoxCategories = map[string]BrowserCategory{
	"cache":    {ID: "cache", Label: "缓存文件", Paths: []string{"cache2", "offlinecache"}},
	"cookies":  {ID: "cookies", Label: "Cookies", Paths: []string{"cookies.sqlite", "cookies.sqlite-wal"}},
	"history":  {ID: "history", Label: "浏览历史", Paths: []string{"places.sqlite", "places.sqlite-wal"}},
	"sessions": {ID: "sessions", Label: "会话数据", Paths: []string{"sessionstore.jsonlz4", "sessionstore-backups"}},
}

// DetectBrowsers 检测已安装的浏览器
func DetectBrowsers() []Browser {
	var detected []Browser
	for _, b := range browsers {
		path := expandEnv(b.Path)
		if _, err := os.Stat(path); err == nil {
			b.Path = path
			b.Enabled = true
			detected = append(detected, b)
		}
	}
	return detected
}

// ScanBrowserData 扫描浏览器可清理数据
func ScanBrowserData() []CleanItem {
	var mu sync.Mutex
	var items []CleanItem
	var wg sync.WaitGroup

	for _, b := range DetectBrowsers() {
		wg.Add(1)
		go func(b Browser) {
			defer wg.Done()
			catMap := categories
			if b.ID == "firefox" {
				catMap = firefoxCategories
			}
			for _, cat := range catMap {
				item := scanBrowserCategory(b, cat)
				mu.Lock()
				items = append(items, item)
				mu.Unlock()
			}
		}(b)
	}

	wg.Wait()
	return items
}

func scanBrowserCategory(b Browser, cat BrowserCategory) CleanItem {
	item := CleanItem{
		BrowserID:   b.ID,
		BrowserName: b.Name,
		Category:    cat.ID,
		Label:       cat.Label,
		Checked:     true,
	}

	var totalSize uint64
	var fileCount int

	if b.ID == "firefox" {
		// Firefox: 每个配置文件都扫描
		profiles, _ := filepath.Glob(filepath.Join(b.Path, "*.default*"))
		for _, profile := range profiles {
			for _, p := range cat.Paths {
				fullPath := filepath.Join(profile, p)
				count, size := scanPath(fullPath)
				fileCount += count
				totalSize += size
			}
		}
	} else {
		// Chrome/Edge: 用户数据目录直接扫描
		for _, p := range cat.Paths {
			fullPath := filepath.Join(b.Path, p)
			count, size := scanPath(fullPath)
			fileCount += count
			totalSize += size
		}
	}

	item.FileCount = fileCount
	item.TotalSize = totalSize
	if fileCount == 0 {
		item.Checked = false
	}
	return item
}

func scanPath(path string) (int, uint64) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, 0
	}
	if !info.IsDir() {
		return 1, uint64(info.Size())
	}

	count := 0
	var total uint64
	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				count++
				total += uint64(info.Size())
			}
		}
		return nil
	})
	return count, total
}

// CleanBrowserData 清理选中项
func CleanBrowserData(items []CleanItem) []CleanResult {
	var mu sync.Mutex
	var results []CleanResult
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(item CleanItem) {
			defer wg.Done()
			result := cleanItem(item)
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(item)
	}

	wg.Wait()
	return results
}

func cleanItem(item CleanItem) CleanResult {
	result := CleanResult{
		BrowserID:  item.BrowserID,
		Category:   item.Category,
		Label:      item.Label,
	}

	info, err := os.Stat(item.Path)
	if err != nil {
		result.Success = true // 不存在也算成功
		return result
	}

	if !info.IsDir() {
		size := uint64(info.Size())
		if err := os.Remove(item.Path); err != nil {
			result.Success = false
			result.Error = err.Error()
			return result
		}
		result.FileCount = 1
		result.FreedBytes = size
		result.Success = true
		return result
	}

	count := 0
	var freed uint64
	filepath.WalkDir(item.Path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				s := uint64(info.Size())
				if err := os.Remove(p); err == nil {
					count++
					freed += s
				}
			}
		}
		return nil
	})

	result.FileCount = count
	result.FreedBytes = freed
	result.Success = true
	return result
}

func expandEnv(path string) string {
	return common.ExpandEnv(path)
}
