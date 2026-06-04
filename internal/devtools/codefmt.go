package devtools

type CodeBeautifyResult struct {
	Input   string `json:"input"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func BeautifyHTML(input string) CodeBeautifyResult {
	return CodeBeautifyResult{Input: input, Output: input, Success: true}
}

func BeautifyCSS(input string) CodeBeautifyResult {
	return CodeBeautifyResult{Input: input, Output: input, Success: true}
}

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
				out += string(c) + "\n" + stringsRepeat("  ", depth)
				continue
			}
			if c == ')' {
				depth--
				if depth < 0 {
					depth = 0
				}
				out += "\n" + stringsRepeat("  ", depth) + string(c)
				continue
			}
			if c == ',' && !inString {
				out += ",\n" + stringsRepeat("  ", depth)
				continue
			}
		}
		out += string(c)
	}
	return CodeBeautifyResult{Input: input, Output: out, Success: true}
}

func stringsRepeat(s string, count int) string {
	r := ""
	for i := 0; i < count; i++ {
		r += s
	}
	return r
}
