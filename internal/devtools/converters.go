package devtools

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ============================================================
//  Base32/Hex 互转
// ============================================================

// Base32ToHex Base32 转 Hex
func Base32ToHex(input string) (string, error) {
	data, err := base32.StdEncoding.DecodeString(strings.ToUpper(strings.TrimSpace(input)))
	if err != nil {
		return "", fmt.Errorf("Base32 解码失败: %w", err)
	}
	return hex.EncodeToString(data), nil
}

// HexToBase32 Hex 转 Base32
func HexToBase32(input string) (string, error) {
	data, err := hex.DecodeString(strings.TrimSpace(input))
	if err != nil {
		return "", fmt.Errorf("Hex 解码失败: %w", err)
	}
	return base32.StdEncoding.EncodeToString(data), nil
}

// ============================================================
//  时间戳转换
// ============================================================

// TimestampConvert 时间戳转换结果
type TimestampConvert struct {
	Timestamp int64  `json:"timestamp"`
	DateTime  string `json:"dateTime"`
	Unix10    int64  `json:"unix10"`
	Unix13    int64  `json:"unix13"`
	ISO8601   string `json:"iso8601"`
}

// TimestampToDate 时间戳转日期
func TimestampToDate(ts int64) (*TimestampConvert, error) {
	// 自动判断是秒级还是毫秒级
	var t time.Time
	if ts > 9999999999 {
		// 毫秒级
		t = time.Unix(0, ts*int64(time.Millisecond))
	} else {
		// 秒级
		t = time.Unix(ts, 0)
	}

	return &TimestampConvert{
		Timestamp: ts,
		DateTime:  t.Format("2006-01-02 15:04:05"),
		Unix10:    t.Unix(),
		Unix13:    t.UnixMilli(),
		ISO8601:   t.Format(time.RFC3339),
	}, nil
}

// DateToTimestamp 日期转时间戳
func DateToTimestamp(dateStr string) (*TimestampConvert, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
		"2006-01-02T15:04:05Z07:00",
	}

	var t time.Time
	var err error
	for _, format := range formats {
		t, err = time.ParseInLocation(format, strings.TrimSpace(dateStr), time.Local)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("无法解析日期: %s", dateStr)
	}

	return &TimestampConvert{
		Timestamp: t.Unix(),
		DateTime:  t.Format("2006-01-02 15:04:05"),
		Unix10:    t.Unix(),
		Unix13:    t.UnixMilli(),
		ISO8601:   t.Format(time.RFC3339),
	}, nil
}

// GetNowTimestamp 获取当前时间戳
func GetNowTimestamp() *TimestampConvert {
	now := time.Now()
	return &TimestampConvert{
		Timestamp: now.Unix(),
		DateTime:  now.Format("2006-01-02 15:04:05"),
		Unix10:    now.Unix(),
		Unix13:    now.UnixMilli(),
		ISO8601:   now.Format(time.RFC3339),
	}
}

// ============================================================
//  颜色代码转换
// ============================================================

// ColorConvert 颜色转换结果
type ColorConvert struct {
	Hex  string `json:"hex"`
	RGB  string `json:"rgb"`
	HSL  string `json:"hsl"`
	R    int    `json:"r"`
	G    int    `json:"g"`
	B    int    `json:"b"`
	Hue  int    `json:"hue"`
	Sat  int    `json:"sat"`
	Lit  int    `json:"lit"`
}

// HexToColor Hex 转 RGB/HSL
func HexToColor(hexStr string) (*ColorConvert, error) {
	hexStr = strings.TrimPrefix(hexStr, "#")

	var r, g, b uint8
	var err error

	if len(hexStr) == 3 {
		// 短格式 #RGB -> #RRGGBB
		hexStr = string([]byte{hexStr[0], hexStr[0], hexStr[1], hexStr[1], hexStr[2], hexStr[2]})
	}

	if len(hexStr) != 6 {
		return nil, fmt.Errorf("无效的 HEX 颜色: %s", hexStr)
	}

	r64, err := strconv.ParseUint(hexStr[0:2], 16, 8)
	if err != nil {
		return nil, err
	}
	g64, err := strconv.ParseUint(hexStr[2:4], 16, 8)
	if err != nil {
		return nil, err
	}
	b64, err := strconv.ParseUint(hexStr[4:6], 16, 8)
	if err != nil {
		return nil, err
	}
	r = uint8(r64)
	g = uint8(g64)
	b = uint8(b64)

	h, s, l := rgbToHsl(int(r), int(g), int(b))

	return &ColorConvert{
		Hex: fmt.Sprintf("#%02X%02X%02X", r, g, b),
		RGB: fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
		HSL: fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l),
		R:   int(r), G: int(g), B: int(b),
		Hue: h, Sat: s, Lit: l,
	}, nil
}

// RGBToColor RGB 转 HEX/HSL
func RGBToColor(r, g, b int) (*ColorConvert, error) {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return nil, fmt.Errorf("RGB 值范围 0-255")
	}

	h, s, l := rgbToHsl(r, g, b)

	return &ColorConvert{
		Hex: fmt.Sprintf("#%02X%02X%02X", r, g, b),
		RGB: fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
		HSL: fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l),
		R:   r, G: g, B: b,
		Hue: h, Sat: s, Lit: l,
	}, nil
}

// HSLToColor HSL 转 HEX/RGB
func HSLToColor(h, s, l int) (*ColorConvert, error) {
	if h < 0 || h > 360 || s < 0 || s > 100 || l < 0 || l > 100 {
		return nil, fmt.Errorf("HSL 值范围: H(0-360) S(0-100) L(0-100)")
	}

	r, g, b := hslToRgb(h, s, l)

	return &ColorConvert{
		Hex: fmt.Sprintf("#%02X%02X%02X", r, g, b),
		RGB: fmt.Sprintf("rgb(%d, %d, %d)", r, g, b),
		HSL: fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l),
		R:   r, G: g, B: b,
		Hue: h, Sat: s, Lit: l,
	}, nil
}

func rgbToHsl(r, g, b int) (int, int, int) {
	rf := float64(r) / 255
	gf := float64(g) / 255
	bf := float64(b) / 255

	max := rf
	if gf > max {
		max = gf
	}
	if bf > max {
		max = bf
	}
	min := rf
	if gf < min {
		min = gf
	}
	if bf < min {
		min = bf
	}

	l := (max + min) / 2
	var h, s float64

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
			h = (gf - bf) / d
			if gf < bf {
				h += 6
			}
		case gf:
			h = (bf-rf)/d + 2
		case bf:
			h = (rf-gf)/d + 4
		}
		h /= 6
	}

	return int(h * 360), int(s * 100), int(l * 100)
}

func hslToRgb(h, s, l int) (int, int, int) {
	hf := float64(h) / 360
	sf := float64(s) / 100
	lf := float64(l) / 100

	var r, g, b float64

	if sf == 0 {
		r, g, b = lf, lf, lf
	} else {
		var q float64
		if lf < 0.5 {
			q = lf * (1 + sf)
		} else {
			q = lf + sf - lf*sf
		}
		p := 2*lf - q

		r = hue2rgb(p, q, hf+1.0/3)
		g = hue2rgb(p, q, hf)
		b = hue2rgb(p, q, hf-1.0/3)
	}

	return int(r * 255), int(g * 255), int(b * 255)
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}
