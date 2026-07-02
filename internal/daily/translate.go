package daily

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"pc-toolbox/internal/common"
)

// TranslateResult 翻译结果
type TranslateResult struct {
	Source     string `json:"source"`
	Target     string `json:"target"`
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
	Provider   string `json:"provider"`
}

// Translate 翻译文本（自动检测语言）
func Translate(text string, fromLang string, toLang string) (*TranslateResult, error) {
	if text == "" {
		return nil, fmt.Errorf("文本不能为空")
	}

	if fromLang == "" {
		fromLang = detectLanguage(text)
	}
	if toLang == "" {
		if fromLang == "zh" {
			toLang = "en"
		} else {
			toLang = "zh"
		}
	}

	result, err := myMemoryTranslate(text, fromLang, toLang)
	if err != nil {
		result, err = googleTranslate(text, fromLang, toLang)
		if err != nil {
			return nil, fmt.Errorf("翻译失败: %w", err)
		}
	}

	return result, nil
}

func detectLanguage(text string) string {
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FFF {
			return "zh"
		}
	}
	for _, r := range text {
		if (r >= 0x3040 && r <= 0x309F) || (r >= 0x30A0 && r <= 0x30FF) {
			return "ja"
		}
	}
	for _, r := range text {
		if r >= 0xAC00 && r <= 0xD7AF {
			return "ko"
		}
	}
	return "en"
}

func googleTranslate(text string, from, to string) (*TranslateResult, error) {
	apiURL := "https://translate.googleapis.com/translate_a/single"
	params := url.Values{}
	params.Set("client", "gtx")
	params.Set("sl", from)
	params.Set("tl", to)
	params.Set("dt", "t")
	params.Set("q", text)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw [][][]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("解析响应失败")
	}

	var translated strings.Builder
	for _, segment := range raw {
		for _, s := range segment {
			if len(s) > 0 {
				if str, ok := s[0].(string); ok {
					translated.WriteString(str)
				}
			}
		}
	}

	return &TranslateResult{
		Source:     text,
		Target:     translated.String(),
		SourceLang: from,
		TargetLang: to,
		Provider:   "Google",
	}, nil
}

func myMemoryTranslate(text string, from, to string) (*TranslateResult, error) {
	// MyMemory API 限制 500 字符，长文本需要分段
	const maxChunkSize = 450
	if len(text) <= maxChunkSize {
		return myMemoryTranslateSingle(text, from, to)
	}

	// 按句子分段
	chunks := splitText(text, maxChunkSize)
	var translated strings.Builder

	for _, chunk := range chunks {
		result, err := myMemoryTranslateSingle(chunk, from, to)
		if err != nil {
			// 某段失败则返回已翻译的部分
			if translated.Len() > 0 {
				break
			}
			return nil, err
		}
		translated.WriteString(result.Target)
	}

	return &TranslateResult{
		Source:     text,
		Target:     translated.String(),
		SourceLang: from,
		TargetLang: to,
		Provider:   "MyMemory",
	}, nil
}

func myMemoryTranslateSingle(text string, from, to string) (*TranslateResult, error) {
	apiURL := "https://api.mymemory.translated.net/get"
	params := url.Values{}
	params.Set("q", text)
	params.Set("langpair", from+"|"+to)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败")
	}

	if result.ResponseData.TranslatedText == "" {
		return nil, fmt.Errorf("翻译服务返回空结果")
	}

	return &TranslateResult{
		Source:     text,
		Target:     result.ResponseData.TranslatedText,
		SourceLang: from,
		TargetLang: to,
		Provider:   "MyMemory",
	}, nil
}

// splitText 按句子边界分段
func splitText(text string, maxSize int) []string {
	var chunks []string
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?' || r == '。' || r == '！' || r == '？' || r == '\n'
	})

	current := ""
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if len(current)+len(s)+1 > maxSize && current != "" {
			chunks = append(chunks, current)
			current = s
		} else {
			if current != "" {
				current += " " + s
			} else {
				current = s
			}
		}
	}
	if current != "" {
		chunks = append(chunks, current)
	}
	return chunks
}

// GetClipboardAndTranslate 获取剪贴板内容并翻译
func GetClipboardAndTranslate(toLang string) (*TranslateResult, error) {
	text := getClipboardText()
	if text == "" {
		return nil, fmt.Errorf("剪贴板为空或不包含文本")
	}
	return Translate(text, "", toLang)
}

func getClipboardText() string {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "Get-Clipboard")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(common.GbkToUtf8(string(out)))
}
