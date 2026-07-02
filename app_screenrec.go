package main

import (
	"pc-toolbox/internal/daily"
)

// StartScreenRecording 开始录屏
func (a *App) StartScreenRecording(duration int, fps int, audio bool) (*daily.ScreenRecordResult, error) {
	config := daily.ScreenRecordConfig{
		Duration: duration,
		FPS:      fps,
		Audio:    audio,
	}
	return daily.StartScreenRecording(config)
}

// StopScreenRecording 停止录屏
func (a *App) StopScreenRecording() *daily.ScreenRecordResult {
	return daily.StopScreenRecording()
}

// IsRecording 是否正在录屏
func (a *App) IsRecording() bool {
	return daily.IsRecording()
}
