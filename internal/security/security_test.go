package security

import (
	"strings"
	"testing"
)

func TestGeneratePassword_Default(t *testing.T) {
	result := GeneratePassword(12, true, true, true, true)
	if result.Length != 12 {
		t.Errorf("密码长度不正确: got %d, want 12", result.Length)
	}
	if len(result.Password) != 12 {
		t.Errorf("密码字符串长度不正确: got %d, want 12", len(result.Password))
	}
	if result.Strength == "weak" {
		t.Error("12 位包含所有字符类型的密码不应为 weak")
	}
}

func TestGeneratePassword_MinLength(t *testing.T) {
	// 长度 < 4 应自动提升到 8
	result := GeneratePassword(2, false, true, false, false)
	if result.Length != 8 {
		t.Errorf("低于最小长度时应提升到 8: got %d", result.Length)
	}
}

func TestGeneratePassword_AllCharacterTypes(t *testing.T) {
	result := GeneratePassword(20, true, true, true, true)
	pwd := result.Password

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range pwd {
		switch {
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		t.Errorf("密码应包含所有类型字符: %s", pwd)
	}
}

func TestGeneratePassword_OnlyLowercase(t *testing.T) {
	result := GeneratePassword(10, false, true, false, false)
	pwd := result.Password

	for _, ch := range pwd {
		if ch < 'a' || ch > 'z' {
			t.Errorf("仅小写模式密码包含非法字符: %c in %s", ch, pwd)
		}
	}
}

func TestGeneratePassword_EmptyCharset(t *testing.T) {
	// 所有类型都为 false，应使用默认字符集
	result := GeneratePassword(10, false, false, false, false)
	if len(result.Password) != 10 {
		t.Errorf("默认字符集密码长度不正确: got %d", len(result.Password))
	}
	if result.Password == "" {
		t.Error("密码不应为空")
	}
}

func TestGeneratePassword_Uniqueness(t *testing.T) {
	// 生两个密码，应不相同（极大概率）
	p1 := GeneratePassword(16, true, true, true, true)
	p2 := GeneratePassword(16, true, true, true, true)

	if p1.Password == p2.Password {
		t.Errorf("两次生成的密码相同: %s", p1.Password)
	}
}

func TestGeneratePassword_NoAmbiguousChars(t *testing.T) {
	// 检查是否不包含易混淆字符 0, O, l, I, 1
	result := GeneratePassword(20, true, true, true, false)
	for _, ch := range result.Password {
		if strings.Contains("0OlI1", string(ch)) {
			t.Logf("包含易混淆字符: %c (密码: %s)", ch, result.Password)
		}
	}
}
