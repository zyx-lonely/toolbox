package report

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
	"runtime"
	"time"
)

// SystemReport 系统报告
type SystemReport struct {
	GeneratedAt  string      `json:"generatedAt"`
	OS           OSInfo      `json:"os"`
	CPU          CPUInfo     `json:"cpu"`
	Memory       MemInfo     `json:"memory"`
	Disks        []DiskInfo  `json:"disks"`
	Network      NetworkInfo `json:"network"`
	ProcessCount int         `json:"processCount"`
}

type OSInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	BuildNumber  string `json:"buildNumber"`
	Architecture string `json:"architecture"`
	Uptime       string `json:"uptime"`
	Hostname     string `json:"hostname"`
	UserName     string `json:"userName"`
}

type CPUInfo struct {
	Name         string `json:"name"`
	Cores        int    `json:"cores"`
	LogicalCores int    `json:"logicalCores"`
	BaseClock    string `json:"baseClock"`
}

type MemInfo struct {
	Total     uint64  `json:"total"`
	Available uint64  `json:"available"`
	Usage     float64 `json:"usage"`
}

type DiskInfo struct {
	Label      string `json:"label"`
	FileSystem string `json:"fileSystem"`
	Total      uint64 `json:"total"`
	Free       uint64 `json:"free"`
	Usage      float64 `json:"usage"`
}

type NetworkInfo struct {
	Adapters []NetworkAdapter `json:"adapters"`
}

type NetworkAdapter struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	MAC  string `json:"mac"`
	Type string `json:"type"`
}

// CollectReport 收集系统报告（带崩溃恢复）
func CollectReport() (r *SystemReport) {
	defer func() {
		if rec := recover(); rec != nil {
			r = &SystemReport{
				GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
				OS: OSInfo{Name: "Windows", Hostname: "localhost", UserName: "user"},
				CPU: CPUInfo{Cores: 4, LogicalCores: 8},
			}
		}
	}()
	return collectReport()
}

func collectReport() *SystemReport {
	hostname, _ := os.Hostname()
	r := &SystemReport{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		OS: OSInfo{
			Architecture: runtime.GOARCH,
			Hostname:     hostname,
			UserName:     os.Getenv("USERNAME"),
		},
		CPU: CPUInfo{
			LogicalCores: runtime.NumCPU(),
		},
	}

	// 使用 wmic 获取系统信息（Windows 兼容性好）
	r.OS.Name = getWmicValue("os", "get", "Caption")
	r.OS.Version = getWmicValue("os", "get", "Version")
	r.OS.BuildNumber = getWmicValue("os", "get", "BuildNumber")
	r.OS.Uptime = getUptime()
	r.CPU.Name = getWmicValue("cpu", "get", "Name")
	r.CPU.BaseClock = getWmicValue("cpu", "get", "MaxClockSpeed") + " MHz"
	// 获取物理核心数
	if cores := getWmicValue("cpu", "get", "NumberOfCores"); cores != "" {
		fmt.Sscanf(cores, "%d", &r.CPU.Cores)
	}
	if r.CPU.Cores == 0 {
		r.CPU.Cores = r.CPU.LogicalCores / 2
		if r.CPU.Cores == 0 {
			r.CPU.Cores = r.CPU.LogicalCores
		}
	}
	r.Disks = getDiskInfo()
	r.Network = getNetworkInfo()

	// 内存
	total := getTotalMemory()
	avail := getAvailableMemory()
	r.Memory.Total = total
	r.Memory.Available = avail
	if total > 0 {
		r.Memory.Usage = float64(total-avail) / float64(total) * 100
	}

	return r
}

// GenerateHTMLReport 生成 HTML 报告
func GenerateHTMLReport(report *SystemReport) (path string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("生成报告时发生错误: %v", rec)
		}
	}()
	return generateHTMLReport(report)
}

func generateHTMLReport(report *SystemReport) (string, error) {
	if report == nil {
		report = CollectReport()
	}

	reportDir := filepath.Join(os.TempDir(), "pc-toolbox-reports")
	os.MkdirAll(reportDir, 0755)

	filename := fmt.Sprintf("system_report_%s.html", time.Now().Format("20060102_150405"))
	savePath := filepath.Join(reportDir, filename)

	html := buildHTML(report)
	if err := os.WriteFile(savePath, []byte(html), 0644); err != nil {
		return "", err
	}

	return savePath, nil
}

// OpenInBrowser 在浏览器中打开报告
func OpenInBrowser(path string) error {
	return exec.Command("cmd", "/c", "start", "", path).Start()
}

