package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"sort"
	"strings"
	"sync"
	"time"

	"pc-toolbox/internal/common"
)

// PingResult Ping 结果
type PingResult struct {
	Target   string `json:"target"`
	Success  bool   `json:"success"`
	Latency  string `json:"latency"`
	TTL      int    `json:"ttl"`
	Sequence int    `json:"sequence"`
	Error    string `json:"error,omitempty"`
}

// PingSummary Ping 统计
type PingSummary struct {
	Target    string       `json:"target"`
	Results   []PingResult `json:"results"`
	Sent      int          `json:"sent"`
	Received  int          `json:"received"`
	LossRate  float64      `json:"lossRate"`
	MinLatency string      `json:"minLatency"`
	MaxLatency string      `json:"maxLatency"`
	AvgLatency string      `json:"avgLatency"`
}

// PortResult 端口扫描结果
type PortResult struct {
	Port    int    `json:"port"`
	State   string `json:"state"` // "open", "closed", "filtered"
	Service string `json:"service"`
}

// DNSResult DNS 查询结果
type DNSResult struct {
	Hostname string   `json:"hostname"`
	Type     string   `json:"type"`
	Answers  []string `json:"answers"`
	TTL      uint32   `json:"ttl"`
	Error    string   `json:"error,omitempty"`
}

// ConnectionInfo 网络连接信息
type ConnectionInfo struct {
	Protocol  string `json:"protocol"`
	LocalAddr string `json:"localAddr"`
	LocalPort int    `json:"localPort"`
	RemoteAddr string `json:"remoteAddr"`
	RemotePort int    `json:"remotePort"`
	State     string `json:"state"`
	PID       int    `json:"pid"`
	Process   string `json:"process"`
}

// 常用端口到服务名的映射
var commonPorts = map[int]string{
	20: "FTP-DATA", 21: "FTP", 22: "SSH", 23: "Telnet",
	25: "SMTP", 53: "DNS", 80: "HTTP", 110: "POP3",
	143: "IMAP", 443: "HTTPS", 445: "SMB", 993: "IMAPS",
	995: "POP3S", 1433: "MSSQL", 1521: "Oracle", 3306: "MySQL",
	3389: "RDP", 5432: "PostgreSQL", 6379: "Redis", 8080: "HTTP-Alt",
	8443: "HTTPS-Alt", 27017: "MongoDB",
}

// Ping 执行 Ping 操作
func Ping(host string, count int, timeout int) (*PingSummary, error) {
	if count <= 0 {
		count = 4
	}
	if timeout <= 0 {
		timeout = 1000
	}

	summary := &PingSummary{
		Target:  host,
		Results: make([]PingResult, 0, count),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	timeoutMs := timeout

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			result := pingOnce(host, seq, timeoutMs)
			mu.Lock()
			summary.Results = append(summary.Results, result)
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	// 排序
	sort.Slice(summary.Results, func(i, j int) bool {
		return summary.Results[i].Sequence < summary.Results[j].Sequence
	})

	// 统计
	summary.Sent = count
	var totalLatency time.Duration
	var minLatency, maxLatency time.Duration

	for _, r := range summary.Results {
		if r.Success {
			summary.Received++
			d, _ := time.ParseDuration(r.Latency)
			totalLatency += d
			if summary.Received == 1 || d < minLatency {
				minLatency = d
			}
			if d > maxLatency {
				maxLatency = d
			}
		}
	}
	summary.LossRate = float64(summary.Sent-summary.Received) * 100.0 / float64(summary.Sent)

	if summary.Received > 0 {
		summary.MinLatency = fmt.Sprintf("%.0fms", float64(minLatency.Microseconds())/1000)
		summary.MaxLatency = fmt.Sprintf("%.0fms", float64(maxLatency.Microseconds())/1000)
		summary.AvgLatency = fmt.Sprintf("%.0fms", float64(totalLatency.Microseconds())/1000/float64(summary.Received))
	}

	return summary, nil
}

func pingOnce(host string, seq int, timeoutMs int) PingResult {
	start := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", host, time.Duration(timeoutMs)*time.Millisecond)
	if err != nil {
		// 回退到 TCP Ping
		return tcpPingFallback(host, seq, timeoutMs)
	}
	defer conn.Close()
	latency := time.Since(start)
	return PingResult{
		Target:   host,
		Success:  true,
		Latency:  latency.String(),
		Sequence: seq,
	}
}

func tcpPingFallback(host string, seq int, timeoutMs int) PingResult {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"),
		time.Duration(timeoutMs)*time.Millisecond)
	if err != nil {
		return PingResult{
			Target:   host,
			Success:  false,
			Sequence: seq,
			Error:    fmt.Sprintf("连接超时或失败: %v", err),
		}
	}
	conn.Close()
	latency := time.Since(start)
	return PingResult{
		Target:   host,
		Success:  true,
		Latency:  latency.String(),
		Sequence: seq,
	}
}

