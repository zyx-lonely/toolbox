package optimize

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"pc-toolbox/internal/common"
)

// HealthReport 健康体检报告
type HealthReport struct {
	GeneratedAt  string        `json:"generatedAt"`
	Score        int           `json:"score"`        // 0-100
	Items        []HealthItem  `json:"items"`
	Suggestions  []string      `json:"suggestions"`
}

// HealthItem 单项检查结果
type HealthItem struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "good", "warning", "bad"
	Value   string `json:"value"`
	Message string `json:"message"`
}

// RunHealthCheck 执行系统健康体检（带崩溃恢复）
func RunHealthCheck() (r *HealthReport) {
	defer func() {
		if rec := recover(); rec != nil {
			r = &HealthReport{
				GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
				Score:       0,
				Items: []HealthItem{{
					Name: "系统体检", Status: "bad", Value: "崩溃",
					Message: fmt.Sprintf("体检过程发生错误: %v", rec),
				}},
				Suggestions: []string{"请以管理员身份运行本程序后重试"},
			}
		}
	}()
	
	r = doHealthCheck()
	return r
}

func doHealthCheck() *HealthReport {
	report := &HealthReport{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	checks := []struct {
		name string
		fn   func() HealthItem
	}{
		{"磁盘空间", checkDiskSpace},
		{"内存使用", checkMemoryUsage},
		{"CPU 负载", checkCPULoad},
		{"启动项数量", checkStartupItems},
		{"系统运行时间", checkUptime},
		{"临时文件", checkTempFiles},
		{"Windows 更新", checkWindowsUpdate},
	}

	score := 100
	for _, c := range checks {
		item := c.fn()
		report.Items = append(report.Items, item)
		switch item.Status {
		case "bad":
			score -= 20
			report.Suggestions = append(report.Suggestions, item.Message)
		case "warning":
			score -= 10
			report.Suggestions = append(report.Suggestions, item.Message)
		}
	}

	if score < 0 {
		score = 0
	}
	report.Score = score
	return report
}

func checkDiskSpace() HealthItem {
	// 检查 C 盘可用空间
	free := getDiskFreeSpace("C")
	if free > 50*1024*1024*1024 { // > 50GB
		return HealthItem{Name: "C盘可用空间", Status: "good", Value: formatBytes(free), Message: "磁盘空间充足"}
	} else if free > 10*1024*1024*1024 {
		return HealthItem{Name: "C盘可用空间", Status: "warning", Value: formatBytes(free), Message: "建议清理 C 盘空间，可用不足 50GB"}
	}
	return HealthItem{Name: "C盘可用空间", Status: "bad", Value: formatBytes(free), Message: "C 盘空间严重不足，请立即清理！"}
}

func checkMemoryUsage() HealthItem {
	total := getTotalMem()
	avail := getAvailMem()
	if total == 0 {
		return HealthItem{Name: "内存使用", Status: "good", Value: "未知", Message: "无法检测"}
	}
	used := total - avail
	pct := float64(used) / float64(total) * 100
	if pct < 60 {
		return HealthItem{Name: "内存使用", Status: "good", Value: fmt.Sprintf("%.0f%%", pct), Message: "内存使用正常"}
	} else if pct < 85 {
		return HealthItem{Name: "内存使用", Status: "warning", Value: fmt.Sprintf("%.0f%%", pct), Message: "内存使用偏高，建议关闭未使用的程序"}
	}
	return HealthItem{Name: "内存使用", Status: "bad", Value: fmt.Sprintf("%.0f%%", pct), Message: "内存严重不足！建议增加物理内存"}
}

func checkCPULoad() HealthItem {
	// 使用 Windows API 获取 CPU 使用率
	total1, idle1 := getCPUTimes()
	time.Sleep(500 * time.Millisecond)
	total2, idle2 := getCPUTimes()

	deltaTotal := total2 - total1
	deltaIdle := idle2 - idle1

	var usage float64
	if deltaTotal > 0 {
		usage = (1.0 - float64(deltaIdle)/float64(deltaTotal)) * 100
	}
	if usage < 0 {
		usage = 0
	}
	if usage < 50 {
		return HealthItem{Name: "CPU 负载", Status: "good", Value: fmt.Sprintf("%.0f%%", usage), Message: "CPU 使用正常"}
	} else if usage < 80 {
		return HealthItem{Name: "CPU 负载", Status: "warning", Value: fmt.Sprintf("%.0f%%", usage), Message: "CPU 负载较高"}
	}
	return HealthItem{Name: "CPU 负载", Status: "bad", Value: fmt.Sprintf("%.0f%%", usage), Message: "CPU 负载过高"}
}

func checkStartupItems() HealthItem {
	// 简化检查：读取注册表启动项
	count := countStartupRegistry()
	if count <= 5 {
		return HealthItem{Name: "启动项数量", Status: "good", Value: fmt.Sprintf("%d 项", count), Message: "启动项数量合理"}
	} else if count <= 15 {
		return HealthItem{Name: "启动项数量", Status: "warning", Value: fmt.Sprintf("%d 项", count), Message: "启动项较多，建议打开「启动项管理」禁用不必要的项"}
	}
	return HealthItem{Name: "启动项数量", Status: "bad", Value: fmt.Sprintf("%d 项", count), Message: "启动项过多！请立即清理"}
}

func checkUptime() HealthItem {
	days := getSystemUptimeDays()
	if days < 7 {
		return HealthItem{Name: "系统运行时间", Status: "good", Value: fmt.Sprintf("%d 天", days), Message: "系统已定期重启"}
	} else if days < 30 {
		return HealthItem{Name: "系统运行时间", Status: "warning", Value: fmt.Sprintf("%d 天", days), Message: "建议重启系统以释放缓存"}
	}
	return HealthItem{Name: "系统运行时间", Status: "bad", Value: fmt.Sprintf("%d 天", days), Message: "系统长期未重启，建议立即重启"}
}

func checkTempFiles() HealthItem {
	tempDir := os.Getenv("TEMP")
	total := countTempSize(tempDir)
	if total < 500*1024*1024 {
		return HealthItem{Name: "临时文件", Status: "good", Value: formatBytes(total), Message: "临时文件较少"}
	} else if total < 2*1024*1024*1024 {
		return HealthItem{Name: "临时文件", Status: "warning", Value: formatBytes(total), Message: "临时文件较多，建议打开「磁盘清理」清理"}
	}
	return HealthItem{Name: "临时文件", Status: "bad", Value: formatBytes(total), Message: "临时文件过多！请立即清理"}
}

func checkWindowsUpdate() HealthItem {
	systemRoot := os.Getenv("SystemRoot")
	powershellPath := filepath.Join(systemRoot, "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
	c1 := &exec.Cmd{
		Path: powershellPath,
		Args: []string{powershellPath, "-NoProfile", "(Get-WUApiVersion).ToString()"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out, _ := c1.Output()
	if len(out) > 0 {
		return HealthItem{Name: "Windows 更新", Status: "good", Value: "正常", Message: "Windows 更新服务运行正常"}
	}
	// 检查更新服务状态
	scPath := filepath.Join(systemRoot, "System32", "sc.exe")
	c2 := &exec.Cmd{
		Path: scPath,
		Args: []string{scPath, "query", "wuauserv"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out, _ = c2.Output()
	if strings.Contains(common.GbkToUtf8(string(out)), "RUNNING") {
		return HealthItem{Name: "Windows 更新", Status: "good", Value: "运行中", Message: "更新服务运行正常"}
	}
	return HealthItem{Name: "Windows 更新", Status: "warning", Value: "未运行", Message: "Windows Update 服务未运行，可能无法接收更新"}
}

// -- 辅助函数 --

func getDiskFreeSpace(drive string) uint64 {
	defer func() { recover() }()
	dll := windows.NewLazySystemDLL("kernel32.dll")
	proc := dll.NewProc("GetDiskFreeSpaceExW")
	if proc == nil {
		return 0
	}
	var free uint64
	drivePath, err := windows.UTF16PtrFromString(drive + ":\\")
	if err != nil {
		return 0
	}
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(drivePath)), 0, 0, uintptr(unsafe.Pointer(&free)))
	if ret == 0 {
		return 0
	}
	return free
}

func getTotalMem() uint64 {
	defer func() { recover() }()
	dll := windows.NewLazySystemDLL("kernel32.dll")
	proc := dll.NewProc("GlobalMemoryStatusEx")
	if proc == nil {
		return 16 * 1024 * 1024 * 1024
	}
	var ms memoryStatusEx
	ms.length = uint32(unsafe.Sizeof(ms))
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&ms)))
	if ret == 0 {
		return 16 * 1024 * 1024 * 1024
	}
	return ms.TotalPhys
}

func getAvailMem() uint64 {
	defer func() { recover() }()
	dll := windows.NewLazySystemDLL("kernel32.dll")
	proc := dll.NewProc("GlobalMemoryStatusEx")
	if proc == nil {
		return 8 * 1024 * 1024 * 1024
	}
	var ms memoryStatusEx
	ms.length = uint32(unsafe.Sizeof(ms))
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&ms)))
	if ret == 0 {
		return 8 * 1024 * 1024 * 1024
	}
	return ms.AvailPhys
}

