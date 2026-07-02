package clipboard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"pc-toolbox/internal/common"
)

// ClipItem 剪贴板记录
type ClipItem struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Type      string `json:"type"` // "text", "image"
	Time      string `json:"time"`
	Size      int    `json:"size"`
}

var (
	history []ClipItem
	mu      sync.Mutex
	nextID  = 1
	maxSize = 50
)

func init() {
	// 启动时加载历史记录
	loadHistory()
}

// getHistoryPath 获取历史记录文件路径
func getHistoryPath() string {
	appData, err := os.UserConfigDir()
	if err != nil {
		appData = "."
	}
	dir := filepath.Join(appData, "pc-toolbox")
	if err := os.MkdirAll(dir, 0755); err != nil {
		// 目录创建失败，使用当前目录作为回退
		dir = "."
	}
	return filepath.Join(dir, "clipboard_history.json")
}

// loadHistory 从文件加载剪贴板历史
func loadHistory() {
	mu.Lock()
	defer mu.Unlock()

	path := getHistoryPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return // 文件不存在则忽略
	}

	var items []ClipItem
	if err := json.Unmarshal(data, &items); err != nil {
		return // 文件损坏则忽略
	}

	history = items
	if len(items) > 0 {
		// 恢复 nextID 为最大 ID + 1
		maxID := 0
		for _, item := range items {
			if item.ID > maxID {
				maxID = item.ID
			}
		}
		nextID = maxID + 1
	}
}

// saveHistory 保存剪贴板历史到文件
func saveHistory() {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return
	}

	path := getHistoryPath()
	if err := os.WriteFile(path, data, 0644); err != nil {
		// 日志记录保存失败，但不中断程序
		fmt.Fprintf(os.Stderr, "保存剪贴板历史失败: %v\n", err)
	}
}

// AddItem 添加剪贴板记录
func AddItem(content string, contentType string) ClipItem {
	mu.Lock()
	defer mu.Unlock()

	item := ClipItem{
		ID:      nextID,
		Content: content,
		Type:    contentType,
		Time:    time.Now().Format("15:04:05"),
		Size:    len(content),
	}
	nextID++

	history = append([]ClipItem{item}, history...)
	if len(history) > maxSize {
		history = history[:maxSize]
	}

	// 异步保存到文件
	common.GoSafe(saveHistory)

	return item
}

// GetHistory 获取剪贴板历史
func GetHistory() []ClipItem {
	mu.Lock()
	defer mu.Unlock()

	result := make([]ClipItem, len(history))
	copy(result, history)
	return result
}

// ClearHistory 清空剪贴板历史
func ClearHistory() {
	mu.Lock()
	defer mu.Unlock()
	history = nil

	common.GoSafe(saveHistory)
}

// RemoveItem 删除单条记录
func RemoveItem(id int) {
	mu.Lock()
	defer mu.Unlock()

	for i, item := range history {
		if item.ID == id {
			history = append(history[:i], history[i+1:]...)
			common.GoSafe(saveHistory)
			return
		}
	}
}
