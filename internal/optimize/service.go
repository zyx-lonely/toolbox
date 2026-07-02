package optimize

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// ServiceBackup 服务配置备份
type ServiceBackup struct {
	Name      string `json:"name"`
	StartType uint32 `json:"startType"`
	Status    string `json:"status"`
}

// 全局服务备份存储
var (
	serviceBackups   = make(map[string]ServiceBackup)
	serviceBackupMu  sync.Mutex
)

// SaveServiceBackup 保存服务原始配置
func SaveServiceBackup(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("打开服务 %s 失败: %w", name, err)
	}
	defer s.Close()

	config, err := s.Config()
	if err != nil {
		return err
	}
	status, _ := s.Query()

	state := "stopped"
	if status.State == svc.Running {
		state = "running"
	}

	serviceBackupMu.Lock()
	serviceBackups[name] = ServiceBackup{
		Name:      name,
		StartType: config.StartType,
		Status:    state,
	}
	serviceBackupMu.Unlock()
	return nil
}

// RestoreService 还原单个服务
func RestoreService(name string) error {
	serviceBackupMu.Lock()
	backup, ok := serviceBackups[name]
	serviceBackupMu.Unlock()
	if !ok {
		return fmt.Errorf("未找到服务 %s 的备份", name)
	}

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("打开服务 %s 失败: %w", name, err)
	}
	defer s.Close()

	config, err := s.Config()
	if err != nil {
		return err
	}
	config.StartType = backup.StartType
	if err := s.UpdateConfig(config); err != nil {
		return err
	}

	// 恢复运行状态
	if backup.Status == "running" {
		if err := s.Start(); err != nil {
			return fmt.Errorf("启动服务失败: %w", err)
		}
	}

	return nil
}

// GetServiceBackups 获取所有已备份的服务
func GetServiceBackups() []ServiceBackup {
	serviceBackupMu.Lock()
	defer serviceBackupMu.Unlock()
	var backups []ServiceBackup
	for _, b := range serviceBackups {
		backups = append(backups, b)
	}
	return backups
}

// ClearServiceBackups 清除服务备份
func ClearServiceBackups() {
	serviceBackupMu.Lock()
	serviceBackups = make(map[string]ServiceBackup)
	serviceBackupMu.Unlock()
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Status      string `json:"status"`      // "running", "stopped", "paused"
	StartType   string `json:"startType"`   // "auto", "manual", "disabled"
	Recommended string `json:"recommended"` // "auto", "manual", "disabled"
}

// OptimizationProfile 优化方案
type OptimizationProfile struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Services    []ServiceRecommendation `json:"services"`
}

// ServiceRecommendation 服务推荐配置
type ServiceRecommendation struct {
	Name   string `json:"name"`
	Action string `json:"action"` // "disable", "manual", "auto"
}

