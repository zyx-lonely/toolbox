package main

import (
	"pc-toolbox/internal/filetools"
)

// ============================================================
//  文件工具模块
// ============================================================

// FindDuplicateFiles 查找重复文件
func (a *App) FindDuplicateFiles(rootPath string, mode string) ([]filetools.DuplicateGroup, error) {
	return filetools.FindDuplicates(rootPath, mode)
}

// ComputeFileHash 计算文件哈希
func (a *App) ComputeFileHash(path string, algorithm string) (string, error) {
	return filetools.ComputeFileHash(path, algorithm)
}

// BatchRenamePreview 批量重命名预览
func (a *App) BatchRenamePreview(dir string, rule filetools.RenameRule) ([]filetools.RenamePreview, error) {
	return filetools.BatchRenamePreview(dir, rule)
}

// BatchRename 执行批量重命名
func (a *App) BatchRename(dir string, rule filetools.RenameRule) ([]filetools.RenamePreview, error) {
	return filetools.BatchRename(dir, rule)
}

// PreviewOrganize 预览归类
func (a *App) PreviewOrganize(dir string, rule filetools.OrganizeRule) ([]filetools.OrganizePreview, error) {
	return filetools.PreviewOrganize(dir, rule)
}

// ExecuteOrganize 执行归类
func (a *App) ExecuteOrganize(dir string, rule filetools.OrganizeRule) ([]filetools.OrganizeResult, error) {
	return filetools.ExecuteOrganize(dir, rule)
}

// FindLargeFiles 查找大文件
func (a *App) FindLargeFiles(rootPath string, minSizeMB int, maxCount int) ([]filetools.LargeFile, error) {
	return filetools.FindLargeFiles(rootPath, minSizeMB, maxCount)
}

// BatchCompressImages 批量压缩图片
func (a *App) BatchCompressImages(dir string, quality int, targetFormat string, maxWidth int) []filetools.BatchCompressResult {
	return filetools.BatchCompressImages(dir, quality, targetFormat, maxWidth)
}

// ConvertToPDF 转换文档到 PDF
func (a *App) ConvertToPDF(inputPath string) filetools.DocConvertResult {
	return filetools.ConvertToPDF(inputPath)
}

// SearchFileContent 搜索文件内容
func (a *App) SearchFileContent(rootDir string, keyword string, fileTypes string) []filetools.SearchResult {
	return filetools.SearchFileContent(rootDir, keyword, fileTypes)
}

// SearchAndReplace 文件内容替换
func (a *App) SearchAndReplace(dir string, search string, replace string, fileTypes string) []filetools.ReplaceResult {
	return filetools.SearchAndReplace(dir, search, replace, fileTypes)
}

// AnalyzeFolderSizes 文件夹大小分析
func (a *App) AnalyzeFolderSizes(rootPath string, depth int) []filetools.FolderSize {
	return filetools.AnalyzeFolderSizes(rootPath, depth)
}

// GetRecycleBinInfo 获取回收站信息
func (a *App) GetRecycleBinInfo() filetools.RecycleBinInfo {
	return filetools.GetRecycleBinInfo()
}

// EmptyRecycleBin 清空回收站
func (a *App) EmptyRecycleBin() error {
	return filetools.EmptyRecycleBin()
}

// DiffFiles 文件差异对比
func (a *App) DiffFiles(oldPath string, newPath string) ([]filetools.DiffLine, error) {
	return filetools.DiffFiles(oldPath, newPath)
}

// ConvertImage 转换图片格式
func (a *App) ConvertImage(inputPath string, targetFormat string) filetools.ConvertResult {
	return filetools.ConvertImage(inputPath, targetFormat)
}

// ConvertEncoding 转换文本编码
func (a *App) ConvertEncoding(text string, fromCharset string, toCharset string) filetools.EncodingResult {
	return filetools.ConvertEncoding(text, fromCharset, toCharset)
}

// PreviewFile 预览文件内容
func (a *App) PreviewFile(path string) filetools.FilePreviewResult {
	return filetools.PreviewFile(path)
}
