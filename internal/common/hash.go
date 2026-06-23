package common

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
)

// ComputeFileHash 计算指定文件的哈希值
// algorithm: "md5", "sha1", "sha256", "sha512"
func ComputeFileHash(path string, algorithm string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, fmt.Errorf("打开文件失败: %w", err)
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
		return "", 0, fmt.Errorf("不支持的哈希算法: %s", algorithm)
	}

	if _, err := io.Copy(h, f); err != nil {
		return "", 0, fmt.Errorf("计算哈希失败: %w", err)
	}

	info, _ := os.Stat(path)
	return hex.EncodeToString(h.Sum(nil)), info.Size(), nil
}
