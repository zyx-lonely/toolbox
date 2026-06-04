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
	"time"

	"pc-toolbox/internal/common"
)

// BatchPingResult 批量 Ping 结果
type BatchPingResult struct {
	IP       string `json:"ip"`
	Alive    bool   `json:"alive"`
	Latency  string `json:"latency"`
}

// BatchPing 批量 Ping IP 段
// cidr: 如 "192.168.1.1-100" 或 "192.168.1.0/24"
func BatchPing(cidr string, timeout int) []BatchPingResult {
	var results []BatchPingResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	ips := expandCIDR(cidr)
	sem := make(chan struct{}, 20) // 并发限制
	timeoutMs := timeout
	if timeoutMs <= 0 {
		timeoutMs = 500
	}

	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			r := BatchPingResult{IP: ip}
			pingPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "ping.exe")
			c := &exec.Cmd{
				Path: pingPath,
				Args: []string{pingPath, "-n", "1", "-w", strconv.Itoa(timeoutMs), ip},
				SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
			}
			start := time.Now()
			out, err := c.Output()
			latency := time.Since(start)

			if err == nil {
				r.Alive = true
				output := common.GbkToUtf8(string(out))
				if strings.Contains(output, "时间=") || strings.Contains(output, "time=") {
					r.Latency = fmt.Sprintf("%dms", latency.Milliseconds())
				} else {
					r.Latency = "0ms"
				}
			}

			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(ip)
	}
	wg.Wait()
	return results
}

func expandCIDR(cidr string) []string {
	var ips []string
	if strings.Contains(cidr, "-") {
		parts := strings.Split(cidr, "-")
		if len(parts) == 2 {
			startIP := parts[0]
			endNum, _ := strconv.Atoi(parts[1])
			lastDot := strings.LastIndex(startIP, ".")
			if lastDot > 0 {
				base := startIP[:lastDot+1]
				startNum, _ := strconv.Atoi(startIP[lastDot+1:])
				for i := startNum; i <= endNum; i++ {
					ips = append(ips, fmt.Sprintf("%s%d", base, i))
				}
			}
		}
	}
	return ips
}
