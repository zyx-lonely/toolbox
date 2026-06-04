package common

import (
	"os/exec"
	"syscall"
)

// CmdHidden 创建隐藏窗口的命令
func CmdHidden(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	return cmd
}

// CmdOutput 执行命令并获取输出（隐藏窗口）
func CmdOutput(name string, args ...string) ([]byte, error) {
	return CmdHidden(name, args...).CombinedOutput()
}
