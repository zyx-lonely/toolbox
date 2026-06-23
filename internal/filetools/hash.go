package filetools

import (
	"os"
	"path/filepath"

	"pc-toolbox/internal/common"
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
	hashVal, size, err := common.ComputeFileHash(path, algorithm)
	if err != nil {
		return nil, err
	}
	return &HashResult{
		Path:      path,
		Algorithm: algorithm,
		Hash:      hashVal,
		FileSize:  size,
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
