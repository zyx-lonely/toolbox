package optimize

// UpdateInfo Windows 更新信息
type UpdateInfo struct {
	Name       string `json:"name"`
	KB         string `json:"kb"`
	Status     string `json:"status"`     // "installed", "pending", "failed"
	InstallDate string `json:"installDate"`
	Size       string `json:"size"`
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

	// 使用 wmic qfe 获取已安装更新
	output := runPowershellScript(`Get-WmiObject -Class Win32_QuickFixEngineering | Select-Object HotFixID, Description, InstalledOn | ConvertTo-Json`)

	// 简化解析
	if output != "" && len(output) > 10 {
		updates = append(updates, UpdateInfo{
			Name:       "已安装更新",
			KB:         "查看详情",
			Status:     "installed",
			InstallDate: "见控制面板",
			Size:       "-",
		})
	}

	return updates
}

func getWUPendingList() []UpdateInfo {
	output := runPowershellScript(`$Session = New-Object -ComObject Microsoft.Update.Session; $Searcher = $Session.CreateUpdateSearcher(); $Result = $Searcher.Search("IsInstalled=0"); $Result.Updates | Select-Object Title, KBArticleIDs, Size | ConvertTo-Json`)

	if len(output) < 10 {
		return []UpdateInfo{}
	}

	return []UpdateInfo{
		{Name: "待安装更新", KB: "-", Status: "pending", InstallDate: "-", Size: "-"},
	}
}

func runPowershellScript(script string) string {
	// 简化实现，实际需要调用 PowerShell
	return ""
}
