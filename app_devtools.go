package main

import (
	"pc-toolbox/internal/common"
	"pc-toolbox/internal/devtools"
)

// ============================================================
//  开发工具模块
// ============================================================

// FormatJSON 格式化 JSON
func (a *App) FormatJSON(input string) devtools.JSONResult {
	return devtools.FormatJSON(input)
}

// MinifyJSON 压缩 JSON
func (a *App) MinifyJSON(input string) devtools.JSONResult {
	return devtools.MinifyJSON(input)
}

// DiffText 文本差异对比
func (a *App) DiffText(oldText, newText string) []devtools.DiffResult {
	return devtools.DiffText(oldText, newText)
}

// EncodeBase64 Base64 编码
func (a *App) EncodeBase64(input string) devtools.CodecResult {
	return devtools.EncodeBase64(input)
}

// DecodeBase64 Base64 解码
func (a *App) DecodeBase64(input string) devtools.CodecResult {
	return devtools.DecodeBase64(input)
}

// ReadFileAsBase64 读取文件并返回 Base64 编码
func (a *App) ReadFileAsBase64(path string) string {
	return common.ReadFileAsBase64(path)
}

// EncodeURL URL 编码
func (a *App) EncodeURL(input string) devtools.CodecResult {
	return devtools.EncodeURL(input)
}

// DecodeURL URL 解码
func (a *App) DecodeURL(input string) devtools.CodecResult {
	return devtools.DecodeURL(input)
}

// TestRegex 测试正则表达式
func (a *App) TestRegex(pattern, text string) devtools.RegexTestResult {
	return devtools.TestRegex(pattern, text)
}

// ConvertTimestamp 转换时间戳
func (a *App) ConvertTimestamp(timestamp int64, fromUnit string) devtools.TimestampResult {
	return devtools.ConvertTimestamp(timestamp, fromUnit)
}

// ConvertColor 颜色值转换
func (a *App) ConvertColor(hex string) devtools.ColorResult {
	return devtools.ConvertColor(hex)
}

// GenerateUUID 生成 UUID
func (a *App) GenerateUUID() string {
	return devtools.GenerateUUID()
}

// SendHTTPRequest 发送 HTTP 请求
func (a *App) SendHTTPRequest(req devtools.HTTPRequest) devtools.HTTPResponse {
	return devtools.SendHTTPRequest(req)
}

// DecodeJWT 解码 JWT
func (a *App) DecodeJWT(token string) devtools.JWTResult {
	return devtools.DecodeJWT(token)
}

// FormatYAML 格式化 YAML
func (a *App) FormatYAML(input string) devtools.FormatResult {
	return devtools.FormatYAML(input)
}

// FormatTOML 格式化 TOML
func (a *App) FormatTOML(input string) devtools.FormatResult {
	return devtools.FormatTOML(input)
}

// GenerateUUIDs 批量生成 UUID
func (a *App) GenerateUUIDs(count int, version int) devtools.UUIDGenResult {
	return devtools.GenerateUUIDs(count, version)
}

// CheckUpdate 检查更新
func (a *App) CheckUpdate(currentVersion string) devtools.ReleaseInfo {
	return devtools.CheckUpdate(currentVersion)
}

// BeautifyHTML 美化 HTML
func (a *App) BeautifyHTML(input string) devtools.CodeBeautifyResult {
	return devtools.BeautifyHTML(input)
}

// BeautifyCSS 美化 CSS
func (a *App) BeautifyCSS(input string) devtools.CodeBeautifyResult {
	return devtools.BeautifyCSS(input)
}

// BeautifySQL 美化 SQL
func (a *App) BeautifySQL(input string) devtools.CodeBeautifyResult {
	return devtools.BeautifySQL(input)
}

// GenerateQRCode 生成二维码
func (a *App) GenerateQRCode(content string, size int) devtools.QRResult {
	return devtools.GenerateQRCode(content, size)
}
