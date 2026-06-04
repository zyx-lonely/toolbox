package network

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"pc-toolbox/internal/common"
)

// SignalInfo WiFi 信号信息
type SignalInfo struct {
	SSID     string `json:"ssid"`
	BSSID    string `json:"bssid"`
	Signal   int    `json:"signal"`    // dBm
	Channel  int    `json:"channel"`
	Auth     string `json:"auth"`
	MHz      int    `json:"mhz"`
}

// ScanWiFiSignal 扫描 WiFi 信号强度
func ScanWiFiSignal() []SignalInfo {
	var signals []SignalInfo
	netshPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "netsh.exe")
	c := &exec.Cmd{
		Path: netshPath,
		Args: []string{netshPath, "wlan", "show", "networks", "mode=bssid"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out, err := c.Output()
	if err != nil {
		return signals
	}
	// 转换 GBK 到 UTF-8
	output := common.GbkToUtf8(string(out))
	lines := strings.Split(output, "\n")

	var current SignalInfo
	parsingBSSID := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检测 SSID 行: "SSID 1 : MyWiFi" 或 "SSID 1: MyWiFi"
		if strings.HasPrefix(line, "SSID") && !strings.Contains(line, "BSSID") {
			// 保存上一个
			if current.SSID != "" && current.Channel > 0 {
				signals = append(signals, current)
			}
			idx := strings.Index(line, ":")
			if idx >= 0 {
				current = SignalInfo{SSID: strings.TrimSpace(line[idx+1:])}
				parsingBSSID = false
			}
			continue
		}

		if current.SSID == "" {
			continue
		}

		if strings.Contains(line, "BSSID") || strings.Contains(line, "信号") || strings.Contains(line, "Signal") ||
			strings.Contains(line, "信道") || strings.Contains(line, "Channel") ||
			strings.Contains(line, "MHz") || strings.Contains(line, "认证") || strings.Contains(line, "Auth") {

			idx := strings.Index(line, ":")
			if idx < 0 {
				continue
			}
			val := strings.TrimSpace(line[idx+1:])

			if strings.Contains(line, "BSSID") {
				if parsingBSSID && current.BSSID != "" {
					// 同网络不同BSSID
				}
				current.BSSID = val
				parsingBSSID = true
			} else if strings.Contains(line, "信号") || strings.Contains(line, "Signal") {
				// 提取信号百分比
				pct := extractPercent(val)
				if pct > 0 {
					current.Signal = pct
				}
			} else if strings.Contains(line, "信道") || strings.Contains(line, "Channel") {
				if ch, err := strconv.Atoi(strings.Fields(val)[0]); err == nil {
					current.Channel = ch
				}
			} else if strings.Contains(line, "认证") || strings.Contains(line, "Auth") {
				current.Auth = val
			}
		}
	}

	// 添加最后一个
	if current.SSID != "" && current.Channel > 0 {
		signals = append(signals, current)
	}

	return signals
}

func extractPercent(s string) int {
	fields := strings.Fields(s)
	for _, f := range fields {
		f = strings.TrimSuffix(f, "%")
		if v, err := strconv.Atoi(f); err == nil && v > 0 && v <= 100 {
			return v
		}
	}
	return 0
}

// CheckIPConflict 检测 IP 冲突
func CheckIPConflict(localIP string) []string {
	var conflicts []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	ips := []string{localIP}
	if localIP != "" {
		ips = append(ips, localIP)
	}

	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			// 通过 arp 检查是否有多个 MAC 对应同一 IP
			arpPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "arp.exe")
			ac := &exec.Cmd{
				Path: arpPath,
				Args: []string{arpPath, "-a"},
				SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
			}
			out, _ := ac.Output()
			if strings.Count(string(out), ip) > 1 {
				mu.Lock()
				conflicts = append(conflicts, fmt.Sprintf("IP %s 可能存在冲突", ip))
				mu.Unlock()
			}
		}(ip)
	}
	wg.Wait()
	return conflicts
}

var _ = fmt.Sprintf
