package system

import (
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// PerfSnapshot 性能快照
type PerfSnapshot struct {
	Timestamp    int64   `json:"timestamp"`
	CPUUsage     float64 `json:"cpuUsage"`
	MemoryUsage  float64 `json:"memoryUsage"`
	MemoryTotal  uint64  `json:"memoryTotal"`
	MemoryUsed   uint64  `json:"memoryUsed"`
	MemoryAvail  uint64  `json:"memoryAvail"`
	DiskReadBps  uint64  `json:"diskReadBps"`
	DiskWriteBps uint64  `json:"diskWriteBps"`
	NetRecvBps   uint64  `json:"netRecvBps"`
	NetSendBps   uint64  `json:"netSendBps"`
	ProcessCount int     `json:"processCount"`
	ThreadCount  int     `json:"threadCount"`
	HandleCount  int     `json:"handleCount"`
}

// PerfHistory 性能历史
type PerfHistory struct {
	mu       sync.Mutex
	snapshots []PerfSnapshot
	maxSize   int
}

// NewPerfHistory 创建性能历史记录
func NewPerfHistory(maxSize int) *PerfHistory {
	if maxSize <= 0 {
		maxSize = 360 // 默认保存 60 秒（每秒采样）
	}
	return &PerfHistory{
		snapshots: make([]PerfSnapshot, 0, maxSize),
		maxSize:   maxSize,
	}
}

// Add 添加快照
func (h *PerfHistory) Add(snap PerfSnapshot) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.snapshots = append(h.snapshots, snap)
	if len(h.snapshots) > h.maxSize {
		h.snapshots = h.snapshots[1:]
	}
}

// GetHistory 获取历史数据
func (h *PerfHistory) GetHistory() []PerfSnapshot {
	h.mu.Lock()
	defer h.mu.Unlock()
	result := make([]PerfSnapshot, len(h.snapshots))
	copy(result, h.snapshots)
	return result
}

// GetLatest 获取最新快照
func (h *PerfHistory) GetLatest() PerfSnapshot {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.snapshots) == 0 {
		return PerfSnapshot{}
	}
	return h.snapshots[len(h.snapshots)-1]
}

// Clear 清空历史
func (h *PerfHistory) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.snapshots = h.snapshots[:0]
}

var (
	globalHistory = NewPerfHistory(360)
	prevCPUTimes  [2]uint64
)

// GetPerfSnapshot 获取当前性能快照
func GetPerfSnapshot() PerfSnapshot {
	snap := PerfSnapshot{
		Timestamp: time.Now().UnixMilli(),
	}

	// CPU 使用率
	total1, idle1 := getSystemTimes()
	if prevCPUTimes[0] > 0 {
		dTotal := total1 - prevCPUTimes[0]
		dIdle := idle1 - prevCPUTimes[1]
		if dTotal > 0 {
			snap.CPUUsage = (1.0 - float64(dIdle)/float64(dTotal)) * 100
		}
	}
	prevCPUTimes[0] = total1
	prevCPUTimes[1] = idle1

	// 内存信息
	snap.MemoryTotal, snap.MemoryAvail, snap.MemoryUsed, snap.MemoryUsage = getMemoryInfoDetailed()

	// 进程/线程/句柄数
	snap.ProcessCount, snap.ThreadCount, snap.HandleCount = getSystemCounts()

	// 记录到历史
	globalHistory.Add(snap)

	return snap
}

// GetPerfHistory 获取性能历史
func GetPerfHistory() []PerfSnapshot {
	return globalHistory.GetHistory()
}

// ClearPerfHistory 清空性能历史
func ClearPerfHistory() {
	globalHistory.Clear()
}

// getSystemTimes 获取系统 CPU 时间
func getSystemTimes() (total, idle uint64) {
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
	idle = uint64(idleTime.Nanoseconds())
	total = uint64(kernelTime.Nanoseconds() + userTime.Nanoseconds())
	return
}

