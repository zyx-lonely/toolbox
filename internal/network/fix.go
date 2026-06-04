package network

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"strings"

	"pc-toolbox/internal/common"
)

// FixResult 网络修复结果
type FixResult struct {
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// FlushDNS 刷新 DNS 缓存
func FlushDNS() FixResult {
	return runCmd("ipconfig", []string{"/flushdns"}, "刷新 DNS 缓存")
}

// ReleaseIP 释放 IP 地址
func ReleaseIP() FixResult {
	return runCmd("ipconfig", []string{"/release"}, "释放 IP 地址")
}

// RenewIP 续租 IP 地址
func RenewIP() FixResult {
	return runCmd("ipconfig", []string{"/renew"}, "续租 IP 地址")
}

// ResetWinsock 重置 Winsock 目录
func ResetWinsock() FixResult {
	return runCmd("netsh", []string{"winsock", "reset"}, "重置 Winsock 目录")
}

// ResetTCPIP 重置 TCP/IP 协议栈
func ResetTCPIP() FixResult {
	return runCmd("netsh", []string{"int", "ip", "reset"}, "重置 TCP/IP 协议栈")
}

// ResetFirewall 重置防火墙规则
func ResetFirewall() FixResult {
	return runCmd("netsh", []string{"advfirewall", "reset"}, "重置防火墙规则")
}

// ResetProxy 清除代理设置
func ResetProxy() FixResult {
	return runCmd("reg", []string{"delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/f"}, "清除代理设置")
}

// ResetArp 清空 ARP 缓存
func ResetArp() FixResult {
	return runCmd("arp", []string{"-d"}, "清空 ARP 缓存")
}

// DiagnoseNetwork 快速诊断
func DiagnoseNetwork() []FixResult {
	return []FixResult{
		FlushDNS(),
		ResetArp(),
		ResetWinsock(),
	}
}

// FixAll 一键网络修复
func FixAll() []FixResult {
	var results []FixResult

	actions := []struct {
		label string
		fn    func() FixResult
	}{
		{"刷新 DNS 缓存", FlushDNS},
		{"重置 Winsock", ResetWinsock},
		{"重置 TCP/IP", ResetTCPIP},
		{"清除 ARP 缓存", ResetArp},
	}

	for _, a := range actions {
		r := a.fn()
		results = append(results, r)
	}

	return results
}

func runCmd(name string, args []string, action string) FixResult {
	// 使用完整路径，避免 PATH 环境问题
	sysDir := filepath.Join(os.Getenv("SystemRoot"), "System32")
	fullPath := filepath.Join(sysDir, name+".exe")
	cmdPath := name
	if _, err := os.Stat(fullPath); err == nil {
		cmdPath = fullPath
	}
	cmd := &exec.Cmd{
		Path: cmdPath,
		Args: append([]string{cmdPath}, args...),
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	output, err := cmd.CombinedOutput()

	result := FixResult{
		Action: action,
		Output: common.GbkToUtf8(strings.TrimSpace(string(output))),
	}

	if err != nil {
		result.Success = false
		result.Error = common.GbkToUtf8(fmt.Sprintf("%v", err))
	} else {
		result.Success = true
	}

	return result
}
