package filetools

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type ConvertResult struct {
	InputPath  string `json:"inputPath"`
	OutputPath string `json:"outputPath"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

func ConvertImage(inputPath string, targetFormat string) ConvertResult {
	ext := strings.ToLower(filepath.Ext(inputPath))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".bmp" {
		return ConvertResult{InputPath: inputPath, Error: "不支持的输入格式: " + ext}
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return ConvertResult{InputPath: inputPath, Error: err.Error()}
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return ConvertResult{InputPath: inputPath, Error: "解码失败: " + err.Error()}
	}

	outPath := inputPath[:len(inputPath)-len(ext)] + "." + targetFormat

	outFile, err := os.Create(outPath)
	if err != nil {
		return ConvertResult{InputPath: inputPath, Error: err.Error()}
	}
	defer outFile.Close()

	switch targetFormat {
	case "png":
		err = png.Encode(outFile, img)
	case "jpg", "jpeg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	default:
		return ConvertResult{InputPath: inputPath, Error: "不支持的输出格式: " + targetFormat}
	}

	if err != nil {
		return ConvertResult{InputPath: inputPath, Error: "编码失败: " + err.Error()}
	}

	return ConvertResult{InputPath: inputPath, OutputPath: outPath, Success: true}
}

var _ = fmt.Sprintf
