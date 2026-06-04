package clipboard

import (
	"sync"
	"time"
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
}

// RemoveItem 删除单条记录
func RemoveItem(id int) {
	mu.Lock()
	defer mu.Unlock()

	for i, item := range history {
		if item.ID == id {
			history = append(history[:i], history[i+1:]...)
			return
		}
	}
}
