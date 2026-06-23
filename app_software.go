package main

import (
	"encoding/json"
	"fmt"

	"pc-toolbox/internal/common"
	"pc-toolbox/internal/software"
	"pc-toolbox/internal/tools"
)

// ============================================================
//  软件管理
// ============================================================

// GetInstalledSoftware 获取已安装的软件列表
func (a *App) GetInstalledSoftware() ([]software.SoftwareInfo, error) {
	return software.GetInstalledSoftware()
}

// UninstallSoftware 卸载软件
func (a *App) UninstallSoftware(uninstallCmd string) common.APIResponse {
	tools.AddLog("卸载软件", uninstallCmd, "info")
	err := software.UninstallSoftware(uninstallCmd)
	if err != nil {
		tools.AddLog("卸载软件失败", err.Error(), "error")
		return common.NewErrorResponseStr(err.Error())
	}
	tools.AddLog("卸载软件成功", uninstallCmd, "success")
	return common.NewSuccessResponse(nil)
}

// BatchUninstallSoftware 批量卸载软件
func (a *App) BatchUninstallSoftware(uninstallCmds []string) common.APIResponse {
	tools.AddLog("批量卸载软件", fmt.Sprintf("共 %d 个", len(uninstallCmds)), "info")
	failed, err := software.BatchUninstallSoftware(uninstallCmds)
	if err != nil {
		tools.AddLog("批量卸载失败", err.Error(), "error")
		return common.NewErrorResponseStr(err.Error())
	}

	if len(failed) > 0 {
		tools.AddLog("部分卸载失败", joinStrings(failed, "\n"), "warning")
		return common.NewErrorResponseStr(joinStrings(failed, "\n"))
	}
	tools.AddLog("批量卸载成功", "全部完成", "success")
	return common.NewSuccessResponse(nil)
}

// ExportSoftwareList 导出软件列表
func (a *App) ExportSoftwareList(filePath string, softwareListJSON string) common.APIResponse {
	// 解析 JSON
	var softwareList []software.SoftwareInfo
	err := json.Unmarshal([]byte(softwareListJSON), &softwareList)
	if err != nil {
		return common.NewErrorResponseStr(err.Error())
	}

	err = software.ExportSoftwareList(softwareList, filePath)
	if err != nil {
		return common.NewErrorResponseStr(err.Error())
	}
	return common.NewSuccessResponse(nil)
}
