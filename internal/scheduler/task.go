package scheduler

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
)

// TaskAction 定时任务动作
type TaskAction string

const (
	ActionCleanup  TaskAction = "cleanup"
	ActionScreenshot TaskAction = "screenshot"
	ActionShutdown TaskAction = "shutdown"
)

// TaskInfo 定时任务信息
type TaskInfo struct {
	Name     string `json:"name"`
	Action   string `json:"action"`
	Time     string `json:"time"`     // HH:mm 格式
	Enabled  bool   `json:"enabled"`
	TaskPath string `json:"taskPath,omitempty"`
}

// CreateTask 创建计划任务
func CreateTask(action TaskAction, hour, minute int) TaskInfo {
	exePath, _ := os.Executable()

	taskName := fmt.Sprintf("PCToolbox_%s_%02d%02d", action, hour, minute)
	timeStr := fmt.Sprintf("%02d:%02d", hour, minute)

	var args string
	switch action {
	case ActionCleanup:
		args = "--clean-temp"
	case ActionScreenshot:
		args = "--screenshot"
	case ActionShutdown:
		args = "--shutdown"
	}

	cmdStr := fmt.Sprintf(
		`schtasks /create /tn "%s" /tr "\"%s\" %s" /sc daily /st %s /f`,
		taskName, exePath, args, timeStr)

	c := exec.Command("cmd", "/c", cmdStr)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c.Run()

	return TaskInfo{
		Name:     taskName,
		Action:   string(action),
		Time:     timeStr,
		Enabled:  true,
		TaskPath: filepath.Join(`\PCToolbox`, taskName),
	}
}

// ListTasks 列出已创建的定时任务
func ListTasks() []TaskInfo {
	c := exec.Command("schtasks", "/query", "/fo", "CSV", "/nh")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, _ := c.Output()
	_ = output
	return []TaskInfo{}
}

// DeleteTask 删除定时任务
func DeleteTask(name string) bool {
	c := exec.Command("schtasks", "/delete", "/tn", name, "/f")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := c.Run()
	return err == nil
}
