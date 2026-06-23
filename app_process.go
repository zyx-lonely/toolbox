package main

import (
	"fmt"

	"pc-toolbox/internal/common"
	"pc-toolbox/internal/process"
)

// ============================================================
//  进程管理
// ============================================================

// GetProcessList 获取进程列表
func (a *App) GetProcessList() ([]process.ProcessInfo, error) {
	return process.GetProcessList()
}

// KillProcess 结束进程
func (a *App) KillProcess(pid int) error {
	return process.KillProcess(pid)
}

// SetProcessPriority 设置进程优先级
// priority: "idle", "below_normal", "normal", "above_normal", "high", "realtime"
func (a *App) SetProcessPriority(pid int, priority string) common.APIResponse {
	var priorityClass uint32
	switch priority {
	case "idle":
		priorityClass = process.IDLE_PRIORITY_CLASS
	case "below_normal":
		priorityClass = process.BELOW_NORMAL_PRIORITY_CLASS
	case "normal":
		priorityClass = process.NORMAL_PRIORITY_CLASS
	case "above_normal":
		priorityClass = process.ABOVE_NORMAL_PRIORITY_CLASS
	case "high":
		priorityClass = process.HIGH_PRIORITY_CLASS
	case "realtime":
		priorityClass = process.REALTIME_PRIORITY_CLASS
	default:
		return common.NewErrorResponseStr(fmt.Sprintf("无效的优先级: %s", priority))
	}

	err := process.SetProcessPriority(pid, priorityClass)
	if err != nil {
		return common.NewErrorResponseStr(err.Error())
	}
	return common.NewSuccessResponse(nil)
}

// BatchKillProcesses 批量结束进程
func (a *App) BatchKillProcesses(pids []int) common.APIResponse {
	var failed []string
	for _, pid := range pids {
		err := process.KillProcess(pid)
		if err != nil {
			failed = append(failed, fmt.Sprintf("PID %d: %v", pid, err))
		}
	}

	if len(failed) > 0 {
		return common.NewErrorResponseStr(fmt.Sprintf("部分进程结束失败:\n%s", joinStrings(failed, "\n")))
	}
	return common.NewSuccessResponse(nil)
}

// joinStrings 连接字符串数组
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
