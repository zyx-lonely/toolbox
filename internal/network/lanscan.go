package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

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

// ScanLAN 扫描局域网设备（主动 ping 子网 1-254）
func ScanLAN() ([]LANDevice, error) {
	localIP := getLocalIP()
	if localIP == "" {
		return nil, fmt.Errorf("无法获取本机 IP")
	}

	subnet := localIP[:strings.LastIndex(localIP, ".")+1]
	var devices []LANDevice
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 限制并发数为 50，避免网络拥塞
	sem := make(chan struct{}, 50)

	for i := 1; i <= 254; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(suffix int) {
			defer wg.Done()
			defer func() { <-sem }()
			ip := fmt.Sprintf("%s%d", subnet, suffix)
			if pingTest(ip, 1) {
				mac := getMACFromARP(ip)
				hostname := lookupHostname(ip)
				mu.Lock()
				devices = append(devices, LANDevice{
					IP:       ip,
					MAC:      mac,
					Hostname: hostname,
					Alive:    true,
				})
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return devices, nil
}

// getLocalIP 获取本机局域网 IP（优先返回私有地址）
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	var fallback string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ip := ipnet.IP.String()
			// 优先返回私有地址
			if strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") ||
				(strings.HasPrefix(ip, "172.") && len(ip) > 4) {
				return ip
			}
			if fallback == "" {
				fallback = ip
			}
		}
	}
	return fallback
}

// getMACFromARP 从 ARP 表获取 MAC 地址
func getMACFromARP(ip string) string {
	cmd := exec.Command("arp", "-a", ip)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	output := common.GbkToUtf8(string(out))
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) >= 3 && parts[0] == ip {
			mac := parts[1]
			if isMAC(mac) {
				return mac
			}
		}
	}
	return ""
}

// pingTest ping 测试
func pingTest(ip string, count int) bool {
	cmd := exec.Command("ping", "-n", fmt.Sprintf("%d", count), "-w", "200", ip)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run() == nil
}

// lookupHostname 反向 DNS 查询
func lookupHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return ""
	}
	return strings.TrimSuffix(names[0], ".")
}

// isMAC 验证 MAC 地址格式
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

// GetLANTopology 获取局域网拓扑（同 ScanLAN，返回相同数据）
func GetLANTopology() ([]LANDevice, error) {
	return ScanLAN()
}

func init() {
	_ = time.Now // 避免未使用导入
}