// ChangeResult 变更结果
type ChangeResult struct {
	Name    string `json:"name"`
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

const (
	serviceAutoStart  = 2
	serviceDemandStart = 3
	serviceDisabled   = 4
)

// GetServices 获取所有服务
func GetServices() ([]ServiceInfo, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	names, err := m.ListServices()
	if err != nil {
		return nil, fmt.Errorf("列出服务失败: %w", err)
	}

	var services []ServiceInfo
	for _, name := range names {
		s, err := m.OpenService(name)
		if err != nil {
			continue
		}

		config, err := s.Config()
		if err != nil {
			s.Close()
			continue
		}

		status, err := s.Query()
		s.Close()
		if err != nil {
			continue
		}

		svcInfo := ServiceInfo{
			Name:        name,
			DisplayName: config.DisplayName,
			Description: config.Description,
			StartType:   startTypeToString(config.StartType),
			Recommended: getRecommendedAction(name),
		}

		switch status.State {
		case svc.Running:
			svcInfo.Status = "running"
		case svc.Stopped:
			svcInfo.Status = "stopped"
		case svc.Paused:
			svcInfo.Status = "paused"
		default:
			svcInfo.Status = "unknown"
		}

		services = append(services, svcInfo)
	}

	return services, nil
}

// ChangeService 修改服务状态或启动类型
// action: "start", "stop", "restart", "disable", "manual", "auto"
func ChangeService(name string, action string) error {
	// 备份服务原始配置（仅在修改启动类型时）
	if action != "start" && action != "stop" && action != "restart" {
		if err := SaveServiceBackup(name); err != nil {
			// 备份失败只记录警告，不阻止操作
			fmt.Printf("警告: 备份服务配置失败: %v\n", err)
		}
	}

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("打开服务 %s 失败: %w", name, err)
	}
	defer s.Close()

	switch action {
	case "start":
		return s.Start()
	case "stop":
		_, err := s.Control(svc.Stop)
		return err
	case "restart":
		_, err := s.Control(svc.Stop)
		if err != nil {
			return err
		}
		return s.Start()
	default:
		config, err := s.Config()
		if err != nil {
			return err
		}
		config.StartType = stringToStartType(action)
		return s.UpdateConfig(config)
	}
}

// GetOptimizationProfiles 获取优化方案列表
func GetOptimizationProfiles() []OptimizationProfile {
	return []OptimizationProfile{
		{
			Name:        "游戏模式",
			Description: "关闭非游戏必需的后台服务，释放系统资源",
			Services:    getGameProfile(),
		},
		{
			Name:        "办公模式",
			Description: "保留办公必需服务，关闭娱乐相关服务",
			Services:    getOfficeProfile(),
		},
		{
			Name:        "开发模式",
			Description: "保留开发相关服务（Docker/IIS/远程等）",
			Services:    getDevProfile(),
		},
		{
			Name:        "节能模式",
			Description: "最大程度降低后台活动，适合笔记本用户",
			Services:    getPowerSaverProfile(),
		},
		{
			Name:        "Win10/Win11 优化",
			Description: "针对 Windows 10/11 系统优化，禁用更新、Defender、Cortana、Xbox 等，提升系统流畅度",
			Services:    getWin10OptimizeProfile(),
		},
	}
}

// ApplyOptimizationProfile 应用优化方案
func ApplyOptimizationProfile(profileName string) ([]ChangeResult, error) {
	profiles := GetOptimizationProfiles()
	var profile *OptimizationProfile
	for _, p := range profiles {
		if p.Name == profileName {
			profile = &p
			break
		}
	}
	if profile == nil {
		return nil, fmt.Errorf("未找到优化方案: %s", profileName)
	}

	var results []ChangeResult
	for _, rec := range profile.Services {
		err := ChangeService(rec.Name, rec.Action)
		results = append(results, ChangeResult{
			Name:    rec.Name,
			Action:  rec.Action,
			Success: err == nil,
			Error:   getErrorMsg(err),
		})
	}
	return results, nil
}

func getRecommendedAction(name string) string {
	recommendations := map[string]string{
		"XblAuthManager": "disabled",
		"XblGameSave":    "disabled",
		"XboxNetApiSvc":  "disabled",
		"XboxGipSvc":     "disabled",
		"DiagTrack":      "disabled",
		"WSearch":        "manual",
		"Spooler":        "auto",
		"RemoteRegistry": "disabled",
		"WMPNetworkSvc":  "manual",
		"TabletInputService": "manual",
		"lfsvc":          "manual",
		"MapsBroker":     "manual",
		"PcaSvc":         "manual",
		"WlanSvc":        "auto",
		"Dhcp":           "auto",
		"Dnscache":       "auto",
		"EventLog":       "auto",
		"Audiosrv":       "auto",
		"Themes":         "auto",
		"UsoSvc":         "manual",
		"wuauserv":       "manual",
	}
	if action, ok := recommendations[name]; ok {
		return action
	}
	return "keep"
}

