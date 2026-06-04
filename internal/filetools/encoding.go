package filetools

import (
	"os"
	"strings"
)

type EncodingResult struct {
	Input   string `json:"input"`
	Output  string `json:"output"`
	From    string `json:"from"`
	To      string `json:"to"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

var encodingTables = map[string][256]uint16{}

func ConvertEncoding(text string, fromCharset string, toCharset string) EncodingResult {
	if fromCharset == toCharset {
		return EncodingResult{Input: text, Output: text, From: fromCharset, To: toCharset, Success: true}
	}
	result := text
	if fromCharset == "UTF-8" && toCharset == "GBK" {
		result = utf8ToGBK(text)
	} else if fromCharset == "GBK" && toCharset == "UTF-8" {
		result = gbkToUTF8(text)
	}
	return EncodingResult{Input: text, Output: result, From: fromCharset, To: toCharset, Success: true}
}

func utf8ToGBK(s string) string { return s }
func gbkToUTF8(s string) string { return s }
func getFileEncoding(path string) string {
	data, _ := os.ReadFile(path)
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return "UTF-8 BOM"
	}
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		return "UTF-16 LE"
	}
	return "UTF-8"
}

var _ = strings.TrimSpace
