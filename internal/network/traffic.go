package network

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"strconv"
	"strings"
	"time"

	"pc-toolbox/internal/common"
)

// TrafficSample 流量采样点
type TrafficSample struct {
	Time     string `json:"time"`
	Download uint64 `json:"download"` // bytes/s
	Upload   uint64 `json:"upload"`   // bytes/s
}

// GetTrafficSamples 获取一组网络流量采样（通过 netstat 近似）
func GetTrafficSamples(durationSec int) ([]TrafficSample, error) {
	if durationSec <= 0 {
		durationSec = 5
	}

	before, err := getNetstatReceived()
	if err != nil {
		return nil, err
	}

	var samples []TrafficSample
	for i := 0; i < durationSec; i++ {
		time.Sleep(1 * time.Second)
		after, err := getNetstatReceived()
		if err != nil {
			continue
		}

		diff := after - before
		if diff < 0 {
			diff = 0
		}

		samples = append(samples, TrafficSample{
			Time:     time.Now().Format("15:04:05"),
			Download: diff,
			Upload:   diff / 2, // 近似
		})
		before = after
	}

	return samples, nil
}

func getNetstatReceived() (uint64, error) {
	// 使用 netstat -e 获取网络接口统计
	netstatPath := filepath.Join(os.Getenv("SystemRoot"), "System32", "netstat.exe")
	cmd := &exec.Cmd{
		Path: netstatPath,
		Args: []string{netstatPath, "-e"},
		SysProcAttr: &syscall.SysProcAttr{HideWindow: true},
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("获取网络统计失败: %w", err)
	}

	lines := strings.Split(common.GbkToUtf8(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "已接收字节数") || strings.HasPrefix(line, "Bytes Received") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)
			}
		}
	}

	return 0, fmt.Errorf("未找到网络统计信息")
}
