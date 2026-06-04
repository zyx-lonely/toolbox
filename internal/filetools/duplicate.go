package filetools

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// FileMatch 文件匹配信息（用于重复文件检测）
type FileMatch struct {
	Size      int64  `json:"size"`
	ModTime   string `json:"modTime"`
	Path      string `json:"path"`
	Hash      string `json:"hash"`
	MatchType string `json:"matchType"` // "exact" - hash匹配, "fuzzy" - 大小+时间匹配
}

// DuplicateGroup 重复文件组
type DuplicateGroup struct {
	Hash      string       `json:"hash"`
	FileCount int          `json:"fileCount"`
	TotalSize uint64       `json:"totalSize"`
	Files     []FileMatch  `json:"files"`
}

// FindDuplicates 查找重复文件
// mode: "quick" (大小+修改时间) 或 "exact" (MD5 哈希)
func FindDuplicates(rootPath string, mode string) ([]DuplicateGroup, error) {
	// 第一阶段：按大小分组
	sizeMap := make(map[int64][]string)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Size() > 0 {
			sizeMap[info.Size()] = append(sizeMap[info.Size()], path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 第二阶段：过滤出有重复大小的组，计算哈希
	var groups []DuplicateGroup
	for size, paths := range sizeMap {
		if len(paths) < 2 {
			continue
		}

		if mode == "quick" {
			groups = append(groups, DuplicateGroup{
				FileCount: len(paths),
				TotalSize: uint64(size) * uint64(len(paths)-1),
				Files:     pathsToFileMatches(paths, size),
			})
			continue
		}

		// 精确模式：按哈希分组
		hashMap := make(map[string][]FileMatch)
		for _, path := range paths {
			h, err := computeFileHash(path, "md5")
			if err != nil {
				continue
			}

			info, _ := os.Stat(path)
			hashMap[h] = append(hashMap[h], FileMatch{
				Size:      size,
				ModTime:   info.ModTime().Format("2006-01-02 15:04:05"),
				Path:      path,
				Hash:      h,
				MatchType: "exact",
			})
		}

		for h, files := range hashMap {
			if len(files) >= 2 {
				groups = append(groups, DuplicateGroup{
					Hash:      h,
					FileCount: len(files),
					TotalSize: uint64(size) * uint64(len(files)-1),
					Files:     files,
				})
			}
		}
	}

	return groups, nil
}

// ComputeFileHash 计算文件哈希
func ComputeFileHash(path string, algorithm string) (string, error) {
	return computeFileHash(path, algorithm)
}

func computeFileHash(path string, algorithm string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return "", fmt.Errorf("不支持的哈希算法: %s", algorithm)
	}

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func pathsToFileMatches(paths []string, size int64) []FileMatch {
	var matches []FileMatch
	for _, p := range paths {
		info, _ := os.Stat(p)
		matches = append(matches, FileMatch{
			Size:      size,
			ModTime:   info.ModTime().Format("2006-01-02 15:04:05"),
			Path:      p,
			MatchType: "fuzzy",
		})
	}
	return matches
}
