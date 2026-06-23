package common

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	// DefaultHTTPClient 带超时的默认 HTTP 客户端（全局复用连接池）
	DefaultHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  false,
		},
	}
)

// FetchString 获取 URL 内容并返回字符串
func FetchString(url string) (string, error) {
	resp, err := DefaultHTTPClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 限制 1MB
	if err != nil {
		return "", err
	}
	return string(body), nil
}
