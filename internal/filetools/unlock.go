package filetools

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"strings"
)

// LockInfo 文件锁定信息
type LockInfo struct {
	FilePath  string `json:"filePath"`
	ProcessID int    `json:"processId"`
	ProcessName string `json:"processName"`
}

// UnlockResult 解锁结果
type UnlockResult struct {
	FilePath    string `json:"filePath"`
	Success     bool   `json:"success"`
	ReleasedBy []string `json:"releasedBy,omitempty"`
	Error       string `json:"error,omitempty"`
}

// CheckLocks 检查哪些进程占用了文件
func CheckLocks(path string) ([]LockInfo, error) {
	// 使用 PowerShell 的 Get-Process + 文件句柄检测
	psScript := fmt.Sprintf(`
$file = "%s"
$handle = Get-Process | Where-Object { $_.Modules.FileName -like $file }
if ($handle) {
	$handle | Select-Object Id, ProcessName | ConvertTo-Json
}
`, escapePath(path))

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 非致命，返回空
		return nil, nil
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return nil, nil
	}

	return parseLockInfo(result), nil
}

// TryUnlock 尝试解除文件占用
func TryUnlock(path string) UnlockResult {
	result := UnlockResult{FilePath: path}

	// 方法1: 尝试直接删除
	if err := os.Remove(path); err == nil {
		result.Success = true
		result.ReleasedBy = []string{"直接删除"}
		return result
	}

	// 方法2: 重命名文件（有时可解除占用）
	tempPath := path + ".old"
	if err := os.Rename(path, tempPath); err == nil {
		os.Remove(tempPath)
		result.Success = true
		result.ReleasedBy = []string{"重命名后删除"}
		return result
	}

	// 方法3: 使用 PowerShell 终止占用进程
	psScript := fmt.Sprintf(`
$file = "%s"
$procs = Get-Process | Where-Object { $_.Modules.FileName -eq $file }
$procs | ForEach-Object { Stop-Process -Id $_.Id -Force }
`, escapePath(path))

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if output, err := cmd.CombinedOutput(); err == nil {
		_ = output
		// 再次尝试删除
		if err := os.Remove(path); err == nil {
			result.Success = true
			result.ReleasedBy = []string{"终止进程后删除"}
			return result
		}
	}

	// 方法4: 使用系统工具解除
	// 调用 sysinternals handle.exe（如果存在）
	handlePath := `C:\Program Files\Sysinternals\handle64.exe`
	if _, err := os.Stat(handlePath); err == nil {
		cmd := exec.Command(handlePath, "-a", "-u", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if output, err := cmd.Output(); err == nil {
			_ = output
			// handle 输出分析略
		}
	}

	result.Error = "无法解除文件占用（文件正在被使用）"
	return result
}

func parseLockInfo(jsonOutput string) []LockInfo {
	var locks []LockInfo
	lines := strings.Split(jsonOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Id") {
			// 简化解析
			locks = append(locks, LockInfo{
				FilePath:    "",
				ProcessID:   0,
				ProcessName: "未知进程",
			})
		}
	}
	return locks
}

func escapePath(path string) string {
	return strings.ReplaceAll(path, "'", "''")
}