// getMemoryInfoDetailed 获取内存详细信息
func getMemoryInfoDetailed() (total, avail, used uint64, usage float64) {
	dll := windows.NewLazySystemDLL("kernel32.dll")
	proc := dll.NewProc("GlobalMemoryStatusEx")
	if proc == nil {
		return
	}

	type memoryStatusEx struct {
		Length               uint32
		MemoryLoad           uint32
		TotalPhys            uint64
		AvailPhys            uint64
		TotalPageFile        uint64
		AvailPageFile        uint64
		TotalVirtual         uint64
		AvailVirtual         uint64
		AvailExtendedVirtual uint64
	}

	var ms memoryStatusEx
	ms.Length = uint32(unsafe.Sizeof(ms))
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&ms)))
	if ret == 0 {
		return
	}

	total = ms.TotalPhys
	avail = ms.AvailPhys
	used = total - avail
	if total > 0 {
		usage = float64(used) / float64(total) * 100
	}
	return
}

// getSystemCounts 获取系统进程/线程/句柄数
func getSystemCounts() (processes, threads, handles int) {
	// 使用 NtQuerySystemInformation 获取系统信息（快速，无外部进程）
	modNt := windows.NewLazySystemDLL("ntdll.dll")
	procNtQuery := modNt.NewProc("NtQuerySystemInformation")

	// SystemBasicInformation (class 0) 包含进程数
	type systemBasicInformation struct {
		Reserved             uint16
		ProcessorArchitecture uint16
		PageSize             uint32
		MinUserAddress       uintptr
		MaxUserAddress       uintptr
		ActiveProcessorsMask uintptr
		NumberOfProcessors   uint32
		ProcessorType        uint32
		AllocationGranularity uint32
		ProcessorLevel       uint16
		ProcessorRevision    uint16
		NumberOfProcessors32 uint32
	}
	var basicInfo systemBasicInformation
	ret, _, _ := procNtQuery.Call(
		0, // SystemBasicInformation
		uintptr(unsafe.Pointer(&basicInfo)),
		unsafe.Sizeof(basicInfo),
		0,
	)
	if ret == 0 {
		processes = int(basicInfo.NumberOfProcessors)
	}

	// SystemProcessInformation (class 5) 遍历所有进程累加线程和句柄
	type unicodeString struct {
		Length        uint16
		MaxLength     uint16
		Buffer        uintptr
	}
	type systemProcessInformation struct {
		NextEntryOffset              uint32
		NumberOfThreads              uint32
		CreateTime                   int64
		UserTime                     int64
		KernelTime                   int64
		ImageName                    unicodeString
		BasePriority                 int32
		UniqueProcessId              uintptr
		InheritedFromUniqueProcessId uintptr
		HandleCount                  uint32
		SessionId                    uint32
		PeakVirtualSize              uintptr
		VirtualSize                  uintptr
		PageFaultCount               uint32
		PeakWorkingSetSize           uintptr
		WorkingSetSize               uintptr
		QuotaPagedPoolUsage          uintptr
		QuotaNonPagedPoolUsage       uintptr
		PagefileUsage                uintptr
		PeakPagefileUsage            uintptr
		PrivatePageCount             uintptr
		ReadOperationCount           int64
		WriteOperationCount          int64
		OtherOperationCount          int64
		ReadTransferCount            int64
		WriteTransferCount           int64
		OtherTransferCount           int64
	}

	bufSize := uint32(1024 * 1024) // 1MB buffer
	buf := make([]byte, bufSize)
	var returnLen uint32
	ret, _, _ = procNtQuery.Call(
		5, // SystemProcessInformation
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(bufSize),
		uintptr(unsafe.Pointer(&returnLen)),
	)

	if ret == 0 {
		offset := 0
		for {
			proc := (*systemProcessInformation)(unsafe.Pointer(&buf[offset]))
			// 累加线程数（包括 System Idle Process 的线程）
			threads += int(proc.NumberOfThreads)
			handles += int(proc.HandleCount)
			if proc.NextEntryOffset == 0 {
				break
			}
			offset += int(proc.NextEntryOffset)
		}
		// 减去 System Idle Process 的句柄数（它不算真实句柄）
		handles--
	}

	return
}
