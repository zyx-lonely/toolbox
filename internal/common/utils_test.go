package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0 B"},
		{1023, "1023 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
		{1099511627776, "1.00 TB"},
	}

	for _, tt := range tests {
		result := FormatBytes(tt.input)
		if result != tt.expected {
			t.Errorf("FormatBytes(%d) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestFormatBytes_Int64(t *testing.T) {
	result := FormatBytes(int64(2048))
	if result != "2.00 KB" {
		t.Errorf("FormatBytes(int64) = %s, want 2.00 KB", result)
	}
}

func TestIsDir(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建临时文件和目录
	dir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(dir, 0755)

	file := filepath.Join(tmpDir, "file.txt")
	os.WriteFile(file, []byte("test"), 0644)

	if !IsDir(dir) {
		t.Error("IsDir 对目录应返回 true")
	}

	if IsDir(file) {
		t.Error("IsDir 对文件应返回 false")
	}

	if IsDir("/nonexistent/path") {
		t.Error("IsDir 对不存在的路径应返回 false")
	}
}

func TestSafeFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal.txt", "normal.txt"},
		{"file:name.txt", "file_name.txt"},
		{"a/b\\c*d?e\"f<g>h|i", "a_b_c_d_e_f_g_h_i"},
		{"no_invalid", "no_invalid"},
		{"***", "___"},
	}

	for _, tt := range tests {
		result := SafeFileName(tt.input)
		if result != tt.expected {
			t.Errorf("SafeFileName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "existing.txt")
	os.WriteFile(file, []byte("test"), 0644)

	if !FileExists(file) {
		t.Error("FileExists 对存在的文件应返回 true")
	}

	if FileExists("/nonexistent/file.txt") {
		t.Error("FileExists 对不存在的文件应返回 false")
	}
}

func TestGetTotalSize(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建文件结构
	os.WriteFile(filepath.Join(tmpDir, "a.txt"), []byte("12345"), 0644)       // 5 bytes
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("1234567890"), 0644) // 10 bytes
	subDir := filepath.Join(tmpDir, "sub")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "c.txt"), []byte("abc"), 0644) // 3 bytes

	total, err := GetTotalSize(tmpDir)
	if err != nil {
		t.Fatalf("GetTotalSize 失败: %v", err)
	}

	expected := uint64(5 + 10 + 3)
	if total != expected {
		t.Errorf("GetTotalSize = %d, want %d", total, expected)
	}
}

func TestOpenURL_EmptyURL(t *testing.T) {
	err := OpenURL("")
	if err == nil {
		t.Error("空 URL 应返回错误")
	}
}
