package system

import (
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

type ShortcutKey struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Modifiers string `json:"modifiers"`
	Key       string `json:"key"`
	FullPath  string `json:"fullPath"`
	AppPath   string `json:"appPath"`
	AppName   string `json:"appName"`
	Enabled   bool   `json:"enabled"`
	Source    string `json:"source"`
}

var (
	modUser32            = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey   = modUser32.NewProc("RegisterHotKey")
	procUnregisterHotKey = modUser32.NewProc("UnregisterHotKey")
)

func GetShortcutKeys() []ShortcutKey {
	var keys []ShortcutKey
	id := 1

	hotkeyKeys := []string{
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\AutoplayHandlers`,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\Explorer`,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced`,
	}

	for _, keyPath := range hotkeyKeys {
		for _, root := range []registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
			key, err := registry.OpenKey(root, keyPath, registry.READ)
			if err != nil {
				continue
			}

			val, _, err := key.GetStringValue("Hotkey")
			if err == nil && val != "" {
				keys = append(keys, ShortcutKey{
					ID:        id,
					Name:      keyPath,
					FullPath:  val,
					Enabled:   true,
					Source:    "registry",
					Modifiers: parseModifiers(val),
					Key:       parseKey(val),
				})
				id++
			}
			key.Close()
		}
	}

	windowsHotkeys := getWindowsHotkeys()
	for _, hk := range windowsHotkeys {
		hk.ID = id
		id++
		keys = append(keys, hk)
	}

	return keys
}

func getWindowsHotkeys() []ShortcutKey {
	var keys []ShortcutKey

	keyCombinations := []struct {
		modifiers string
		key       string
		name      string
		action    string
	}{
		{"Ctrl+Shift", "Esc", "任务管理器", "打开任务管理器"},
		{"Win", "L", "锁定计算机", "锁定计算机"},
		{"Win", "D", "显示桌面", "最小化所有窗口"},
		{"Win", "E", "文件资源管理器", "打开资源管理器"},
		{"Win", "R", "运行", "打开运行对话框"},
		{"Win", "I", "设置", "打开Windows设置"},
		{"Win", "X", "快速链接菜单", "打开高级用户菜单"},
		{"Win", "P", "投影", "切换投影模式"},
		{"Win", "A", "操作中心", "打开操作中心"},
		{"Win", "S", "搜索", "打开搜索"},
		{"Win", "Tab", "任务视图", "切换任务视图"},
		{"Alt", "Tab", "切换窗口", "切换活动窗口"},
		{"Alt", "F4", "关闭窗口", "关闭当前窗口"},
		{"Ctrl", "C", "复制", "复制选中内容"},
		{"Ctrl", "V", "粘贴", "粘贴剪贴板内容"},
		{"Ctrl", "X", "剪切", "剪切选中内容"},
		{"Ctrl", "Z", "撤销", "撤销上一步操作"},
		{"Ctrl", "Y", "重做", "重做上一步操作"},
		{"Ctrl", "A", "全选", "全选所有内容"},
		{"Ctrl", "F", "查找", "查找内容"},
		{"Ctrl", "H", "替换", "查找并替换"},
		{"Ctrl", "N", "新建", "新建窗口/文档"},
		{"Ctrl", "S", "保存", "保存当前文件"},
		{"Ctrl", "P", "打印", "打印当前文件"},
		{"Ctrl", "W", "关闭标签", "关闭当前标签"},
		{"Ctrl", "Tab", "切换标签", "切换标签页"},
		{"Ctrl", "+", "放大", "放大页面"},
		{"Ctrl", "-", "缩小", "缩小页面"},
		{"Ctrl", "0", "重置缩放", "重置页面缩放"},
		{"F1", "", "帮助", "打开帮助"},
		{"F2", "", "重命名", "重命名选中项"},
		{"F3", "", "搜索", "打开搜索"},
		{"F5", "", "刷新", "刷新当前页面"},
		{"F11", "", "全屏", "切换全屏模式"},
		{"PrintScreen", "", "截屏", "截取整个屏幕"},
		{"Alt", "PrintScreen", "活动窗口截图", "截取活动窗口"},
	}

	for _, combo := range keyCombinations {
		fullPath := combo.modifiers + "+" + combo.key
		keys = append(keys, ShortcutKey{
			Name:      combo.name,
			Modifiers: combo.modifiers,
			Key:       combo.key,
			FullPath:  fullPath,
			AppName:   combo.action,
			Enabled:   true,
			Source:    "system",
		})
	}

	return keys
}

func parseModifiers(val string) string {
	val = strings.ToUpper(val)
	var mods []string
	if strings.Contains(val, "CTRL") {
		mods = append(mods, "Ctrl")
	}
	if strings.Contains(val, "ALT") {
		mods = append(mods, "Alt")
	}
	if strings.Contains(val, "SHIFT") {
		mods = append(mods, "Shift")
	}
	if strings.Contains(val, "WIN") {
		mods = append(mods, "Win")
	}
	return strings.Join(mods, "+")
}

func parseKey(val string) string {
	val = strings.ToUpper(val)
	keyMap := map[string]string{
		"ESC": "Esc", "ESCAPE": "Esc",
		"TAB": "Tab", "ENTER": "Enter", "RETURN": "Enter",
		"SPACE": "Space", "DEL": "Delete", "DELETE": "Delete",
		"INS": "Insert", "INSERT": "Insert",
		"HOME": "Home", "END": "End",
		"PAGEUP": "PageUp", "PAGEDOWN": "PageDown",
		"UP": "Up", "DOWN": "Down", "LEFT": "Left", "RIGHT": "Right",
		"F1": "F1", "F2": "F2", "F3": "F3", "F4": "F4",
		"F5": "F5", "F6": "F6", "F7": "F7", "F8": "F8",
		"F9": "F9", "F10": "F10", "F11": "F11", "F12": "F12",
		"PRINTSCREEN": "PrintScreen", "PRTSC": "PrintScreen",
		"NUMLOCK": "NumLock", "CAPSLOCK": "CapsLock", "SCROLLLOCK": "ScrollLock",
	}

	for k, v := range keyMap {
		if strings.Contains(val, k) {
			return v
		}
	}

	re := regexp.MustCompile(`\d+`)
	if match := re.FindString(val); match != "" {
		return match
	}

	if len(val) == 1 {
		return strings.ToUpper(val)
	}

	return val
}
