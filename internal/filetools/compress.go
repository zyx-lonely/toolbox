package filetools

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// BatchCompressResult 批量压缩结果
type BatchCompressResult struct {
	InputPath    string `json:"inputPath"`
	OutputPath   string `json:"outputPath"`
	OriginalSize int64  `json:"originalSize"`
	NewSize      int64  `json:"newSize"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
}

// BatchCompressImages 批量压缩图片
func BatchCompressImages(dir string, quality int, targetFormat string, maxWidth int) []BatchCompressResult {
	var results []BatchCompressResult
	extensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".bmp": true}
	if quality <= 0 { quality = 80 }
	if quality > 100 { quality = 100 }

	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.IsDir() { continue }
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if !extensions[ext] { continue }

		inputPath := filepath.Join(dir, entry.Name())
		result := compressSingleImage(inputPath, quality, targetFormat, maxWidth)
		results = append(results, result)
	}

	return results
}

func compressSingleImage(inputPath string, quality int, targetFormat string, maxWidth int) BatchCompressResult {
	info, _ := os.Stat(inputPath)
	originalSize := int64(0)
	if info != nil { originalSize = info.Size() }

	f, err := os.Open(inputPath)
	if err != nil {
		return BatchCompressResult{InputPath: inputPath, Error: err.Error()}
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return BatchCompressResult{InputPath: inputPath, Error: "解码失败: " + err.Error()}
	}

	// 调整尺寸
	if maxWidth > 0 {
		bounds := img.Bounds()
		w := bounds.Dx()
		h := bounds.Dy()
		if w > maxWidth {
			newW := maxWidth
			newH := int(float64(h) * float64(newW) / float64(w))
			if newH < 1 {
				newH = 1
			}
			resized := image.NewRGBA(image.Rect(0, 0, newW, newH))
			draw.ApproxBiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)
			img = resized
		}
	}

	ext := filepath.Ext(inputPath)
	outPath := inputPath[:len(inputPath)-len(ext)] + "_compressed." + targetFormat

	outFile, err := os.Create(outPath)
	if err != nil {
		return BatchCompressResult{InputPath: inputPath, Error: err.Error()}
	}
	defer outFile.Close()

	switch targetFormat {
	case "jpg", "jpeg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(outFile, img)
	default:
		return BatchCompressResult{InputPath: inputPath, Error: "不支持格式: " + targetFormat}
	}

	if err != nil {
		return BatchCompressResult{InputPath: inputPath, Error: "编码失败: " + err.Error()}
	}

	outInfo, _ := os.Stat(outPath)
	newSize := int64(0)
	if outInfo != nil { newSize = outInfo.Size() }

	return BatchCompressResult{
		InputPath: inputPath, OutputPath: outPath,
		OriginalSize: originalSize, NewSize: newSize,
		Success: true,
	}
}

var _ = fmt.Sprintf
