package filetools

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// WatermarkOptions 水印选项
type WatermarkOptions struct {
	Text     string `json:"text"`     // 水印文字
	Position string `json:"position"` // "bottom-right", "bottom-left", "top-right", "top-left", "center"
	Opacity  int    `json:"opacity"`  // 0-100
	FontSize int    `json:"fontSize"` // 字体大小
}

// WatermarkResult 水印处理结果
type WatermarkResult struct {
	File   string `json:"file"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Size   int64  `json:"size"`
}

// AddWatermarkToImages 批量添加水印
func AddWatermarkToImages(files []string, opts WatermarkOptions) ([]WatermarkResult, error) {
	if opts.Text == "" {
		return nil, fmt.Errorf("水印文字不能为空")
	}
	if opts.Opacity <= 0 || opts.Opacity > 100 {
		opts.Opacity = 50
	}
	if opts.FontSize <= 0 {
		opts.FontSize = 16
	}
	if opts.Position == "" {
		opts.Position = "bottom-right"
	}

	var results []WatermarkResult
	var errs []error

	for _, file := range files {
		result, err := addWatermarkSingle(file, opts)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", filepath.Base(file), err))
			continue
		}
		results = append(results, result)
	}

	if len(errs) > 0 && len(results) == 0 {
		return nil, fmt.Errorf("处理失败: %v", errs[0])
	}

	return results, nil
}

func addWatermarkSingle(filePath string, opts WatermarkOptions) (WatermarkResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return WatermarkResult{}, err
	}
	defer f.Close()

	// 解码图片
	var img image.Image
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	default:
		return WatermarkResult{}, fmt.Errorf("不支持的格式: %s", ext)
	}
	if err != nil {
		return WatermarkResult{}, err
	}

	// 创建可绘制的图片副本
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// 计算水印位置
	x, y := calcWatermarkPosition(bounds, opts)

	// 绘制水印
	drawWatermark(rgba, x, y, opts)

	// 保存到同目录，加 _watermark 后缀
	outPath := addSuffix(filePath, "_watermark")
	outFile, err := os.Create(outPath)
	if err != nil {
		return WatermarkResult{}, err
	}
	defer outFile.Close()

	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outFile, rgba, &jpeg.Options{Quality: 95})
	case ".png":
		err = png.Encode(outFile, rgba)
	}
	if err != nil {
		return WatermarkResult{}, err
	}

	info, _ := os.Stat(outPath)
	return WatermarkResult{
		File:   outPath,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
		Size:   info.Size(),
	}, nil
}

func calcWatermarkPosition(bounds image.Rectangle, opts WatermarkOptions) (int, int) {
	margin := 10
	textWidth := len(opts.Text) * 8
	textHeight := opts.FontSize

	switch opts.Position {
	case "top-left":
		return bounds.Min.X + margin, bounds.Min.Y + margin + textHeight
	case "top-right":
		return bounds.Max.X - textWidth - margin, bounds.Min.Y + margin + textHeight
	case "bottom-left":
		return bounds.Min.X + margin, bounds.Max.Y - margin
	case "center":
		return (bounds.Dx() - textWidth) / 2, bounds.Dy() / 2
	default: // bottom-right
		return bounds.Max.X - textWidth - margin, bounds.Max.Y - margin
	}
}

func drawWatermark(img *image.RGBA, x, y int, opts WatermarkOptions) {
	// 创建半透明颜色
	alpha := uint8(float64(opts.Opacity) / 100.0 * 255)
	col := color.RGBA{R: 255, G: 255, B: 255, A: alpha}

	face := basicfont.Face7x13
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(opts.Text)
}

func addSuffix(filePath, suffix string) string {
	ext := filepath.Ext(filePath)
	name := strings.TrimSuffix(filePath, ext)
	return name + suffix + ext
}
