package main

import (
	"pc-toolbox/internal/scheduler"
)

// ============================================================
//  定时任务
// ============================================================

// CreateScheduledTask 创建定时任务
func (a *App) CreateScheduledTask(action string, hour int, minute int) scheduler.TaskInfo {
	return scheduler.CreateTask(scheduler.TaskAction(action), hour, minute)
}

// ListScheduledTasks 列出定时任务
func (a *App) ListScheduledTasks() []scheduler.TaskInfo {
	return scheduler.ListTasks()
}

// DeleteScheduledTask 删除定时任务
func (a *App) DeleteScheduledTask(name string) bool {
	return scheduler.DeleteTask(name)
}
