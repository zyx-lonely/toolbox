package filetools

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FolderChangeEvent 文件变化事件
type FolderChangeEvent struct {
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"` // "created", "modified", "deleted"
	Path      string `json:"path"`
	Size      int64  `json:"size"`
}

// FolderMonitor 文件夹监控器
type FolderMonitor struct {
	mu        sync.Mutex
	path      string
	events    []FolderChangeEvent
	lastState map[string]fileInfo
	running   bool
	stopCh    chan struct{}
	maxEvents int
}

type fileInfo struct {
	size    int64
	modTime time.Time
}

// NewFolderMonitor 创建文件夹监控器
func NewFolderMonitor(path string, maxEvents int) *FolderMonitor {
	if maxEvents <= 0 {
		maxEvents = 1000
	}
	return &FolderMonitor{
		path:      path,
		events:    make([]FolderChangeEvent, 0),
		lastState: make(map[string]fileInfo),
		stopCh:    make(chan struct{}),
		maxEvents: maxEvents,
	}
}

// Start 开始监控
func (m *FolderMonitor) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("已在监控中")
	}

	if _, err := os.Stat(m.path); err != nil {
		return fmt.Errorf("目录不存在: %s", m.path)
	}

	// 记录初始状态
	m.lastState = m.snapshotDir()
	m.running = true

	go m.monitorLoop()
	return nil
}

// Stop 停止监控
func (m *FolderMonitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	m.running = false
	close(m.stopCh)
}

// GetEvents 获取变化事件
func (m *FolderMonitor) GetEvents() []FolderChangeEvent {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]FolderChangeEvent, len(m.events))
	copy(result, m.events)
	return result
}

// ClearEvents 清空事件记录
func (m *FolderMonitor) ClearEvents() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = m.events[:0]
}

// IsRunning 是否正在监控
func (m *FolderMonitor) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

func (m *FolderMonitor) monitorLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.checkChanges()
		case <-m.stopCh:
			return
		}
	}
}

func (m *FolderMonitor) checkChanges() {
	newState := m.snapshotDir()

	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05")

	// 检测新增和修改
	for path, newInfo := range newState {
		oldInfo, exists := m.lastState[path]
		if !exists {
			m.addEvent(FolderChangeEvent{
				Timestamp: now,
				Type:      "created",
				Path:      path,
				Size:      newInfo.size,
			})
		} else if newInfo.modTime.After(oldInfo.modTime) || newInfo.size != oldInfo.size {
			m.addEvent(FolderChangeEvent{
				Timestamp: now,
				Type:      "modified",
				Path:      path,
				Size:      newInfo.size,
			})
		}
	}

	// 检测删除
	for path := range m.lastState {
		if _, exists := newState[path]; !exists {
			m.addEvent(FolderChangeEvent{
				Timestamp: now,
				Type:      "deleted",
				Path:      path,
			})
		}
	}

	m.lastState = newState
}

func (m *FolderMonitor) snapshotDir() map[string]fileInfo {
	state := make(map[string]fileInfo)

	filepath.WalkDir(m.path, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(m.path, p)
		state[relPath] = fileInfo{
			size:    info.Size(),
			modTime: info.ModTime(),
		}
		return nil
	})

	return state
}

func (m *FolderMonitor) addEvent(event FolderChangeEvent) {
	m.events = append(m.events, event)
	if len(m.events) > m.maxEvents {
		m.events = m.events[len(m.events)-m.maxEvents:]
	}
}
