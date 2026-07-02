package optimize

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// FileAssociation 文件关联信息
type FileAssociation struct {
	Extension    string `json:"extension"`    // 文件扩展名
	ContentType  string `json:"contentType"`  // MIME 类型
	ProgID       string `json:"progId"`       // 程序标识符
	Description  string `json:"description"`  // 文件类型描述
	OpenCommand  string `json:"openCommand"`  // 打开命令
	Icon         string `json:"icon"`         // 图标路径
	DefaultApp   string `json:"defaultApp"`   // 默认打开程序
}

// GetFileAssociations 获取所有文件关联
func GetFileAssociations() []FileAssociation {
	var associations []FileAssociation

	// 常见文件扩展名
	commonExtensions := []string{
		".txt", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".pdf", ".jpg", ".jpeg", ".png", ".gif", ".bmp",
		".mp3", ".mp4", ".avi", ".mkv", ".mov",
		".zip", ".rar", ".7z", ".tar", ".gz",
		".exe", ".msi", ".bat", ".cmd", ".ps1",
		".html", ".htm", ".css", ".js", ".json", ".xml",
		".py", ".java", ".cpp", ".c", ".h", ".go", ".rs",
	}

	for _, ext := range commonExtensions {
		assoc := getFileAssociation(ext)
		if assoc != nil {
			associations = append(associations, *assoc)
		}
	}

	return associations
}

// getFileAssociation 获取单个文件关联
func getFileAssociation(extension string) *FileAssociation {
	assoc := &FileAssociation{
		Extension: extension,
	}

	// 从注册表获取文件关联信息
	k, err := registry.OpenKey(registry.CLASSES_ROOT, extension, registry.READ)
	if err != nil {
		return nil
	}
	defer k.Close()

	// 获取 ProgID
	progID, _, err := k.GetStringValue("")
	if err != nil || progID == "" {
		return nil
	}
	assoc.ProgID = progID

	// 获取文件类型描述
	descKey, err := registry.OpenKey(registry.CLASSES_ROOT, progID, registry.READ)
	if err == nil {
		desc, _, _ := descKey.GetStringValue("")
		assoc.Description = desc
		descKey.Close()
	}

	// 获取打开命令
	commandKey, err := registry.OpenKey(registry.CLASSES_ROOT, progID+`\shell\open\command`, registry.READ)
	if err == nil {
		command, _, _ := commandKey.GetStringValue("")
		assoc.OpenCommand = command
		commandKey.Close()
	}

	// 获取图标
	iconKey, err := registry.OpenKey(registry.CLASSES_ROOT, progID+`\DefaultIcon`, registry.READ)
	if err == nil {
		icon, _, _ := iconKey.GetStringValue("")
		assoc.Icon = icon
		iconKey.Close()
	}

	// 获取默认打开程序
	assoc.DefaultApp = getDefaultAppForExtension(extension)

	return assoc
}

// getDefaultAppForExtension 获取默认打开程序
func getDefaultAppForExtension(extension string) string {
	// 使用 assoc 和 ftype 命令获取
	cmd := exec.Command("cmd", "/c", "assoc", extension)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// 解析输出: .txt=txtfile
	parts := strings.SplitN(strings.TrimSpace(string(output)), "=", 2)
	if len(parts) != 2 {
		return ""
	}

	// 获取 ftype
	cmd2 := exec.Command("cmd", "/c", "ftype", parts[1])
	cmd2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		return ""
	}

	// 解析输出: txtfile="%1" %*
	ftypeParts := strings.SplitN(strings.TrimSpace(string(output2)), "=", 2)
	if len(ftypeParts) == 2 {
		return ftypeParts[1]
	}

	return ""
}

// SetFileAssociation 设置文件关联
func SetFileAssociation(extension, progID, command string) error {
	// 设置扩展名到 ProgID 的映射
	k, _, err := registry.CreateKey(registry.CLASSES_ROOT, extension, registry.WRITE)
	if err != nil {
		return fmt.Errorf("创建注册表键失败: %w", err)
	}
	defer k.Close()

	if err := k.SetStringValue("", progID); err != nil {
		return fmt.Errorf("设置 ProgID 失败: %w", err)
	}

	// 设置 ProgID 的打开命令
	commandKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, progID+`\shell\open\command`, registry.WRITE)
	if err != nil {
		return fmt.Errorf("创建命令注册表键失败: %w", err)
	}
	defer commandKey.Close()

	if err := commandKey.SetStringValue("", command); err != nil {
		return fmt.Errorf("设置打开命令失败: %w", err)
	}

	return nil
}

