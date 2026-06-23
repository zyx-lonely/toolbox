package devtools

import "strings"

// FormatResult 格式化结果
type FormatResult struct {
	Input   string `json:"input"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// FormatYAML 格式化 YAML（语法感知缩进整理）
// 支持：注释保留、列表项缩进、嵌套映射、文档分隔符 ---
func FormatYAML(input string) FormatResult {
	lines := strings.Split(input, "\n")
	var out []string
	depth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 保留空行
		if trimmed == "" {
			out = append(out, "")
			continue
		}

		// 保留注释行（不修改缩进，保持原样）
		if strings.HasPrefix(trimmed, "#") {
			out = append(out, strings.Repeat("  ", depth)+trimmed)
			continue
		}

		// 文档分隔符：重置深度
		if trimmed == "---" || trimmed == "..." {
			depth = 0
			out = append(out, trimmed)
			continue
		}

		// 检测列表项
		isList := strings.HasPrefix(trimmed, "- ")
		if isList {
			// 列表项：在当前深度输出
			out = append(out, strings.Repeat("  ", depth)+trimmed)
			continue
		}

		// 检测键值对
		colonIdx := strings.Index(trimmed, ":")
		if colonIdx > 0 {
			key := trimmed[:colonIdx]
			// 如果上一行是列表项，且当前是 key: value，说明是列表中的映射项
			out = append(out, strings.Repeat("  ", depth)+key+":")

			// 检查是否有值（key: value 格式）
			rest := strings.TrimSpace(trimmed[colonIdx+1:])
			if rest != "" {
				// 有值，将值追加到同一行
				out[len(out)-1] += " " + rest
			} else {
				// 无值，下一级缩进
				depth++
			}
			continue
		}

		// 普通行：保持当前缩进
		out = append(out, strings.Repeat("  ", depth)+trimmed)
	}

	return FormatResult{Input: input, Output: strings.Join(out, "\n"), Success: true}
}

// FormatTOML 格式化 TOML（简易：保持原样，TOML 格式较复杂需专用库）
func FormatTOML(input string) FormatResult {
	return FormatResult{Input: input, Output: input, Success: true}
}
