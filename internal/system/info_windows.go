//go:build windows

package system

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procGetNativeSystemInfo    = modkernel32.NewProc("GetNativeSystemInfo")
	procGlobalMemoryStatusEx   = modkernel32.NewProc("GlobalMemoryStatusEx")
	procGetSystemTimes         = modkernel32.NewProc("GetSystemTimes")
	procGetTickCount64         = modkernel32.NewProc("GetTickCount64")
	procGetDiskFreeSpaceExW    = modkernel32.NewProc("GetDiskFreeSpaceExW")
	procGetLogicalDrives       = modkernel32.NewProc("GetLogicalDrives")
	procGetDriveTypeW          = modkernel32.NewProc("GetDriveTypeW")
	procGetVolumeInformationW  = modkernel32.NewProc("GetVolumeInformationW")
)

type systemInfo struct {
	wProcessorArchitecture      uint16
	wReserved                   uint16
	dwPageSize                  uint32
	lpMinimumApplicationAddress uintptr
	lpMaximumApplicationAddress uintptr
	dwActiveProcessorMask       uintptr
	dwNumberOfProcessors        uint32
	dwProcessorType             uint32
	dwAllocationGranularity     uint32
	wProcessorLevel             uint16
	wProcessorRevision          uint16
}

type memoryStatusEx struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var (
	// 用于 CPU 使用率计算的先前值
	prevIdleTime   uint64
	prevKernelTime uint64
	prevUserTime   uint64
	firstCPUSample = true
)

// GetSystemInfo 获取完整系统信息
func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{}

	info.OS = getOSInfo()
	info.CPU = getCPUInfo()
	info.Memory = getMemoryInfo()
	info.Disks = getDiskInfo()
	info.Network = getNetworkInfo()

	return info, nil
}

// GetHardwareMonitor 获取实时硬件监控数据
func GetHardwareMonitor() (*HardwareMonitor, error) {
	mem := getMemoryInfo()
	cpu := getCPUUsage()
	uptime := getUptime()

	return &HardwareMonitor{
		CPUUsage:    cpu,
		MemoryUsage: mem.Usage,
		MemoryUsed:  mem.Used,
		MemoryTotal: mem.Total,
		Uptime:      uptime,
	}, nil
}

func getOSInfo() OSInfo {
	osInfo := OSInfo{
		Architecture: "Unknown",
	}

	// 检测 Windows 版本
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.READ)
	if err == nil {
		defer k.Close()

		osInfo.Name, _, _ = k.GetStringValue("ProductName")
		osInfo.Version, _, _ = k.GetStringValue("DisplayVersion")
		osInfo.BuildNumber, _, _ = k.GetStringValue("CurrentBuild")
		installDateStr, _, _ := k.GetStringValue("InstallDate")
		if installDateStr != "" {
			osInfo.InstallDate = installDateStr
		}
	}

	// 检测架构
	var si systemInfo
	procGetNativeSystemInfo.Call(uintptr(unsafe.Pointer(&si)))
	switch si.wProcessorArchitecture {
	case 0:
		osInfo.Architecture = "x86"
	case 9:
		osInfo.Architecture = "x64"
	case 12:
		osInfo.Architecture = "ARM64"
	case 5:
		osInfo.Architecture = "ARM"
	}

	osInfo.Uptime = getUptime()

	return osInfo
}

func getCPUInfo() CPUInfo {
	cpu := CPUInfo{
		Cores:        0,
		LogicalCores: 0,
	}

	// 获取逻辑核心数
	var si systemInfo
	procGetNativeSystemInfo.Call(uintptr(unsafe.Pointer(&si)))
	cpu.LogicalCores = int(si.dwNumberOfProcessors)

	// 获取 CPU 名称（通过注册表）
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`HARDWARE\DESCRIPTION\System\CentralProcessor\0`, registry.READ)
	if err == nil {
		defer k.Close()
		name, _, _ := k.GetStringValue("ProcessorNameString")
		cpu.Name = name
		cpu.Cores = int(si.dwNumberOfProcessors) // 简化处理
	}

	cpu.Usage = getCPUUsage()

	return cpu
}

