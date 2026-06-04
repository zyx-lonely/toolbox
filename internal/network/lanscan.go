package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"strings"
	"sync"

	"pc-toolbox/internal/common"
)

// LANDevice 局域网设备
type LANDevice struct {
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
	Vendor   string `json:"vendor"`
	Alive    bool   `json:"alive"`
}

// ScanLAN 扫描局域网设备
func ScanLAN() ([]LANDevice, error) {
	// 先通过 arp -a 获取现有表
	arpDevices := getARPEntries()

	// 并发 Ping 测试存活
	var wg sync.WaitGroup
	var mu sync.Mutex
	var alive []LANDevice

	for _, d := range arpDevices {
		wg.Add(1)
		go func(dev LANDevice) {
			defer wg.Done()
			if pingTest(dev.IP, 1) {
				dev.Alive = true
				if dev.Hostname == "" {
					dev.Hostname = lookupHostname(dev.IP)
				}
				mu.Lock()
				alive = append(alive, dev)
				mu.Unlock()
			}
		}(d)
	}
	wg.Wait()

	return alive, nil
}

func getARPEntries() []LANDevice {
	var devices []LANDevice
	sysDir := filepath.Join(os.Getenv("SystemRoot"), "System32")
	arpPath := filepath.Join(sysDir, "arp.exe")
	cmd := &exec.Cmd{
		Path: arpPath,
		Args: []string{arpPath, "-a"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out, err := cmd.Output()
	if err != nil {
		return devices
	}

	// 转换 GBK 到 UTF-8
	output := common.GbkToUtf8(string(out))
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 跳过接口头部: "接口: 192.168.1.1 --- 0x1" 或 "Interface: 192.168.1.1 --- 0x1"
		if strings.Contains(line, "---") || strings.Contains(line, "接口") || strings.Contains(line, "Interface") {
			continue
		}
		// Internet 地址 / 物理地址 / 类型 或 Internet Address / Physical Address / Type
		if strings.Contains(line, "Internet") || strings.Contains(line, "物理") || strings.Contains(line, "地址") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 3 {
			ip := parts[0]
			mac := parts[1]
			if isIP(ip) && isMAC(mac) {
				devices = append(devices, LANDevice{
					IP:  ip,
					MAC: mac,
				})
			}
		}
	}
	return devices
}

func pingTest(ip string, count int) bool {
	sysDir := filepath.Join(os.Getenv("SystemRoot"), "System32")
	pingPath := filepath.Join(sysDir, "ping.exe")
	cmd := &exec.Cmd{
		Path: pingPath,
		Args: []string{pingPath, "-n", fmt.Sprintf("%d", count), "-w", "500", ip},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func lookupHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return ""
	}
	return strings.TrimSuffix(names[0], ".")
}

func isIP(s string) bool {
	return net.ParseIP(s) != nil
}

func isMAC(s string) bool {
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, ":", "")
	if len(s) != 12 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
