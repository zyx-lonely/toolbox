package filetools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// RenameRule 重命名规则
type RenameRule struct {
	Pattern     string `json:"pattern"`     // 命名模式: {index} {name} {date} {random}
	StartIndex  int    `json:"startIndex"`  // 起始序号
	Padding     int    `json:"padding"`     // 序号位数补零
	ReplaceFrom string `json:"replaceFrom"` // 替换原名称中的文本
	ReplaceTo   string `json:"replaceTo"`   // 替换为
	FileFilter  string `json:"fileFilter"`  // 文件扩展名过滤(逗号分隔)
}

// RenamePreview 重命名预览项
type RenamePreview struct {
	OriginalPath string `json:"originalPath"`
	NewPath      string `json:"newPath"`
	OriginalName string `json:"originalName"`
	NewName      string `json:"newName"`
	Index        int    `json:"index"`
}

// BatchRenamePreview 批量重命名预览
func BatchRenamePreview(dir string, rule RenameRule) ([]RenamePreview, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// 解析过滤条件
	var allowedExts []string
	if rule.FileFilter != "" {
		for _, ext := range strings.Split(rule.FileFilter, ",") {
			allowedExts = append(allowedExts, strings.TrimSpace(strings.ToLower(ext)))
		}
	}

	var previews []RenamePreview
	index := rule.StartIndex

	for i, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 过滤扩展名
		if len(allowedExts) > 0 {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if !contains(allowedExts, ext) {
				continue
			}
		}

		newName := applyPattern(entry.Name(), rule, index)
		previews = append(previews, RenamePreview{
			OriginalPath: filepath.Join(dir, entry.Name()),
			NewPath:      filepath.Join(dir, newName),
			OriginalName: entry.Name(),
			NewName:      newName,
			Index:        index,
		})
		_ = i
		index++
	}

	return previews, nil
}

// BatchRename 执行批量重命名
func BatchRename(dir string, rule RenameRule) ([]RenamePreview, error) {
	previews, err := BatchRenamePreview(dir, rule)
	if err != nil {
		return nil, err
	}

	// 先重命名为临时名，避免文件名冲突
	var tempNames []string
	for i, p := range previews {
		tempName := fmt.Sprintf(".temp_%d_%s", i, p.NewName)
		tempPath := filepath.Join(dir, tempName)
		if err := os.Rename(p.OriginalPath, tempPath); err != nil {
			// 回滚
			rollbackRenames(dir, previews[:i], tempNames)
			return nil, fmt.Errorf("重命名失败: %w", err)
		}
		tempNames = append(tempNames, tempName)
	}

	// 再从临时名改为目标名
	for i, p := range previews {
		tempPath := filepath.Join(dir, tempNames[i])
		if err := os.Rename(tempPath, p.NewPath); err != nil {
			return nil, fmt.Errorf("最终重命名失败: %w", err)
		}
		previews[i].NewPath = p.NewPath
	}

	return previews, nil
}

func applyPattern(originalName string, rule RenameRule, index int) string {
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)

	// 替换模式
	result := rule.Pattern
	result = strings.ReplaceAll(result, "{index}", padNumber(index, rule.Padding))
	result = strings.ReplaceAll(result, "{name}", baseName)
	result = strings.ReplaceAll(result, "{ext}", ext)

	// 替换文本
	if rule.ReplaceFrom != "" {
		result = strings.ReplaceAll(result, rule.ReplaceFrom, rule.ReplaceTo)
	}

	return result
}

func padNumber(n int, padding int) string {
	return fmt.Sprintf("%0*d", padding, n)
}

func rollbackRenames(dir string, previews []RenamePreview, tempNames []string) {
	for i, p := range previews {
		tempPath := filepath.Join(dir, tempNames[i])
		os.Rename(tempPath, p.OriginalPath)
	}
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
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
