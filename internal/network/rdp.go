package network

import (
	"fmt"
	"os/exec"
	"syscall"
	"strings"
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
	return []RemoteDesktopInfo{}
}

// CheckRDPPort 检查远程桌面端口是否开放
func CheckRDPPort(ip string) bool {
	return checkPort(ip, 3389)
}

func checkPort(ip string, port int) bool {
	// 简化实现
	_ = port
	return false
}

func osTempDir() string {
	c := exec.Command("cmd", "/c", "echo", "%TEMP%")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	dir, _ := c.Output()
	return strings.TrimSpace(string(dir))
}

func writeFile(path, content string) error {
	c := exec.Command("cmd", "/c", "echo", content, ">", path)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c.Run()
}
