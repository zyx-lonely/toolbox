package main

import (
	"pc-toolbox/internal/security"
)

// ============================================================
//  安全与隐私模块
// ============================================================

// ShredFile 安全删除文件
func (a *App) ShredFile(path string, passes int) security.ShredResult {
	return security.ShredFile(path, passes)
}

// ShredDir 安全删除目录中的所有文件
func (a *App) ShredDir(dir string, passes int) []security.ShredResult {
	return security.ShredDir(dir, passes)
}

// GeneratePassword 生成密码
func (a *App) GeneratePassword(length int, useUpper bool, useLower bool, useDigits bool, useSpecial bool) security.PasswordResult {
	return security.GeneratePassword(length, useUpper, useLower, useDigits, useSpecial)
}

// EncryptFile 加密文件
func (a *App) EncryptFile(inputPath string, password string) security.EncryptResult {
	return security.EncryptFile(inputPath, password)
}

// DecryptFile 解密文件
func (a *App) DecryptFile(inputPath string, password string) security.EncryptResult {
	return security.DecryptFile(inputPath, password)
}

// ClearRecentDocs 清理最近文档记录
func (a *App) ClearRecentDocs() error {
	return security.ClearRecentDocs()
}
