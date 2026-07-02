package main

import (
	"pc-toolbox/internal/security"
)

// GetFirewallRules 获取防火墙规则
func (a *App) GetFirewallRules(direction string) ([]security.FirewallRule, error) {
	return security.GetFirewallRules(direction)
}

// ToggleFirewallRule 启用/禁用防火墙规则
func (a *App) ToggleFirewallRule(name string, enabled bool) error {
	return security.ToggleFirewallRule(name, enabled)
}

// BlockPort 阻止端口
func (a *App) BlockPort(port int, protocol string, direction string) error {
	return security.BlockPort(port, protocol, direction)
}

// DeleteFirewallRule 删除防火墙规则
func (a *App) DeleteFirewallRule(name string) error {
	return security.DeleteFirewallRule(name)
}

// GetFirewallStatus 获取防火墙状态
func (a *App) GetFirewallStatus() (map[string]bool, error) {
	return security.GetFirewallStatus()
}

// BlockProgram 阻止程序联网
func (a *App) BlockProgram(programPath string, direction string) error {
	return security.BlockProgram(programPath, direction)
}
