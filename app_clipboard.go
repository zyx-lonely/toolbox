package main

import (
	"pc-toolbox/internal/clipboard"
	"pc-toolbox/internal/common"
)

// ============================================================
//  剪贴板历史
// ============================================================

// AddClipboardItem 添加剪贴板记录
func (a *App) AddClipboardItem(content string, contentType string) common.APIResponse {
	item := clipboard.AddItem(content, contentType)
	return common.NewSuccessResponse(item)
}

// GetClipboardHistory 获取剪贴板历史
func (a *App) GetClipboardHistory() common.APIResponse {
	items := clipboard.GetHistory()
	return common.NewSuccessResponse(items)
}

// ClearClipboardHistory 清空剪贴板历史
func (a *App) ClearClipboardHistory() common.APIResponse {
	clipboard.ClearHistory()
	return common.NewSuccessResponse(nil)
}

// RemoveClipboardItem 删除单条记录
func (a *App) RemoveClipboardItem(id int) common.APIResponse {
	clipboard.RemoveItem(id)
	return common.NewSuccessResponse(nil)
}

// ReadClipboardImage 读取剪贴板图片并返回 base64
func (a *App) ReadClipboardImage() common.APIResponse {
	img, err := clipboard.ReadClipboardImage()
	if err != nil {
		return common.NewErrorResponseStr(err.Error())
	}
	return common.NewSuccessResponse(img)
}

// GetClipboardText 读取剪贴板文本
func (a *App) GetClipboardText() string {
	return clipboard.ReadClipboardText()
}

// SaveClipboardImage 保存剪贴板图片到文件
func (a *App) SaveClipboardImage(filePath string) common.APIResponse {
	err := clipboard.SaveClipboardImageToFile(filePath)
	if err != nil {
		return common.NewErrorResponseStr(err.Error())
	}
	return common.NewSuccessResponse(nil)
}
