package common

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

// GbkToUtf8 将 GBK 编码的字符串转换为 UTF-8
func GbkToUtf8(s string) string {
	if s == "" {
		return ""
	}
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	data, err := io.ReadAll(reader)
	if err != nil {
		return s // 如果转换失败，返回原字符串
	}
	return string(data)
}

// Utf8ToGbk 将 UTF-8 编码的字符串转换为 GBK
func Utf8ToGbk(s string) string {
	if s == "" {
		return ""
	}
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	data, err := io.ReadAll(reader)
	if err != nil {
		return s
	}
	return string(data)
}
