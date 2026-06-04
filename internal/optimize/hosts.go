package optimize

import (
	"os"
	"os/exec"
	"syscall"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// HostsEntry hosts 条目
type HostsEntry struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Comment  string `json:"comment"`
	Enabled  bool   `json:"enabled"`
	Line     string `json:"line,omitempty"`
}

const hostsFilePath = `C:\Windows\System32\drivers\etc\hosts`

// GetHostsEntries 读取 hosts 文件所有条目
func GetHostsEntries() ([]HostsEntry, error) {
	data, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return nil, err
	}

	var entries []HostsEntry
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n ")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			entries = append(entries, HostsEntry{
				Comment: strings.TrimPrefix(strings.TrimPrefix(line, "# "), "#"),
				Enabled: true,
				Line:    line,
			})
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			entries = append(entries, HostsEntry{
				IP:       parts[0],
				Hostname: strings.Join(parts[1:], " "),
				Enabled:  true,
				Line:     line,
			})
		}
	}

	return entries, nil
}

// FlushDNSCache 刷新 DNS 缓存
func FlushDNSCache() error {
	cmd := exec.Command("ipconfig", "/flushdns")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

// SetDNSServer 设置 DNS 服务器
func SetDNSServer(adapterName string, primaryDNS string, secondaryDNS string) error {
	keyPath := `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces`

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
	if err != nil {
		return err
	}
	defer k.Close()

	keys, _ := k.ReadSubKeyNames(100)
	k.Close()

	for _, subKey := range keys {
		fullPath := keyPath + `\` + subKey
		sk, err := registry.OpenKey(registry.LOCAL_MACHINE, fullPath, registry.READ)
		if err != nil {
			continue
		}
		desc, _, _ := sk.GetStringValue("Description")
		if strings.Contains(desc, adapterName) || adapterName == "" {
			sk.Close()

			wk, err := registry.OpenKey(registry.LOCAL_MACHINE, fullPath, registry.WRITE)
			if err != nil {
				continue
			}

			dnsList := primaryDNS
			if secondaryDNS != "" {
				dnsList += "," + secondaryDNS
			}

			if err := wk.SetStringValue("NameServer", dnsList); err != nil {
				wk.Close()
				return err
			}
			wk.Close()

			return FlushDNSCache()
		}
		sk.Close()
	}

	return nil
}

// SaveHostsEntries 保存 hosts 文件
func SaveHostsEntries(entries []HostsEntry) error {
	var lines []string
	for _, entry := range entries {
		if entry.Line != "" {
			lines = append(lines, entry.Line)
		} else if entry.Comment != "" {
			lines = append(lines, "# "+entry.Comment)
		} else if entry.IP != "" && entry.Hostname != "" {
			lines = append(lines, entry.IP+" "+entry.Hostname)
		}
	}

	return os.WriteFile(hostsFilePath, []byte(strings.Join(lines, "\r\n")), 0644)
}
