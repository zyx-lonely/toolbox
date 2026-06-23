package common

import (
	"os"
	"strings"
)

// ExpandEnv 展开 Windows 环境变量（支持 %VAR% 格式）
func ExpandEnv(path string) string {
	result := path
	for {
		start := strings.Index(result, "%")
		if start == -1 {
			break
		}
		end := strings.Index(result[start+1:], "%")
		if end == -1 {
			break
		}
		key := result[start+1 : start+1+end]
		val := os.Getenv(key)
		if val != "" {
			result = result[:start] + val + result[start+1+end+1:]
		} else if key == "WINDIR" {
			result = result[:start] + os.Getenv("SystemRoot") + result[start+1+end+1:]
		} else {
			// 环境变量不存在，移除占位符避免路径中出现空段
			result = result[:start] + result[start+1+end+1:]
		}
	}
	return result
}

// ExpandPath 展开路径中的环境变量并规范化
func ExpandPath(path string) string {
	p := ExpandEnv(path)
	p = strings.ReplaceAll(p, "\\", string(os.PathSeparator))
	return p
}