func startTypeToString(st uint32) string {
	switch st {
	case serviceDisabled:
		return "disabled"
	case serviceDemandStart:
		return "manual"
	case serviceAutoStart:
		return "auto"
	default:
		return "unknown"
	}
}

func stringToStartType(s string) uint32 {
	switch strings.ToLower(s) {
	case "disabled":
		return serviceDisabled
	case "manual":
		return serviceDemandStart
	case "auto":
		return serviceAutoStart
	default:
		return serviceDemandStart
	}
}

func getErrorMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func getGameProfile() []ServiceRecommendation {
	return []ServiceRecommendation{
		{"XblAuthManager", "disabled"},
		{"XblGameSave", "disabled"},
		{"XboxNetApiSvc", "disabled"},
		{"XboxGipSvc", "disabled"},
		{"WSearch", "disabled"},
		{"DiagTrack", "disabled"},
		{"TabletInputService", "disabled"},
		{"PcaSvc", "manual"},
		{"SysMain", "disabled"},
	}
}

func getOfficeProfile() []ServiceRecommendation {
	return []ServiceRecommendation{
		{"XblAuthManager", "disabled"},
		{"XblGameSave", "disabled"},
		{"XboxNetApiSvc", "disabled"},
		{"DiagTrack", "disabled"},
	}
}

func getDevProfile() []ServiceRecommendation {
	return []ServiceRecommendation{
		{"XblAuthManager", "disabled"},
		{"RemoteRegistry", "manual"},
	}
}

func getPowerSaverProfile() []ServiceRecommendation {
	return []ServiceRecommendation{
		{"XblAuthManager", "disabled"},
		{"WSearch", "disabled"},
		{"DiagTrack", "disabled"},
		{"SysMain", "disabled"},
		{"Themes", "disabled"},
		{"TabletInputService", "disabled"},
		{"PcaSvc", "disabled"},
	}
}

func getWin10OptimizeProfile() []ServiceRecommendation {
	return []ServiceRecommendation{
		// Windows 更新服务
		{"wuauserv", "disabled"},           // Windows Update
		{"UsoSvc", "disabled"},             // Update Orchestrator
		{"WaaSMedicSvc", "disabled"},       // Windows Update Medic Service
		{"TrustedInstaller", "manual"},     // Windows Modules Installer
		// Windows Defender 相关
		{"WinDefend", "disabled"},          // Windows Defender
		{"SecurityHealthService", "disabled"}, // 安全健康服务
		{"Sense", "disabled"},              // Windows Defender 高级防护
		{"WdBoot", "disabled"},             // Defender 启动驱动
		{"WdFilter", "disabled"},           // Defender 文件过滤
		{"WdNisSvc", "disabled"},           // Defender 网络检查
		// 遥测与诊断
		{"DiagTrack", "disabled"},          // 连接用户体验和遥测
		{"dmwappushservice", "disabled"},   // 设备管理 WAP 推送
		// Xbox 相关
		{"XblAuthManager", "disabled"},
		{"XblGameSave", "disabled"},
		{"XboxNetApiSvc", "disabled"},
		{"XboxGipSvc", "disabled"},
		// 其他无用服务
		{"WSearch", "disabled"},            // Windows Search
		{"SysMain", "disabled"},            // SysMain (Superfetch)
		{"PcaSvc", "disabled"},             // 程序兼容性助手
		{"TabletInputService", "manual"},   // 触摸键盘
		{"WalletService", "disabled"},      // 钱包服务
		{"MixedRealityOpenXRSvc", "disabled"}, // 混合现实
		{"VacSvc", "disabled"},             // 音量自适应
		{"lfsvc", "manual"},                // 地理位置
		{"MapsBroker", "manual"},           // 地图下载
		{"RetailDemo", "disabled"},         // 零售演示
		{"MessagingService", "manual"},     // 消息服务
	}
}
