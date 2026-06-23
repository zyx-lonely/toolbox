package devtools

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// JSONResult JSON 处理结果
type JSONResult struct {
	Formatted string `json:"formatted"`
	Valid     bool   `json:"valid"`
	Error     string `json:"error,omitempty"`
	Size      int    `json:"size"`
}

// DiffResult 差异对比结果
type DiffResult struct {
	Type     string `json:"type"` // "equal", "added", "removed", "modified"
	OldLine  string `json:"oldLine"`
	NewLine  string `json:"newLine"`
	OldNum   int    `json:"oldNum"`
	NewNum   int    `json:"newNum"`
}

// CodecResult 编解码结果
type CodecResult struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// ColorResult 颜色转换结果
type ColorResult struct {
	HEX   string `json:"hex"`
	RGB   string `json:"rgb"`
	HSL   string `json:"hsl"`
	HSV   string `json:"hsv"`
	Name  string `json:"name,omitempty"`
	Error string `json:"error,omitempty"`
}

// FormatJSON 格式化 JSON 字符串
func FormatJSON(input string) JSONResult {
	var raw interface{}
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return JSONResult{
			Valid: false,
			Error: fmt.Sprintf("JSON 格式错误: %v", err),
		}
	}

	formatted, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return JSONResult{
			Valid: false,
			Error: fmt.Sprintf("格式化失败: %v", err),
		}
	}

	return JSONResult{
		Formatted: string(formatted),
		Valid:     true,
		Size:      len(formatted),
	}
}

// MinifyJSON 压缩 JSON
func MinifyJSON(input string) JSONResult {
	var raw interface{}
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return JSONResult{
			Valid: false,
			Error: fmt.Sprintf("JSON 格式错误: %v", err),
		}
	}

	compacted, err := json.Marshal(raw)
	if err != nil {
		return JSONResult{
			Valid: false,
			Error: fmt.Sprintf("压缩失败: %v", err),
		}
	}

	return JSONResult{
		Formatted: string(compacted),
		Valid:     true,
		Size:      len(compacted),
	}
}

// DiffText 文本差异对比（简单行对比）
func DiffText(oldText, newText string) []DiffResult {
	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")

	var results []DiffResult
	maxLen := len(oldLines)
	if len(newLines) > maxLen {
		maxLen = len(newLines)
	}

	for i := 0; i < maxLen; i++ {
		switch {
		case i >= len(oldLines):
			results = append(results, DiffResult{
				Type:    "added",
				NewLine: newLines[i],
				NewNum:  i + 1,
			})
		case i >= len(newLines):
			results = append(results, DiffResult{
				Type:   "removed",
				OldLine: oldLines[i],
				OldNum: i + 1,
			})
		case oldLines[i] == newLines[i]:
			results = append(results, DiffResult{
				Type:   "equal",
				OldLine: oldLines[i],
				NewLine: newLines[i],
				OldNum:  i + 1,
				NewNum:  i + 1,
			})
		default:
			results = append(results, DiffResult{
				Type:    "modified",
				OldLine: oldLines[i],
				NewLine: newLines[i],
				OldNum:  i + 1,
				NewNum:  i + 1,
			})
		}
	}

	return results
}

// EncodeBase64 编码 Base64
func EncodeBase64(input string) CodecResult {
	return CodecResult{
		Input:  input,
		Output: base64.StdEncoding.EncodeToString([]byte(input)),
	}
}

// DecodeBase64 解码 Base64
func DecodeBase64(input string) CodecResult {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return CodecResult{
			Input: input,
			Error: fmt.Sprintf("Base64 解码失败: %v", err),
		}
	}
	return CodecResult{
		Input:  input,
		Output: string(decoded),
	}
}

// EncodeURL URL 编码
func EncodeURL(input string) CodecResult {
	return CodecResult{
		Input:  input,
		Output: url.QueryEscape(input),
	}
}

// DecodeURL URL 解码
func DecodeURL(input string) CodecResult {
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return CodecResult{
			Input: input,
			Error: fmt.Sprintf("URL 解码失败: %v", err),
		}
	}
	return CodecResult{
		Input:  input,
		Output: decoded,
	}
}

// TestRegex 测试正则表达式
type RegexTestResult struct {
	Pattern  string   `json:"pattern"`
	Text     string   `json:"text"`
	Matches  []string `json:"matches"`
	Count    int      `json:"count"`
	Error    string   `json:"error,omitempty"`
}

// TestRegex 测试正则匹配
func TestRegex(pattern, text string) RegexTestResult {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return RegexTestResult{
			Pattern: pattern,
			Text:    text,
			Error:   fmt.Sprintf("正则表达式错误: %v", err),
		}
	}

	matches := re.FindAllString(text, -1)
	if matches == nil {
		matches = []string{}
	}

	return RegexTestResult{
		Pattern: pattern,
		Text:    text,
		Matches: matches,
		Count:   len(matches),
	}
}

