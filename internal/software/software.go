package software

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// SoftwareInfo 软件信息
type SoftwareInfo struct {
	Name             string `json:"name"`
	Version          string `json:"version"`
	Publisher        string `json:"publisher"`
	InstallDate      string `json:"installDate"`
	Size             string `json:"size"`
	Uninstall        string `json:"uninstall"`
	RegistryKey      string `json:"registryKey"`
	InstallLocation  string `json:"installLocation"`
	QuietUninstall   string `json:"quietUninstall"`
	SystemComponent  int    `json:"systemComponent"`
	WindowsInstaller int    `json:"windowsInstaller"`
}

var registryKeys = []string{
	`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
	`SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
}

// GetInstalledSoftware 获取已安装的软件列表
func GetInstalledSoftware() ([]SoftwareInfo, error) {
	var allSoftware []SoftwareInfo

	// HKLM 64位
	allSoftware = append(allSoftware, readFromRegistryKey(registry.LOCAL_MACHINE, registryKeys[0])...)
	// HKLM 32位
	allSoftware = append(allSoftware, readFromRegistryKey(registry.LOCAL_MACHINE, registryKeys[1])...)
	// HKCU
	allSoftware = append(allSoftware, readFromRegistryKey(registry.CURRENT_USER, registryKeys[0])...)

	// 去重
	seen := make(map[string]bool)
	var result []SoftwareInfo
	for _, s := range allSoftware {
		if s.SystemComponent == 1 {
			continue
		}
		key := s.Name + "|" + s.Publisher
		if !seen[key] {
			seen[key] = true
			result = append(result, s)
		}
	}

	return result, nil
}

func readFromRegistryKey(root registry.Key, subKey string) []SoftwareInfo {
	key, err := registry.OpenKey(root, subKey, registry.READ)
	if err != nil {
		return nil
	}
	defer key.Close()

	names, err := key.ReadSubKeyNames(0)
	if err != nil {
		return nil
	}

	var list []SoftwareInfo
	for _, name := range names {
		sub, err := registry.OpenKey(key, name, registry.READ)
		if err != nil {
			continue
		}

		displayName, _, err := sub.GetStringValue("DisplayName")
		if err != nil || displayName == "" {
			sub.Close()
			continue
		}

		uninstallString, _, _ := sub.GetStringValue("UninstallString")
		if uninstallString == "" {
			sub.Close()
			continue
		}

		version, _, _ := sub.GetStringValue("DisplayVersion")
		publisher, _, _ := sub.GetStringValue("Publisher")
		installDate, _, _ := sub.GetStringValue("InstallDate")
		installLocation, _, _ := sub.GetStringValue("InstallLocation")
		quietUninstall, _, _ := sub.GetStringValue("QuietUninstallString")

		var estimatedSize int64
		if v, _, err := sub.GetIntegerValue("EstimatedSize"); err == nil {
			estimatedSize = int64(v) * 1024 // KB -> bytes
		}

		var sysComp int64
		if v, _, err := sub.GetIntegerValue("SystemComponent"); err == nil {
			sysComp = int64(v)
		}

		var winInst int64
		if v, _, err := sub.GetIntegerValue("WindowsInstaller"); err == nil {
			winInst = int64(v)
		}

		sizeStr := ""
		if estimatedSize > 0 {
			sizeStr = formatSize(estimatedSize)
		}

		dateStr := installDate
		if len(dateStr) >= 8 {
			dateStr = dateStr[:4] + "-" + dateStr[4:6] + "-" + dateStr[6:8]
		}

		list = append(list, SoftwareInfo{
			Name:             displayName,
			Version:          version,
			Publisher:        publisher,
			InstallDate:      dateStr,
			Size:             sizeStr,
			Uninstall:        uninstallString,
			RegistryKey:      subKey,
			InstallLocation:  installLocation,
			QuietUninstall:   quietUninstall,
			SystemComponent:  int(sysComp),
			WindowsInstaller: int(winInst),
		})

		sub.Close()
	}

	return list
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	if bytes >= GB {
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	} else if bytes >= MB {
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	} else if bytes >= KB {
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	}
	return fmt.Sprintf("%d B", bytes)
}

// UninstallSoftware 卸载软件
func UninstallSoftware(uninstallCmd string) error {
	uninstallCmd = strings.TrimSpace(uninstallCmd)

	if strings.Contains(strings.ToLower(uninstallCmd), "msiexec") {
		cmd := exec.Command("cmd", "/c", uninstallCmd)
		if runtime.GOOS == "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		}
		return cmd.Start()
	}

	if strings.HasSuffix(strings.ToLower(uninstallCmd), ".exe") {
		var exePath, args string
		if strings.HasPrefix(uninstallCmd, "\"") {
			endQuote := strings.Index(uninstallCmd[1:], "\"")
			if endQuote != -1 {
				exePath = uninstallCmd[1 : endQuote+1]
				args = strings.TrimSpace(uninstallCmd[endQuote+2:])
			}
		} else {
			spaceIndex := strings.Index(uninstallCmd, " ")
			if spaceIndex != -1 {
				exePath = uninstallCmd[:spaceIndex]
				args = strings.TrimSpace(uninstallCmd[spaceIndex+1:])
			} else {
				exePath = uninstallCmd
			}
		}

		cmd := exec.Command(exePath, strings.Fields(args)...)
		if runtime.GOOS == "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		}
		return cmd.Start()
	}

	cmd := exec.Command("cmd", "/c", uninstallCmd)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	return cmd.Start()
}

// BatchUninstallSoftware 批量卸载软件
func BatchUninstallSoftware(uninstallCmds []string) ([]string, error) {
	var failed []string
	for _, cmd := range uninstallCmds {
		if err := UninstallSoftware(cmd); err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", cmd, err))
		}
	}
	return failed, nil
}

// ExportSoftwareList 导出软件列表
func ExportSoftwareList(softwareList []SoftwareInfo, filePath string) error {
	data, err := json.MarshalIndent(softwareList, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
