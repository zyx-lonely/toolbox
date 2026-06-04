package filetools

import (
	"os"
	"path/filepath"
	"sort"
)

type FolderSize struct {
	Path  string  `json:"path"`
	Name  string  `json:"name"`
	Size  uint64  `json:"size"`
	Files int     `json:"files"`
	Pct   float64 `json:"pct"`
}

func AnalyzeFolderSizes(rootPath string, depth int) []FolderSize {
	if depth <= 0 { depth = 1 }
	if depth > 3 { depth = 3 }

	var folders []FolderSize

	entries, _ := os.ReadDir(rootPath)

	// 先计算所有条目（文件+文件夹）的总大小
	var rootTotal uint64
	for _, entry := range entries {
		if !entry.IsDir() {
			info, _ := entry.Info()
			if info != nil { rootTotal += uint64(info.Size()) }
		} else {
			path := filepath.Join(rootPath, entry.Name())
			var size uint64
			var count int
			filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
				if err != nil { return nil }
				if !d.IsDir() {
					info, _ := d.Info()
					if info != nil { size += uint64(info.Size()); count++ }
				}
				return nil
			})
			folders = append(folders, FolderSize{
				Path: path, Name: entry.Name(), Size: size, Files: count,
			})
			rootTotal += size
		}
	}

	if rootTotal == 0 { rootTotal = 1 }
	for i := range folders {
		folders[i].Pct = float64(folders[i].Size) / float64(rootTotal) * 100
	}
	sort.Slice(folders, func(i, j int) bool { return folders[i].Size > folders[j].Size })
	return folders
}
