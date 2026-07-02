package optimize

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/svc/mgr"
)

// ServiceDependency 服务依赖关系
type ServiceDependency struct {
	Name         string   `json:"name"`
	DisplayName  string   `json:"displayName"`
	DependsOn    []string `json:"dependsOn"`    // 依赖的服务
	DependedBy   []string `json:"dependedBy"`   // 被依赖的服务
	DepType      string   `json:"depType"`      // "required", "optional"
}

// ServiceDependencyGraph 服务依赖关系图
type ServiceDependencyGraph struct {
	Services    []ServiceDependency `json:"services"`
	Isolated    []string            `json:"isolated"`    // 无依赖的服务
	Critical    []string            `json:"critical"`    // 被多个服务依赖的关键服务
	Orphaned    []string            `json:"orphaned"`    // 依赖不存在的服务
}

// GetServiceDependencies 获取服务依赖关系
func GetServiceDependencies(serviceName string) (*ServiceDependency, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("打开服务 %s 失败: %w", serviceName, err)
	}
	defer s.Close()

	config, err := s.Config()
	if err != nil {
		return nil, fmt.Errorf("获取服务配置失败: %w", err)
	}

	dep := &ServiceDependency{
		Name:        serviceName,
		DisplayName: config.DisplayName,
		DependsOn:   config.Dependencies,
		DepType:     "required",
	}

	return dep, nil
}

// GetServiceDependencyGraph 获取所有服务的依赖关系图
func GetServiceDependencyGraph() (*ServiceDependencyGraph, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	names, err := m.ListServices()
	if err != nil {
		return nil, fmt.Errorf("列出服务失败: %w", err)
	}

	graph := &ServiceDependencyGraph{
		Services: make([]ServiceDependency, 0),
	}

	// 收集所有服务名称
	serviceSet := make(map[string]bool)
	for _, name := range names {
		serviceSet[name] = true
	}

	// 构建依赖关系
	for _, name := range names {
		s, err := m.OpenService(name)
		if err != nil {
			continue
		}

		config, err := s.Config()
		s.Close()
		if err != nil {
			continue
		}

		dep := ServiceDependency{
			Name:        name,
			DisplayName: config.DisplayName,
			DependsOn:   config.Dependencies,
			DepType:     "required",
		}

		graph.Services = append(graph.Services, dep)
	}

	// 分析依赖关系
	graph.Isolated = findIsolatedServices(graph.Services, serviceSet)
	graph.Critical = findCriticalServices(graph.Services)
	graph.Orphaned = findOrphanedServices(graph.Services, serviceSet)

	return graph, nil
}

// findIsolatedServices 查找无依赖的服务
func findIsolatedServices(services []ServiceDependency, serviceSet map[string]bool) []string {
	var isolated []string
	for _, svc := range services {
		if len(svc.DependsOn) == 0 {
			isolated = append(isolated, svc.Name)
		}
	}
	return isolated
}

// findCriticalServices 查找被多个服务依赖的关键服务
func findCriticalServices(services []ServiceDependency) []string {
	depCount := make(map[string]int)
	for _, svc := range services {
		for _, dep := range svc.DependsOn {
			depCount[dep]++
		}
	}

	var critical []string
	for name, count := range depCount {
		if count >= 3 { // 被3个或更多服务依赖
			critical = append(critical, name)
		}
	}
	return critical
}

// findOrphanedServices 查找依赖不存在服务的孤儿服务
func findOrphanedServices(services []ServiceDependency, serviceSet map[string]bool) []string {
	var orphaned []string
	for _, svc := range services {
		for _, dep := range svc.DependsOn {
			if !serviceSet[dep] {
				orphaned = append(orphaned, svc.Name)
				break
			}
		}
	}
	return orphaned
}

