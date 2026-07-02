package main

import (
	"pc-toolbox/internal/network"
)

// ============================================================
//  网络工具模块
// ============================================================

// Ping 执行 Ping
func (a *App) Ping(host string, count int, timeout int) (*network.PingSummary, error) {
	return network.Ping(host, count, timeout)
}

// PortScan 端口扫描
func (a *App) PortScan(host string, ports string) ([]network.PortResult, error) {
	return network.PortScan(host, ports)
}

// DNSLookup DNS 查询
func (a *App) DNSLookup(hostname string) (*network.DNSResult, error) {
	return network.DNSLookup(hostname)
}

// GetNetworkConnections 获取网络连接
func (a *App) GetNetworkConnections() ([]network.ConnectionInfo, error) {
	return network.GetNetworkConnections()
}

// ============================================================
//  网络修复
// ============================================================

// FlushDNS 刷新 DNS 缓存
func (a *App) FlushDNSCache() network.FixResult {
	return network.FlushDNS()
}

// ResetWinsock 重置 Winsock
func (a *App) ResetWinsock() network.FixResult {
	return network.ResetWinsock()
}

// ResetTCPIP 重置 TCP/IP
func (a *App) ResetTCPIP() network.FixResult {
	return network.ResetTCPIP()
}

// ReleaseIP 释放 IP
func (a *App) ReleaseIP() network.FixResult {
	return network.ReleaseIP()
}

// RenewIP 续租 IP
func (a *App) RenewIP() network.FixResult {
	return network.RenewIP()
}

// DiagnoseNetwork 快速诊断
func (a *App) DiagnoseNetwork() []network.FixResult {
	return network.DiagnoseNetwork()
}

// FixAllNetwork 一键网络修复
func (a *App) FixAllNetwork() []network.FixResult {
	return network.FixAll()
}

// ResetProxy 清除代理
func (a *App) ResetProxy() network.FixResult {
	return network.ResetProxy()
}

// ResetArpCache 清空 ARP 缓存
func (a *App) ResetArpCache() network.FixResult {
	return network.ResetArp()
}

// GetTrafficSamples 获取流量采样
func (a *App) GetTrafficSamples(durationSec int) ([]network.TrafficSample, error) {
	return network.GetTrafficSamples(durationSec)
}

// QueryGeoIP 查询 IP 地理位置
func (a *App) QueryGeoIP(ip string) network.GeoIPResult {
	return network.QueryGeoIP(ip)
}

// ScanLAN 扫描局域网设备
func (a *App) ScanLAN() ([]network.LANDevice, error) {
	return network.ScanLAN()
}

// GetWiFiPasswords 获取 WiFi 密码列表
func (a *App) GetWiFiPasswords() []network.WiFiProfile {
	return network.GetWiFiPasswords()
}

// LaunchMSTSC 启动远程桌面
func (a *App) LaunchMSTSC(computer string, address string, port int) error {
	return network.LaunchMSTSC(computer, address, port)
}

// BatchPing 批量 Ping
func (a *App) BatchPing(cidr string, timeout int) []network.BatchPingResult {
	return network.BatchPing(cidr, timeout)
}

// ScanWiFiSignal 扫描 WiFi 信号
func (a *App) ScanWiFiSignal() []network.SignalInfo {
	return network.ScanWiFiSignal()
}

// CheckIPConflict 检测 IP 冲突
func (a *App) CheckIPConflict(localIP string) []string {
	return network.CheckIPConflict(localIP)
}

// ============================================================
//  网络连接查看器
// ============================================================

// GetAllNetConnections 获取所有网络连接
func (a *App) GetAllNetConnections() []network.NetConnection {
	conns, _ := network.GetAllConnections()
	return conns
}

// GetEstablishedConnections 获取已建立的连接
func (a *App) GetEstablishedConnections() []network.NetConnection {
	return network.GetEstablishedConnections()
}

// GetListeningConnections 获取监听中的连接
func (a *App) GetListeningConnections() []network.NetConnection {
	return network.GetListeningConnections()
}

// SearchNetConnections 搜索网络连接
func (a *App) SearchNetConnections(keyword string) []network.NetConnection {
	return network.SearchConnections(keyword)
}

// KillNetConnection 终止网络连接
func (a *App) KillNetConnection(pid int) error {
	return network.KillConnection(pid)
}

// GetNetConnectionStats 获取网络连接统计
func (a *App) GetNetConnectionStats() map[string]int {
	return network.GetConnectionStats()
}
