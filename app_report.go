package main

import (
	"time"

	"pc-toolbox/internal/optimize"
	"pc-toolbox/internal/process"
	"pc-toolbox/internal/system"
)

// GenerateSystemReport 生成系统报告（HTML 格式，可打印为 PDF）
func (a *App) GenerateSystemReport(reportType string) (string, error) {
	var data system.ExportReportData
	data.Title = "电脑工具箱 - 系统报告"
	data.GeneratedAt = time.Now().Format("2006-01-02 15:04:05")

	switch reportType {
	case "system":
		info, err := system.GetSystemInfo()
		if err != nil {
			return "", err
		}
		data.SystemInfo = info

	case "monitor":
		monitor, err := system.GetHardwareMonitor()
		if err != nil {
			return "", err
		}
		data.Monitor = monitor

	case "full":
		info, _ := system.GetSystemInfo()
		data.SystemInfo = info

		monitor, _ := system.GetHardwareMonitor()
		data.Monitor = monitor

		temps, _ := system.GetTemperatures()
		data.Temperatures = temps

		// 启动项
		startupItems := optimize.GetStartupItems()
		for _, item := range startupItems {
			data.StartupItems = append(data.StartupItems, system.ReportStartupItem{
				Name:    item.Name,
				Publisher: item.Publisher,
				Enabled:  item.Enabled,
				Impact:   item.Impact,
			})
		}

		// 进程 (Top 50)
		processes, _ := process.GetProcessList()
		count := len(processes)
		if count > 50 {
			count = 50
		}
		for i := 0; i < count; i++ {
			p := processes[i]
			data.Processes = append(data.Processes, system.ReportProcessInfo{
				PID:    p.PID,
				Name:   p.Name,
				CPU:    p.CPU,
				Memory: p.Memory,
				Status: p.Status,
			})
		}
	}

	return system.GenerateHTMLReport(data)
}
