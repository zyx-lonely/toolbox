package system

import (
	"bytes"
	"fmt"
)

// ReportStartupItem 报告中的启动项数据（避免跨包依赖）
type ReportStartupItem struct {
	Name    string `json:"name"`
	Publisher string `json:"publisher"`
	Enabled  bool   `json:"enabled"`
	Impact   string `json:"impact"`
}

// ReportProcessInfo 报告中的进程数据（避免跨包依赖）
type ReportProcessInfo struct {
	PID      int    `json:"pid"`
	Name     string `json:"name"`
	CPU      string `json:"cpu"`      // 保持字符串格式，如 "5%"
	Memory   string `json:"memory"`  // 保持字符串格式，如 "12,340 K"
	Status   string `json:"status"`
}

// ExportReportData 导出报告所需的数据
type ExportReportData struct {
	Title         string                `json:"title"`
	GeneratedAt   string                `json:"generated_at"`
	SystemInfo    *SystemInfo          `json:"system_info,omitempty"`
	Monitor       *HardwareMonitor    `json:"monitor,omitempty"`
	Temperatures  []TemperatureInfo   `json:"temperatures,omitempty"`
	StartupItems  []ReportStartupItem `json:"startup_items,omitempty"`
	Processes     []ReportProcessInfo  `json:"processes,omitempty"`
}

