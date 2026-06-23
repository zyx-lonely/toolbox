package main

import (
	"pc-toolbox/internal/upload"
)

// ============================================================
//  文件上传
// ============================================================

// UploadFileToServer 将 Base64 文件上传到指定服务器
func (a *App) UploadFileToServer(fileData string, fileName string, serverURL string, fieldName string) upload.UploadResult {
	return upload.UploadFileToServer(fileData, fileName, serverURL, fieldName)
}
