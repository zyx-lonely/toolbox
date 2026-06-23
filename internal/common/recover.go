package common

import (
	"runtime/debug"

	"pc-toolbox/internal/logger"
)

// GoSafe 安全启动一个 goroutine，自动捕获 panic 并记录日志
func GoSafe(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				logger.Error("goroutine panic: %v\n%s", r, stack)
			}
		}()
		fn()
	}()
}

// GoSafeWithRecover 安全启动 goroutine，支持自定义 recover 回调
func GoSafeWithRecover(fn func(), onRecover func(interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				logger.Error("goroutine panic: %v\n%s", r, stack)
				if onRecover != nil {
					onRecover(r)
				}
			}
		}()
		fn()
	}()
}
