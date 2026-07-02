package daily

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// CalcResult 计算结果
type CalcResult struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
	Error      string  `json:"error,omitempty"`
}

// Calculate 计算表达式
func Calculate(expr string) (*CalcResult, error) {
	if expr == "" {
		return nil, fmt.Errorf("表达式不能为空")
	}

	result, err := evalExpression(expr)
	if err != nil {
		return &CalcResult{
			Expression: expr,
			Error:      err.Error(),
		}, nil
	}

	return &CalcResult{
		Expression: expr,
		Result:     result,
	}, nil
}

func evalExpression(expr string) (float64, error) {
	// 预处理：替换中文符号和常用函数
	expr = strings.ReplaceAll(expr, "×", "*")
	expr = strings.ReplaceAll(expr, "÷", "/")
	expr = strings.ReplaceAll(expr, "（", "(")
	expr = strings.ReplaceAll(expr, "）", ")")
	expr = strings.ReplaceAll(expr, "^", "**")
	expr = strings.ReplaceAll(expr, ",", "")
	expr = strings.TrimSpace(expr)

	// 解析并计算
	tokens := tokenize(expr)
	result, _, err := parseExpression(tokens, 0)
	if err != nil {
		return 0, err
	}

	return result, nil
}

type tokenType int

const (
	tokenNumber tokenType = iota
	tokenOperator
	tokenFunction
	tokenLeftParen
	tokenRightParen
	tokenComma
)

type token struct {
	typ    tokenType
	value  string
	numVal float64
}

func tokenize(expr string) []token {
	var tokens []token
	i := 0
	for i < len(expr) {
		ch := expr[i]

		// 空格跳过
		if ch == ' ' {
			i++
			continue
		}

		// 数字
		if ch >= '0' && ch <= '9' || ch == '.' {
			j := i
			for j < len(expr) && ((expr[j] >= '0' && expr[j] <= '9') || expr[j] == '.') {
				j++
			}
			num, _ := strconv.ParseFloat(expr[i:j], 64)
			tokens = append(tokens, token{typ: tokenNumber, value: expr[i:j], numVal: num})
			i = j
			continue
		}

		// 字母（函数名）
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			j := i
			for j < len(expr) && ((expr[j] >= 'a' && expr[j] <= 'z') || (expr[j] >= 'A' && expr[j] <= 'Z')) {
				j++
			}
			name := strings.ToLower(expr[i:j])
			// 检查是否是函数
			if _, ok := mathFunctions[name]; ok {
				tokens = append(tokens, token{typ: tokenFunction, value: name})
			} else {
				// 跳过未知标识符
			}
			i = j
			continue
		}

		// 运算符
		if strings.Contains("+-*/%", string(ch)) {
			tokens = append(tokens, token{typ: tokenOperator, value: string(ch)})
			i++
			continue
		}

		// 括号
		if ch == '(' {
			tokens = append(tokens, token{typ: tokenLeftParen, value: "("})
			i++
			continue
		}
		if ch == ')' {
			tokens = append(tokens, token{typ: tokenRightParen, value: ")"})
			i++
			continue
		}

		// 逗号
		if ch == ',' {
			tokens = append(tokens, token{typ: tokenComma, value: ","})
			i++
			continue
		}

		i++
	}
	return tokens
}

var mathFunctions = map[string]func(args ...float64) float64{
	"sin":   func(args ...float64) float64 { return math.Sin(args[0] * math.Pi / 180) },
	"cos":   func(args ...float64) float64 { return math.Cos(args[0] * math.Pi / 180) },
	"tan":   func(args ...float64) float64 { return math.Tan(args[0] * math.Pi / 180) },
	"asin":  func(args ...float64) float64 { return math.Asin(args[0]) * 180 / math.Pi },
	"acos":  func(args ...float64) float64 { return math.Acos(args[0]) * 180 / math.Pi },
	"atan":  func(args ...float64) float64 { return math.Atan(args[0]) * 180 / math.Pi },
	"sqrt":  func(args ...float64) float64 { return math.Sqrt(args[0]) },
	"cbrt":  func(args ...float64) float64 { return math.Cbrt(args[0]) },
	"abs":   func(args ...float64) float64 { return math.Abs(args[0]) },
	"log":   func(args ...float64) float64 { return math.Log10(args[0]) },
	"ln":    func(args ...float64) float64 { return math.Log(args[0]) },
	"log2":  func(args ...float64) float64 { return math.Log2(args[0]) },
	"exp":   func(args ...float64) float64 { return math.Exp(args[0]) },
	"pow":   func(args ...float64) float64 { return math.Pow(args[0], args[1]) },
	"floor": func(args ...float64) float64 { return math.Floor(args[0]) },
	"ceil":  func(args ...float64) float64 { return math.Ceil(args[0]) },
	"round": func(args ...float64) float64 { return math.Round(args[0]) },
	"pi":    func(args ...float64) float64 { return math.Pi },
	"e":     func(args ...float64) float64 { return math.E },
}

