package optimize

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

// RegistryScanResult 注册表扫描结果
type RegistryScanResult struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Issue      string `json:"issue"`
	Category   string `json:"category"`
	BackupPath string `json:"backupPath,omitempty"`
}

// RegistryBackup 注册表备份信息
type RegistryBackup struct {
	Path      string `json:"path"`
	CreatedAt string `json:"createdAt"`
	Size      uint64 `json:"size"`
}

// RegistryFixResult 修复结果
type RegistryFixResult struct {
	Key        string `json:"key"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	BackupPath string `json:"backupPath,omitempty"`
}

// ScanRegistry 扫描无效注册表项
func ScanRegistry() []RegistryScanResult {
	var results []RegistryScanResult

	results = append(results, scanUninstallResidues()...)
	results = append(results, scanEmptyKeys()...)

	return results
}

func scanUninstallResidues() []RegistryScanResult {
	var results []RegistryScanResult
	paths := []string{
		`Software\Microsoft\Windows\CurrentVersion\Uninstall`,
		`Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	for _, path := range paths {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
		if err != nil {
			continue
		}
		keys, _ := k.ReadSubKeyNames(500)
		k.Close()

		for _, key := range keys {
			sk, err := registry.OpenKey(registry.LOCAL_MACHINE,
				path+`\`+key, registry.READ)
			if err != nil {
				continue
			}

			displayName, _, _ := sk.GetStringValue("DisplayName")
			installLocation, _, _ := sk.GetStringValue("InstallLocation")

			if displayName == "" {
				results = append(results, RegistryScanResult{
					Key:      path + `\` + key,
					Issue:    "缺少 DisplayName 的 Uninstall 项",
					Category: "uninstall_residue",
				})
				sk.Close()
				continue
			}

			if installLocation != "" {
				if _, err := os.Stat(installLocation); os.IsNotExist(err) {
					results = append(results, RegistryScanResult{
						Key:      path + `\` + key,
						Value:    displayName,
						Issue:    "安装目录已不存在: " + installLocation,
						Category: "uninstall_residue",
					})
				}
			}
			sk.Close()
		}
	}

	return results
}

func scanEmptyKeys() []RegistryScanResult {
	var results []RegistryScanResult

	scanPaths := []struct {
		root registry.Key
		path string
	}{
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts`},
	}

	for _, sp := range scanPaths {
		k, err := registry.OpenKey(sp.root, sp.path, registry.READ)
		if err != nil {
			continue
		}
		scanEmptyKeyRecursive(k, sp.path, &results)
		k.Close()
	}

	return results
}

func scanEmptyKeyRecursive(k registry.Key, currentPath string, results *[]RegistryScanResult) {
	keys, err := k.ReadSubKeyNames(100)
	if err != nil {
		return
	}

	for _, keyName := range keys {
		subKey, err := registry.OpenKey(k, keyName, registry.READ)
		if err != nil {
			continue
		}

		vals, _ := subKey.ReadValueNames(50)
		subKeys, _ := subKey.ReadSubKeyNames(50)

		if len(vals) == 0 && len(subKeys) == 0 {
			*results = append(*results, RegistryScanResult{
				Key:      currentPath + `\` + keyName,
				Issue:    "空注册表键",
				Category: "empty_key",
			})
		}

		if len(subKeys) > 0 {
			scanEmptyKeyRecursive(subKey, currentPath+`\`+keyName, results)
		}
		subKey.Close()
	}
}

// FixRegistryItems 修复指定的注册表问题
func FixRegistryItems(items []RegistryScanResult) []RegistryFixResult {
	var results []RegistryFixResult

	backupPath, err := createRegistryBackup()
	if err != nil {
		for _, item := range items {
			results = append(results, RegistryFixResult{
				Key:     item.Key,
				Success: false,
				Error:   "备份失败: " + err.Error(),
			})
		}
		return results
	}

	for _, item := range items {
		result := RegistryFixResult{
			Key:        item.Key,
			BackupPath: backupPath,
		}

		// 解析注册表路径
		parts := strings.SplitN(item.Key, `\`, 2)
		if len(parts) != 2 {
			result.Success = false
			result.Error = "无效的注册表路径格式"
			results = append(results, result)
			continue
		}

		var root registry.Key
		switch parts[0] {
		case "HKLM", "HKEY_LOCAL_MACHINE":
			root = registry.LOCAL_MACHINE
		case "HKCU", "HKEY_CURRENT_USER":
			root = registry.CURRENT_USER
		case "HKCR", "HKEY_CLASSES_ROOT":
			root = registry.CLASSES_ROOT
		default:
			result.Success = false
			result.Error = "不支持的注册表根键: " + parts[0]
			results = append(results, result)
			continue
		}

		keyPath := parts[1]
		parentPath := filepath.Dir(keyPath)
		keyName := filepath.Base(keyPath)

		// 打开父键
		parentKey, err := registry.OpenKey(root, parentPath, registry.WRITE)
		if err != nil {
			result.Success = false
			result.Error = "打开父键失败: " + err.Error()
			results = append(results, result)
			continue
		}

		// 删除子键（递归）
		err = deleteRegistryKeyTree(parentKey, keyName)
		parentKey.Close()

		if err != nil {
			result.Success = false
			result.Error = "删除注册表项失败: " + err.Error()
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results
}

func createRegistryBackup() (string, error) {
	backupDir := filepath.Join(os.TempDir(), "pc-toolbox-regbackup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	backupFile := filepath.Join(backupDir,
		fmt.Sprintf("registry_backup_%d.reg", time.Now().Unix()))
	return backupFile, nil
}

// deleteRegistryKeyTree 递归删除注册表键树
func deleteRegistryKeyTree(parent registry.Key, keyName string) error {
	// 打开要删除的键
	key, err := registry.OpenKey(parent, keyName, registry.READ|registry.WRITE)
	if err != nil {
		return err
	}

	// 获取所有子键
	subKeys, _ := key.ReadSubKeyNames(1000)
	key.Close()

	// 递归删除子键
	for _, sub := range subKeys {
		if err := deleteRegistryKeyTree(parent, keyName+`\`+sub); err != nil {
			// 忽略单个子键删除失败
			continue
		}
	}

	// 删除当前键
	return parent.DeleteValue(keyName)
}