// AnalyzeServiceImpact 分析禁用服务的影响
func AnalyzeServiceImpact(serviceName string) ([]string, error) {
	graph, err := GetServiceDependencyGraph()
	if err != nil {
		return nil, err
	}

	// 找到所有直接依赖此服务的服务
	var affected []string
	for _, svc := range graph.Services {
		for _, dep := range svc.DependsOn {
			if strings.EqualFold(dep, serviceName) {
				affected = append(affected, svc.Name)
			}
		}
	}

	// 递归查找间接依赖
	visited := make(map[string]bool)
	var dfs func(name string)
	dfs = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true

		for _, svc := range graph.Services {
			for _, dep := range svc.DependsOn {
				if strings.EqualFold(dep, name) {
					affected = append(affected, svc.Name)
					dfs(svc.Name)
				}
			}
		}
	}

	for _, name := range affected {
		dfs(name)
	}

	// 去重
	seen := make(map[string]bool)
	var unique []string
	for _, name := range affected {
		if !seen[name] {
			seen[name] = true
			unique = append(unique, name)
		}
	}

	return unique, nil
}

// CanDisableService 检查是否可以安全禁用服务
func CanDisableService(serviceName string) (bool, []string, error) {
	affected, err := AnalyzeServiceImpact(serviceName)
	if err != nil {
		return false, nil, err
	}

	if len(affected) == 0 {
		return true, nil, nil
	}

	return false, affected, nil
}

// GetServiceStartupType 获取服务启动类型
func GetServiceStartupType(serviceName string) (string, error) {
	m, err := mgr.Connect()
	if err != nil {
		return "", fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return "", fmt.Errorf("打开服务 %s 失败: %w", serviceName, err)
	}
	defer s.Close()

	config, err := s.Config()
	if err != nil {
		return "", fmt.Errorf("获取服务配置失败: %w", err)
	}

	switch config.StartType {
	case 0:
		return "boot", nil
	case 1:
		return "system", nil
	case 2:
		return "auto", nil
	case 3:
		return "manual", nil
	case 4:
		return "disabled", nil
	default:
		return "unknown", nil
	}
}

// SetServiceStartupType 设置服务启动类型
func SetServiceStartupType(serviceName string, startupType string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("打开服务 %s 失败: %w", serviceName, err)
	}
	defer s.Close()

	config, err := s.Config()
	if err != nil {
		return fmt.Errorf("获取服务配置失败: %w", err)
	}

	switch strings.ToLower(startupType) {
	case "auto", "automatic":
		config.StartType = 2
	case "manual":
		config.StartType = 3
	case "disabled":
		config.StartType = 4
	case "boot":
		config.StartType = 0
	case "system":
		config.StartType = 1
	default:
		return fmt.Errorf("无效的启动类型: %s", startupType)
	}

	if err := s.UpdateConfig(config); err != nil {
		return fmt.Errorf("更新服务配置失败: %w", err)
	}

	return nil
}

// GetServiceDescription 获取服务描述
func GetServiceDescription(serviceName string) (string, error) {
	m, err := mgr.Connect()
	if err != nil {
		return "", fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return "", fmt.Errorf("打开服务 %s 失败: %w", serviceName, err)
	}
	defer s.Close()

	config, err := s.Config()
	if err != nil {
		return "", fmt.Errorf("获取服务配置失败: %w", err)
	}

	return config.Description, nil
}

// GetServiceBinaryPath 获取服务二进制路径
func GetServiceBinaryPath(serviceName string) (string, error) {
	m, err := mgr.Connect()
	if err != nil {
		return "", fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return "", fmt.Errorf("打开服务 %s 失败: %w", serviceName, err)
	}
	defer s.Close()

	config, err := s.Config()
	if err != nil {
		return "", fmt.Errorf("获取服务配置失败: %w", err)
	}

	return config.BinaryPathName, nil
}

// IsServiceRunning 检查服务是否运行中
func IsServiceRunning(serviceName string) (bool, error) {
	m, err := mgr.Connect()
	if err != nil {
		return false, fmt.Errorf("连接服务管理器失败: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return false, fmt.Errorf("打开服务 %s 失败: %w", serviceName, err)
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return false, fmt.Errorf("查询服务状态失败: %w", err)
	}

	return status.State == 4, nil // 4 = Running
}