func getCPUUsage() float64 {
	var idle, kernel, user windows.Filetime

	procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&idle)),
		uintptr(unsafe.Pointer(&kernel)),
		uintptr(unsafe.Pointer(&user)),
	)

	idleTime := uint64(idle.LowDateTime) | (uint64(idle.HighDateTime) << 32)
	kernelTime := uint64(kernel.LowDateTime) | (uint64(kernel.HighDateTime) << 32)
	userTime := uint64(user.LowDateTime) | (uint64(user.HighDateTime) << 32)

	if firstCPUSample {
		prevIdleTime = idleTime
		prevKernelTime = kernelTime
		prevUserTime = userTime
		firstCPUSample = false
		return 0
	}

	idleDelta := idleTime - prevIdleTime
	kernelDelta := kernelTime - prevKernelTime
	userDelta := userTime - prevUserTime

	totalDelta := kernelDelta + userDelta

	prevIdleTime = idleTime
	prevKernelTime = kernelTime
	prevUserTime = userTime

	if totalDelta == 0 {
		return 0
	}

	return 100.0 * float64(totalDelta-idleDelta) / float64(totalDelta)
}

func getMemoryInfo() MemoryInfo {
	var ms memoryStatusEx
	ms.dwLength = uint32(unsafe.Sizeof(ms))

	procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&ms)))

	mem := MemoryInfo{
		Total:     ms.ullTotalPhys,
		Available: ms.ullAvailPhys,
	}
	mem.Used = mem.Total - mem.Available
	if mem.Total > 0 {
		mem.Usage = float64(mem.Used) * 100.0 / float64(mem.Total)
	}

	return mem
}

func getDiskInfo() []DiskInfo {
	var disks []DiskInfo

	driveMask, _, _ := procGetLogicalDrives.Call()
	for i := 0; i < 26; i++ {
		if driveMask&(1<<uint(i)) != 0 {
			disk := getSingleDiskInfo(i)
			if disk != nil {
				disks = append(disks, *disk)
			}
		}
	}

	return disks
}

func getSingleDiskInfo(driveIndex int) *DiskInfo {
	rootPath := string(rune('A'+driveIndex)) + ":\\"

	// 检查驱动器类型（只保留固定驱动器）
	driveType, _, _ := procGetDriveTypeW.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(rootPath))))
	if driveType != 3 { // DRIVE_FIXED = 3
		return nil
	}

	var freeBytes, totalBytes, totalFree uint64
	rootPtr, _ := syscall.UTF16PtrFromString(rootPath)

	ret, _, _ := procGetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(rootPtr)),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFree)),
	)
	if ret == 0 {
		return nil
	}

	// 获取文件系统类型
	var fsNameBuf [260]uint16
	var volNameBuf [260]uint16
	var fsFlags, serialNum uint32

	procGetVolumeInformationW.Call(
		uintptr(unsafe.Pointer(rootPtr)),
		uintptr(unsafe.Pointer(&volNameBuf)),
		260,
		uintptr(unsafe.Pointer(&serialNum)),
		0,
		uintptr(unsafe.Pointer(&fsFlags)),
		uintptr(unsafe.Pointer(&fsNameBuf)),
		260,
	)

	disk := &DiskInfo{
		Label:      string(rune('A'+driveIndex)) + ":",
		FileSystem: syscall.UTF16ToString(fsNameBuf[:]),
		Total:      totalBytes,
		Used:       totalBytes - totalFree,
		Free:       totalFree,
	}
	if totalBytes > 0 {
		disk.Usage = float64(disk.Used) * 100.0 / float64(disk.Total)
	}

	return disk
}

func getNetworkInfo() NetworkInfo {
	// 简化实现：读取注册表中的网络适配器信息
	var adapters []NetworkAdapter

	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`,
		registry.READ)
	if err != nil {
		return NetworkInfo{Adapters: adapters}
	}
	defer k.Close()

	keys, _ := k.ReadSubKeyNames(100)
	for _, subKey := range keys {
		sk, err := registry.OpenKey(k, subKey, registry.READ)
		if err != nil {
			continue
		}
		ip, _, _ := sk.GetStringValue("IPAddress")
		if ip != "" {
			adapters = append(adapters, NetworkAdapter{
				Name:   subKey,
				Type:   "Ethernet",
				IP:     ip,
				IsConnected: len(ip) > 0,
			})
		}
		sk.Close()
	}

	return NetworkInfo{
		Adapters: adapters,
	}
}

func getUptime() string {
	millis, _, _ := procGetTickCount64.Call()
	duration := time.Duration(millis) * time.Millisecond

	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%d天 %d小时 %d分钟", days, hours, minutes)
}
