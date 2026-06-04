package devtools

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// JWTResult JWT 解码结果
type JWTResult struct {
	Raw       string `json:"raw"`
	Header    string `json:"header"`
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
	Valid     bool   `json:"valid"`
	Error     string `json:"error,omitempty"`
}

// DecodeJWT 解码 JWT Token
func DecodeJWT(token string) JWTResult {
	if token == "" {
		return JWTResult{Error: "Token 不能为空"}
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return JWTResult{Raw: token, Error: "无效的 JWT 格式（需要 3 段）"}
	}

	result := JWTResult{Raw: token, Signature: parts[2], Valid: true}

	// 解码 Header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err == nil {
		result.Header = prettyJSON(string(headerBytes))
	} else {
		result.Header = parts[0] + " (解码失败)"
	}

	// 解码 Payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err == nil {
		result.Payload = prettyJSON(string(payloadBytes))
	} else {
		result.Payload = parts[1] + " (解码失败)"
	}

	return result
}

func prettyJSON(s string) string {
	var data interface{}
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return s + "\n(非 JSON 格式)"
	}
	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return s
	}
	return string(pretty)
}
