package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level 日志级别
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger 日志记录器
type Logger struct {
	level  Level
	file   *os.File
	mu     sync.Mutex
	prefix string
}

var (
	defaultLogger *Logger
	once          sync.Once
	initMu        sync.Mutex
)

// Init 初始化日志系统
func Init(logDir string, level Level) error {
	once.Do(func() {
		if logDir == "" {
			logDir = "."
		}

		// 创建日志目录
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "创建日志目录失败: %v\n", err)
			return
		}

		// 清理过期日志
		go cleanOldLogs(logDir, 7)

		// 打开日志文件
		logFile := filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "打开日志文件失败: %v\n", err)
			return
		}

		defaultLogger = &Logger{
			level: level,
			file:  f,
		}
	})

	return nil
}

// cleanOldLogs 自动清理超过 delayDays 天的日志文件
func cleanOldLogs(logDir string, delayDays int) {
	cutoff := time.Now().AddDate(0, 0, -delayDays)
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// 匹配 app_YYYY-MM-DD.log 格式
		if len(name) < 15 || name[:4] != "app_" || name[len(name)-4:] != ".log" {
			continue
		}
		dateStr := name[4 : len(name)-4]
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if t.Before(cutoff) {
			_ = os.Remove(filepath.Join(logDir, name))
		}
	}
}

// GetLogger 获取默认日志实例
func GetLogger() *Logger {
	initMu.Lock()
	defer initMu.Unlock()
	if defaultLogger == nil {
		// 如果未初始化，创建一个控制台输出的 logger
		defaultLogger = &Logger{
			level: LevelInfo,
			file:  os.Stdout,
		}
	}
	return defaultLogger
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetPrefix 设置日志前缀
func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	var levelStr string
	switch level {
	case LevelDebug:
		levelStr = "DEBUG"
	case LevelInfo:
		levelStr = "INFO"
	case LevelWarn:
		levelStr = "WARN"
	case LevelError:
		levelStr = "ERROR"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	prefix := ""
	if l.prefix != "" {
		prefix = "[" + l.prefix + "] "
	}

	logLine := fmt.Sprintf("%s %s %s%s\n", timestamp, levelStr, prefix, msg)

	if l.file != nil {
		l.file.WriteString(logLine)
	}
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.file != nil && l.file != os.Stdout {
		return l.file.Close()
	}
	return nil
}

// 便捷函数
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}