// RemoveFileAssociation 删除文件关联
func RemoveFileAssociation(extension string) error {
	// 删除扩展名的注册表键
	k, err := registry.OpenKey(registry.CLASSES_ROOT, "", registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer k.Close()

	// 递归删除子键
	deleteSubKeyTree(k, extension)

	return nil
}

// deleteSubKeyTree 递归删除注册表子键
func deleteSubKeyTree(parent registry.Key, subKeyPath string) {
	subKey, err := registry.OpenKey(parent, subKeyPath, registry.READ)
	if err != nil {
		return
	}

	// 获取所有子键
	subKeys, _ := subKey.ReadSubKeyNames(1000)
	subKey.Close()

	// 递归删除子键
	for _, sub := range subKeys {
		deleteSubKeyTree(parent, subKeyPath+`\`+sub)
	}

	// 删除当前键
	parent.DeleteValue(subKeyPath)
}

// ResetFileAssociation 重置文件关联为系统默认
func ResetFileAssociation(extension string) error {
	// 使用 assoc 命令重置
	cmd := exec.Command("cmd", "/c", "assoc", extension+"=")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重置文件关联失败: %w", err)
	}

	return nil
}

// GetFileAssociationByExtension 通过扩展名获取文件关联
func GetFileAssociationByExtension(extension string) (*FileAssociation, error) {
	assoc := getFileAssociation(extension)
	if assoc == nil {
		return nil, fmt.Errorf("未找到扩展名 %s 的文件关联", extension)
	}
	return assoc, nil
}

// SetDefaultApp 设置默认打开程序
func SetDefaultApp(extension, appPath string) error {
	// 使用 ftype 设置默认打开程序
	progID := getProgIDForExtension(extension)
	if progID == "" {
		return fmt.Errorf("未找到扩展名 %s 的 ProgID", extension)
	}

	cmd := exec.Command("cmd", "/c", "ftype", progID+"=\""+appPath+"\" \"%1\"")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置默认程序失败: %w", err)
	}

	return nil
}

// getProgIDForExtension 获取扩展名对应的 ProgID
func getProgIDForExtension(extension string) string {
	k, err := registry.OpenKey(registry.CLASSES_ROOT, extension, registry.READ)
	if err != nil {
		return ""
	}
	defer k.Close()

	progID, _, err := k.GetStringValue("")
	if err != nil {
		return ""
	}

	return progID
}

// GetRegisteredExtensions 获取所有已注册的扩展名
func GetRegisteredExtensions() []string {
	var extensions []string

	// 遍历 HKEY_CLASSES_ROOT 获取所有扩展名
	k, err := registry.OpenKey(registry.CLASSES_ROOT, "", registry.READ)
	if err != nil {
		return extensions
	}
	defer k.Close()

	key, err := k.ReadSubKeyNames(10000)
	if err != nil {
		return extensions
	}

	for _, name := range key {
		if strings.HasPrefix(name, ".") {
			extensions = append(extensions, name)
		}
	}

	return extensions
}

// ExportFileAssociations 导出文件关联配置
func ExportFileAssociations() map[string]string {
	result := make(map[string]string)

	extensions := []string{
		".txt", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".pdf", ".jpg", ".jpeg", ".png", ".gif",
		".mp3", ".mp4", ".avi",
		".zip", ".rar", ".7z",
		".html", ".css", ".js",
	}

	for _, ext := range extensions {
		assoc := getFileAssociation(ext)
		if assoc != nil && assoc.DefaultApp != "" {
			result[ext] = assoc.DefaultApp
		}
	}

	return result
}

// RestoreFileAssociations 恢复文件关联配置
func RestoreFileAssociations(associations map[string]string) error {
	for ext, appPath := range associations {
		if err := SetDefaultApp(ext, appPath); err != nil {
			return fmt.Errorf("恢复扩展名 %s 的关联失败: %w", ext, err)
		}
	}

	return nil
}

// GetOpenWithList 获取"打开方式"列表
func GetOpenWithList(extension string) []string {
	var apps []string

	// 从注册表获取"打开方式"列表
	keyPath := `Software\Classes\` + extension + `\OpenWithList`
	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.READ)
	if err != nil {
		return apps
	}
	defer k.Close()

	subkeys, err := k.ReadSubKeyNames(100)
	if err != nil {
		return apps
	}

	return subkeys
}

// AddToOpenWithList 添加到"打开方式"列表
func AddToOpenWithList(extension, appPath string) error {
	keyPath := `Software\Classes\` + extension + `\OpenWithList`
	k, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.WRITE)
	if err != nil {
		return fmt.Errorf("创建注册表键失败: %w", err)
	}
	defer k.Close()

	// 创建应用程序条目
	appKey, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath+`\`+appPath, registry.WRITE)
	if err != nil {
		return fmt.Errorf("创建应用程序条目失败: %w", err)
	}
	defer appKey.Close()

	return nil
}

// RemoveFromOpenWithList 从"打开方式"列表移除
func RemoveFromOpenWithList(extension, appPath string) error {
	keyPath := `Software\Classes\` + extension + `\OpenWithList`
	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.WRITE)
	if err != nil {
		return fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer k.Close()

	if err := k.DeleteValue(appPath); err != nil {
		return fmt.Errorf("移除应用程序失败: %w", err)
	}

	return nil
}
