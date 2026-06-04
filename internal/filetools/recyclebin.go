package filetools

import (
	"os/exec"
	"syscall"
)

type RecycleBinInfo struct {
	ItemCount int    `json:"itemCount"`
	Size      string `json:"size"`
	Path      string `json:"path"`
}

func GetRecycleBinInfo() RecycleBinInfo {
	return RecycleBinInfo{
		ItemCount: 0,
		Size:      "未知",
		Path:      "$Recycle.Bin",
	}
}

func EmptyRecycleBin() error {
	c := exec.Command("cmd", "/c", "rd", "/s", "/q", "C:\\$Recycle.Bin")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c.Run()
}
