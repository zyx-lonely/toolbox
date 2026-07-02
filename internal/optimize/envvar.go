package optimize

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	modUser32             = syscall.NewLazyDLL("user32.dll")
	procSendMessageTimeout = modUser32.NewProc("SendMessageTimeoutW")
)

const (
	hWndBroadcast      = 0xFFFF
	wmSettingChange    = 0x001A
	smtoAbortIfHung    = 0x0002
)

// EnvVar 环境变量信息
type EnvVar struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Scope     string `json:"scope"`     // "user" or "system"
	Expanded  string `json:"expanded"`  // 展开后的值
}

// GetEnvironmentVariables 获取所有环境变量
func GetEnvironmentVariables(scope string) []EnvVar {
	var vars []EnvVar

	// 获取用户环境变量
	if scope == "" || scope == "user" {
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				vars = append(vars, EnvVar{
					Name:     parts[0],
					Value:    parts[1],
					Scope:    "user",
					Expanded: os.ExpandEnv(parts[1]),
				})
			}
		}
	}

	// 获取系统环境变量
	if scope == "" || scope == "system" {
		systemVars := getSystemEnvVars()
		vars = append(vars, systemVars...)
	}

	return vars
}

// getSystemEnvVars 从注册表获取系统环境变量
func getSystemEnvVars() []EnvVar {
	var vars []EnvVar

	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ)
	if err != nil {
		return vars
	}
	defer k.Close()

	// 获取所有值
	valNames, err := k.ReadValueNames(1000)
	if err != nil {
		return vars
	}

	for _, name := range valNames {
		val, _, err := k.GetStringValue(name)
		if err != nil {
			continue
		}
		vars = append(vars, EnvVar{
			Name:     name,
			Value:    val,
			Scope:    "system",
			Expanded: os.ExpandEnv(val),
		})
	}

	return vars
}

// GetEnvironmentVariable 获取单个环境变量
func GetEnvironmentVariable(name string, scope string) (*EnvVar, error) {
	if scope == "user" {
		val := os.Getenv(name)
		return &EnvVar{
			Name:     name,
			Value:    val,
			Scope:    "user",
			Expanded: os.ExpandEnv(val),
		}, nil
	}

	// 系统环境变量
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("打开注册表失败: %w", err)
	}
	defer k.Close()

	val, _, err := k.GetStringValue(name)
	if err != nil {
		return nil, fmt.Errorf("环境变量 %s 不存在", name)
	}

	return &EnvVar{
		Name:     name,
		Value:    val,
		Scope:    "system",
		Expanded: os.ExpandEnv(val),
	}, nil
}

// SetEnvironmentVariable 设置环境变量
func SetEnvironmentVariable(name, value string, scope string) error {
	if scope == "user" {
		// 设置用户环境变量（需要通知系统）
		return setUserEnvVar(name, value)
	}

	// 设置系统环境变量
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer k.Close()

	if err := k.SetStringValue(name, value); err != nil {
		return fmt.Errorf("设置环境变量失败: %w", err)
	}

	// 通知系统环境变量已更改
	notifyEnvironmentChange()
	return nil
}

// DeleteEnvironmentVariable 删除环境变量
func DeleteEnvironmentVariable(name string, scope string) error {
	if scope == "user" {
		return deleteUserEnvVar(name)
	}

	// 删除系统环境变量
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer k.Close()

	if err := k.DeleteValue(name); err != nil {
		return fmt.Errorf("删除环境变量失败: %w", err)
	}

	notifyEnvironmentChange()
	return nil
}

// setUserEnvVar 设置用户环境变量
func setUserEnvVar(name, value string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Environment`, registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer k.Close()

	if err := k.SetStringValue(name, value); err != nil {
		return fmt.Errorf("设置环境变量失败: %w", err)
	}

	notifyEnvironmentChange()
	return nil
}

// deleteUserEnvVar 删除用户环境变量
func deleteUserEnvVar(name string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Environment`, registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer k.Close()

	if err := k.DeleteValue(name); err != nil {
		return fmt.Errorf("删除环境变量失败: %w", err)
	}

	notifyEnvironmentChange()
	return nil
}

// notifyEnvironmentChange 通知系统环境变量已更改
func notifyEnvironmentChange() {
	lParam, _ := windows.UTF16PtrFromString("Environment")
	procSendMessageTimeout.Call(
		uintptr(hWndBroadcast),
		uintptr(wmSettingChange),
		0,
		uintptr(unsafe.Pointer(lParam)),
		uintptr(smtoAbortIfHung),
		5000,
		0,
	)
}

// ExpandEnvPath 展开环境变量路径
func ExpandEnvPath(path string) string {
	return os.ExpandEnv(path)
}

// ValidateEnvVarName 验证环境变量名是否合法
func ValidateEnvVarName(name string) error {
	if name == "" {
		return fmt.Errorf("环境变量名不能为空")
	}

	// 环境变量名不能包含空格或特殊字符
	for _, c := range name {
		if c == ' ' || c == '=' || c == '\t' {
			return fmt.Errorf("环境变量名包含非法字符")
		}
	}

	return nil
}

// GetPathVariable 获取 PATH 环境变量
func GetPathVariable(scope string) []string {
	var pathValue string
	if scope == "user" {
		pathValue = os.Getenv("PATH")
	} else {
		// 从注册表获取系统 PATH
		k, err := registry.OpenKey(registry.LOCAL_MACHINE,
			`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.READ)
		if err != nil {
			return nil
		}
		defer k.Close()
		pathValue, _, _ = k.GetStringValue("Path")
	}

	return strings.Split(pathValue, ";")
}

// SetPathVariable 设置 PATH 环境变量
func SetPathVariable(paths []string, scope string) error {
	newPath := strings.Join(paths, ";")
	return SetEnvironmentVariable("PATH", newPath, scope)
}

// AddToPathVariable 向 PATH 添加新路径
func AddToPathVariable(newPath string, scope string) error {
	paths := GetPathVariable(scope)

	// 检查是否已存在
	for _, p := range paths {
		if strings.EqualFold(p, newPath) {
			return nil // 已存在，无需添加
		}
	}

	paths = append(paths, newPath)
	return SetPathVariable(paths, scope)
}

// RemoveFromPathVariable 从 PATH 移除路径
func RemoveFromPathVariable(pathToRemove string, scope string) error {
	paths := GetPathVariable(scope)

	var newPaths []string
	for _, p := range paths {
		if !strings.EqualFold(p, pathToRemove) {
			newPaths = append(newPaths, p)
		}
	}

	return SetPathVariable(newPaths, scope)
}
