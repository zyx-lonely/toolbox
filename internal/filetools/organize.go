package filetools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}

// OrganizeRule 归类规则
type OrganizeRule struct {
	Mode     string `json:"mode"`     // "extension", "date", "size"
	Target   string `json:"target"`   // "year", "month", "year-month" (for date mode)
	Move     bool   `json:"move"`     // true=移动, false=复制
	SortInto string `json:"sortInto"` // 归类到的根目录
}

// OrganizePreview 归类预览项
type OrganizePreview struct {
	SourcePath string `json:"sourcePath"`
	DestPath   string `json:"destPath"`
	SourceName string `json:"sourceName"`
	FolderName string `json:"folderName"` // 目标子文件夹名
	FileSize   uint64 `json:"fileSize"`
}

// OrganizeResult 归类结果
type OrganizeResult struct {
	SourcePath string `json:"sourcePath"`
	DestPath   string `json:"destPath"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

// PreviewOrganize 预览归类结果
func PreviewOrganize(dir string, rule OrganizeRule) ([]OrganizePreview, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var previews []OrganizePreview
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}

		folder := classifyFile(entry.Name(), info, rule)
		if folder == "" {
			continue
		}

		destPath := filepath.Join(rule.SortInto, folder, entry.Name())
		previews = append(previews, OrganizePreview{
			SourcePath: filepath.Join(dir, entry.Name()),
			DestPath:   destPath,
			SourceName: entry.Name(),
			FolderName: folder,
			FileSize:   uint64(info.Size()),
		})
	}
	return previews, nil
}

// ExecuteOrganize 执行归类
func ExecuteOrganize(dir string, rule OrganizeRule) ([]OrganizeResult, error) {
	previews, err := PreviewOrganize(dir, rule)
	if err != nil {
		return nil, err
	}

	var results []OrganizeResult
	for _, p := range previews {
		if err := os.MkdirAll(filepath.Dir(p.DestPath), 0755); err != nil {
			results = append(results, OrganizeResult{
				SourcePath: p.SourcePath, DestPath: p.DestPath,
				Success: false, Error: err.Error(),
			})
			continue
		}

		var opErr error
		if rule.Move {
			opErr = os.Rename(p.SourcePath, p.DestPath)
		} else {
			opErr = copyFile(p.SourcePath, p.DestPath)
		}

		results = append(results, OrganizeResult{
			SourcePath: p.SourcePath, DestPath: p.DestPath,
			Success: opErr == nil, Error: getErrMsg(opErr),
		})
	}
	return results, nil
}

func classifyFile(name string, info os.FileInfo, rule OrganizeRule) string {
	switch rule.Mode {
	case "extension":
		ext := strings.TrimPrefix(filepath.Ext(name), ".")
		if ext == "" {
			return "无扩展名"
		}
		return strings.ToUpper(ext)

	case "date":
		modTime := info.ModTime()
		switch rule.Target {
		case "year":
			return fmt.Sprintf("%d", modTime.Year())
		case "month":
			return fmt.Sprintf("%02d月", modTime.Month())
		case "year-month":
			return fmt.Sprintf("%d年%02d月", modTime.Year(), modTime.Month())
		default:
			return fmt.Sprintf("%d", modTime.Year())
		}

	case "size":
		s := info.Size()
		switch {
		case s < 1024:
			return "<1KB"
		case s < 1024*1024:
			return "1KB-1MB"
		case s < 100*1024*1024:
			return "1MB-100MB"
		default:
			return ">100MB"
		}
	}
	return "其他"
}

func getErrMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// 忽略不必要的 import
var _ = time.Second
