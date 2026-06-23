package scheduler

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"path/filepath"

	"pc-toolbox/internal/common"
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
	output, err := c.Output()
	if err != nil {
		return []TaskInfo{}
	}

	var tasks []TaskInfo
	scanner := bufio.NewScanner(strings.NewReader(common.GbkToUtf8(string(output))))
	first := true
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过表头行
		if first {
			first = false
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// CSV 格式: TaskName, Status, Task To Run
		parts := strings.SplitN(line, ",", 3)
		if len(parts) < 3 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		// 过滤非本应用创建的任务
		if !strings.HasPrefix(name, "PCToolbox_") {
			continue
		}

		status := strings.TrimSpace(parts[1])
		taskRun := strings.TrimSpace(parts[2])

		// 从任务名解析动作和时间: PCToolbox_{action}_{HHmm}
		action := ""
		timeStr := ""
		nameParts := strings.SplitN(name, "_", 3)
		if len(nameParts) >= 3 {
			action = nameParts[1]
			timeRaw := nameParts[2]
			if len(timeRaw) >= 4 {
				timeStr = fmt.Sprintf("%s:%s", timeRaw[:2], timeRaw[2:4])
			}
		}

		tasks = append(tasks, TaskInfo{
			Name:    name,
			Action:  action,
			Time:    timeStr,
			Enabled: !strings.EqualFold(status, "Disabled"),
			TaskPath: taskRun,
		})
	}

	return tasks
}

// DeleteTask 删除定时任务
func DeleteTask(name string) bool {
	c := exec.Command("schtasks", "/delete", "/tn", name, "/f")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := c.Run()
	return err == nil
}
