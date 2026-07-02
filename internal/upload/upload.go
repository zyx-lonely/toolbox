package upload

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// UploadResult 上传结果
type UploadResult struct {
	Success    bool   `json:"success"`
	FileName   string `json:"fileName"`
	ServerURL  string `json:"serverUrl"`
	StatusCode int    `json:"statusCode"`
	Response   string `json:"response"`
	Error      string `json:"error,omitempty"`
}

// validateURL 验证 URL 安全性，防止 SSRF
func validateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("无效 URL: %v", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("仅支持 http/https 协议")
	}
	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("URL 缺少主机名")
	}

	// 解析主机名，检查是否为内网地址
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("无法解析主机名: %v", err)
	}
	for _, ip := range ips {
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() ||
			ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
			return fmt.Errorf("不允许访问内网地址: %s", host)
		}
	}
	return nil
}

// UploadFileToServer 将 Base64 编码的文件上传到指定服务器
// fileData: Base64 编码的文件数据
// fileName: 文件名（如 "photo_20260604.jpg"）
// serverURL: 服务器上传接口地址
// fieldName: 表单字段名（如 "file"）
func UploadFileToServer(fileData string, fileName string, serverURL string, fieldName string) UploadResult {
	if fileData == "" {
		return UploadResult{Error: "文件数据为空"}
	}
	if serverURL == "" {
		return UploadResult{Error: "服务器地址为空"}
	}
	if fieldName == "" {
		fieldName = "file"
	}

	if err := validateURL(serverURL); err != nil {
		return UploadResult{Error: fmt.Sprintf("URL 验证失败: %v", err)}
	}

	// 解析 Base64 数据
	// 支持 data:image/png;base64,xxxx 格式或纯 base64
	data := fileData
	if idx := strings.Index(data, ","); idx >= 0 && strings.HasPrefix(data, "data:") {
		data = data[idx+1:]
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return UploadResult{Error: fmt.Sprintf("Base64 解码失败: %v", err)}
	}

	// 构造 multipart/form-data 请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return UploadResult{Error: fmt.Sprintf("创建表单文件失败: %v", err)}
	}
	if _, err := part.Write(decoded); err != nil {
		return UploadResult{Error: fmt.Sprintf("写入文件数据失败: %v", err)}
	}
	writer.Close()

	// 发送请求
	httpReq, err := http.NewRequest("POST", serverURL, body)
	if err != nil {
		return UploadResult{Error: fmt.Sprintf("创建请求失败: %v", err)}
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return UploadResult{Error: fmt.Sprintf("上传请求失败: %v", err)}
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	return UploadResult{
		Success:    resp.StatusCode >= 200 && resp.StatusCode < 300,
		FileName:   fileName,
		ServerURL:  serverURL,
		StatusCode: resp.StatusCode,
		Response:   strings.TrimSpace(string(respBody)),
	}
}
