package optimize

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"pc-toolbox/internal/logger"
)

// CleanupTask 清理任务配置
type CleanupTask struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Enabled     bool     `json:"enabled"`
	Schedule    string   `json:"schedule"`    // "daily", "weekly", "monthly", "interval"
	DayOfWeek   int      `json:"dayOfWeek"`   // 0-6, 周日到周六
	DayOfMonth  int      `json:"dayOfMonth"`  // 1-31
	Hour        int      `json:"hour"`        // 0-23
	Minute      int      `json:"minute"`      // 0-59
	IntervalMin int      `json:"intervalMin"` // 间隔分钟数
	Targets     []string `json:"targets"`     // 清理目标路径
	AutoDelete  bool     `json:"autoDelete"`  // 自动删除过期文件
	MaxAgeDays  int      `json:"maxAgeDays"`  // 最大保留天数
	LastRun     string   `json:"lastRun"`
	NextRun     string   `json:"nextRun"`
	RunCount    int      `json:"runCount"`
}

// CleanupLog 清理日志
type CleanupLog struct {
	TaskID    string `json:"taskId"`
	TaskName  string `json:"taskName"`
	RunTime   string `json:"runTime"`
	FreedSize uint64 `json:"freedSize"`
	FileCount int    `json:"fileCount"`
	Status    string `json:"status"` // "success", "partial", "failed"
	Error     string `json:"error,omitempty"`
}

// CleanupScheduler 清理调度器
type CleanupScheduler struct {
	Tasks []CleanupTask `json:"tasks"`
	Logs  []CleanupLog  `json:"logs"`
}

var (
	schedulerPath string
	scheduler     *CleanupScheduler
)

func init() {
	schedulerPath = filepath.Join(os.Getenv("APPDATA"), "pc-toolbox", "cleanup_scheduler.json")
	scheduler = loadScheduler()
}

// loadScheduler 加载调度器配置
func loadScheduler() *CleanupScheduler {
	data, err := os.ReadFile(schedulerPath)
	if err != nil {
		return &CleanupScheduler{
			Tasks: getDefaultTasks(),
			Logs:  make([]CleanupLog, 0),
		}
	}

	var s CleanupScheduler
	if err := json.Unmarshal(data, &s); err != nil {
		return &CleanupScheduler{
			Tasks: getDefaultTasks(),
			Logs:  make([]CleanupLog, 0),
		}
	}

	return &s
}

