package system

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// BatteryInfo 电池信息
type BatteryInfo struct {
	Name            string  `json:"name"`
	Status          string  `json:"status"`
	ChargeLevel     int     `json:"chargeLevel"`
	DesignCapacity  int     `json:"designCapacity"`
	FullCapacity    int     `json:"fullCapacity"`
	Voltage         float64 `json:"voltage"`
	Temperature     float64 `json:"temperature"`
	HealthPercent   float64 `json:"healthPercent"`
	IsPresent       bool    `json:"isPresent"`
}

// GetBatteryInfo 获取电池信息
func GetBatteryInfo() ([]BatteryInfo, error) {
	psScript := `Get-WmiObject -Class Win32_Battery | ForEach-Object {
		$name = $_.Name
		$status = $_.BatteryStatus
		$charge = $_.EstimatedChargeRemaining
		$design = $_.DesignCapacity
		$full = $_.FullChargeCapacity
		$voltage = $_.DesignVoltage
		"$name|$status|$charge|$design|$full|$voltage"
	}`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取电池信息失败: %w", err)
	}

	output := strings.TrimSpace(string(out))
	if output == "" {
		return nil, nil
	}

	var batteries []BatteryInfo
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "|") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 6 {
			continue
		}

		b := BatteryInfo{
			Name: strings.TrimSpace(parts[0]),
		}

		// BatteryStatus: 1=Discharging, 2=Charging, 3=Full
		statusCode, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		switch statusCode {
		case 1:
			b.Status = "放电中"
		case 2:
			b.Status = "充电中"
		case 3:
			b.Status = "已充满"
		default:
			b.Status = "未知"
		}

		b.ChargeLevel, _ = strconv.Atoi(strings.TrimSpace(parts[2]))
		b.DesignCapacity, _ = strconv.Atoi(strings.TrimSpace(parts[3]))
		b.FullCapacity, _ = strconv.Atoi(strings.TrimSpace(parts[4]))

		volt, _ := strconv.ParseFloat(strings.TrimSpace(parts[5]), 64)
		b.Voltage = volt / 1000

		if b.DesignCapacity > 0 && b.FullCapacity > 0 {
			b.HealthPercent = float64(b.FullCapacity) / float64(b.DesignCapacity) * 100
		} else {
			b.HealthPercent = 100
		}

		b.IsPresent = true
		batteries = append(batteries, b)
	}

	return batteries, nil
}

// HasBattery 检测是否有电池
func HasBattery() bool {
	batteries, err := GetBatteryInfo()
	return err == nil && len(batteries) > 0
}
