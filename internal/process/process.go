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

// GetProcessList 获取进程列表（使用 tasklist）
func GetProcessList() ([]ProcessInfo, error) {
	// 直接使用完整路径，避免 PATH 找不到
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

		processes = append(processes, ProcessInfo{
			PID:     pid,
			Name:    name,
			Memory:  memStr,
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
