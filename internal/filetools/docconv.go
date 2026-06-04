package filetools

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// DocConvertResult 文档转换结果
type DocConvertResult struct {
	InputPath  string `json:"inputPath"`
	OutputPath string `json:"outputPath"`
	TargetType string `json:"targetType"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

// ConvertToPDF 将 Office 文档转换为 PDF（需要安装 Office 或 LibreOffice）
func ConvertToPDF(inputPath string) DocConvertResult {
	ext := strings.ToLower(filepath.Ext(inputPath))
	supported := map[string]bool{".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true}

	if !supported[ext] {
		return DocConvertResult{InputPath: inputPath, TargetType: "pdf", Error: "不支持的文件格式: " + ext}
	}

	outPath := inputPath[:len(inputPath)-len(ext)] + ".pdf"

	// 尝试使用 LibreOffice（如果安装）
	c := exec.Command("soffice", "--headless", "--convert-to", "pdf", "--outdir", filepath.Dir(outPath), inputPath)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := c.Run(); err == nil {
		return DocConvertResult{InputPath: inputPath, OutputPath: outPath, TargetType: "pdf", Success: true}
	}

	// 尝试使用 Office COM（需要 Windows + Office 安装）
	c2 := exec.Command("powershell", "-NoProfile",
		fmt.Sprintf(`$wd = New-Object -ComObject Word.Application; $wd.Visible = $false; $doc = $wd.Documents.Open('%s'); $doc.SaveAs('%s', 17); $doc.Close(); $wd.Quit()`,
			inputPath, outPath))
	c2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := c2.Run(); err == nil {
		return DocConvertResult{InputPath: inputPath, OutputPath: outPath, TargetType: "pdf", Success: true}
	}

	return DocConvertResult{InputPath: inputPath, TargetType: "pdf", Success: false, Error: "转换失败，请安装 LibreOffice 或 Microsoft Office"}
}
