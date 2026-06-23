package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComputeFileHash_MD5(t *testing.T) {
	// 创建临时文件
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	content := "hello world"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	hash, size, err := ComputeFileHash(file, "md5")
	if err != nil {
		t.Fatalf("计算 MD5 失败: %v", err)
	}

	// "hello world" 的 MD5: 5eb63bbbe01eeed093cb22bb8f5acdc3
	expected := "5eb63bbbe01eeed093cb22bb8f5acdc3"
	if hash != expected {
		t.Errorf("MD5 不匹配: got %s, want %s", hash, expected)
	}

	if size != int64(len(content)) {
		t.Errorf("文件大小不匹配: got %d, want %d", size, len(content))
	}
}

func TestComputeFileHash_SHA256(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	content := "hello world"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}

	hash, _, err := ComputeFileHash(file, "sha256")
	if err != nil {
		t.Fatalf("计算 SHA256 失败: %v", err)
	}

	// "hello world" 的 SHA256
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	if hash != expected {
		t.Errorf("SHA256 不匹配: got %s, want %s", hash, expected)
	}
}

func TestComputeFileHash_UnsupportedAlgorithm(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(file, []byte("test"), 0644)

	_, _, err := ComputeFileHash(file, "invalid")
	if err == nil {
		t.Error("预期不支持的算法应返回错误")
	}
}

func TestComputeFileHash_FileNotFound(t *testing.T) {
	_, _, err := ComputeFileHash("/nonexistent/file.txt", "md5")
	if err == nil {
		t.Error("预期文件不存在应返回错误")
	}
}