// PortScan 执行端口扫描
func PortScan(host string, ports string) ([]PortResult, error) {
	var portList []int

	if ports == "common" {
		for p := range commonPorts {
			portList = append(portList, p)
		}
		sort.Ints(portList)
	} else if strings.Contains(ports, "-") {
		parts := strings.SplitN(ports, "-", 2)
		start, end := 0, 0
		fmt.Sscanf(parts[0], "%d", &start)
		fmt.Sscanf(parts[1], "%d", &end)
		if start > 0 && end >= start {
			for p := start; p <= end; p++ {
				portList = append(portList, p)
			}
		}
	} else if strings.Contains(ports, ",") {
		for _, s := range strings.Split(ports, ",") {
			var p int
			fmt.Sscanf(strings.TrimSpace(s), "%d", &p)
			if p > 0 {
				portList = append(portList, p)
			}
		}
	} else {
		var p int
		fmt.Sscanf(ports, "%d", &p)
		if p > 0 {
			portList = append(portList, p)
		}
	}

	var results []PortResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, 50) // 并发限制

	for _, port := range portList {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := scanPort(host, p)
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(port)
	}
	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results, nil
}

func scanPort(host string, port int) PortResult {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return PortResult{
			Port:    port,
			State:   "closed",
			Service: commonPorts[port],
		}
	}
	conn.Close()
	return PortResult{
		Port:    port,
		State:   "open",
		Service: commonPorts[port],
	}
}

// DNSLookup DNS 查询
func DNSLookup(hostname string) (*DNSResult, error) {
	ips, err := net.LookupHost(hostname)
	if err != nil {
		return &DNSResult{
			Hostname: hostname,
			Type:     "A",
			Error:    err.Error(),
		}, err
	}

	return &DNSResult{
		Hostname: hostname,
		Type:     "A",
		Answers:  ips,
	}, nil
}

// TraceRoute 路由追踪（简化实现）
func TraceRoute(host string, maxHops int) ([]PingResult, error) {
	if maxHops <= 0 {
		maxHops = 30
	}
	// 使用系统的 tracert 命令
	tracertPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "tracert.exe")
	cmd := &exec.Cmd{
		Path: tracertPath,
		Args: []string{tracertPath, "-h", fmt.Sprintf("%d", maxHops), host},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("tracert 执行失败: %w", err)
	}

	// 解析输出
	lines := strings.Split(common.GbkToUtf8(string(output)), "\n")
	var results []PingResult
	for i, line := range lines {
		if i < 3 || strings.TrimSpace(line) == "" {
			continue
		}
		results = append(results, PingResult{
			Target:  host,
			Success: true,
			Latency: fmt.Sprintf("hop %d", i-2),
		})
	}

	return results, nil
}

// GetNetworkConnections 获取网络连接
func GetNetworkConnections() ([]ConnectionInfo, error) {
	// 使用 netstat 命令
	netstatPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "netstat.exe")
	cmd := &exec.Cmd{
		Path: netstatPath,
		Args: []string{netstatPath, "-ano"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var connections []ConnectionInfo
	lines := strings.Split(common.GbkToUtf8(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		protocol := fields[0]
		if protocol != "TCP" && protocol != "UDP" {
			continue
		}

		localAddr := parseAddr(fields[1])
		remoteAddr := parseAddr(fields[2])
		state := ""
		pid := 0

		if protocol == "TCP" && len(fields) >= 5 {
			state = fields[3]
			fmt.Sscanf(fields[len(fields)-1], "%d", &pid)
		} else if len(fields) >= 4 {
			fmt.Sscanf(fields[len(fields)-1], "%d", &pid)
		}

		connections = append(connections, ConnectionInfo{
			Protocol:   protocol,
			LocalAddr:  localAddr.addr,
			LocalPort:  localAddr.port,
			RemoteAddr: remoteAddr.addr,
			RemotePort: remoteAddr.port,
			State:      state,
			PID:        pid,
		})
	}

	return connections, nil
}

type addrInfo struct {
	addr string
	port int
}

func parseAddr(addr string) addrInfo {
	// 处理 IPv4 和 IPv6 地址
	if strings.HasPrefix(addr, "[") {
		// IPv6: [addr]:port
		parts := strings.SplitN(addr, "]:", 2)
		if len(parts) == 2 {
			var port int
			fmt.Sscanf(parts[1], "%d", &port)
			return addrInfo{addr: parts[0][1:], port: port}
		}
	}
	// IPv4: addr:port
	lastColon := strings.LastIndex(addr, ":")
	if lastColon > 0 {
		var port int
		fmt.Sscanf(addr[lastColon+1:], "%d", &port)
		return addrInfo{addr: addr[:lastColon], port: port}
	}
	return addrInfo{addr: addr, port: 0}
}
