package optimize

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"pc-toolbox/internal/common"
)

// CleanupTarget 清理目标
type CleanupTarget struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Risk        string `json:"risk"`
	Browser     string `json:"browser,omitempty"`
	FileCount   int    `json:"fileCount"`
	TotalSize   uint64 `json:"totalSize"`
	Checked     bool   `json:"checked"`
	Error       string `json:"error,omitempty"`
}

// CleanResult 清理结果
type CleanResult struct {
	Target      string `json:"target"`
	Description string `json:"description"`
	FileCount   int    `json:"fileCount"`
	FreedBytes  uint64 `json:"freedBytes"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}

var defaultCleanupPaths = []CleanupTarget{
	{Path: "%TEMP%", Description: "用户临时文件", Risk: "low"},
	{Path: "%WINDIR%\\Temp", Description: "系统临时文件", Risk: "low"},
	{Path: "%WINDIR%\\Prefetch", Description: "预读取文件", Risk: "low"},
	{Path: "%WINDIR%\\SoftwareDistribution\\Download", Description: "Windows 更新缓存", Risk: "medium"},
	{Path: "%LOCALAPPDATA%\\Microsoft\\Windows\\Explorer", Description: "缩略图缓存", Risk: "low"},
	{Path: "%USERPROFILE%\\Recent", Description: "最近文档记录", Risk: "low"},
	{Path: "%LOCALAPPDATA%\\Temp", Description: "本地应用临时文件", Risk: "low"},
	{Path: "%LOCALAPPDATA%\\Google\\Chrome\\User Data\\Default\\Cache", Description: "Chrome 浏览器缓存", Risk: "low", Browser: "chrome"},
	{Path: "%LOCALAPPDATA%\\Microsoft\\Edge\\User Data\\Default\\Cache", Description: "Edge 浏览器缓存", Risk: "low", Browser: "edge"},
}

// ScanCleanupPaths 扫描清理路径，计算可释放空间（并发 + 超时保护）
func ScanCleanupPaths() []CleanupTarget {
	var mu sync.Mutex
	var targets []CleanupTarget
	var wg sync.WaitGroup

	for _, t := range defaultCleanupPaths {
		wg.Add(1)
		go func(t CleanupTarget) {
			defer wg.Done()
			result := scanSinglePath(t)
			mu.Lock()
			targets = append(targets, result)
			mu.Unlock()
		}(t)
	}

	wg.Wait()

	// 按风险等级排序：low 在前
	sortTargets(targets)
	return targets
}

func scanSinglePath(t CleanupTarget) CleanupTarget {
	expanded := expandEnv(t.Path)
	if expanded == "" {
		return CleanupTarget{Path: t.Path, Description: t.Description, Risk: t.Risk, Error: "环境变量解析失败"}
	}

	info, err := os.Stat(expanded)
	if err != nil {
		return CleanupTarget{
			Path:        expanded,
			Description: t.Description,
			Risk:        t.Risk,
			Browser:     t.Browser,
			Checked:     false,
			Error:       "路径不存在或无法访问",
		}
	}

	if !info.IsDir() {
		return CleanupTarget{
			Path:        expanded,
			Description: t.Description,
			Risk:        t.Risk,
			Browser:     t.Browser,
			FileCount:   1,
			TotalSize:   uint64(info.Size()),
			Checked:     t.Risk == "low",
		}
	}

	// 使用 WalkDir（更高效）+ 超时管道
	done := make(chan struct{})
	count := 0
	var totalSize uint64

	go func() {
		// 限制扫描文件数量以提高性能
		const maxFiles = 5000
		filepath.WalkDir(expanded, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				if os.IsPermission(err) {
					return fs.SkipDir
				}
				return nil
			}
			if !d.IsDir() {
				info, err := d.Info()
				if err == nil {
					count++
					totalSize += uint64(info.Size())
					if count >= maxFiles {
						return fs.SkipDir
					}
				}
			}
			return nil
		})
		close(done)
	}()

	// 30秒超时，防止卡死
	select {
	case <-done:
	case <-time.After(30 * time.Second):
	}

	return CleanupTarget{
		Path:        expanded,
		Description: t.Description,
		Risk:        t.Risk,
		Browser:     t.Browser,
		FileCount:   count,
		TotalSize:   totalSize,
		Checked:     t.Risk == "low",
	}
}

// CleanTargets 执行清理操作
func CleanTargets(paths []string) []CleanResult {
	var mu sync.Mutex
	var results []CleanResult
	var wg sync.WaitGroup

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			result := cleanTarget(p)
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(path)
	}

	wg.Wait()
	return results
}

func cleanTarget(path string) CleanResult {
	result := CleanResult{Target: path, Success: true}

	info, err := os.Stat(path)
	if err != nil {
		result.Success = false
		result.Error = "路径不存在: " + err.Error()
		return result
	}

	if !info.IsDir() {
		size := uint64(info.Size())
		err := os.Remove(path)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			return result
		}
		result.FileCount = 1
		result.FreedBytes = size
		return result
	}

	// 顺序执行清理（避免并发竞争）
	count := 0
	var freed uint64

	// 第一遍：删除所有文件
	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return fs.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				s := uint64(info.Size())
				if e := os.Remove(p); e == nil {
					count++
					freed += s
				}
			}
		}
		return nil
	})

	// 第二遍：移除空目录（从下往上）
	removeEmptyDirs(path)

	result.FileCount = count
	result.FreedBytes = freed
	return result
}

func removeEmptyDirs(root string) {
	// 先递归到底，再从底部往上删除空目录
	var dirs []string
	filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && p != root {
			dirs = append([]string{p}, dirs...) // 反转顺序，先删子目录
		}
		return nil
	})

	for _, dir := range dirs {
		entries, _ := os.ReadDir(dir)
		if len(entries) == 0 {
			os.Remove(dir) // 忽略错误
		}
	}
}

func expandEnv(path string) string {
	return common.ExpandEnv(path)
}

func sortTargets(targets []CleanupTarget) {
	riskOrder := map[string]int{"low": 0, "medium": 1, "high": 2}
	for i := 0; i < len(targets); i++ {
		for j := i + 1; j < len(targets); j++ {
			if riskOrder[targets[i].Risk] > riskOrder[targets[j].Risk] {
				targets[i], targets[j] = targets[j], targets[i]
			}
		}
	}
}
