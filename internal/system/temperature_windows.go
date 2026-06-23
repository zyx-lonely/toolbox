package system

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"pc-toolbox/internal/common"
)

// TemperatureInfo 温度信息
type TemperatureInfo struct {
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
	Unit        string  `json:"unit"`
}

// GetTemperatures 获取硬件温度（Windows）
// 尝试多种方法获取温度信息
func GetTemperatures() ([]TemperatureInfo, error) {
	var temps []TemperatureInfo

	// 方法1：使用 PowerShell 查询 WMI
	psTemps, err := getTemperaturesViaPowerShell()
	if err == nil && len(psTemps) > 0 {
		return psTemps, nil
	}

	// 方法2：使用 wmic 命令（已弃用但可能仍可用）
	cmdTemps, err := getTemperaturesViaWmic()
	if err == nil && len(cmdTemps) > 0 {
		return cmdTemps, nil
	}

	// 如果都失败，返回空数组和提示
	if len(temps) == 0 {
		return temps, fmt.Errorf("无法读取温度信息，可能需要管理员权限或安装第三方工具（如 Open Hardware Monitor）")
	}

	return temps, nil
}

// getTemperaturesViaPowerShell 使用 PowerShell 查询温度
func getTemperaturesViaPowerShell() ([]TemperatureInfo, error) {
	var temps []TemperatureInfo

	// PowerShell 命令：查询 MSAcpi_ThermalZoneTemperature 类
	psCommand := `Get-WmiObject -Namespace "root/wmi" -Class MSAcpi_ThermalZoneTemperature | ForEach-Object { $temp = [math]::Round(($_.CurrentTemperature - 2732) / 10.0, 1); Write-Output "$($_.InstanceName):$temp" }`

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psCommand)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return temps, err
	}

	lines := strings.Split(common.GbkToUtf8(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		tempStr := strings.TrimSpace(parts[1])
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			continue
		}
		temps = append(temps, TemperatureInfo{
			Name:        name,
			Temperature: temp,
			Unit:        "°C",
		})
	}

	return temps, nil
}

// getTemperaturesViaWmic 使用 wmic 命令查询温度
func getTemperaturesViaWmic() ([]TemperatureInfo, error) {
	var temps []TemperatureInfo

	// wmic 命令：查询温度
	cmd := exec.Command("wmic", "/namespace:\\\\root\\wmi", "PATH", "MSAcpi_ThermalZoneTemperature", "GET", "CurrentTemperature,InstanceName", "/FORMAT:CSV")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return temps, err
	}

	lines := strings.Split(common.GbkToUtf8(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Node") {
			continue
		}
		// CSV 格式：Node,InstanceName,CurrentTemperature
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue
		}
		name := strings.Trim(parts[1], "\"")
		tempStr := strings.Trim(parts[2], "\"")
		tempVal, err := strconv.ParseInt(tempStr, 10, 64)
		if err != nil {
			continue
		}
		// WMI 返回的温度是开尔文温度 * 10
		temp := float64(tempVal-2732) / 10.0
		temps = append(temps, TemperatureInfo{
			Name:        name,
			Temperature: temp,
			Unit:        "°C",
		})
	}

	return temps, nil
}