// GenerateHTMLReport 生成 HTML 格式报告（可打印为 PDF）
func GenerateHTMLReport(data ExportReportData) (string, error) {
	var buf bytes.Buffer

	buf.WriteString(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>`)
	buf.WriteString(data.Title)
	buf.WriteString(`</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: 'Microsoft YaHei', 'Segoe UI', sans-serif; font-size: 14px; color: #333; background: #fff; }
  .container { max-width: 900px; margin: 0 auto; padding: 30px; }
  h1 { font-size: 28px; color: #18a058; border-bottom: 3px solid #18a058; padding-bottom: 12px; margin-bottom: 24px; }
  h2 { font-size: 20px; color: #333; border-left: 4px solid #18a058; padding-left: 12px; margin: 24px 0 16px; }
  .meta { color: #888; font-size: 13px; margin-bottom: 20px; }
  table { width: 100%; border-collapse: collapse; margin: 12px 0; }
  th, td { padding: 10px 14px; border: 1px solid #e8e8e8; text-align: left; }
  th { background: #f5f5f5; font-weight: 600; }
  tr:nth-child(even) td { background: #fafafa; }
  .badge { display: inline-block; padding: 2px 10px; border-radius: 10px; font-size: 12px; font-weight: 500; }
  .badge-green { background: #e8f5e9; color: #2e7d32; }
  .badge-red { background: #ffebee; color: #c62828; }
  .badge-yellow { background: #fff8e1; color: #f57f17; }
  .progress-bar { height: 12px; background: #eee; border-radius: 6px; overflow: hidden; width: 120px; display: inline-block; vertical-align: middle; }
  .progress-fill { height: 100%; border-radius: 6px; }
  @media print {
    body { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
  }
</style>
</head>
<body>
<div class="container">
`)
	// 标题
	buf.WriteString(fmt.Sprintf("<h1>%s</h1>\n", data.Title))
	buf.WriteString(fmt.Sprintf(`<div class="meta">生成时间：%s | 由 电脑工具箱 生成</div>`+"\n", data.GeneratedAt))

	// 系统信息
	if data.SystemInfo != nil {
		buf.WriteString("<h2>🖥️ 系统信息</h2>\n")
		buf.WriteString(`<table><tr><th>项目</th><th>值</th></tr>` + "\n")
		si := data.SystemInfo
		buf.WriteString(fmt.Sprintf("<tr><td>操作系统</td><td>%s %s (Build %s)</td></tr>\n", si.OS.Name, si.OS.Version, si.OS.BuildNumber))
		buf.WriteString(fmt.Sprintf("<tr><td>架构</td><td>%s</td></tr>\n", si.OS.Architecture))
		buf.WriteString(fmt.Sprintf("<tr><td>安装日期</td><td>%s</td></tr>\n", si.OS.InstallDate))
		buf.WriteString(fmt.Sprintf("<tr><td>运行时长</td><td>%s</td></tr>\n", si.OS.Uptime))
		buf.WriteString(fmt.Sprintf("<tr><td>CPU</td><td>%s (×%d 核心，×%d 逻辑核心)</td></tr>\n", si.CPU.Name, si.CPU.Cores, si.CPU.LogicalCores))
		buf.WriteString(fmt.Sprintf("<tr><td>CPU 主频</td><td>%s</td></tr>\n", si.CPU.BaseClock))
		memPercent := si.Memory.Usage
		buf.WriteString(fmt.Sprintf("<tr><td>内存</td><td>%s / %s (使用率 %.1f%%)</td></tr>\n", formatBytes(si.Memory.Used), formatBytes(si.Memory.Total), memPercent))
		
		// 磁盘信息
		for _, d := range si.Disks {
			buf.WriteString(fmt.Sprintf("<tr><td>磁盘 %s</td><td>%s / %s (使用率 %.1f%%)</td></tr>\n", d.Label, formatBytes(d.Used), formatBytes(d.Total), d.Usage))
		}
		buf.WriteString("</table>\n")
	}

	// 硬件监控
	if data.Monitor != nil {
		buf.WriteString("<h2>📊 硬件监控</h2>\n")
		m := data.Monitor
		buf.WriteString(`<table><tr><th>指标</th><th>数值</th><th>使用情况</th></tr>` + "\n")
		
		cpuColor := "#4caf50"
		if m.CPUUsage > 80 { cpuColor = "#f44336" }
		buf.WriteString(fmt.Sprintf(`<tr><td>CPU 使用率</td><td>%.1f%%</td><td><div class="progress-bar"><div class="progress-fill" style="width:%.1f%%;background:%s"></div></div></td></tr>`+"\n", m.CPUUsage, m.CPUUsage, cpuColor))

		memColor := "#4caf50"
		if m.MemoryUsage > 80 { memColor = "#f44336" }
		buf.WriteString(fmt.Sprintf(`<tr><td>内存使用率</td><td>%.1f%% (%s / %s)</td><td><div class="progress-bar"><div class="progress-fill" style="width:%.1f%%;background:%s"></div></div></td></tr>`+"\n", m.MemoryUsage, formatBytes(m.MemoryUsed), formatBytes(m.MemoryTotal), m.MemoryUsage, memColor))

		buf.WriteString(fmt.Sprintf(`<tr><td>网络上传</td><td>%s</td><td></td></tr>`+"\n", formatBytes(m.NetUp)))
		buf.WriteString(fmt.Sprintf(`<tr><td>网络下载</td><td>%s</td><td></td></tr>`+"\n", formatBytes(m.NetDown)))
		buf.WriteString(fmt.Sprintf(`<tr><td>运行时长</td><td>%s</td><td></td></tr>`+"\n", m.Uptime))
		buf.WriteString("</table>\n")
	}

	// 温度
	if len(data.Temperatures) > 0 {
		buf.WriteString("<h2>🌡️ 硬件温度</h2>\n")
		buf.WriteString(`<table><tr><th>传感器</th><th>温度</th><th>状态</th></tr>` + "\n")
		for _, t := range data.Temperatures {
			status := "正常"
			cls := "badge-green"
			if t.Temperature > 80 {
				status = "过高"
				cls = "badge-red"
			} else if t.Temperature > 60 {
				status = "偏高"
				cls = "badge-yellow"
			}
			buf.WriteString(fmt.Sprintf(`<tr><td>%s</td><td>%.1f °C</td><td><span class="badge %s">%s</span></td></tr>`+"\n", t.Name, t.Temperature, cls, status))
		}
		buf.WriteString("</table>\n")
	}

	// 启动项
	if len(data.StartupItems) > 0 {
		buf.WriteString("<h2>🚀 启动项</h2>\n")
		buf.WriteString(`<table><tr><th>名称</th><th>发布者</th><th>状态</th><th>影响</th></tr>` + "\n")
		for _, item := range data.StartupItems {
			statusCls := "badge-green"
			if !item.Enabled { statusCls = "badge-red" }
			impactCls := "badge-green"
			if item.Impact == "high" { impactCls = "badge-red" }
			buf.WriteString(fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td><span class="badge %s">%s</span></td><td><span class="badge %s">%s</span></td></tr>`+"\n", item.Name, item.Publisher, statusCls, boolToStatus(item.Enabled), impactCls, item.Impact))
		}
		buf.WriteString("</table>\n")
	}

	// 进程
	if len(data.Processes) > 0 {
		buf.WriteString("<h2>⚙️ 进程列表 (Top 50)</h2>\n")
		buf.WriteString(`<table><tr><th>PID</th><th>名称</th><th>CPU %%</th><th>内存</th><th>状态</th></tr>` + "\n")
		count := len(data.Processes)
		if count > 50 { count = 50 }
		for i := 0; i < count; i++ {
			p := data.Processes[i]
			buf.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n", p.PID, p.Name, p.CPU, p.Memory, p.Status))
		}
		buf.WriteString("</table>\n")
	}

	buf.WriteString(`</div></body></html>`)

	return buf.String(), nil
}

func boolToStatus(b bool) string {
	if b { return "启用" }
	return "禁用"
}

func formatBytes(b uint64) string {
	if b < 1024 { return fmt.Sprintf("%d B", b) }
	if b < 1024*1024 { return fmt.Sprintf("%.1f KB", float64(b)/1024) }
	if b < 1024*1024*1024 { return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024)) }
	return fmt.Sprintf("%.1f GB", float64(b)/(1024*1024*1024))
}
