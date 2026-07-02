package optimize

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// RestorePointInfo 系统还原点信息
type RestorePointInfo struct {
	Name        string `json:"name"`
	CreatedAt   string `json:"createdAt"`
	Description string `json:"description"`
	Sequence    int    `json:"sequence"`
}

// wuRestorePoint WMI 还原点
type wuRestorePoint struct {
	Description    string `json:"Description"`
	CreationTime   string `json:"CreationTime"`
	SequenceNumber int    `json:"SequenceNumber"`
}

// CreateRestorePoint 创建系统还原点
func CreateRestorePoint(description string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		"Checkpoint-Computer", "-Description", description, "-RestorePointType", "MODIFY_SETTINGS")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建还原点失败: %w, 输出: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

// GetRestorePoints 获取系统还原点列表
func GetRestorePoints() ([]RestorePointInfo, error) {
	psScript := `Get-ComputerRestorePoint | Select-Object Description, CreationTime, SequenceNumber | ConvertTo-Json`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取还原点失败: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" || result == "null" {
		return []RestorePointInfo{}, nil
	}

	return parseRestorePoints(result)
}

func parseRestorePoints(jsonOutput string) ([]RestorePointInfo, error) {
	var points []RestorePointInfo

	// 处理单条和数组两种情况
	if !strings.HasPrefix(jsonOutput, "[") {
		jsonOutput = "[" + jsonOutput + "]"
	}

	var items []wuRestorePoint
	if err := json.Unmarshal([]byte(jsonOutput), &items); err != nil {
		return nil, fmt.Errorf("解析还原点 JSON 失败: %w", err)
	}

	for _, item := range items {
		points = append(points, RestorePointInfo{
			Name:        item.Description,
			Description: item.Description,
			CreatedAt:   item.CreationTime,
			Sequence:    item.SequenceNumber,
		})
	}

	return points, nil
}


