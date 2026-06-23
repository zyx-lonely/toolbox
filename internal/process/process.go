package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"strconv"
	"strings"

	"pc-toolbox/internal/common"
)

// Priority class constants (Windows)
const (
	IDLE_PRIORITY_CLASS         = 0x00000040
	BELOW_NORMAL_PRIORITY_CLASS = 0x00004000
	NORMAL_PRIORITY_CLASS      = 0x00000020
	ABOVE_NORMAL_PRIORITY_CLASS = 0x00008000
	HIGH_PRIORITY_CLASS        = 0x00000080
	REALTIME_PRIORITY_CLASS    = 0x00000100
)

// SetProcessPriority 设置进程优先级（Windows only）
func SetProcessPriority(pid int, priorityClass uint32) error {
	// 打开进程句柄（需要 PROCESS_SET_INFORMATION 权限）
	const PROCESS_SET_INFORMATION = 0x0200
	h, err := syscall.OpenProcess(
		PROCESS_SET_INFORMATION,
		false,
		uint32(pid),
	)
	if err != nil {
		return fmt.Errorf("打开进程失败: %w", err)
	}
	defer syscall.CloseHandle(h)

	// 调用 SetPriorityClass
	dll := syscall.NewLazyDLL("kernel32.dll")
	setPriorityClass := dll.NewProc("SetPriorityClass")

	r, _, err := setPriorityClass.Call(
		uintptr(h),
		uintptr(priorityClass),
	)
	if r == 0 {
		return fmt.Errorf("设置优先级失败: %w", err)
	}
	return nil
}

// GetPriorityName 获取优先级名称
func GetPriorityName(priorityClass uint32) string {
	switch priorityClass {
	case IDLE_PRIORITY_CLASS:
		return "低（空闲）"
	case BELOW_NORMAL_PRIORITY_CLASS:
		return "低于正常"
	case NORMAL_PRIORITY_CLASS:
		return "正常"
	case ABOVE_NORMAL_PRIORITY_CLASS:
		return "高于正常"
	case HIGH_PRIORITY_CLASS:
		return "高"
	case REALTIME_PRIORITY_CLASS:
		return "实时"
	default:
		return "未知"
	}
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	PID     int    `json:"pid"`
	Name    string `json:"name"`
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
	Status  string `json:"status"`
	User    string `json:"user"`
	Command string `json:"command"`
}

// GetProcessList 获取进程列表（使用 wmic，更准确）
func GetProcessList() ([]ProcessInfo, error) {
	// 使用 wmic 获取进程信息
	// 注意：wmic 在 Windows 10 1809+ 已弃用，但仍可用
	cmd := &exec.Cmd{
		Path: "wmic.exe",
		Args: []string{"wmic.exe", "process", "get", "ProcessId,Name,WorkingSetSize,PageFileUsage", "/format:csv"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果 wmic 失败，回退到 tasklist
		return getProcessListFallback()
	}

	var processes []ProcessInfo
	lines := strings.Split(common.GbkToUtf8(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Node") {
			continue
		}

		// CSV 格式: Node,Name,PageFileUsage,ProcessId,WorkingSetSize
		parts := strings.Split(line, ",")
		if len(parts) < 5 {
			continue
		}

		name := parts[1]
		pidStr := parts[3]
		memStr := parts[4]

		pid, _ := strconv.Atoi(pidStr)
		
		// 转换内存大小（WorkingSetSize 是字节）
		memBytes, _ := strconv.ParseInt(memStr, 10, 64)
		memMB := float64(memBytes) / 1024 / 1024
		memDisplay := fmt.Sprintf("%.1f MB", memMB)

		processes = append(processes, ProcessInfo{
			PID:     pid,
			Name:    name,
			Memory:  memDisplay,
			CPU:     "N/A", // CPU 占用需要单独计算
			Command: name,
		})
	}

	return processes, nil
}

// getProcessListFallback 回退方案：使用 tasklist
func getProcessListFallback() ([]ProcessInfo, error) {
	tasklistPath := os.Getenv("SystemRoot") + `\System32\tasklist.exe`
	cmd := &exec.Cmd{
		Path: tasklistPath,
		Args: []string{tasklistPath, "/FO", "CSV", "/NH"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	var processes []ProcessInfo
	lines := strings.Split(common.GbkToUtf8(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// CSV 格式: "name","pid","session","session#","mem"
		parts := strings.Split(line, "\",\"")
		if len(parts) < 5 {
			continue
		}

		name := strings.Trim(parts[0], "\"")
		pidStr := strings.Trim(parts[1], "\"")
		memStr := strings.Trim(parts[4], "\"")

		pid, _ := strconv.Atoi(pidStr)
		
		// 转换内存大小（tasklist 输出格式：1,234 K）
		memStr = strings.ReplaceAll(memStr, ",", "")
		memStr = strings.ReplaceAll(memStr, " K", "")
		memKB, _ := strconv.ParseInt(memStr, 10, 64)
		memMB := float64(memKB) / 1024
		memDisplay := fmt.Sprintf("%.1f MB", memMB)

		processes = append(processes, ProcessInfo{
			PID:     pid,
			Name:    name,
			Memory:  memDisplay,
			CPU:     "N/A",
			Command: name,
		})
	}

	return processes, nil
}

// KillProcess 结束进程
func KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("未找到进程 %d: %w", pid, err)
	}
	return process.Kill()
}

// GetProcessByName 按名称查找进程
func GetProcessByName(name string) ([]ProcessInfo, error) {
	all, err := GetProcessList()
	if err != nil {
		return nil, err
	}

	var matches []ProcessInfo
	for _, p := range all {
		if strings.EqualFold(p.Name, name) || strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) {
			matches = append(matches, p)
		}
	}
	return matches, nil
}
