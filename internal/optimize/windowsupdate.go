package optimize

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// UpdateInfo Windows 更新信息
type UpdateInfo struct {
	Name       string `json:"name"`
	KB         string `json:"kb"`
	Status     string `json:"status"`     // "installed", "pending", "failed"
	InstallDate string `json:"installDate"`
	Size       string `json:"size"`
}

// wuUpdateItem WMI 更新项
type wuUpdateItem struct {
	HotFixID    string `json:"HotFixID"`
	Description string `json:"Description"`
	InstalledOn string `json:"InstalledOn"`
}

// wuPendingItem 待安装更新项
type wuPendingItem struct {
	Title        string `json:"Title"`
	KBArticleIDs string `json:"KBArticleIDs"`
	Size         float64 `json:"Size"`
}

// GetWindowsUpdates 获取已安装的更新列表
func GetWindowsUpdates() []UpdateInfo {
	return getWUUpdateList()
}

// GetPendingUpdates 获取待安装的更新
func GetPendingUpdates() []UpdateInfo {
	return getWUPendingList()
}

// PauseUpdates 暂停更新（暂停 7 天）
func PauseUpdates() string {
	return runPowershellScript(`Install-Module PSWindowsUpdate -Force -Confirm:$false; Get-WUInstall -AcceptAll -AutoReboot:$false`)
}

// ResumeUpdates 恢复更新
func ResumeUpdates() string {
	return runPowershellScript(`Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\WindowsUpdate\UX\Settings" -Name "PauseFeatureUpdatesStartTime" -Value ""; Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\WindowsUpdate\UX\Settings" -Name "PauseFeatureUpdatesEndTime" -Value ""; Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\WindowsUpdate\UX\Settings" -Name "PauseUpdatesExpiryTime" -Value ""`)
}

// UninstallUpdate 卸载指定 KB 更新
func UninstallUpdate(kbID string) string {
	return runPowershellScript("wusa /uninstall /kb:" + kbID + " /quiet /norestart")
}

// CheckForUpdates 检查更新
func CheckForUpdates() string {
	return runPowershellScript("(New-Object -ComObject Microsoft.Update.AutoUpdate).DetectNow()")
}

func getWUUpdateList() []UpdateInfo {
	var updates []UpdateInfo

	output := runPowershellScript(`Get-WmiObject -Class Win32_QuickFixEngineering | Select-Object HotFixID, Description, InstalledOn | ConvertTo-Json`)

	if output == "" || output == "null" {
		return updates
	}

	// 处理单条和数组两种情况
	if !strings.HasPrefix(output, "[") {
		output = "[" + output + "]"
	}

	var items []wuUpdateItem
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return updates
	}

	for _, item := range items {
		kb := item.HotFixID
		if !strings.HasPrefix(kb, "KB") {
			kb = "KB" + kb
		}
		updates = append(updates, UpdateInfo{
			Name:       item.Description,
			KB:         kb,
			Status:     "installed",
			InstallDate: item.InstalledOn,
			Size:       "-",
		})
	}

	return updates
}

func getWUPendingList() []UpdateInfo {
	var updates []UpdateInfo

	output := runPowershellScript(`$Session = New-Object -ComObject Microsoft.Update.Session; $Searcher = $Session.CreateUpdateSearcher(); $Result = $Searcher.Search("IsInstalled=0"); $Result.Updates | Select-Object Title, KBArticleIDs, Size | ConvertTo-Json`)

	if output == "" || output == "null" {
		return updates
	}

	// 处理单条和数组两种情况
	if !strings.HasPrefix(output, "[") {
		output = "[" + output + "]"
	}

	var items []wuPendingItem
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return updates
	}

	for _, item := range items {
		sizeMB := item.Size / 1024 / 1024
		updates = append(updates, UpdateInfo{
			Name:       item.Title,
			KB:         item.KBArticleIDs,
			Status:     "pending",
			InstallDate: "-",
			Size:       fmt.Sprintf("%.1f MB", sizeMB),
		})
	}

	return updates
}

func runPowershellScript(script string) string {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
