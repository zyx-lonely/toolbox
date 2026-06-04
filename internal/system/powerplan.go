package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"pc-toolbox/internal/common"
)

type PowerPlan struct {
	GUID   string `json:"guid"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func GetPowerPlans() []PowerPlan {
	var plans []PowerPlan
	powercfgPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "powercfg.exe")
	c := &exec.Cmd{
		Path: powercfgPath,
		Args: []string{powercfgPath, "/list"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	out, _ := c.Output()
	output := common.GbkToUtf8(string(out))
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "---") {
			continue
		}

		active := strings.Contains(line, "*")

		// 提取 GUID: "电源方案 GUID: xxxx-xxxx  (名称)" 或 "Power Scheme GUID: xxxx"
		guidStart := strings.Index(line, "GUID: ")
		if guidStart < 0 {
			continue
		}
		guidPart := line[guidStart+6:]
		// GUID 结束于下一个空格
		endIdx := strings.Index(guidPart, " ")
		if endIdx < 0 {
			endIdx = len(guidPart)
		}
		guid := guidPart[:endIdx]

		// 提取名称：最后括号内的内容
		name := ""
		parenStart := strings.LastIndex(line, "(")
		parenEnd := strings.LastIndex(line, ")")
		if parenStart >= 0 && parenEnd > parenStart {
			name = line[parenStart+1 : parenEnd]
		}

		if guid != "" {
			plans = append(plans, PowerPlan{GUID: guid, Name: name, Active: active})
		}
	}
	return plans
}

func SetPowerPlan(guid string) error {
	powercfgPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "powercfg.exe")
	c := &exec.Cmd{
		Path: powercfgPath,
		Args: []string{powercfgPath, "/s", guid},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	return c.Run()
}
