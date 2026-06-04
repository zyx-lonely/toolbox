package devtools

import "strings"

// FormatResult 格式化结果
type FormatResult struct {
	Input   string `json:"input"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// FormatYAML 格式化 YAML（简易缩进整理）
func FormatYAML(input string) FormatResult {
	lines := strings.Split(input, "\n")
	var out []string
	depth := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			out = append(out, "")
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			// 列表项
		}
		if strings.HasPrefix(trimmed, "---") {
			depth = 0
		}
		out = append(out, strings.Repeat("  ", depth)+trimmed)
		if strings.HasSuffix(trimmed, ":") || strings.HasSuffix(trimmed, ": ") {
			depth++
		}
		if strings.HasPrefix(trimmed, "- ") && !strings.Contains(trimmed, ":") {
		} else if strings.HasPrefix(trimmed, "- ") {
		}
	}
	_ = depth
	return FormatResult{Input: input, Output: strings.Join(out, "\n"), Success: true}
}

// FormatTOML 格式化 TOML（简易）
func FormatTOML(input string) FormatResult {
	return FormatResult{Input: input, Output: input, Success: true}
}
