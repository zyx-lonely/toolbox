package security

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"pc-toolbox/internal/common"
)

// FirewallRule 防火墙规则
type FirewallRule struct {
	Name       string `json:"name"`
	DisplayName string `json:"displayName"`
	Direction  string `json:"direction"`
	Action     string `json:"action"`
	Protocol   string `json:"protocol"`
	LocalPort  string `json:"localPort"`
	Program    string `json:"program"`
	Enabled    bool   `json:"enabled"`
	Profile    string `json:"profile"`
}

// GetFirewallRules 获取防火墙规则
func GetFirewallRules(direction string) ([]FirewallRule, error) {
	psScript := fmt.Sprintf(`Get-NetFirewallRule -Direction %s | Select-Object -First 200 | ForEach-Object {
		"$($_.Name)|$($_.DisplayName)|$($_.Direction)|$($_.Action)|$($_.Enabled)|$($_.Profile)"
	}`, direction)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取防火墙规则失败: %w", err)
	}

	output := strings.TrimSpace(common.GbkToUtf8(string(out)))
	if output == "" {
		return nil, nil
	}

	var rules []FirewallRule
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "|") {
			continue
		}

		parts := strings.SplitN(line, "|", 6)
		if len(parts) < 4 {
			continue
		}

		rule := FirewallRule{
			Name:        strings.TrimSpace(parts[0]),
			DisplayName: strings.TrimSpace(parts[1]),
			Direction:   strings.TrimSpace(parts[2]),
			Action:      strings.TrimSpace(parts[3]),
		}

		if len(parts) > 4 {
			rule.Enabled = strings.TrimSpace(parts[4]) == "True"
		}
		if len(parts) > 5 {
			rule.Profile = strings.TrimSpace(parts[5])
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// ToggleFirewallRule 启用/禁用防火墙规则
func ToggleFirewallRule(name string, enabled bool) error {
	action := "Enable"
	if !enabled {
		action = "Disable"
	}

	psScript := fmt.Sprintf(`%s-NetFirewallRule -Name "%s"`, action, name)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("修改防火墙规则失败: %w", err)
	}
	return nil
}

// BlockPort 阻止端口
func BlockPort(port int, protocol string, direction string) error {
	if protocol == "" {
		protocol = "TCP"
	}
	if direction == "" {
		direction = "Inbound"
	}

	name := fmt.Sprintf("Block-%s-%d", protocol, port)
	psScript := fmt.Sprintf(`New-NetFirewallRule -DisplayName "%s" -Direction %s -Protocol %s -LocalPort %d -Action Block -Enabled True`, name, direction, protocol, port)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("阻止端口失败: %w", err)
	}
	return nil
}

// BlockProgram 阻止程序联网
func BlockProgram(programPath string, direction string) error {
	if direction == "" {
		direction = "Both"
	}

	name := fmt.Sprintf("Block-%s", programPath)
	psScript := fmt.Sprintf(`New-NetFirewallRule -DisplayName "%s" -Direction %s -Program "%s" -Action Block -Enabled True`, name, direction, programPath)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("阻止程序失败: %w", err)
	}
	return nil
}

// DeleteFirewallRule 删除防火墙规则
func DeleteFirewallRule(name string) error {
	psScript := fmt.Sprintf(`Remove-NetFirewallRule -Name "%s"`, name)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("删除防火墙规则失败: %w", err)
	}
	return nil
}

// GetFirewallStatus 获取防火墙状态
func GetFirewallStatus() (map[string]bool, error) {
	psScript := `Get-NetFirewallProfile | ForEach-Object { "$($_.Name)|$($_.Enabled)" }`
	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	status := make(map[string]bool)
	output := strings.TrimSpace(string(out))
	if output == "" {
		return status, nil
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "|") {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			enabled := strings.TrimSpace(parts[1]) == "True"
			status[name] = enabled
		}
	}

	return status, nil
}
