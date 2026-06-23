package devtools

import (
	"strings"
)

type CodeBeautifyResult struct {
	Input   string `json:"input"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// BeautifyHTML 美化 HTML（缩进对齐标签）
func BeautifyHTML(input string) CodeBeautifyResult {
	defer func() {
		if r := recover(); r != nil {
			// 防止解析异常导致崩溃
		}
	}()

	// 预处理：去除多余空白
	input = strings.TrimSpace(input)
	if input == "" {
		return CodeBeautifyResult{Input: input, Output: "", Success: true}
	}

	var output strings.Builder
	depth := 0
	inPre := false // 是否在 <pre> 或 <textarea> 内
	inComment := false
	inScript := false
	inStyle := false
	i := 0
	runes := []rune(input)

	for i < len(runes) {
		ch := runes[i]

		// 处理注释
		if !inComment && !inPre && !inScript && !inStyle && i+4 <= len(runes) && string(runes[i:i+4]) == "<!--" {
			inComment = true
			indent := strings.Repeat("  ", depth)
			output.WriteString(indent + "<!--")
			i += 4
			continue
		}
		if inComment && i+3 <= len(runes) && string(runes[i:i+3]) == "-->" {
			inComment = false
			output.WriteString("-->\n")
			i += 3
			continue
		}
		if inComment {
			output.WriteRune(ch)
			i++
			continue
		}

		// 检测进入 pre/script/style 区域
		if !inPre && !inScript && !inStyle && ch == '<' {
			lower := strings.ToLower(string(runes[i:]))
			if strings.HasPrefix(lower, "<pre") {
				inPre = true
			}
			if strings.HasPrefix(lower, "<script") {
				inScript = true
			}
			if strings.HasPrefix(lower, "<style") {
				inStyle = true
			}
		}

		// 检测离开 pre/script/style 区域
		if inPre && ch == '<' {
			lower := strings.ToLower(string(runes[i:]))
			if strings.HasPrefix(lower, "</pre>") {
				inPre = false
			}
		}
		if inScript && ch == '<' {
			lower := strings.ToLower(string(runes[i:]))
			if strings.HasPrefix(lower, "</script>") {
				inScript = false
			}
		}
		if inStyle && ch == '<' {
			lower := strings.ToLower(string(runes[i:]))
			if strings.HasPrefix(lower, "</style>") {
				inStyle = false
			}
		}

		// 在 pre/script/style 内保持原样
		if inPre || inScript || inStyle {
			output.WriteRune(ch)
			i++
			continue
		}

		// 处理标签
		if ch == '<' {
			// 判断是否为闭合标签
			if i+1 < len(runes) && runes[i+1] == '/' {
				depth--
				if depth < 0 {
					depth = 0
				}
				indent := strings.Repeat("  ", depth)
				// 输出到下一个 '>'
				endIdx := strings.Index(string(runes[i:]), ">")
				if endIdx == -1 {
					output.WriteString(indent + string(runes[i:]))
					break
				}
				tag := string(runes[i : i+endIdx+1])
				// 自闭合标签不缩进（如 <br/> <img/>）
				if strings.HasSuffix(strings.TrimRight(tag, " "), "/>") {
					depth++ // 抵消之前的 depth--
				}
				output.WriteString(indent + tag + "\n")
				i += endIdx + 1
				continue
			}

			// 处理 <?xml ... ?> 等
			if i+1 < len(runes) && runes[i+1] == '?' {
				endIdx := strings.Index(string(runes[i:]), "?>")
				if endIdx == -1 {
					output.WriteRune(ch)
					i++
					continue
				}
				indent := strings.Repeat("  ", depth)
				output.WriteString(indent + string(runes[i:i+endIdx+2]) + "\n")
				i += endIdx + 2
				continue
			}

			// 检测自闭合标签 <br/> <img/> <input/> <hr/> <meta/> <link/>
			tagContent := string(runes[i:])
			selfClose := isSelfClosingTag(tagContent)

			indent := strings.Repeat("  ", depth)
			// 找到标签结束 '>'
			endIdx := findTagEnd(runes[i:])
			if endIdx == -1 {
				output.WriteRune(ch)
				i++
				continue
			}

			tag := string(runes[i : i+endIdx+1])
			output.WriteString(indent + tag)

			if !selfClose {
				// 跳过纯内容节点中的空白
			}
			output.WriteString("\n")

			if !selfClose {
				depth++
			}
			i += endIdx + 1
			continue
		}

		// 文本内容：跳过纯空白
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			// 只在非空行首空白有意义
			allSpace := true
			for j := i; j < len(runes); j++ {
				next := runes[j]
				if next == '<' {
					break
				}
				if next != ' ' && next != '\t' && next != '\n' && next != '\r' {
					allSpace = false
					break
				}
			}
			if allSpace {
				i = skipSpaces(runes, i)
				continue
			}
			// 保留单个空格
			output.WriteRune(' ')
			i = skipSpaces(runes, i)
			continue
		}

		output.WriteRune(ch)
		i++
	}

	result := strings.TrimSpace(output.String())
	return CodeBeautifyResult{Input: input, Output: result, Success: true}
}

// BeautifyCSS 美化 CSS（格式化缩进）
func BeautifyCSS(input string) CodeBeautifyResult {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	input = strings.TrimSpace(input)
	if input == "" {
		return CodeBeautifyResult{Input: input, Output: "", Success: true}
	}

	var output strings.Builder
	depth := 0
	inComment := false
	inString := false
	stringCh := rune(0)
	i := 0
	runes := []rune(input)

	for i < len(runes) {
		ch := runes[i]

		// 注释处理
		if !inComment && !inString && i+2 <= len(runes) && runes[i] == '/' && runes[i+1] == '*' {
			inComment = true
			indent := strings.Repeat("  ", depth)
			output.WriteString(indent + "/*")
			i += 2
			continue
		}
		if inComment && i+2 <= len(runes) && runes[i] == '*' && runes[i+1] == '/' {
			inComment = false
			output.WriteString("*/\n")
			i += 2
			continue
		}
		if inComment {
			if ch == '\n' {
				output.WriteString("\n")
			} else {
				output.WriteRune(ch)
			}
			i++
			continue
		}

		// 字符串处理
		if !inString && (ch == '"' || ch == '\'') {
			inString = true
			stringCh = ch
			output.WriteRune(ch)
			i++
			continue
		}
		if inString {
			if ch == '\\' && i+1 < len(runes) {
				output.WriteRune(ch)
				output.WriteRune(runes[i+1])
				i += 2
				continue
			}
			if ch == stringCh {
				inString = false
			}
			output.WriteRune(ch)
			i++
			continue
		}

		// 格式化
		switch ch {
		case '{':
			_ = strings.Repeat("  ", depth) // 预留缩进
			output.WriteString(" {\n")
			depth++
			i++
			// 跳过 { 后的空白
			for i < len(runes) && (runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\n') {
				i++
			}
			continue
		case '}':
			depth--
			if depth < 0 {
				depth = 0
			}
			indent := strings.Repeat("  ", depth)
			output.WriteString(indent + "}\n")
			i++
			continue
		case ';':
			indent := strings.Repeat("  ", depth)
			output.WriteString(indent + ";\n")
			i++
			// 跳过 ; 后的空白
			for i < len(runes) && (runes[i] == ' ' || runes[i] == '\t') {
				i++
			}
			continue
		default:
			// 在属性值上下文中输出
			if depth > 0 {
				indent := strings.Repeat("  ", depth)
				// 检查是否需要换行
				if output.Len() > 0 {
					lastBytes := output.String()
					lastCh := rune(lastBytes[len(lastBytes)-1])
					if lastCh == '\n' {
						output.WriteString(indent)
					} else if lastCh == ';' || lastCh == '{' {
						// 已有缩进
					}
				}
			}
			output.WriteRune(ch)
			i++
		}
	}

	result := strings.TrimSpace(output.String())
	if result == "" {
		result = input
	}
	return CodeBeautifyResult{Input: input, Output: result, Success: true}
}

// BeautifyJS 美化 JavaScript（简单缩进美化）
func BeautifyJS(input string) CodeBeautifyResult {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	if strings.TrimSpace(input) == "" {
		return CodeBeautifyResult{Input: input, Output: input, Success: true}
	}

	var output strings.Builder
	depth := 0
	inString := false
	stringCh := rune(0)
	inComment := false
	inLineComment := false
	i := 0
	runes := []rune(input)

	for i < len(runes) {
		ch := runes[i]

		// 行注释
		if !inLineComment && !inComment && !inString && i+2 <= len(runes) && runes[i] == '/' && runes[i+1] == '/' {
			inLineComment = true
			output.WriteString("//")
			i += 2
			continue
		}
		if inLineComment && ch == '\n' {
			inLineComment = false
			output.WriteString("\n")
			i++
			continue
		}
		if inLineComment {
			output.WriteRune(ch)
			i++
			continue
		}

		// 块注释
		if !inComment && !inString && i+2 <= len(runes) && runes[i] == '/' && runes[i+1] == '*' {
			inComment = true
			output.WriteString("/*")
			i += 2
			continue
		}
		if inComment && i+2 <= len(runes) && runes[i] == '*' && runes[i+1] == '/' {
			inComment = false
			output.WriteString("*/\n")
			i += 2
			continue
		}
		if inComment {
			output.WriteRune(ch)
			i++
			continue
		}

		// 字符串
		if !inString && (ch == '"' || ch == '\'' || ch == '`') {
			inString = true
			stringCh = ch
			output.WriteRune(ch)
			i++
			continue
		}
		if inString {
			if ch == '\\' && i+1 < len(runes) {
				output.WriteRune(ch)
				output.WriteRune(runes[i+1])
				i += 2
				continue
			}
			if ch == stringCh {
				inString = false
			}
			output.WriteRune(ch)
			i++
			continue
		}

		// 缩进控制
		switch ch {
		case '{':
			indent := strings.Repeat("  ", depth)
			output.WriteString(indent + "{\n")
			depth++
			i++
			// 跳过空白
			for i < len(runes) && (runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\n') {
				i++
			}
			continue
		case '}':
			depth--
			if depth < 0 {
				depth = 0
			}
			indent := strings.Repeat("  ", depth)
			output.WriteString(indent + "}\n")
			i++
			for j := i; j < len(runes); j++ {
				if runes[j] == ' ' || runes[j] == '\t' || runes[j] == '\n' {
					i++
				} else {
					break
				}
			}
			continue
		case ';':
			output.WriteString(";\n")
			i++
			for j := i; j < len(runes); j++ {
				if runes[j] == ' ' || runes[j] == '\t' || runes[j] == '\n' {
					i++
				} else {
					break
				}
			}
			continue
		default:
			// 换行后加缩进
			if ch == '\n' {
				output.WriteRune(ch)
				i++
				// 下一行添加缩进
				for i < len(runes) && (runes[i] == ' ' || runes[i] == '\t') {
					i++
				}
				if i < len(runes) && runes[i] != '\n' && depth > 0 {
					output.WriteString(strings.Repeat("  ", depth))
				}
				continue
			}
			// 跳过行首空白（缩进已加）
			if (ch == ' ' || ch == '\t') && output.Len() > 0 {
				lastOutput := output.String()
				if lastOutput[len(lastOutput)-1] == '\n' {
					i++
					continue
				}
			}
			output.WriteRune(ch)
			i++
		}
	}

	result := strings.TrimSpace(output.String())
	if result == "" {
		result = input
	}
	return CodeBeautifyResult{Input: input, Output: result, Success: true, Error: ""}
}

