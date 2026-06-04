package optimize

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
)

// ContextMenuAction 右键菜单操作
type ContextMenuAction struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	IconPath    string `json:"iconPath,omitempty"`
}

// ContextMenuStatus 右键菜单状态
type ContextMenuStatus struct {
	Installed bool   `json:"installed"`
	Path      string `json:"path,omitempty"`
	Error     string `json:"error,omitempty"`
}

// InstallContextMenu 安装右键菜单（文件）
func InstallContextMenu() ContextMenuStatus {
	exePath, err := os.Executable()
	if err != nil {
		return ContextMenuStatus{Error: "获取程序路径失败: " + err.Error()}
	}

	// 注册表路径: HKEY_CLASSES_ROOT\*\shell\PCToolbox
	regCmd := fmt.Sprintf(
		`REG ADD "HKCR\*\shell\PCToolbox" /ve /t REG_SZ /d "发送到电脑工具箱" /f && `+
			`REG ADD "HKCR\*\shell\PCToolbox\command" /ve /t REG_SZ /d "\"%s\" \"%%1\"" /f`,
		exePath)

	cmd := exec.Command("cmd", "/c", regCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if output, err := cmd.CombinedOutput(); err != nil {
		return ContextMenuStatus{
			Error: fmt.Sprintf("注册失败: %v\n输出: %s", err, string(output)),
		}
	}

	return ContextMenuStatus{Installed: true, Path: exePath}
}

// InstallContextMenuDir 安装右键菜单（文件夹）
func InstallContextMenuDir() ContextMenuStatus {
	exePath, err := os.Executable()
	if err != nil {
		return ContextMenuStatus{Error: "获取程序路径失败: " + err.Error()}
	}

	regCmd := fmt.Sprintf(
		`REG ADD "HKCR\Directory\shell\PCToolbox" /ve /t REG_SZ /d "电脑工具箱打开" /f && `+
			`REG ADD "HKCR\Directory\shell\PCToolbox\command" /ve /t REG_SZ /d "\"%s\" \"%%1\"" /f`,
		exePath)

	cmd := exec.Command("cmd", "/c", regCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if output, err := cmd.CombinedOutput(); err != nil {
		return ContextMenuStatus{
			Error: fmt.Sprintf("注册失败: %v\n输出: %s", err, string(output)),
		}
	}

	return ContextMenuStatus{Installed: true, Path: exePath}
}

// UninstallContextMenu 卸载右键菜单
func UninstallContextMenu() ContextMenuStatus {
	regCmd := `REG DELETE "HKCR\*\shell\PCToolbox" /f && REG DELETE "HKCR\Directory\shell\PCToolbox" /f`
	cmd := exec.Command("cmd", "/c", regCmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if output, err := cmd.CombinedOutput(); err != nil {
		return ContextMenuStatus{
			Error: fmt.Sprintf("卸载失败: %v\n输出: %s", err, string(output)),
		}
	}
	return ContextMenuStatus{Installed: false}
}

// CheckContextMenu 检查右键菜单是否已安装
func CheckContextMenu() ContextMenuStatus {
	exePath, _ := os.Executable()
	appDir := filepath.Dir(exePath)

	// 检查注册表项是否存在
	cmd := exec.Command("reg", "query", `HKCR\*\shell\PCToolbox`)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return ContextMenuStatus{Installed: false, Path: appDir}
	}

	return ContextMenuStatus{Installed: true, Path: appDir}
}

// 清除未使用的 import 警告
var _ = fmt.Sprintf
