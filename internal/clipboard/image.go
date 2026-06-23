package clipboard

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"syscall"
	"unsafe"
)

var (
	user32                    = syscall.NewLazyDLL("user32.dll")
	procOpenClipboard         = user32.NewProc("OpenClipboard")
	procCloseClipboard        = user32.NewProc("CloseClipboard")
	procIsClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable")
	procGetClipboardData      = user32.NewProc("GetClipboardData")
	procEmptyClipboard        = user32.NewProc("EmptyClipboard")
	procSetClipboardData      = user32.NewProc("SetClipboardData")

	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	procGlobalLock   = kernel32.NewProc("GlobalLock")
	procGlobalUnlock = kernel32.NewProc("GlobalUnlock")
	procGlobalSize   = kernel32.NewProc("GlobalSize")
	procGlobalAlloc  = kernel32.NewProc("GlobalAlloc")
	procCopyMemory   = kernel32.NewProc("RtlMoveMemory")
)

const (
	CF_UNICODETEXT = 13
	CF_DIB         = 8
	GMEM_MOVEABLE  = 0x0002
)

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

// dibToImage 将 DIB 数据转换为 image.RGBA
func dibToImage(data []byte) (*image.RGBA, error) {
	if len(data) < 40 {
		return nil, fmt.Errorf("无效的 DIB 数据")
	}

	infoHeader := *(*BITMAPINFOHEADER)(unsafe.Pointer(&data[0]))

	width := int(infoHeader.BiWidth)
	height := int(infoHeader.BiHeight)
	if height < 0 {
		height = -height
	}

	bpp := int(infoHeader.BiBitCount)
	rowSize := ((bpp*width + 31) / 32) * 4

	pixelOffset := int(infoHeader.BiSize)
	if int(infoHeader.BiClrUsed) > 0 {
		pixelOffset += int(infoHeader.BiClrUsed) * 4
	} else if bpp <= 8 {
		pixelOffset += (1 << uint(bpp)) * 4
	}

	if pixelOffset >= len(data) {
		return nil, fmt.Errorf("DIB 数据不完整")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	switch bpp {
	case 32:
		for y := 0; y < height; y++ {
			srcRow := pixelOffset + (height-1-y)*rowSize
			for x := 0; x < width; x++ {
				srcIdx := srcRow + x*4
				if srcIdx+3 >= len(data) {
					continue
				}
				b, g, r, a := data[srcIdx], data[srcIdx+1], data[srcIdx+2], data[srcIdx+3]
				img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
			}
		}
	case 24:
		for y := 0; y < height; y++ {
			srcRow := pixelOffset + (height-1-y)*rowSize
			for x := 0; x < width; x++ {
				srcIdx := srcRow + x*3
				if srcIdx+2 >= len(data) {
					continue
				}
				b, g, r := data[srcIdx], data[srcIdx+1], data[srcIdx+2]
				img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}
	default:
		return nil, fmt.Errorf("不支持的位深度: %d", bpp)
	}

	return img, nil
}

// readDIBFromClipboard 从剪贴板读取 DIB 数据
func readDIBFromClipboard() ([]byte, error) {
	r, _, _ := procOpenClipboard.Call(0)
	if r == 0 {
		return nil, fmt.Errorf("无法打开剪贴板")
	}
	defer procCloseClipboard.Call()

	r, _, _ = procIsClipboardFormatAvailable.Call(CF_DIB)
	if r == 0 {
		return nil, fmt.Errorf("剪贴板中没有图片")
	}

	hMem, _, _ := procGetClipboardData.Call(CF_DIB)
	if hMem == 0 {
		return nil, fmt.Errorf("无法获取剪贴板图片数据")
	}

	pMem, _, _ := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return nil, fmt.Errorf("无法锁定内存")
	}
	defer procGlobalUnlock.Call(hMem)

	size, _, _ := procGlobalSize.Call(hMem)
	if size == 0 {
		return nil, fmt.Errorf("无法获取内存大小")
	}

	data := make([]byte, size)
	procCopyMemory.Call(
		uintptr(unsafe.Pointer(&data[0])),
		pMem,
		size,
	)

	return data, nil
}

func ReadClipboardImage() (string, error) {
	dibData, err := readDIBFromClipboard()
	if err != nil {
		return "", err
	}

	img, err := dibToImage(dibData)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("编码图片失败: %w", err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func ReadClipboardText() string {
	r, _, _ := procOpenClipboard.Call(0)
	if r == 0 {
		return ""
	}
	defer procCloseClipboard.Call()

	hMem, _, _ := procGetClipboardData.Call(CF_UNICODETEXT)
	if hMem == 0 {
		return ""
	}

	pMem, _, _ := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return ""
	}
	defer procGlobalUnlock.Call(hMem)

	size, _, _ := procGlobalSize.Call(hMem)
	if size == 0 {
		return ""
	}

	data := make([]byte, size)
	procCopyMemory.Call(
		uintptr(unsafe.Pointer(&data[0])),
		pMem,
		size,
	)

	var result []byte
	for i := 0; i+1 < len(data); i += 2 {
		lo := uint16(data[i])
		hi := uint16(data[i+1])
		ch := lo | (hi << 8)
		if ch == 0 {
			break
		}
		if ch < 0x80 {
			result = append(result, byte(ch))
		} else if ch < 0x800 {
			result = append(result, byte(0xC0|(ch>>6)), byte(0x80|(ch&0x3F)))
		} else {
			result = append(result, byte(0xE0|(ch>>12)), byte(0x80|((ch>>6)&0x3F)), byte(0x80|(ch&0x3F)))
		}
	}

	return string(result)
}

func SetClipboardText(text string) error {
	r, _, _ := procOpenClipboard.Call(0)
	if r == 0 {
		return fmt.Errorf("无法打开剪贴板")
	}
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	utf16, _ := syscall.UTF16FromString(text)
	dataSize := len(utf16) * 2

	hMem, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, uintptr(dataSize+2))
	if hMem == 0 {
		return fmt.Errorf("无法分配内存")
	}

	pMem, _, _ := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return fmt.Errorf("无法锁定内存")
	}

	procCopyMemory.Call(pMem, uintptr(unsafe.Pointer(&utf16[0])), uintptr(dataSize))
	procGlobalUnlock.Call(hMem)

	ret, _, _ := procSetClipboardData.Call(CF_UNICODETEXT, hMem)
	if ret == 0 {
		return fmt.Errorf("设置剪贴板失败")
	}

	return nil
}

func SaveClipboardImageToFile(filePath string) error {
	dibData, err := readDIBFromClipboard()
	if err != nil {
		return err
	}

	img, err := dibToImage(dibData)
	if err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
