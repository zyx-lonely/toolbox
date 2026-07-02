package filetools

import (
	"os/exec"
	"syscall"
)

type RecycleBinInfo struct {
	ItemCount int    `json:"itemCount"`
	Size      string `json:"size"`
	Path      string `json:"path"`
}

// GetRecycleBinInfo 获取回收站信息
func GetRecycleBinInfo() RecycleBinInfo {
	// 使用 PowerShell 获取回收站内容
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		`(New-Object -ComObject Shell.Application).NameSpace(0x0a).Items() | Measure-Object | Select-Object Count | ConvertTo-Json`)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, _ := cmd.Output()

	count := 0
	// 简单解析 Count 字段
	if len(out) > 0 {
		for _, line := range splitLines(string(out)) {
			if containsString(line, "Count") {
				// 提取数字
				for i := 0; i < len(line); i++ {
					if line[i] >= '0' && line[i] <= '9' {
						start := i
						for i < len(line) && line[i] >= '0' && line[i] <= '9' {
							i++
						}
						count = parseInt(line[start:i])
						break
					}
				}
			}
		}
	}

	return RecycleBinInfo{
		ItemCount: count,
		Size:      "未知",
		Path:      "$Recycle.Bin",
	}
}

// EmptyRecycleBin 清空回收站（所有驱动器）
func EmptyRecycleBin() error {
	// 使用 Shell32 API 清空回收站（更安全）
	c := exec.Command("powershell", "-NoProfile", "-Command",
		`(New-Object -ComObject Shell.Application).NameSpace(0x0a).Items() | ForEach-Object { $_.InvokeVerb("delete") }`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c.Run()
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func containsString(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsAt(s, sub))
}

func containsAt(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
