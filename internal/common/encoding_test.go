package common

import (
	"testing"
)

func TestGbkToUtf8_EmptyString(t *testing.T) {
	result := GbkToUtf8("")
	if result != "" {
		t.Errorf("空字符串应返回空字符串, got %q", result)
	}
}

func TestUtf8ToGbk_EmptyString(t *testing.T) {
	result := Utf8ToGbk("")
	if result != "" {
		t.Errorf("空字符串应返回空字符串, got %q", result)
	}
}

func TestGbkToUtf8_RoundTrip(t *testing.T) {
	// UTF-8 → GBK → UTF-8 往返测试
	original := "你好世界，这是一段中文测试文本。"
	gbkBytes := Utf8ToGbk(original)
	roundTripped := GbkToUtf8(gbkBytes)

	if roundTripped != original {
		t.Errorf("往返转换不一致: got %q, want %q", roundTripped, original)
	}
}
