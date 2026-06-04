package devtools

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

type QRResult struct {
	Content string `json:"content"`
	DataURI string `json:"dataUri"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func GenerateQRCode(content string, size int) QRResult {
	if content == "" {
		return QRResult{Error: "内容不能为空"}
	}
	if size <= 0 {
		size = 256
	}
	png, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		return QRResult{Content: content, Error: err.Error()}
	}
	dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	return QRResult{Content: content, DataURI: dataURI, Success: true}
}

var _ = fmt.Sprintf
