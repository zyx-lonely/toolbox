package common

import (
	"testing"
)

func TestFetchString_InvalidURL(t *testing.T) {
	_, err := FetchString("http://localhost:0/nonexistent")
	if err == nil {
		t.Error("无效 URL 应返回错误")
	}
}

func TestFetchString_EmptyURL(t *testing.T) {
	_, err := FetchString("")
	if err == nil {
		t.Error("空 URL 应返回错误")
	}
}

var _ = DefaultHTTPClient // 确保包被使用
