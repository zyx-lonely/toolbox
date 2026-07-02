package network

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

// RemoteDesktopInfo 远程桌面信息
type RemoteDesktopInfo struct {
	Computer string `json:"computer"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
}

// LaunchMSTSC 启动远程桌面连接
func LaunchMSTSC(computer, address string, port int) error {
	// 创建临时 RDP 文件
	rdpContent := fmt.Sprintf(
		"full address:s:%s:%d\n"+
			"prompt for credentials:i:1\n"+
			"administrative session:i:1\n"+
			"connection type:i:2\n"+
			"session bpp:i:32\n"+
			"screen mode id:i:2\n"+
			"desktopwidth:i:1920\n"+
			"desktopheight:i:1080\n"+
			"autoreconnection enabled:i:1\n",
		address, port)

	// 保存到临时文件
	tmpFile := fmt.Sprintf("%s\\pc-toolbox-rdp.rdp", osTempDir())
	if err := writeFile(tmpFile, rdpContent); err != nil {
		return err
	}

	c := exec.Command("mstsc", tmpFile)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c.Start()
}

// GetRDPClients 获取局域网中开启了 RDP 的设备
func GetRDPClients() []RemoteDesktopInfo {
	var clients []RemoteDesktopInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 获取本机 IP 和子网
	localIP := getLocalIP()
	if localIP == "" {
		return clients
	}

	// 计算子网范围 (假设 /24)
	subnet := localIP[:strings.LastIndex(localIP, ".")+1]

	// 并发扫描 1-254
	for i := 1; i <= 254; i++ {
		wg.Add(1)
		go func(ipSuffix int) {
			defer wg.Done()
			ip := fmt.Sprintf("%s%d", subnet, ipSuffix)
			if checkPort(ip, 3389) {
				hostname, _ := net.LookupAddr(ip)
				name := ip
				if len(hostname) > 0 {
					name = strings.TrimSuffix(hostname[0], ".")
				}
				mu.Lock()
				clients = append(clients, RemoteDesktopInfo{
					Computer: name,
					Address:  ip,
					Port:     3389,
				})
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return clients
}

// CheckRDPPort 检查远程桌面端口是否开放
func CheckRDPPort(ip string) bool {
	return checkPort(ip, 3389)
}

// checkPort 检查指定 IP 和端口是否开放
func checkPort(ip string, port int) bool {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func osTempDir() string {
	return os.TempDir()
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0600)
}
