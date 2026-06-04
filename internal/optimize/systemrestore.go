package optimize

import (
	"fmt"
	"os/exec"
	"syscall"
	"strings"
)

// RestorePointInfo 系统还原点信息
type RestorePointInfo struct {
	Name        string `json:"name"`
	CreatedAt   string `json:"createdAt"`
	Description string `json:"description"`
}

// CreateRestorePoint 创建系统还原点
func CreateRestorePoint(description string) error {
	// 使用 PowerShell 创建还原点
	psScript := fmt.Sprintf(`
$description = "%s"
Checkpoint-Computer -Description $description -RestorePointType MODIFY_SETTINGS
`, escapePowerShellString(description))

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建还原点失败: %w, 输出: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

// GetRestorePoints 获取系统还原点列表
func GetRestorePoints() ([]RestorePointInfo, error) {
	psScript := `
Get-ComputerRestorePoint | Select-Object Description, CreationTime, SequenceNumber | ConvertTo-Json
`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取还原点失败: %w", err)
	}

	// 解析 JSON 输出
	result := string(output)
	if strings.TrimSpace(result) == "" || strings.Contains(result, "null") {
		return []RestorePointInfo{}, nil
	}

	// 简化处理：返回基本信息
	return parseRestorePoints(result), nil
}

func parseRestorePoints(jsonOutput string) []RestorePointInfo {
	// 简化解析逻辑
	var points []RestorePointInfo

	// 按行扫描
	lines := strings.Split(jsonOutput, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Description") {
			desc := extractJSONValue(line)
			if i+1 < len(lines) {
				timeLine := strings.TrimSpace(lines[i+1])
				if strings.Contains(timeLine, "CreationTime") {
					ct := extractJSONValue(timeLine)
					points = append(points, RestorePointInfo{
						Name:        desc,
						Description: desc,
						CreatedAt:   ct,
					})
				}
			}
		}
	}

	return points
}

func extractJSONValue(line string) string {
	// 提取 "key": "value" 中的 value
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return ""
	}
	val := strings.TrimSpace(parts[1])
	val = strings.Trim(val, "\",")
	return val
}

func escapePowerShellString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
