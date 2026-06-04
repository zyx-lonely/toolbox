package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

// ShredResult 安全删除结果
type ShredResult struct {
	Path      string `json:"path"`
	Success   bool   `json:"success"`
	Passes    int    `json:"passes"`
	Error     string `json:"error,omitempty"`
}

// PasswordResult 密码生成结果
type PasswordResult struct {
	Password string `json:"password"`
	Strength string `json:"strength"` // "weak", "medium", "strong"
	Length   int    `json:"length"`
}

// EncryptResult 加密结果
type EncryptResult struct {
	InputPath      string `json:"inputPath"`
	OutputPath     string `json:"outputPath"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
	Algorithm      string `json:"algorithm,omitempty"`
	FileSize       int64  `json:"fileSize,omitempty"`
	EncryptedSize  int64  `json:"encryptedSize,omitempty"`
}

// ShredFile 安全删除文件（覆写后删除）
func ShredFile(path string, passes int) ShredResult {
	if passes <= 0 {
		passes = 3 // DoD 5220.22-M 标准
	}

	info, err := os.Stat(path)
	if err != nil {
		return ShredResult{
			Path:    path,
			Success: false,
			Error:   fmt.Sprintf("文件访问失败: %v", err),
		}
	}

	size := info.Size()

	for pass := 0; pass < passes; pass++ {
		if err := overwriteFile(path, size, pass); err != nil {
			return ShredResult{
				Path:    path,
				Success: false,
				Passes:  pass,
				Error:   fmt.Sprintf("第 %d 遍覆写失败: %v", pass+1, err),
			}
		}
	}

	// 最后删除文件
	if err := os.Remove(path); err != nil {
		return ShredResult{
			Path:    path,
			Success: false,
			Passes:  passes,
			Error:   fmt.Sprintf("删除失败: %v", err),
		}
	}

	return ShredResult{
		Path:    path,
		Success: true,
		Passes:  passes,
	}
}

// ShredDir 安全删除目录中的所有文件
func ShredDir(dir string, passes int) []ShredResult {
	var results []ShredResult

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			result := ShredFile(path, passes)
			results = append(results, result)
		}
		return nil
	})

	return results
}

func overwriteFile(path string, size int64, pass int) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	// 生成覆写数据
	buf := make([]byte, 4096)
	for i := range buf {
		switch pass {
		case 0:
			buf[i] = 0x00 // 第一遍: 全 0
		case 1:
			buf[i] = 0xFF // 第二遍: 全 1
		default:
			n, _ := rand.Int(rand.Reader, big.NewInt(256))
			buf[i] = byte(n.Int64()) // 后续: 随机数据
		}
	}

	remaining := size
	for remaining > 0 {
		writeSize := int64(len(buf))
		if writeSize > remaining {
			writeSize = remaining
		}
		if _, err := f.Write(buf[:writeSize]); err != nil {
			return err
		}
		remaining -= writeSize
	}

	return f.Sync()
}

// GeneratePassword 生成随机密码
func GeneratePassword(length int, useUpper bool, useLower bool, useDigits bool, useSpecial bool) PasswordResult {
	if length < 4 {
		length = 8
	}

	var charset string
	if useUpper {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if useLower {
		charset += "abcdefghijklmnopqrstuvwxyz"
	}
	if useDigits {
		charset += "0123456789"
	}
	if useSpecial {
		charset += "!@#$%^&*()-_=+[]{}|;:,.<>?"
	}

	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	}

	var password strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			continue
		}
		password.WriteByte(charset[n.Int64()])
	}

	pwd := password.String()
	strength := evaluateStrength(pwd)

	return PasswordResult{
		Password: pwd,
		Strength: strength,
		Length:   length,
	}
}

func evaluateStrength(pwd string) string {
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, c := range pwd {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	score := 0
	if len(pwd) >= 12 {
		score++
	}
	if len(pwd) >= 16 {
		score++
	}
	if hasUpper && hasLower {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}

	switch {
	case score >= 4:
		return "strong"
	case score >= 2:
		return "medium"
	default:
		return "weak"
	}
}

// EncryptFile 加密文件
func EncryptFile(inputPath string, password string) EncryptResult {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	outputPath := inputPath + ".enc"
	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	return EncryptResult{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Success:    true,
	}
}

// DecryptFile 解密文件
func DecryptFile(inputPath string, password string) EncryptResult {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	key := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return EncryptResult{InputPath: inputPath, Error: "数据损坏"}
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: "解密失败（密码错误或数据损坏）"}
	}

	outputPath := strings.TrimSuffix(inputPath, ".enc")
	outputPath = outputPath + ".decrypted"

	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return EncryptResult{InputPath: inputPath, Error: err.Error()}
	}

	return EncryptResult{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Success:    true,
	}
}

// ClearRecentDocs 清理最近文档记录
func ClearRecentDocs() error {
	recentDir := filepath.Join(os.Getenv("USERPROFILE"), "Recent")
	if _, err := os.Stat(recentDir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(recentDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(recentDir, entry.Name())
		os.Remove(path)
	}

	return nil
}

// EncodeBase64 Base64 编码
func EncodeBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// DecodeBase64 Base64 解码
func DecodeBase64(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
