package network

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"strings"

	"pc-toolbox/internal/common"
)

// WiFiProfile WiFi 配置文件信息
type WiFiProfile struct {
	SSID     string `json:"ssid"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

// GetWiFiPasswords 获取所有已保存的 WiFi 密码
func GetWiFiPasswords() []WiFiProfile {
	var profiles []WiFiProfile
	sysDir := filepath.Join(os.Getenv("SystemRoot"), "System32")
	netshPath := filepath.Join(sysDir, "netsh.exe")

	// 获取所有 WiFi 配置文件
	c := &exec.Cmd{
		Path: netshPath,
		Args: []string{netshPath, "wlan", "show", "profiles"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out, err := c.Output()
	if err != nil {
		return profiles
	}

	// 转换 GBK 到 UTF-8
	output := common.GbkToUtf8(string(out))
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		ssid := strings.TrimSpace(parts[1])

		if ssid == "" {
			continue
		}

		// SSID 行的特征是包含"所有用户配置文件"或"All User Profile"或"用户配置文件"或"User Profile"
		if !strings.Contains(key, "配置文件") && !strings.Contains(key, "Profile") {
			continue
		}

		password := getWiFiPassword(netshPath, ssid)
		profiles = append(profiles, WiFiProfile{
			SSID:     ssid,
			Password: password,
			Auth:     "WPA2",
		})
	}

	return profiles
}

func getWiFiPassword(netshPath, ssid string) string {
	c := &exec.Cmd{
		Path: netshPath,
		Args: []string{netshPath, "wlan", "show", "profile", "name="+ssid, "key=clear"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out, err := c.Output()
	if err != nil {
		return "需要管理员权限"
	}

	output := common.GbkToUtf8(string(out))
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 匹配: "关键内容", "Key Content", "Contenu de la clé"
		if containsAny(line, "关键内容", "Key Content", "Contenu") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				pwd := strings.TrimSpace(parts[1])
				if pwd != "" && pwd != "<无>" {
					return pwd
				}
				return "（空密码）"
			}
		}
	}

	return "（空）"
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