// TimestampConvert 时间戳转换
type TimestampResult struct {
	UnixTimestamp int64  `json:"unixTimestamp"`
	DateTime      string `json:"dateTime"`
	ISO8601       string `json:"iso8601"`
}

// ConvertTimestamp 转换时间戳
func ConvertTimestamp(timestamp int64, fromUnit string) TimestampResult {
	var t time.Time
	switch fromUnit {
	case "s":
		t = time.Unix(timestamp, 0)
	case "ms":
		t = time.UnixMilli(timestamp)
	default:
		t = time.Unix(timestamp, 0)
	}

	return TimestampResult{
		UnixTimestamp: t.Unix(),
		DateTime:      t.Format("2006-01-02 15:04:05"),
		ISO8601:       t.Format(time.RFC3339),
	}
}

// ParseTimeString 从时间字符串转换
func ParseTimeString(timeStr string) TimestampResult {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006/01/02 15:04:05",
		"2006-01-02",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return ConvertTimestamp(t.Unix(), "s")
		}
	}

	return TimestampResult{
		DateTime: fmt.Sprintf("无法解析: %s", timeStr),
	}
}

// ConvertColor 颜色值转换
func ConvertColor(hex string) ColorResult {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return ColorResult{HEX: "#" + hex, Error: "无效的颜色值"}
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)

	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := max(rf, max(gf, bf))
	min := min(rf, min(gf, bf))

	var h, s, l, v float64

	// HSL
	l = (max + min) / 2

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case rf:
			h = 60 * (((gf - bf) / d) + 0)
		case gf:
			h = 60 * (((bf - rf) / d) + 2)
		case bf:
			h = 60 * (((rf - gf) / d) + 4)
		}
		if h < 0 {
			h += 360
		}
	}

	// HSV
	v = max
	if max == 0 {
		s = 0
	} else {
		s = 1 - min/max
	}

	return ColorResult{
		HEX:  "#" + hex,
		RGB:  fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
		HSL:  fmt.Sprintf("hsl(%.0f, %.1f%%, %.1f%%)", h, s*100, l*100),
		HSV:  fmt.Sprintf("hsv(%.0f, %.1f%%, %.1f%%)", h, s*100, v*100),
	}
}

// GenerateUUID 生成标准 UUID v4（使用 crypto/rand）
func GenerateUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// 极端情况下降级到时间戳方案
		return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			big.NewInt(0).SetBytes(b[:4]).Uint64(),
			big.NewInt(0).SetBytes(b[4:6]).Uint64(),
			big.NewInt(0).SetBytes(b[6:8]).Uint64(),
			big.NewInt(0).SetBytes(b[8:10]).Uint64(),
			big.NewInt(0).SetBytes(b[10:16]).Uint64())
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// CronExpression Cron 表达式结构
type CronExpression struct {
	Expression  string `json:"expression"`
	Description string `json:"description"`
	NextRun     string `json:"nextRun"`
}

// ParseCron 解析 Cron 表达式，计算下次执行时间
func ParseCron(expr string) CronExpression {
	// 尝试解析 5 位或 6 位（含秒）cron 表达式
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.SecondOptional)
	schedule, err := parser.Parse(expr)
	if err != nil {
		return CronExpression{
			Expression:  expr,
			Description: fmt.Sprintf("无效的 Cron 表达式: %v", err),
		}
	}

	desc := describeCron(expr)
	next := schedule.Next(time.Now()).Format("2006-01-02 15:04:05")

	return CronExpression{
		Expression:  expr,
		Description: desc,
		NextRun:     next,
	}
}

func describeCron(expr string) string {
	parts := strings.Fields(expr)
	if len(parts) < 5 {
		return "自定义间隔: " + expr
	}

	// 生成人类可读描述
	minute, hour := parts[0], parts[1]
	dom, month, dow := parts[2], parts[3], parts[4]

	switch {
	case minute == "0" && hour == "0":
		return "每天午夜执行"
	case minute == "0" && hour == "*":
		return "每小时整点执行"
	case minute == "*/5":
		return "每 5 分钟执行"
	case minute == "*/15":
		return "每 15 分钟执行"
	case minute == "*/30":
		return "每 30 分钟执行"
	case minute == "0" && dom == "*" && month == "*" && dow == "*":
		return "每天 " + hour + ":00 执行"
	case dow != "*" && dow != "?":
		return "每周 " + dow + " 的 " + hour + ":" + minute + " 执行"
	default:
		return "按计划: " + expr
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
