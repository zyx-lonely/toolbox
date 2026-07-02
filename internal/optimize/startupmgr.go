package optimize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// StartupProgram 开机自启程序
type StartupProgram struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Command     string `json:"command"`
	Location    string `json:"location"`    // "HKCU-Run", "HKLM-Run", "StartupFolder", "TaskScheduler"
	Enabled     bool   `json:"enabled"`
	Publisher   string `json:"publisher"`
	StartTime   string `json:"startTime"`   // 启动耗时
	Size        int64  `json:"size"`        // 文件大小
	LastModified string `json:"lastModified"`
}

// StartupGroup 启动项分组统计
type StartupGroup struct {
	Total      int `json:"total"`
	Enabled    int `json:"enabled"`
	Disabled   int `json:"disabled"`
	HKCUCount  int `json:"hkcuCount"`
	HKLMCount  int `json:"hklmCount"`
	FolderCount int `json:"folderCount"`
}

// GetAllStartupPrograms 获取所有开机自启程序
func GetAllStartupPrograms() ([]StartupProgram, StartupGroup) {
	var programs []StartupProgram
	group := StartupGroup{}

	// 1. HKCU Run
	programs = append(programs, getStartupFromRegistry(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		"HKCU-Run")...)
	group.HKCUCount = len(programs)

	// 2. HKCU RunOnce
	programs = append(programs, getStartupFromRegistry(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\RunOnce`,
		"HKCU-RunOnce")...)

	// 3. HKLM Run
	hkLMStart := len(programs)
	programs = append(programs, getStartupFromRegistry(
		registry.LOCAL_MACHINE,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		"HKLM-Run")...)
	group.HKLMCount = len(programs) - hkLMStart

	// 4. HKLM RunOnce
	programs = append(programs, getStartupFromRegistry(
		registry.LOCAL_MACHINE,
		`Software\Microsoft\Windows\CurrentVersion\RunOnce`,
		"HKLM-RunOnce")...)

	// 5. 启动文件夹
	folderStart := len(programs)
	programs = append(programs, getStartupFromFolder()...)
	group.FolderCount = len(programs) - folderStart

	// 统计
	group.Total = len(programs)
	for _, p := range programs {
		if p.Enabled {
			group.Enabled++
		} else {
			group.Disabled++
		}
	}

	return programs, group
}

// getStartupFromRegistry 从注册表获取启动项
func getStartupFromRegistry(root registry.Key, subKey, location string) []StartupProgram {
	var programs []StartupProgram

	k, err := registry.OpenKey(root, subKey, registry.READ)
	if err != nil {
		return programs
	}
	defer k.Close()

	names, err := k.ReadValueNames(200)
	if err != nil {
		return programs
	}

	for _, name := range names {
		val, valType, err := k.GetStringValue(name)
		if err != nil || valType != registry.SZ {
			continue
		}

		prog := StartupProgram{
			Name:     name,
			Command:  val,
			Location: location,
			Enabled:  true,
		}

		// 提取可执行文件路径
		prog.Path = extractExePath(val)

		// 获取文件信息
		if info, err := os.Stat(prog.Path); err == nil {
			prog.Size = info.Size()
			prog.LastModified = info.ModTime().Format("2006-01-02 15:04:05")
		}

		programs = append(programs, prog)
	}

	return programs
}

// getStartupFromFolder 从启动文件夹获取启动项
func getStartupFromFolder() []StartupProgram {
	var programs []StartupProgram

	startupDir := filepath.Join(os.Getenv("APPDATA"),
		"Microsoft", "Windows", "Start Menu", "Programs", "Startup")

	entries, err := os.ReadDir(startupDir)
	if err != nil {
		return programs
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(startupDir, entry.Name())
		prog := StartupProgram{
			Name:        entry.Name(),
			Path:        fullPath,
			Command:     fullPath,
			Location:    "StartupFolder",
			Enabled:     true,
			Size:        info.Size(),
			LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
		}

		programs = append(programs, prog)
	}

	return programs
}

// ToggleStartupProgram 启用/禁用开机自启程序
func ToggleStartupProgram(name, location string, enable bool) error {
	switch location {
	case "HKCU-Run":
		return toggleRegistryStartup(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\Run`, name, enable)
	case "HKCU-RunOnce":
		return toggleRegistryStartup(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\RunOnce`, name, enable)
	case "HKLM-Run":
		return toggleRegistryStartup(registry.LOCAL_MACHINE,
			`Software\Microsoft\Windows\CurrentVersion\Run`, name, enable)
	case "HKLM-RunOnce":
		return toggleRegistryStartup(registry.LOCAL_MACHINE,
			`Software\Microsoft\Windows\CurrentVersion\RunOnce`, name, enable)
	case "StartupFolder":
		return toggleFolderStartup(name, enable)
	default:
		return fmt.Errorf("不支持的启动项位置: %s", location)
	}
}

// toggleRegistryStartup 切换注册表启动项
func toggleRegistryStartup(root registry.Key, subKey, name string, enable bool) error {
	k, err := registry.OpenKey(root, subKey, registry.READ|registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer k.Close()

	if enable {
		// 启用：恢复备份或保持原样
		return nil
	}

	// 禁用：备份后删除
	val, _, err := k.GetStringValue(name)
	if err != nil {
		return fmt.Errorf("获取启动项值失败: %w", err)
	}

	// 确定位置
	location := "HKCU-Run"
	if root == registry.LOCAL_MACHINE {
		location = "HKLM-Run"
	}
	if strings.Contains(subKey, "RunOnce") {
		location += "Once"
	}

	// 保存备份到禁用列表
	saveDisabledBackup(name, val, location)

	// 删除启动项
	return k.DeleteValue(name)
}

// toggleFolderStartup 切换文件夹启动项
func toggleFolderStartup(name string, enable bool) error {
	startupDir := filepath.Join(os.Getenv("APPDATA"),
		"Microsoft", "Windows", "Start Menu", "Programs", "Startup")

	filePath := filepath.Join(startupDir, name)

	if enable {
		// 启用：从备份恢复
		backupPath := getDisabledBackupPath(name)
		if _, err := os.Stat(backupPath); err == nil {
			return os.Rename(backupPath, filePath)
		}
		return nil
	}

	// 禁用：移动到备份目录
	backupDir := filepath.Join(os.Getenv("APPDATA"), "pc-toolbox", "disabled_startup")
	os.MkdirAll(backupDir, 0755)
	backupPath := filepath.Join(backupDir, name)
	return os.Rename(filePath, backupPath)
}

// saveDisabledBackup 保存禁用的启动项备份
func saveDisabledBackup(name, command, location string) {
	backupDir := filepath.Join(os.Getenv("APPDATA"), "pc-toolbox", "disabled_startup")
	os.MkdirAll(backupDir, 0755)

	backupFile := filepath.Join(backupDir, name+".json")
	data := fmt.Sprintf(`{"name":"%s","command":"%s","location":"%s"}`, name, command, location)
	os.WriteFile(backupFile, []byte(data), 0644)
}

// getDisabledBackupPath 获取禁用备份路径
func getDisabledBackupPath(name string) string {
	return filepath.Join(os.Getenv("APPDATA"), "pc-toolbox", "disabled_startup", name)
}

// DeleteStartupProgram 删除开机自启程序
func DeleteStartupProgram(name, location string) error {
	switch location {
	case "HKCU-Run":
		return deleteRegistryStartup(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\Run`, name)
	case "HKCU-RunOnce":
		return deleteRegistryStartup(registry.CURRENT_USER,
			`Software\Microsoft\Windows\CurrentVersion\RunOnce`, name)
	case "HKLM-Run":
		return deleteRegistryStartup(registry.LOCAL_MACHINE,
			`Software\Microsoft\Windows\CurrentVersion\Run`, name)
	case "HKLM-RunOnce":
		return deleteRegistryStartup(registry.LOCAL_MACHINE,
			`Software\Microsoft\Windows\CurrentVersion\RunOnce`, name)
	case "StartupFolder":
		startupDir := filepath.Join(os.Getenv("APPDATA"),
			"Microsoft", "Windows", "Start Menu", "Programs", "Startup")
		return os.Remove(filepath.Join(startupDir, name))
	default:
		return fmt.Errorf("不支持的启动项位置: %s", location)
	}
}

