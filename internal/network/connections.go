package network

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"pc-toolbox/internal/common"
)

// NetConnection 网络连接信息
type NetConnection struct {
	Protocol    string `json:"protocol"`
	LocalAddr   string `json:"localAddr"`
	LocalPort   int    `json:"localPort"`
	RemoteAddr  string `json:"remoteAddr"`
	RemotePort  int    `json:"remotePort"`
	State       string `json:"state"`
	PID         int    `json:"pid"`
	ProcessName string `json:"processName"`
}

// ConnGroup 连接分组统计
type ConnGroup struct {
	Total        int `json:"total"`
	TCPCount     int `json:"tcpCount"`
	UDPCount     int `json:"udpCount"`
	Established  int `json:"established"`
	Listening    int `json:"listening"`
	TimeWait     int `json:"timeWait"`
	CloseWait    int `json:"closeWait"`
	ForeignCount int `json:"foreignCount"`
}

var (
	processCache   = make(map[int]string)
	processCacheMu sync.Mutex
)

// GetAllConnections 获取所有网络连接
func GetAllConnections() ([]NetConnection, ConnGroup) {
	// 先加载一次进程名缓存
	loadProcessCache()

	var connections []NetConnection
	group := ConnGroup{}

	connections = append(connections, getConnectionsByProto("TCP")...)
	connections = append(connections, getConnectionsByProto("TCPv6")...)
	connections = append(connections, getConnectionsByProto("UDP")...)
	connections = append(connections, getConnectionsByProto("UDPv6")...)

	group.Total = len(connections)
	for _, c := range connections {
		if c.Protocol == "TCP" || c.Protocol == "TCPv6" {
			group.TCPCount++
		} else {
			group.UDPCount++
		}
		switch c.State {
		case "ESTABLISHED":
			group.Established++
		case "LISTENING":
			group.Listening++
		case "TIME_WAIT":
			group.TimeWait++
		case "CLOSE_WAIT":
			group.CloseWait++
		}
		if c.RemoteAddr != "" && c.RemoteAddr != "0.0.0.0" && c.RemoteAddr != "*" && c.RemoteAddr != "::" {
			group.ForeignCount++
		}
	}

	return connections, group
}

// loadProcessCache 一次性加载所有进程名到缓存
func loadProcessCache() {
	processCacheMu.Lock()
	defer processCacheMu.Unlock()

	cmd := exec.Command("tasklist", "/fo", "CSV", "/nh")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(common.GbkToUtf8(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) >= 2 {
			name := strings.Trim(parts[0], "\"")
			name = strings.TrimSuffix(name, ".exe")
			pidStr := strings.Trim(parts[1], "\"")
			if pid, err := strconv.Atoi(pidStr); err == nil {
				processCache[pid] = name
			}
		}
	}
}

// getConnectionsByProto 按协议获取连接
func getConnectionsByProto(proto string) []NetConnection {
	var connections []NetConnection

	netstatPath := fmt.Sprintf("%s\\System32\\netstat.exe", getSystemRoot())
	cmd := &exec.Cmd{
		Path: netstatPath,
		Args: []string{netstatPath, "-ano", "-p", proto},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return connections
	}

	// 转换 GBK 到 UTF-8
	output := common.GbkToUtf8(string(out))
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 跳过头部行（中英文）
		if strings.HasPrefix(line, "Protocol") || strings.HasPrefix(line, "Active") ||
			strings.HasPrefix(line, "协议") || strings.HasPrefix(line, "活动") {
			continue
		}
		// 只处理以 TCP 或 UDP 开头的数据行
		if !strings.HasPrefix(line, "TCP") && !strings.HasPrefix(line, "UDP") {
			continue
		}

		conn := parseNetstatLine(line, proto)
		if conn != nil {
			connections = append(connections, *conn)
		}
	}

	return connections
}

