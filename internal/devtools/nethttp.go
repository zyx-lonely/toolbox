package devtools

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPRequest HTTP 请求参数
type HTTPRequest struct {
	URL     string `json:"url"`
	Method  string `json:"method"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
}

// HTTPResponse HTTP 响应结果
type HTTPResponse struct {
	StatusCode int    `json:"statusCode"`
	StatusText string `json:"statusText"`
	Headers    string `json:"headers"`
	Body       string `json:"body"`
	Duration   string `json:"duration"`
	Size       int64  `json:"size"`
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
	ips, err := net.LookupIP(host)
	if err == nil {
		for _, ip := range ips {
			if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
				return fmt.Errorf("不允许访问内网/本地地址: %s", host)
			}
		}
	}
	return nil
}

// SendHTTPRequest 发送 HTTP 请求
func SendHTTPRequest(req HTTPRequest) HTTPResponse {
	if req.URL == "" {
		return HTTPResponse{Error: "URL 不能为空"}
	}
	if req.Method == "" {
		req.Method = "GET"
	}

	if err := validateURL(req.URL); err != nil {
		return HTTPResponse{Error: fmt.Sprintf("URL 安全验证失败: %v", err)}
	}

	start := time.Now()

	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return HTTPResponse{Error: fmt.Sprintf("创建请求失败: %v", err)}
	}

	// 设置请求头
	if req.Headers != "" {
		for _, line := range strings.Split(req.Headers, "\n") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				httpReq.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return HTTPResponse{Error: fmt.Sprintf("请求失败: %v", err)}
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	body, _ := io.ReadAll(resp.Body)

	// 限制响应体大小
	bodyStr := string(body)
	if len(bodyStr) > 50000 {
		bodyStr = bodyStr[:50000] + "\n... (响应过长已截断)"
	}

	// 收集响应头
	var headerLines []string
	for k, v := range resp.Header {
		headerLines = append(headerLines, fmt.Sprintf("%s: %s", k, strings.Join(v, ", ")))
	}

	return HTTPResponse{
		StatusCode: resp.StatusCode,
		StatusText: resp.Status,
		Headers:    strings.Join(headerLines, "\n"),
		Body:       bodyStr,
		Duration:   fmt.Sprintf("%dms", duration.Milliseconds()),
		Size:       int64(len(body)),
	}
}