// deleteRegistryStartup 删除注册表启动项
func deleteRegistryStartup(root registry.Key, subKey, name string) error {
	k, err := registry.OpenKey(root, subKey, registry.READ|registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer k.Close()
	return k.DeleteValue(name)
}

// extractExePath 从命令行提取可执行文件路径
func extractExePath(command string) string {
	cmd := strings.TrimSpace(command)
	if cmd == "" {
		return ""
	}

	// 处理带引号的路径
	if strings.HasPrefix(cmd, `"`) {
		end := strings.Index(cmd[1:], `"`)
		if end >= 0 {
			return cmd[1 : end+1]
		}
	}

	// 处理不带引号的路径
	parts := strings.Fields(cmd)
	if len(parts) > 0 {
		return parts[0]
	}

	return cmd
}

// GetStartupGroupStats 获取启动项分组统计
func GetStartupGroupStats() StartupGroup {
	_, group := GetAllStartupPrograms()
	return group
}

// SearchStartupPrograms 搜索启动项
func SearchStartupPrograms(keyword string) []StartupProgram {
	programs, _ := GetAllStartupPrograms()
	var results []StartupProgram

	keyword = strings.ToLower(keyword)
	for _, p := range programs {
		if strings.Contains(strings.ToLower(p.Name), keyword) ||
			strings.Contains(strings.ToLower(p.Command), keyword) {
			results = append(results, p)
		}
	}

	return results
}

// GetStartupProgramDetail 获取启动项详细信息
func GetStartupProgramDetail(name, location string) (*StartupProgram, error) {
	programs, _ := GetAllStartupPrograms()

	for _, p := range programs {
		if p.Name == name && p.Location == location {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("未找到启动项: %s", name)
}