// BeautifySQL 美化 SQL
func BeautifySQL(input string) CodeBeautifyResult {
	out := ""
	depth := 0
	inString := false
	for _, c := range input {
		if c == '\'' {
			inString = !inString
		}
		if !inString {
			if c == '(' {
				depth++
				out += string(c) + "\n" + strings.Repeat("  ", depth)
				continue
			}
			if c == ')' {
				depth--
				if depth < 0 {
					depth = 0
				}
				out += "\n" + strings.Repeat("  ", depth) + string(c)
				continue
			}
			if c == ',' && !inString {
				out += ",\n" + strings.Repeat("  ", depth)
				continue
			}
		}
		out += string(c)
	}
	return CodeBeautifyResult{Input: input, Output: out, Success: true}
}

// 辅助函数
func isSelfClosingTag(s string) bool {
	selfClosing := []string{"br", "hr", "img", "input", "meta", "link", "area", "base", "col", "embed", "source", "track", "wbr"}
	lower := strings.ToLower(s)
	for _, tag := range selfClosing {
		if strings.HasPrefix(lower, "<"+tag+" ") || strings.HasPrefix(lower, "<"+tag+"/") || strings.HasPrefix(lower, "<"+tag+">") {
			return true
		}
	}
	return strings.Contains(lower, "/>")
}

func findTagEnd(runes []rune) int {
	for i, ch := range runes {
		if ch == '>' {
			return i
		}
	}
	return -1
}

func skipSpaces(runes []rune, i int) int {
	for i < len(runes) && (runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\n' || runes[i] == '\r') {
		i++
	}
	return i
}