// saveScheduler 保存调度器配置
func saveScheduler() error {
	data, err := json.MarshalIndent(scheduler, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(schedulerPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(schedulerPath, data, 0644)
}

// getDefaultTasks 获取默认清理任务
func getDefaultTasks() []CleanupTask {
	return []CleanupTask{
		{
			ID:       "temp_cleanup",
			Name:     "临时文件清理",
			Enabled:  true,
			Schedule: "daily",
			Hour:     3,
			Minute:   0,
			Targets: []string{
				os.Getenv("TEMP"),
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Temp"),
			},
			AutoDelete: true,
			MaxAgeDays: 7,
		},
		{
			ID:       "browser_cache",
			Name:     "浏览器缓存清理",
			Enabled:  false,
			Schedule: "weekly",
			DayOfWeek: 0,
			Hour:     4,
			Minute:   0,
			Targets: []string{
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data", "Default", "Cache"),
				filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Edge", "User Data", "Default", "Cache"),
			},
			AutoDelete: true,
			MaxAgeDays: 30,
		},
		{
			ID:       "windows_update_cache",
			Name:     "Windows 更新缓存清理",
			Enabled:  false,
			Schedule: "monthly",
			DayOfMonth: 1,
			Hour:     5,
			Minute:   0,
			Targets: []string{
				filepath.Join(os.Getenv("WINDIR"), "SoftwareDistribution", "Download"),
			},
			AutoDelete: true,
			MaxAgeDays: 30,
		},
		{
			ID:       "recycle_bin",
			Name:     "回收站清理",
			Enabled:  false,
			Schedule: "weekly",
			DayOfWeek: 0,
			Hour:     3,
			Minute:   30,
			Targets:  []string{"$RECYCLE.BIN"},
			AutoDelete: true,
			MaxAgeDays: 14,
		},
	}
}

// GetAllTasks 获取所有清理任务
func GetAllTasks() []CleanupTask {
	return scheduler.Tasks
}

// GetTask 获取单个任务
func GetTask(id string) *CleanupTask {
	for _, t := range scheduler.Tasks {
		if t.ID == id {
			return &t
		}
	}
	return nil
}

// AddTask 添加清理任务
func AddTask(task CleanupTask) error {
	// 生成 ID
	if task.ID == "" {
		task.ID = fmt.Sprintf("task_%d", time.Now().UnixMilli())
	}

	// 计算下次运行时间
	task.NextRun = calculateNextRun(task)

	scheduler.Tasks = append(scheduler.Tasks, task)
	return saveScheduler()
}

// UpdateTask 更新清理任务
func UpdateTask(task CleanupTask) error {
	for i, t := range scheduler.Tasks {
		if t.ID == task.ID {
			task.NextRun = calculateNextRun(task)
			scheduler.Tasks[i] = task
			return saveScheduler()
		}
	}
	return fmt.Errorf("任务不存在: %s", task.ID)
}

// DeleteTask 删除清理任务
func DeleteTask(id string) error {
	for i, t := range scheduler.Tasks {
		if t.ID == id {
			scheduler.Tasks = append(scheduler.Tasks[:i], scheduler.Tasks[i+1:]...)
			return saveScheduler()
		}
	}
	return fmt.Errorf("任务不存在: %s", id)
}

// ToggleTask 启用/禁用任务
func ToggleTask(id string, enabled bool) error {
	for i, t := range scheduler.Tasks {
		if t.ID == id {
			scheduler.Tasks[i].Enabled = enabled
			return saveScheduler()
		}
	}
	return fmt.Errorf("任务不存在: %s", id)
}

// RunTask 立即执行任务
func RunTask(id string) (*CleanupLog, error) {
	task := GetTask(id)
	if task == nil {
		return nil, fmt.Errorf("任务不存在: %s", id)
	}

	log := CleanupLog{
		TaskID:   task.ID,
		TaskName: task.Name,
		RunTime:  time.Now().Format("2006-01-02 15:04:05"),
		Status:   "success",
	}

	var totalFreed uint64
	var totalFiles int

	for _, target := range task.Targets {
		expanded := os.ExpandEnv(target)
		freed, count, err := cleanupPath(expanded, task.AutoDelete, task.MaxAgeDays)
		if err != nil {
			logger.Warn("清理路径失败: %s, 错误: %v", expanded, err)
			log.Status = "partial"
			log.Error = err.Error()
		}
		totalFreed += freed
		totalFiles += count
	}

	log.FreedSize = totalFreed
	log.FileCount = totalFiles

	// 更新任务状态
	for i, t := range scheduler.Tasks {
		if t.ID == id {
			scheduler.Tasks[i].LastRun = log.RunTime
			scheduler.Tasks[i].NextRun = calculateNextRun(t)
			scheduler.Tasks[i].RunCount++
			break
		}
	}

	// 保存日志
	scheduler.Logs = append(scheduler.Logs, log)
	if len(scheduler.Logs) > 100 {
		scheduler.Logs = scheduler.Logs[len(scheduler.Logs)-100:]
	}

	if err := saveScheduler(); err != nil {
		return nil, err
	}

	return &log, nil
}

// cleanupPath 清理指定路径
func cleanupPath(path string, autoDelete bool, maxAgeDays int) (freed uint64, count int, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}

	if !info.IsDir() {
		if autoDelete && maxAgeDays > 0 {
			cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
			if info.ModTime().Before(cutoff) {
				size := uint64(info.Size())
				if err := os.Remove(path); err == nil {
					return size, 1, nil
				}
			}
		}
		return 0, 0, nil
	}

	cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		if autoDelete && maxAgeDays > 0 && info.ModTime().Before(cutoff) {
			size := uint64(info.Size())
			if e := os.Remove(p); e == nil {
				freed += size
				count++
			}
		} else {
			freed += uint64(info.Size())
			count++
		}
		return nil
	})

	return freed, count, err
}

// calculateNextRun 计算下次运行时间
func calculateNextRun(task CleanupTask) string {
	now := time.Now()
	var next time.Time

	switch task.Schedule {
	case "daily":
		next = time.Date(now.Year(), now.Month(), now.Day(),
			task.Hour, task.Minute, 0, 0, now.Location())
		if next.Before(now) {
			next = next.AddDate(0, 0, 1)
		}

	case "weekly":
		// 找到下一个指定星期几
		daysUntil := task.DayOfWeek - int(now.Weekday())
		if daysUntil < 0 {
			daysUntil += 7
		}
		next = now.AddDate(0, 0, daysUntil)
		next = time.Date(next.Year(), next.Month(), next.Day(),
			task.Hour, task.Minute, 0, 0, next.Location())
		if next.Before(now) {
			next = next.AddDate(0, 0, 7)
		}

	case "monthly":
		next = time.Date(now.Year(), now.Month(), task.DayOfMonth,
			task.Hour, task.Minute, 0, 0, now.Location())
		if next.Before(now) {
			next = next.AddDate(0, 1, 0)
		}

	case "interval":
		if task.IntervalMin > 0 {
			next = now.Add(time.Duration(task.IntervalMin) * time.Minute)
		} else {
			next = now.Add(time.Hour)
		}

	default:
		next = now.AddDate(0, 0, 1)
	}

	return next.Format("2006-01-02 15:04:05")
}

// GetLogs 获取清理日志
func GetLogs(limit int) []CleanupLog {
	if limit <= 0 || limit > len(scheduler.Logs) {
		limit = len(scheduler.Logs)
	}
	return scheduler.Logs[len(scheduler.Logs)-limit:]
}

// ClearLogs 清空日志
func ClearLogs() error {
	scheduler.Logs = make([]CleanupLog, 0)
	return saveScheduler()
}

// GetSchedulerStats 获取调度器统计
func GetSchedulerStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["totalTasks"] = len(scheduler.Tasks)

	enabled := 0
	for _, t := range scheduler.Tasks {
		if t.Enabled {
			enabled++
		}
	}
	stats["enabledTasks"] = enabled
	stats["totalRuns"] = len(scheduler.Logs)

	var totalFreed uint64
	for _, log := range scheduler.Logs {
		totalFreed += log.FreedSize
	}
	stats["totalFreed"] = totalFreed

	return stats
}
