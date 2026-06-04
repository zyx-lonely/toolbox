package filetools

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"sync"
)

// LargeFile 大文件信息
type LargeFile struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     uint64 `json:"size"`
	Modified string `json:"modified"`
	Type     string `json:"type"` // 文件类型
}

// FindLargeFiles 查找大文件
// rootPath: 扫描根目录
// minSize: 最小大小（字节），默认 100MB
// maxCount: 最多返回数量，默认 100
func FindLargeFiles(rootPath string, minSizeMB int, maxCount int) ([]LargeFile, error) {
	if minSizeMB <= 0 {
		minSizeMB = 100
	}
	if maxCount <= 0 {
		maxCount = 100
	}
	minSize := uint64(minSizeMB) * 1024 * 1024

	var mu sync.Mutex
	var files []LargeFile

	walkDir := rootPath
	err := filepath.WalkDir(walkDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // 跳过无法访问的目录
		}
		if d.IsDir() {
			// 跳过常见系统目录
			name := d.Name()
			if name == "Windows" || name == "System32" || name == "WinSxS" ||
				name == "$Recycle.Bin" || name == "System Volume Information" ||
				stringsHasPrefix(name, ".") {
				return fs.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		size := uint64(info.Size())
		if size >= minSize {
			mu.Lock()
			files = append(files, LargeFile{
				Path:     path,
				Name:     d.Name(),
				Size:     size,
				Modified: info.ModTime().Format("2006-01-02 15:04"),
				Type:     filepath.Ext(d.Name()),
			})
			mu.Unlock()
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 按大小降序排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	if len(files) > maxCount {
		files = files[:maxCount]
	}

	return files, nil
}

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

var _ = fmt.Sprintf