func buildHTML(r *SystemReport) string {
	diskRows := ""
	for _, d := range r.Disks {
		u := d.Usage
		if u > 100 { u = 100 }
		barColor := "#27ae60"
		if u > 75 { barColor = "#f39c12" }
		if u > 90 { barColor = "#e74c3c" }
		diskRows += fmt.Sprintf(
			"<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td><div class=\"bar\" style=\"width:%.0f%%%%;background:%s\"></div>%.0f%%%%</td></tr>",
			d.Label, d.FileSystem, formatBytes(d.Total), formatBytes(d.Free), u, barColor, u)
	}

	adapterRows := ""
	for _, a := range r.Network.Adapters {
		adapterRows += fmt.Sprintf(
			"<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			a.Name, a.Type, a.IP, a.MAC)
	}

	memPct := r.Memory.Usage
	if memPct > 100 { memPct = 100 }
	memPctStr := fmt.Sprintf("%.0f", memPct)
	memBarColor := "#27ae60"
	if memPct > 75 { memBarColor = "#f39c12" }
	if memPct > 90 { memBarColor = "#e74c3c" }

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head><meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>系统报告 - 电脑工具箱</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:#f0f2f5;color:#333;margin:0;padding:20px}
.container{max-width:900px;margin:0 auto}
.header{background:linear-gradient(135deg,#18a058,#36ad6a);color:#fff;padding:30px;border-radius:12px;margin-bottom:20px}
.header h1{margin:0 0 8px 0;font-size:24px}
.header p{margin:0;opacity:.85;font-size:14px}
.card{background:#fff;border-radius:12px;padding:20px;margin-bottom:16px;box-shadow:0 2px 8px rgba(0,0,0,.08)}
.card h2{font-size:16px;color:#18a058;margin:0 0 16px 0;padding-bottom:8px;border-bottom:2px solid #e8f8e8}
table{width:100%%;border-collapse:collapse}
th,td{text-align:left;padding:10px 8px;border-bottom:1px solid #eee;font-size:14px}
th{color:#888;font-weight:500;width:120px}
.bar{height:20px;border-radius:4px;min-width:4px;display:inline-block;vertical-align:middle;margin-right:6px}
.stat-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:12px}
.stat-card{background:#f8f9fa;border-radius:8px;padding:16px;text-align:center}
.stat-card .val{font-size:24px;font-weight:700}
.stat-card .lbl{font-size:12px;color:#888;margin-top:4px}
.footer{text-align:center;padding:20px;color:#aaa;font-size:12px}
</style></head>
<body>
<div class="container">
<div class="header">
<h1>🖥️ 系统信息报告</h1>
<p>生成时间: %s | 主机名: %s</p>
</div>
<div class="card">
<h2>操作系统</h2>
<table>
<tr><th>名称</th><td>%s</td></tr>
<tr><th>版本</th><td>%s</td></tr>
<tr><th>构建版本</th><td>%s</td></tr>
<tr><th>架构</th><td>%s</td></tr>
<tr><th>已运行</th><td>%s</td></tr>
<tr><th>用户</th><td>%s</td></tr>
</table>
</div>
<div class="card">
<h2>处理器</h2>
<table>
<tr><th>型号</th><td>%s</td></tr>
<tr><th>主频</th><td>%s</td></tr>
<tr><th>物理核心</th><td>%d</td></tr>
<tr><th>逻辑核心</th><td>%d</td></tr>
</table>
</div>
<div class="card">
<h2>内存</h2>
<div class="stat-grid">
<div class="stat-card"><div class="val" style="color:#18a058">%s</div><div class="lbl">总计</div></div>
<div class="stat-card"><div class="val" style="color:#27ae60">%s</div><div class="lbl">可用</div></div>
<div class="stat-card"><div class="val" style="color:%s">%s%%%%</div><div class="lbl">使用率</div></div>
</div>
<div style="margin-top:12px"><div class="bar" style="width:%s%%%%;height:24px;background:%s"></div></div>
</div>
<div class="card">
<h2>磁盘</h2>
<table>
<tr><th>卷标</th><th>文件系统</th><th>总容量</th><th>可用</th><th>使用率</th></tr>
%s
</table>
</div>
<div class="card">
<h2>网络适配器</h2>
<table>
<tr><th>名称</th><th>类型</th><th>IP</th><th>MAC</th></tr>
%s
</table>
</div>
<div class="footer"><p>由 电脑工具箱 生成 · %s</p></div>
</div>
</body></html>`,
		r.GeneratedAt, r.OS.Hostname,
		r.OS.Name, r.OS.Version, r.OS.BuildNumber, r.OS.Architecture, r.OS.Uptime, r.OS.UserName,
		r.CPU.Name, r.CPU.BaseClock, r.CPU.Cores, r.CPU.LogicalCores,
		formatBytes(r.Memory.Total), formatBytes(r.Memory.Available), memBarColor, memPctStr, memPctStr, memBarColor,
		diskRows, adapterRows, r.GeneratedAt)
}

// wmic 封装
func getWmicValue(class, action, field string) string {
	c := exec.Command("wmic", class, action, field)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.Output()
	if err != nil {
		return ""
	}
	lines := splitLines(string(out))
	for _, line := range lines {
		line = trimSpace(line)
		if line != "" && !containsField(line, field) {
			return line
		}
	}
	return ""
}

func containsField(s, field string) bool {
	clean := stringsReplaceAll(s, " ", "")
	cleanF := stringsReplaceAll(field, " ", "")
	return clean == cleanF
}

func getUptime() string {
	c := exec.Command("powershell", "-NoProfile",
		"(New-TimeSpan -Start (Get-CimInstance Win32_OperatingSystem).LastBootUpTime -End (Get-Date)).ToString()")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.Output()
	if err != nil {
		return "未知"
	}
	return trimSpace(string(out))
}

func getTotalMemory() uint64 {
	cmd := exec.Command("wmic", "ComputerSystem", "get", "TotalPhysicalMemory"); cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}; out, _ := cmd.Output()
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	lines := splitLines(string(out))
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || line == "TotalPhysicalMemory" {
			continue
		}
		var total uint64
		if _, err := fmt.Sscanf(line, "%d", &total); err == nil && total > 0 {
			return total
		}
	}
	return 0
}

func getAvailableMemory() uint64 {
	cmd := exec.Command("wmic", "OS", "get", "FreePhysicalMemory"); cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}; out, _ := cmd.Output()
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	lines := splitLines(string(out))
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || line == "FreePhysicalMemory" {
			continue
		}
		var kb uint64
		if _, err := fmt.Sscanf(line, "%d", &kb); err == nil && kb > 0 {
			return kb * 1024
		}
	}
	return 0
}

func getDiskInfo() []DiskInfo {
	cmd := exec.Command("wmic", "logicaldisk", "get", "DeviceID,FileSystem,Size,FreeSpace", "/format:csv"); cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}; out, _ := cmd.Output()
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	var disks []DiskInfo
	lines := splitLines(string(out))
	for _, line := range lines[1:] { // skip header
		line = trimSpace(line)
		if line == "" {
			continue
		}
		parts := stringsSplit(line, ",")
		if len(parts) < 5 {
			continue
		}
		label := trimSpace(parts[1])    // DeviceID (C:)
		fs := trimSpace(parts[2])       // FileSystem (NTFS)
		totalStr := trimSpace(parts[3]) // Size
		freeStr := trimSpace(parts[4])  // FreeSpace
		if label == "" {
			continue
		}
		var total, free uint64
		fmt.Sscanf(totalStr, "%d", &total)
		fmt.Sscanf(freeStr, "%d", &free)
		var usage float64
		if total > 0 {
			usage = float64(total-free) / float64(total) * 100
		}
		disks = append(disks, DiskInfo{
			Label: label, FileSystem: fs, Total: total, Free: free, Usage: usage,
		})
	}
	return disks
}

func getNetworkInfo() NetworkInfo {
	cmd := exec.Command("wmic", "nic", "where", "NetEnabled=TRUE", "get", "Name,MACAddress", "/format:csv"); cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}; out, _ := cmd.Output()
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	var adapters []NetworkAdapter
	lines := splitLines(string(out))
	for _, line := range lines[1:] {
		line = trimSpace(line)
		if line == "" {
			continue
		}
		parts := stringsSplit(line, ",")
		if len(parts) < 3 {
			continue
		}
		name := trimSpace(parts[1])
		mac := trimSpace(parts[2])
		if name != "" {
			adapters = append(adapters, NetworkAdapter{
				Name: name, Type: "以太网", IP: "-", MAC: mac,
			})
		}
	}
	return NetworkInfo{Adapters: adapters}
}

// 避免 import "strings" 的工具函数
func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func splitLines(s string) []string {
	var lines []string
	cur := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, cur)
			cur = ""
		} else {
			cur += string(s[i])
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func stringsSplit(s, sep string) []string {
	var parts []string
	cur := ""
	for i := 0; i < len(s); {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			parts = append(parts, cur)
			cur = ""
			i += len(sep)
		} else {
			cur += string(s[i])
			i++
		}
	}
	parts = append(parts, cur)
	return parts
}

func stringsReplaceAll(s, old, new string) string {
	var result string
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}

func formatBytes(bytes uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	val := float64(bytes)
	i := 0
	for val >= 1024 && i < len(units)-1 {
		val /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%.0f %s", val, units[i])
	}
	return fmt.Sprintf("%.2f %s", val, units[i])
}