func parseExpression(tokens []token, pos int) (float64, int, error) {
	if pos >= len(tokens) {
		return 0, pos, nil
	}

	result, newPos, err := parseAddSub(tokens, pos)
	if err != nil {
		return 0, newPos, err
	}
	return result, newPos, nil
}

func parseAddSub(tokens []token, pos int) (float64, int, error) {
	left, pos, err := parseMulDiv(tokens, pos)
	if err != nil {
		return 0, pos, err
	}

	for pos < len(tokens) && (tokens[pos].value == "+" || tokens[pos].value == "-") {
		op := tokens[pos].value
		pos++
		right, newPos, err := parseMulDiv(tokens, pos)
		if err != nil {
			return 0, newPos, err
		}
		if op == "+" {
			left += right
		} else {
			left -= right
		}
		pos = newPos
	}

	return left, pos, nil
}

func parseMulDiv(tokens []token, pos int) (float64, int, error) {
	left, pos, err := parseUnary(tokens, pos)
	if err != nil {
		return 0, pos, err
	}

	for pos < len(tokens) && (tokens[pos].value == "*" || tokens[pos].value == "/" || tokens[pos].value == "%") {
		op := tokens[pos].value
		pos++
		right, newPos, err := parseUnary(tokens, pos)
		if err != nil {
			return 0, newPos, err
		}
		switch op {
		case "*":
			left *= right
		case "/":
			if right == 0 {
				return 0, pos, fmt.Errorf("除数不能为零")
			}
			left /= right
		case "%":
			if right == 0 {
				return 0, pos, fmt.Errorf("除数不能为零")
			}
			left = math.Mod(left, right)
		}
		pos = newPos
	}

	return left, pos, nil
}

func parseUnary(tokens []token, pos int) (float64, int, error) {
	if pos < len(tokens) && tokens[pos].value == "-" {
		pos++
		val, newPos, err := parsePrimary(tokens, pos)
		if err != nil {
			return 0, newPos, err
		}
		return -val, newPos, nil
	}
	if pos < len(tokens) && tokens[pos].value == "+" {
		pos++
	}
	return parsePrimary(tokens, pos)
}

func parsePrimary(tokens []token, pos int) (float64, int, error) {
	if pos >= len(tokens) {
		return 0, pos, fmt.Errorf("表达式不完整")
	}

	tok := tokens[pos]

	// 数字
	if tok.typ == tokenNumber {
		return tok.numVal, pos + 1, nil
	}

	// 括号
	if tok.typ == tokenLeftParen {
		pos++
		val, newPos, err := parseExpression(tokens, pos)
		if err != nil {
			return 0, newPos, err
		}
		if newPos < len(tokens) && tokens[newPos].typ == tokenRightParen {
			newPos++
		}
		return val, newPos, nil
	}

	// 函数
	if tok.typ == tokenFunction {
		pos++
		// 跳过左括号
		if pos < len(tokens) && tokens[pos].typ == tokenLeftParen {
			pos++
		}

		// 收集参数
		var args []float64
		for pos < len(tokens) && tokens[pos].typ != tokenRightParen {
			val, newPos, err := parseExpression(tokens, pos)
			if err != nil {
				return 0, newPos, err
			}
			args = append(args, val)
			pos = newPos
			if pos < len(tokens) && tokens[pos].typ == tokenComma {
				pos++
			}
		}

		// 跳过右括号
		if pos < len(tokens) && tokens[pos].typ == tokenRightParen {
			pos++
		}

		if fn, ok := mathFunctions[tok.value]; ok {
			return fn(args...), pos, nil
		}
		return 0, pos, fmt.Errorf("未知函数: %s", tok.value)
	}

	return 0, pos + 1, nil
}
