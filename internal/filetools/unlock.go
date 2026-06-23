package filetools

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// LockInfo 文件锁定信息
type LockInfo struct {
	FilePath    string `json:"filePath"`
	ProcessID   int    `json:"processId"`
	ProcessName string `json:"processName"`
}

// UnlockResult 解锁结果
type UnlockResult struct {
	FilePath   string   `json:"filePath"`
	Success    bool     `json:"success"`
	ReleasedBy []string `json:"releasedBy,omitempty"`
	Error      string   `json:"error,omitempty"`
}

// CheckLocks 检查哪些进程占用了文件
func CheckLocks(path string) ([]LockInfo, error) {
	psScript := `
param($FilePath)
$procs = Get-Process | Where-Object {
    try { $_.Modules.FileName -contains $FilePath } catch { $false }
}
$procs | Select-Object Id, ProcessName | ConvertTo-Json -Compress
`
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript, "-FilePath", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
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

	if err := os.Remove(path); err == nil {
		result.Success = true
		result.ReleasedBy = []string{"直接删除"}
		return result
	}

	tempPath := path + ".old"
	if err := os.Rename(path, tempPath); err == nil {
		os.Remove(tempPath)
		result.Success = true
		result.ReleasedBy = []string{"重命名后删除"}
		return result
	}

	psScript := `
param($FilePath)
$procs = Get-Process | Where-Object {
    try { $_.Modules.FileName -contains $FilePath } catch { $false }
}
$procs | ForEach-Object { Stop-Process -Id $_.Id -Force }
`
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript, "-FilePath", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if _, err := cmd.CombinedOutput(); err == nil {
		if err := os.Remove(path); err == nil {
			result.Success = true
			result.ReleasedBy = []string{"终止进程后删除"}
			return result
		}
	}

	result.Error = "无法解除文件占用（文件正在被使用）"
	return result
}

func parseLockInfo(jsonOutput string) []LockInfo {
	var locks []LockInfo
	jsonOutput = strings.TrimSpace(jsonOutput)
	if jsonOutput == "" {
		return nil
	}
	if strings.HasPrefix(jsonOutput, "[") {
		var arr []struct {
			Id          int    `json:"Id"`
			ProcessName string `json:"ProcessName"`
		}
		if err := json.Unmarshal([]byte(jsonOutput), &arr); err == nil {
			for _, p := range arr {
				locks = append(locks, LockInfo{ProcessID: p.Id, ProcessName: p.ProcessName})
			}
			return locks
		}
	}
	var single struct {
		Id          int    `json:"Id"`
		ProcessName string `json:"ProcessName"`
	}
	if err := json.Unmarshal([]byte(jsonOutput), &single); err == nil && single.Id > 0 {
		locks = append(locks, LockInfo{ProcessID: single.Id, ProcessName: single.ProcessName})
	}
	return locks
}