type memoryStatusEx struct {
	length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

func getCPUTimes() (total, idle uint64) {
	dll := windows.NewLazySystemDLL("kernel32.dll")
	proc := dll.NewProc("GetSystemTimes")
	var idleTime, kernelTime, userTime windows.Filetime
	ret, _, _ := proc.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)))
	if ret == 0 {
		return 0, 0
	}
	// KernelTime 包含 IdleTime
	idle = uint64(idleTime.Nanoseconds())
	total = uint64(kernelTime.Nanoseconds() + userTime.Nanoseconds())
	return total, idle
}

func countStartupRegistry() int {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`, registry.READ)
	if err != nil {
		return 0
	}
	defer k.Close()

	keys, err := k.ReadValueNames(100)
	if err != nil {
		return 0
	}
	return len(keys)
}

func getSystemUptimeDays() int {
	dll := windows.NewLazySystemDLL("kernel32.dll")
	proc := dll.NewProc("GetTickCount64")
	tick, _, _ := proc.Call()
	return int(tick / uintptr(86400000)) // ms to days
}

func countTempSize(dir string) uint64 {
	total, _ := getDirSize(dir)
	return total
}

func getDirSize(path string) (uint64, error) {
	var total uint64
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			total += uint64(info.Size())
		}
		return nil
	})
	return total, err
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
