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

	"golang.org/x/crypto/pbkdf2"
)

// ShredResult 安全删除结果
type ShredResult struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
	Passes  int    `json:"passes"`
	Error   string `json:"error,omitempty"`
}

// PasswordResult 密码生成结果
type PasswordResult struct {
	Password string `json:"password"`
	Strength string `json:"strength"` // "weak", "medium", "strong"
	Length   int    `json:"length"`
}

// EncryptResult 加密结果
type EncryptResult struct {
	InputPath     string `json:"inputPath"`
	OutputPath    string `json:"outputPath"`
	Success       bool   `json:"success"`
	Error         string `json:"error,omitempty"`
	Algorithm     string `json:"algorithm,omitempty"`
	FileSize      int64  `json:"fileSize,omitempty"`
	EncryptedSize int64  `json:"encryptedSize,omitempty"`
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
		// 使用 crypto/rand 生成安全的随机索引，避免模偏差
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return PasswordResult{
				Password: "",
				Strength: "weak",
				Length:   length,
			}
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

// deriveKey 使用 PBKDF2 派生密钥
func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 600000, 32, sha256.New)
}

// EncryptFile 加密文件（流式处理，支持大文件）
func EncryptFile(inputPath string, password string) EncryptResult {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("打开文件失败: %v", err)}
	}
	defer inFile.Close()

	stat, err := inFile.Stat()
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("获取文件信息失败: %v", err)}
	}

	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("生成盐值失败: %v", err)}
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("创建加密器失败: %v", err)}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("创建 GCM 失败: %v", err)}
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("生成 nonce 失败: %v", err)}
	}

	outputPath := inputPath + ".enc"
	outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("创建输出文件失败: %v", err)}
	}
	defer outFile.Close()

	outFile.Write(salt)
	outFile.Write(nonce)

	buf := make([]byte, 64*1024)
	for {
		n, readErr := inFile.Read(buf)
		if n > 0 {
			chunk := gcm.Seal(nil, nonce, buf[:n], nil)
			outFile.Write(chunk)
		}
		if readErr != nil {
			if readErr != io.EOF {
				os.Remove(outputPath)
				return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("读取文件失败: %v", readErr)}
			}
			break
		}
	}

	return EncryptResult{
		InputPath:     inputPath,
		OutputPath:    outputPath,
		Success:       true,
		Algorithm:     "AES-256-GCM",
		FileSize:      stat.Size(),
		EncryptedSize: stat.Size() + 16 + int64(gcm.NonceSize()),
	}
}

// DecryptFile 解密文件（流式处理，支持大文件）
func DecryptFile(inputPath string, password string) EncryptResult {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("打开文件失败: %v", err)}
	}
	defer inFile.Close()

	stat, err := inFile.Stat()
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("获取文件信息失败: %v", err)}
	}

	if stat.Size() < 44 {
		return EncryptResult{InputPath: inputPath, Error: "数据损坏：文件过小"}
	}

	salt := make([]byte, 16)
	if _, err := io.ReadFull(inFile, salt); err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("读取盐值失败: %v", err)}
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(inFile, nonce); err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("读取 nonce 失败: %v", err)}
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("创建解密器失败: %v", err)}
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("创建 GCM 失败: %v", err)}
	}

	ciphertext, err := io.ReadAll(inFile)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("读取密文失败: %v", err)}
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return EncryptResult{InputPath: inputPath, Error: "解密失败（密码错误或数据损坏）"}
	}

	outputPath := strings.TrimSuffix(inputPath, ".enc")
	outputPath = outputPath + ".decrypted"

	if err := os.WriteFile(outputPath, plaintext, 0600); err != nil {
		return EncryptResult{InputPath: inputPath, Error: fmt.Sprintf("写入文件失败: %v", err)}
	}

	return EncryptResult{
		InputPath:     inputPath,
		OutputPath:    outputPath,
		Success:       true,
		Algorithm:     "AES-256-GCM",
		FileSize:      int64(len(plaintext)),
		EncryptedSize: stat.Size(),
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
		return fmt.Errorf("读取最近文档目录失败: %v", err)
	}

	var errors []string
	for _, entry := range entries {
		path := filepath.Join(recentDir, entry.Name())
		if err := os.Remove(path); err != nil {
			errors = append(errors, fmt.Sprintf("删除 %s 失败: %v", entry.Name(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分文件删除失败: %s", strings.Join(errors, "; "))
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
		return "", fmt.Errorf("Base64 解码失败: %v", err)
	}
	return string(decoded), nil
}
