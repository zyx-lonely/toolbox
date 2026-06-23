package main

import (
	"encoding/json"
	"pc-toolbox/internal/common"
	"pc-toolbox/internal/tools"
)

// GetExternalTools 获取外部工具列表
func (a *App) GetExternalTools() common.APIResponse {
	toolList, err := tools.GetExternalTools()
	if err != nil {
		return common.NewErrorResponseStr(err.Error())
	}
	return common.NewSuccessResponse(toolList)
}

// RunExternalTool 运行外部工具（接收 JSON 字符串）
func (a *App) RunExternalTool(toolJSON string) common.APIResponse {
	var tool tools.ExternalTool
	err := json.Unmarshal([]byte(toolJSON), &tool)
	if err != nil {
		return common.NewErrorResponseStr("解析工具配置失败: " + err.Error())
	}

	err = tools.RunExternalTool(tool)
	if err != nil {
		return common.NewErrorResponseStr(err.Error())
	}
	return common.NewSuccessResponse(nil)
}

// GetToolsDir 获取工具目录
func (a *App) GetToolsDir() (string, error) {
	return tools.GetToolsDir(), nil
}

// CheckToolExists 检查工具可执行文件是否存在
func (a *App) CheckToolExists(executable string) bool {
	return tools.CheckToolExists(executable)
}

// GetOperationLogs 获取操作日志
func (a *App) GetOperationLogs() common.APIResponse {
	logs := tools.GetOperationLogs()
	return common.NewSuccessResponse(logs)
}

// ClearOperationLogs 清空操作日志
func (a *App) ClearOperationLogs() common.APIResponse {
	tools.ClearOperationLogs()
	return common.NewSuccessResponse(nil)
}
