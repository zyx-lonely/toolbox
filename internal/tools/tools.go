package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"pc-toolbox/internal/logger"
)

var (
	modShell32    = syscall.NewLazyDLL("shell32.dll")
	procShellExec = modShell32.NewProc("ShellExecuteW")
)

func shellExecuteRunAs(exePath, args, workDir string) error {
	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exePath)
	var params *uint16
	if args != "" {
		params, _ = syscall.UTF16PtrFromString(args)
	}
	var dir *uint16
	if workDir != "" {
		dir, _ = syscall.UTF16PtrFromString(workDir)
	}
	ret, _, err := procShellExec.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(params)),
		uintptr(unsafe.Pointer(dir)),
		1, // SW_SHOWNORMAL
	)
	if ret <= 32 {
		return fmt.Errorf("ShellExecute 失败 (code=%d): %v", ret, err)
	}
	return nil
}

// ExternalTool 外部工具配置
type ExternalTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Executable  string `json:"executable"`
	Args        string `json:"args"`
	WorkingDir  string `json:"workingDir"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
}

// operationLog 操作日志（内存存储）
type operationLog struct {
	Time   string `json:"time"`
	Action string `json:"action"`
	Detail string `json:"detail"`
	Type   string `json:"type"` // success / warning / error / info
}

var (
	logMu  sync.Mutex
	opLogs []operationLog
)

func init() {
	opLogs = loadLogsFromFile()
}

// getLogFilePath 获取日志文件路径
func getLogFilePath() string {
	return filepath.Join(GetToolsDir(), "operation_logs.json")
}

// loadLogsFromFile 从文件加载日志
func loadLogsFromFile() []operationLog {
	filePath := getLogFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return make([]operationLog, 0)
	}
	var logs []operationLog
	if err := json.Unmarshal(data, &logs); err != nil {
		return make([]operationLog, 0)
	}
	return logs
}

// saveLogsToFile 保存日志到文件
func saveLogsToFile() {
	filePath := getLogFilePath()
	data, err := json.MarshalIndent(opLogs, "", "  ")
	if err != nil {
		logger.Info(fmt.Sprintf("序列化日志失败: %v", err))
		return
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		logger.Info(fmt.Sprintf("保存日志文件失败: %v", err))
	}
}

// AddLog 添加操作日志
func AddLog(action, detail, logType string) {
	logMu.Lock()
	defer logMu.Unlock()
	opLogs = append(opLogs, operationLog{
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Action: action,
		Detail: detail,
		Type:   logType,
	})
	// 最多保留 500 条
	if len(opLogs) > 500 {
		opLogs = opLogs[len(opLogs)-500:]
	}
	saveLogsToFile()
}

// GetOperationLogs 获取操作日志
func GetOperationLogs() []operationLog {
	logMu.Lock()
	defer logMu.Unlock()
	// 返回副本，避免竞态
	result := make([]operationLog, len(opLogs))
	copy(result, opLogs)
	return result
}

// ClearOperationLogs 清空操作日志
func ClearOperationLogs() {
	logMu.Lock()
	defer logMu.Unlock()
	opLogs = make([]operationLog, 0)
	saveLogsToFile()
}

// GetExternalTools 获取外部工具列表
func GetExternalTools() ([]ExternalTool, error) {
	configPath := filepath.Join(GetToolsDir(), "tools.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultTools := getDefaultTools()
		if err := saveToolsConfig(defaultTools); err != nil {
			logger.Warn("保存默认工具配置失败: %v", err)
		}
		return defaultTools, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取工具配置失败: %w", err)
	}

	var tools []ExternalTool
	err = json.Unmarshal(data, &tools)
	if err != nil {
		return nil, fmt.Errorf("解析工具配置失败: %w", err)
	}

	return tools, nil
}

// RunExternalTool 运行外部工具
func RunExternalTool(tool ExternalTool) error {
	toolsDir := GetToolsDir()
	var exePath string

	// 如果 executable 已经是绝对路径，直接使用；否则拼接 toolsDir
	if filepath.IsAbs(tool.Executable) {
		exePath = tool.Executable
	} else {
		exePath = filepath.Join(toolsDir, tool.Executable)
	}

	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		AddLog("启动工具失败", fmt.Sprintf("文件不存在: %s", exePath), "error")
		return fmt.Errorf("可执行文件不存在: %s", exePath)
	}

	// 设置工作目录
	workDir := ""
	if tool.WorkingDir != "" {
		if filepath.IsAbs(tool.WorkingDir) {
			workDir = tool.WorkingDir
		} else {
			workDir = filepath.Join(toolsDir, tool.WorkingDir)
		}
	}

	// 使用 ShellExecute runas 提权启动
	err := shellExecuteRunAs(exePath, tool.Args, workDir)
	if err != nil {
		AddLog("启动工具失败", fmt.Sprintf("%s: %v", tool.Name, err), "error")
		logger.Info(fmt.Sprintf("启动工具失败: %s, 错误: %v", tool.Name, err))
		return fmt.Errorf("启动工具失败: %w", err)
	}

	AddLog("启动工具", fmt.Sprintf("已启动: %s (%s)", tool.Name, exePath), "success")
	logger.Info(fmt.Sprintf("已启动外部工具: %s", tool.Name))
	return nil
}

// CheckToolExists 检查工具可执行文件是否存在
func CheckToolExists(executable string) bool {
	toolsDir := GetToolsDir()
	exePath := filepath.Join(toolsDir, executable)
	_, err := os.Stat(exePath)
	return err == nil
}

// GetToolsDir 获取工具目录
func GetToolsDir() string {
	exePath, err := os.Executable()
	if err != nil {
		return "tools"
	}
	toolsDir := filepath.Join(filepath.Dir(exePath), "tools")
	if err := os.MkdirAll(toolsDir, 0755); err != nil {
		logger.Warn("创建工具目录失败: %v", err)
	}
	return toolsDir
}

// getDefaultTools 获取默认工具列表
func getDefaultTools() []ExternalTool {
	return []ExternalTool{
		{
			Name:        "Bulk Crap Uninstaller",
			Description: "强大的软件卸载工具，可以彻底卸载软件并清理残留",
			Executable:  "BCU.exe",
			Args:        "",
			WorkingDir:  "",
			Icon:        "apps-outline",
			Category:    "卸载工具",
		},
		{
			Name:        "Geek Uninstaller",
			Description: "极简强力卸载工具，体积小无残留",
			Executable:  "geek.exe",
			Args:        "",
			WorkingDir:  "",
			Icon:        "trash-outline",
			Category:    "卸载工具",
		},
	}
}

// saveToolsConfig 保存工具配置
func saveToolsConfig(tools []ExternalTool) error {
	configPath := filepath.Join(GetToolsDir(), "tools.json")
	data, err := json.MarshalIndent(tools, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化工具配置失败: %w", err)
	}
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("保存工具配置失败: %w", err)
	}
	logger.Info(fmt.Sprintf("工具配置已保存到: %s", configPath))
	return nil
}
