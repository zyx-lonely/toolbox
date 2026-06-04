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

// HashResult 哈希计算结果
type HashResult struct {
	Path      string `json:"path"`
	Algorithm string `json:"algorithm"`
	Hash      string `json:"hash"`
	FileSize  int64  `json:"fileSize"`
	Error     string `json:"error,omitempty"`
}

// ComputeHash 计算单个文件哈希
func ComputeHash(path string, algorithm string) (*HashResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
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
		return nil, fmt.Errorf("不支持的算法: %s", algorithm)
	}

	if _, err := io.Copy(h, f); err != nil {
		return nil, fmt.Errorf("计算哈希失败: %w", err)
	}

	info, _ := os.Stat(path)
	return &HashResult{
		Path:      path,
		Algorithm: algorithm,
		Hash:      hex.EncodeToString(h.Sum(nil)),
		FileSize:  info.Size(),
	}, nil
}

// ComputeHashes 计算多个文件的哈希
func ComputeHashes(paths []string, algorithm string) []HashResult {
	var results []HashResult
	for _, p := range paths {
		result, err := ComputeHash(p, algorithm)
		if err != nil {
			results = append(results, HashResult{
				Path:      p,
				Algorithm: algorithm,
				Error:     err.Error(),
			})
		} else {
			results = append(results, *result)
		}
	}
	return results
}

// WalkFiles 遍历目录获取所有文件路径
func WalkFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
