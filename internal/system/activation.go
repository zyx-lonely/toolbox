package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"pc-toolbox/internal/common"
)

// ActivationInfo Windows 激活信息
type ActivationInfo struct {
	ActivationStatus string `json:"activationStatus"`
	ProductID        string `json:"productId"`
	Edition          string `json:"edition"`
}

// ActivationTool 激活工具
type ActivationTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	OpenSource  bool   `json:"openSource"`
	Type        string `json:"type"` // "hwid", "kms", "ohook"
}

// GetActivationInfo 获取 Windows 激活信息
func GetActivationInfo() ActivationInfo {
	info := ActivationInfo{}

	// 方法1: 用 slmgr.vbs 检查激活状态
	systemRoot := os.Getenv("SystemRoot")
	slmgrPath := systemRoot + "\\System32\\slmgr.vbs"
	cscriptPath := systemRoot + "\\System32\\cscript.exe"

	c := &exec.Cmd{
		Path: cscriptPath,
		Args: []string{cscriptPath, "//nologo", slmgrPath, "/dli"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out, err := c.Output()
	if err == nil {
		output := common.GbkToUtf8(string(out))
		if strings.Contains(output, "已授权") || strings.Contains(output, "Licensed") {
			info.ActivationStatus = "已激活"
		} else if strings.Contains(output, "初始安装") || strings.Contains(output, "Initial") {
			info.ActivationStatus = "未激活"
		} else {
			info.ActivationStatus = "已授权"
		}
	} else {
		// 方法2: 用 wmic 查激活状态
		wmicPath := filepath.Join(systemRoot, "System32", "wbem", "wmic.exe")
		c2 := &exec.Cmd{
			Path: wmicPath,
			Args: []string{wmicPath, "path", "SoftwareLicensingProduct", "where", "PartialProductKey is not null", "get", "LicenseStatus"},
			SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
		}
		out2, _ := c2.Output()
		if strings.Contains(common.GbkToUtf8(string(out2)), "1") {
			info.ActivationStatus = "已激活"
		} else {
			info.ActivationStatus = "查询失败（建议以管理员运行）"
		}
	}

	// 获取版本
	wmicPath := filepath.Join(systemRoot, "System32", "wbem", "wmic.exe")
	c3 := &exec.Cmd{
		Path: wmicPath,
		Args: []string{wmicPath, "os", "get", "Caption"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out3, _ := c3.Output()
	lines := strings.Split(common.GbkToUtf8(string(out3)), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" && l != "Caption" {
			info.Edition = l
			break
		}
	}

	return info
}

// GetActivationTools 获取推荐的开源激活工具
func GetActivationTools() []ActivationTool {
	return []ActivationTool{
		{
			Name: "Microsoft Activation Scripts (MAS)",
			Description: "最流行的开源激活脚本，支持 HWID（数字许可证永久激活）、Ohook（Office 永久激活）、KMS38（至 2038 年）等多种方式。支持 Win10/Win11/Office。",
			URL: "https://github.com/massgravel/Microsoft-Activation-Scripts",
			OpenSource: true,
			Type: "hwid",
		},
		{
			Name: "HEU KMS Activator",
			Description: "中文知名的 KMS 激活工具，支持 Windows / Office 的 KMS 激活。简单易用，一键激活。注意：下载请从官方 GitHub 仓库获取。",
			URL: "https://github.com/zbezj/HEU_KMS_Activator",
			OpenSource: true,
			Type: "kms",
		},
	}
}

// KMSActivationMethod 返回 KMS 激活步骤
type KMSMethod struct {
	Title       string `json:"title"`
	Steps       []string `json:"steps"`
	Command     string `json:"command,omitempty"`
}

// GetKMSMethods 获取 KMS 激活方法
func GetKMSMethods() []KMSMethod {
	return []KMSMethod{
		{
			Title: "方法一：使用 KMS 工具",
			Steps: []string{
				"从 GitHub 下载 HEU KMS Activator",
				"以管理员身份运行",
				"点击「激活 Windows」或「激活 Office」",
			},
		},
		{
			Title: "方法二：手动 KMS 激活（适用于 VL 版本）",
			Steps: []string{
				"以管理员身份运行 cmd",
				"安装 KMS 密钥: slmgr /ipk XXXXX-XXXXX-XXXXX-XXXXX-XXXXX",
				"设置 KMS 服务器: slmgr /skms kms.xx",
				"激活: slmgr /ato",
			},
			Command: "slmgr /ipk W269N-WFGWX-YVC9B-4J6C9-T83GX\nslmgr /skms kms.03k.org\nslmgr /ato",
		},
		{
			Title: "方法三：使用 MAS（推荐）",
			Steps: []string{
				"以管理员身份运行 PowerShell",
				"执行: irm https://massgrave.dev/get | iex",
				"选择 [1] HWID 获取数字许可证激活",
			},
			Command: "irm https://massgrave.dev/get | iex",
		},
	}
}

var _ = fmt.Sprintf