// parseNetstatLine 解析 netstat 输出行
func parseNetstatLine(line, proto string) *NetConnection {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return nil
	}

	conn := &NetConnection{Protocol: proto}

	// TCP: Proto LocalAddr RemoteAddr State PID
	// UDP: Proto LocalAddr RemoteAddr PID
	if proto == "TCP" || proto == "TCPv6" {
		if len(fields) < 5 {
			return nil
		}
		conn.LocalAddr, conn.LocalPort = parseAddress(fields[1])
		conn.RemoteAddr, conn.RemotePort = parseAddress(fields[2])
		conn.State = fields[3]
		if pid, err := strconv.Atoi(fields[4]); err == nil {
			conn.PID = pid
			conn.ProcessName = getProcessNameCached(pid)
		}
	} else {
		conn.LocalAddr, conn.LocalPort = parseAddress(fields[1])
		conn.RemoteAddr, conn.RemotePort = parseAddress(fields[2])
		if pid, err := strconv.Atoi(fields[3]); err == nil {
			conn.PID = pid
			conn.ProcessName = getProcessNameCached(pid)
		}
	}

	return conn
}

// parseAddress 解析地址:端口
func parseAddress(addr string) (string, int) {
	if strings.HasPrefix(addr, "[") {
		end := strings.Index(addr, "]")
		if end >= 0 {
			ip := addr[1:end]
			portStr := addr[end+2:]
			port, _ := strconv.Atoi(portStr)
			return ip, port
		}
	}

	parts := strings.LastIndex(addr, ":")
	if parts >= 0 {
		ip := addr[:parts]
		portStr := addr[parts+1:]
		port, _ := strconv.Atoi(portStr)
		return ip, port
	}

	return addr, 0
}

// getProcessNameCached 从缓存获取进程名
func getProcessNameCached(pid int) string {
	if pid == 0 || pid == 4 {
		return "System"
	}

	processCacheMu.Lock()
	name, ok := processCache[pid]
	processCacheMu.Unlock()

	if ok {
		return name
	}
	return fmt.Sprintf("PID %d", pid)
}

// getSystemRoot 获取系统根目录
func getSystemRoot() string {
	val, _ := syscall.Getenv("SystemRoot")
	return val
}

// GetConnectionsByPID 获取指定进程的网络连接
func GetConnectionsByPID(pid int) []NetConnection {
	connections, _ := GetAllConnections()
	var result []NetConnection
	for _, c := range connections {
		if c.PID == pid {
			result = append(result, c)
		}
	}
	return result
}

// GetConnectionsByState 按状态筛选连接
func GetConnectionsByState(state string) []NetConnection {
	connections, _ := GetAllConnections()
	var result []NetConnection
	state = strings.ToUpper(state)
	for _, c := range connections {
		if strings.ToUpper(c.State) == state {
			result = append(result, c)
		}
	}
	return result
}

// GetEstablishedConnections 获取已建立的连接
func GetEstablishedConnections() []NetConnection {
	return GetConnectionsByState("ESTABLISHED")
}

// GetListeningConnections 获取监听中的连接
func GetListeningConnections() []NetConnection {
	return GetConnectionsByState("LISTENING")
}

// SearchConnections 搜索连接
func SearchConnections(keyword string) []NetConnection {
	connections, _ := GetAllConnections()
	var result []NetConnection
	keyword = strings.ToLower(keyword)
	for _, c := range connections {
		if strings.Contains(strings.ToLower(c.LocalAddr), keyword) ||
			strings.Contains(strings.ToLower(c.RemoteAddr), keyword) ||
			strings.Contains(strings.ToLower(c.ProcessName), keyword) ||
			strings.Contains(strconv.Itoa(c.PID), keyword) {
			result = append(result, c)
		}
	}
	return result
}

// KillConnection 终止连接
func KillConnection(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("无效的 PID")
	}
	cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

// GetConnectionStats 获取连接统计
func GetConnectionStats() map[string]int {
	stats := make(map[string]int)
	_, group := GetAllConnections()
	stats["total"] = group.Total
	stats["tcp"] = group.TCPCount
	stats["udp"] = group.UDPCount
	stats["established"] = group.Established
	stats["listening"] = group.Listening
	stats["time_wait"] = group.TimeWait
	stats["close_wait"] = group.CloseWait
	return stats
}
