package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"pc-toolbox/internal/common"
	"pc-toolbox/internal/devtools"
	"pc-toolbox/internal/filetools"
	"pc-toolbox/internal/network"
	"pc-toolbox/internal/optimize"
	"pc-toolbox/internal/process"
	"pc-toolbox/internal/screenshot"
	"pc-toolbox/internal/security"
	"pc-toolbox/internal/system"
	"pc-toolbox/internal/browser"
	"pc-toolbox/internal/report"
	"pc-toolbox/internal/tray"
	"pc-toolbox/internal/clipboard"
	"pc-toolbox/internal/scheduler"
	"pc-toolbox/internal/upload"
)

// App struct
type App struct {
	ctx       context.Context
	config    *common.AppConfig
	isQuiting bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg, _ := common.LoadConfig(common.GetDefaultConfigPath())
	if cfg == nil {
		cfg = &common.AppConfig{Theme: "light", Language: "zh-CN"}
	}
	a := &App{config: cfg}
	tray.SetApp(a)
	return a
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go tray.Run()
}

// MenuShowApp 从系统托盘显示窗口
func (a *App) MenuShowApp() {
	if a.ctx != nil {
		runtime.WindowShow(a.ctx)
	}
}

// MenuQuit 退出应用
func (a *App) MenuQuit() {
	a.isQuiting = true
	if a.ctx != nil {
		runtime.Quit(a.ctx)
	}
}

// beforeClose is called when the app is about to close.
func (a *App) beforeClose(ctx context.Context) bool {
	if a.isQuiting {
		return false // 允许退出
	}
	// 关闭时最小化到托盘而不是退出
	runtime.WindowHide(ctx)
	return true // 阻止关闭，改为隐藏
}

// ============================================================
//  系统信息模块
// ============================================================

// GetSystemInfo 获取完整系统信息
func (a *App) GetSystemInfo() (*system.SystemInfo, error) {
	return system.GetSystemInfo()
}

// GetHardwareMonitor 获取实时硬件监控数据
func (a *App) GetHardwareMonitor() (*system.HardwareMonitor, error) {
	return system.GetHardwareMonitor()
}

// ============================================================
//  Windows 优化模块
// ============================================================

// ScanCleanupPaths 扫描可清理的路径
func (a *App) ScanCleanupPaths() []optimize.CleanupTarget {
	return optimize.ScanCleanupPaths()
}

// CleanTargets 执行清理
func (a *App) CleanTargets(paths []string) []optimize.CleanResult {
	return optimize.CleanTargets(paths)
}

// GetStartupItems 获取启动项
func (a *App) GetStartupItems() []optimize.StartupItem {
	return optimize.GetStartupItems()
}

// ToggleStartupItem 切换启动项状态
func (a *App) ToggleStartupItem(name string, enable bool) error {
	return optimize.ToggleStartupItem(name, enable)
}

// GetServices 获取 Windows 服务列表
func (a *App) GetServices() ([]optimize.ServiceInfo, error) {
	return optimize.GetServices()
}

// ChangeService 修改服务
func (a *App) ChangeService(name string, action string) error {
	return optimize.ChangeService(name, action)
}

// GetOptimizationProfiles 获取优化方案
func (a *App) GetOptimizationProfiles() []optimize.OptimizationProfile {
	return optimize.GetOptimizationProfiles()
}

// ApplyOptimizationProfile 应用优化方案
func (a *App) ApplyOptimizationProfile(profileName string) ([]optimize.ChangeResult, error) {
	return optimize.ApplyOptimizationProfile(profileName)
}

// RestoreService 还原服务配置
func (a *App) RestoreService(name string) error {
	return optimize.RestoreService(name)
}

// GetServiceBackups 获取已备份的服务列表
func (a *App) GetServiceBackups() []optimize.ServiceBackup {
	return optimize.GetServiceBackups()
}

// ClearServiceBackups 清除服务备份
func (a *App) ClearServiceBackups() {
	optimize.ClearServiceBackups()
}

// GetHostsEntries 获取 hosts 文件
func (a *App) GetHostsEntries() ([]optimize.HostsEntry, error) {
	return optimize.GetHostsEntries()
}

// SaveHostsEntries 保存 hosts 文件
func (a *App) SaveHostsEntries(entries []optimize.HostsEntry) error {
	return optimize.SaveHostsEntries(entries)
}

// CreateRestorePoint 创建系统还原点
func (a *App) CreateRestorePoint(description string) error {
	return optimize.CreateRestorePoint(description)
}

// ScanRegistry 扫描注册表
func (a *App) ScanRegistry() []optimize.RegistryScanResult {
	return optimize.ScanRegistry()
}

// ============================================================
//  文件工具模块
// ============================================================

// FindDuplicateFiles 查找重复文件
func (a *App) FindDuplicateFiles(rootPath string, mode string) ([]filetools.DuplicateGroup, error) {
	return filetools.FindDuplicates(rootPath, mode)
}

// ComputeFileHash 计算文件哈希
func (a *App) ComputeFileHash(path string, algorithm string) (string, error) {
	return filetools.ComputeFileHash(path, algorithm)
}

// BatchRenamePreview 批量重命名预览
func (a *App) BatchRenamePreview(dir string, rule filetools.RenameRule) ([]filetools.RenamePreview, error) {
	return filetools.BatchRenamePreview(dir, rule)
}

// BatchRename 执行批量重命名
func (a *App) BatchRename(dir string, rule filetools.RenameRule) ([]filetools.RenamePreview, error) {
	return filetools.BatchRename(dir, rule)
}

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

// ============================================================
//  安全与隐私模块
// ============================================================

// ShredFile 安全删除文件
func (a *App) ShredFile(path string, passes int) security.ShredResult {
	return security.ShredFile(path, passes)
}

// ShredDir 安全删除目录中的所有文件
func (a *App) ShredDir(dir string, passes int) []security.ShredResult {
	return security.ShredDir(dir, passes)
}

// GeneratePassword 生成密码
func (a *App) GeneratePassword(length int, useUpper bool, useLower bool, useDigits bool, useSpecial bool) security.PasswordResult {
	return security.GeneratePassword(length, useUpper, useLower, useDigits, useSpecial)
}

// EncryptFile 加密文件
func (a *App) EncryptFile(inputPath string, password string) security.EncryptResult {
	return security.EncryptFile(inputPath, password)
}

// DecryptFile 解密文件
func (a *App) DecryptFile(inputPath string, password string) security.EncryptResult {
	return security.DecryptFile(inputPath, password)
}

// ClearRecentDocs 清理最近文档记录
func (a *App) ClearRecentDocs() error {
	return security.ClearRecentDocs()
}

// ============================================================
//  开发工具模块
// ============================================================

// FormatJSON 格式化 JSON
func (a *App) FormatJSON(input string) devtools.JSONResult {
	return devtools.FormatJSON(input)
}

// MinifyJSON 压缩 JSON
func (a *App) MinifyJSON(input string) devtools.JSONResult {
	return devtools.MinifyJSON(input)
}

// DiffText 文本差异对比
func (a *App) DiffText(oldText, newText string) []devtools.DiffResult {
	return devtools.DiffText(oldText, newText)
}

// EncodeBase64 Base64 编码
func (a *App) EncodeBase64(input string) devtools.CodecResult {
	return devtools.EncodeBase64(input)
}

// DecodeBase64 Base64 解码
func (a *App) DecodeBase64(input string) devtools.CodecResult {
	return devtools.DecodeBase64(input)
}

// ReadFileAsBase64 读取文件并返回 Base64 编码
func (a *App) ReadFileAsBase64(path string) string {
	return common.ReadFileAsBase64(path)
}

// EncodeURL URL 编码
func (a *App) EncodeURL(input string) devtools.CodecResult {
	return devtools.EncodeURL(input)
}

// DecodeURL URL 解码
func (a *App) DecodeURL(input string) devtools.CodecResult {
	return devtools.DecodeURL(input)
}

// TestRegex 测试正则表达式
func (a *App) TestRegex(pattern, text string) devtools.RegexTestResult {
	return devtools.TestRegex(pattern, text)
}

// ConvertTimestamp 转换时间戳
func (a *App) ConvertTimestamp(timestamp int64, fromUnit string) devtools.TimestampResult {
	return devtools.ConvertTimestamp(timestamp, fromUnit)
}

// ConvertColor 颜色值转换
func (a *App) ConvertColor(hex string) devtools.ColorResult {
	return devtools.ConvertColor(hex)
}

// GenerateUUID 生成 UUID
func (a *App) GenerateUUID() string {
	return devtools.GenerateUUID()
}

// ============================================================
//  截屏工具
// ============================================================

// CaptureAllScreens 截取所有屏幕
func (a *App) CaptureAllScreens() screenshot.CaptureResult {
	return screenshot.CaptureAllScreens()
}

// CaptureScreen 截取指定屏幕
func (a *App) CaptureScreen(displayIndex int) screenshot.CaptureResult {
	return screenshot.CaptureScreen(displayIndex)
}

// OpenInExplorer 在资源管理器中打开文件
func (a *App) OpenInExplorer(path string) error {
	return screenshot.OpenInExplorer(path)
}

// ============================================================
//  进程管理
// ============================================================

// GetProcessList 获取进程列表
func (a *App) GetProcessList() ([]process.ProcessInfo, error) {
	return process.GetProcessList()
}

// KillProcess 结束进程
func (a *App) KillProcess(pid int) error {
	return process.KillProcess(pid)
}

// ============================================================
//  浏览器数据管理
// ============================================================

// DetectBrowsers 检测浏览器
func (a *App) DetectBrowsers() []browser.Browser {
	return browser.DetectBrowsers()
}

// ScanBrowserData 扫描浏览器可清理数据
func (a *App) ScanBrowserData() []browser.CleanItem {
	return browser.ScanBrowserData()
}

// CleanBrowserData 清理浏览器数据
func (a *App) CleanBrowserData(items []browser.CleanItem) []browser.CleanResult {
	return browser.CleanBrowserData(items)
}

// ============================================================
//  系统报告
// ============================================================

// CollectReport 收集系统报告
func (a *App) CollectReport() *report.SystemReport {
	return report.CollectReport()
}

// GenerateHTMLReport 生成 HTML 报告
func (a *App) GenerateHTMLReport() (string, error) {
	r := report.CollectReport()
	return report.GenerateHTMLReport(r)
}

// OpenInBrowser 在浏览器中打开
func (a *App) OpenReportInBrowser(path string) error {
	return report.OpenInBrowser(path)
}

// ============================================================
//  网络流量监控
// ============================================================

// GetTrafficSamples 获取流量采样
func (a *App) GetTrafficSamples(durationSec int) ([]network.TrafficSample, error) {
	return network.GetTrafficSamples(durationSec)
}

// ============================================================
//  文件归类
// ============================================================

// PreviewOrganize 预览归类
func (a *App) PreviewOrganize(dir string, rule filetools.OrganizeRule) ([]filetools.OrganizePreview, error) {
	return filetools.PreviewOrganize(dir, rule)
}

// ExecuteOrganize 执行归类
func (a *App) ExecuteOrganize(dir string, rule filetools.OrganizeRule) ([]filetools.OrganizeResult, error) {
	return filetools.ExecuteOrganize(dir, rule)
}

// ============================================================
//  IP 地理位置查询
// ============================================================

// QueryGeoIP 查询 IP 地理位置
func (a *App) QueryGeoIP(ip string) network.GeoIPResult {
	return network.QueryGeoIP(ip)
}

// ============================================================
//  局域网设备扫描
// ============================================================

// ScanLAN 扫描局域网设备
func (a *App) ScanLAN() ([]network.LANDevice, error) {
	return network.ScanLAN()
}

// ============================================================
//  系统健康体检
// ============================================================

// RunHealthCheck 执行健康体检
func (a *App) RunHealthCheck() *optimize.HealthReport {
	return optimize.RunHealthCheck()
}

// ============================================================
//  右键菜单扩展
// ============================================================

// InstallContextMenu 安装右键菜单
func (a *App) InstallContextMenu() optimize.ContextMenuStatus {
	return optimize.InstallContextMenu()
}

// InstallContextMenuDir 安装文件夹右键菜单
func (a *App) InstallContextMenuDir() optimize.ContextMenuStatus {
	return optimize.InstallContextMenuDir()
}

// UninstallContextMenu 卸载右键菜单
func (a *App) UninstallContextMenu() optimize.ContextMenuStatus {
	return optimize.UninstallContextMenu()
}

// CheckContextMenu 检查右键菜单状态
func (a *App) CheckContextMenu() optimize.ContextMenuStatus {
	return optimize.CheckContextMenu()
}

// ============================================================
//  Windows 更新管理
// ============================================================

// GetWindowsUpdates 获取已安装更新
func (a *App) GetWindowsUpdates() []optimize.UpdateInfo {
	return optimize.GetWindowsUpdates()
}

// GetPendingUpdates 获取待安装更新
func (a *App) GetPendingUpdates() []optimize.UpdateInfo {
	return optimize.GetPendingUpdates()
}

// ============================================================
//  剪贴板历史
// ============================================================

// AddClipboardItem 添加剪贴板记录
func (a *App) AddClipboardItem(content string, contentType string) clipboard.ClipItem {
	return clipboard.AddItem(content, contentType)
}

// GetClipboardHistory 获取剪贴板历史
func (a *App) GetClipboardHistory() []clipboard.ClipItem {
	return clipboard.GetHistory()
}

// ClearClipboardHistory 清空剪贴板历史
func (a *App) ClearClipboardHistory() {
	clipboard.ClearHistory()
}

// RemoveClipboardItem 删除单条记录
func (a *App) RemoveClipboardItem(id int) {
	clipboard.RemoveItem(id)
}

// ============================================================
//  定时任务
// ============================================================

// CreateScheduledTask 创建定时任务
func (a *App) CreateScheduledTask(action string, hour int, minute int) scheduler.TaskInfo {
	return scheduler.CreateTask(scheduler.TaskAction(action), hour, minute)
}

// ListScheduledTasks 列出定时任务
func (a *App) ListScheduledTasks() []scheduler.TaskInfo {
	return scheduler.ListTasks()
}

// DeleteScheduledTask 删除定时任务
func (a *App) DeleteScheduledTask(name string) bool {
	return scheduler.DeleteTask(name)
}

// ============================================================
//  字体管理
// ============================================================

// ============================================================
//  WiFi 密码
// ============================================================

// GetWiFiPasswords 获取 WiFi 密码列表
func (a *App) GetWiFiPasswords() []network.WiFiProfile {
	return network.GetWiFiPasswords()
}

// ============================================================
//  远程桌面
// ============================================================

// LaunchMSTSC 启动远程桌面
func (a *App) LaunchMSTSC(computer string, address string, port int) error {
	return network.LaunchMSTSC(computer, address, port)
}

// ============================================================
//  系统备份与还原
// ============================================================

// GetRestorePoints 获取还原点列表
func (a *App) GetRestorePoints() ([]optimize.RestorePointInfo, error) {
	return optimize.GetRestorePoints()
}

// ============================================================
//  大文件查找
// ============================================================

// FindLargeFiles 查找大文件
func (a *App) FindLargeFiles(rootPath string, minSizeMB int, maxCount int) ([]filetools.LargeFile, error) {
	return filetools.FindLargeFiles(rootPath, minSizeMB, maxCount)
}

// ============================================================
//  Windows 激活信息
// ============================================================

// GetActivationInfo 获取激活信息
func (a *App) GetActivationInfo() system.ActivationInfo {
	return system.GetActivationInfo()
}

// GetActivationTools 获取激活工具列表
func (a *App) GetActivationTools() []system.ActivationTool {
	return system.GetActivationTools()
}

// GetKMSMethods 获取 KMS 激活方法
func (a *App) GetKMSMethods() []system.KMSMethod {
	return system.GetKMSMethods()
}

// OpenExternalURL 在系统浏览器中打开 URL
func (a *App) OpenExternalURL(url string) error {
	return common.OpenURL(url)
}

// GetAppVersion 获取应用版本号
func (a *App) GetAppVersion() string {
	return AppVersion
}

// GetBuildDate 获取构建日期
func (a *App) GetBuildDate() string {
	return BuildDate
}

// ============================================================
//  批量 Ping
// ============================================================

// BatchPing 批量 Ping
func (a *App) BatchPing(cidr string, timeout int) []network.BatchPingResult {
	return network.BatchPing(cidr, timeout)
}

// ============================================================
//  图片批量压缩
// ============================================================

// BatchCompressImages 批量压缩图片
func (a *App) BatchCompressImages(dir string, quality int, targetFormat string, maxWidth int) []filetools.BatchCompressResult {
	return filetools.BatchCompressImages(dir, quality, targetFormat, maxWidth)
}

// ConvertToPDF 转换文档到 PDF
func (a *App) ConvertToPDF(inputPath string) filetools.DocConvertResult {
	return filetools.ConvertToPDF(inputPath)
}

// ============================================================
//  文件内容搜索
// ============================================================

// SearchFileContent 搜索文件内容
func (a *App) SearchFileContent(rootDir string, keyword string, fileTypes string) []filetools.SearchResult {
	return filetools.SearchFileContent(rootDir, keyword, fileTypes)
}

// ============================================================
//  HTTP API 调试
// ============================================================

// SendHTTPRequest 发送 HTTP 请求
func (a *App) SendHTTPRequest(req devtools.HTTPRequest) devtools.HTTPResponse {
	return devtools.SendHTTPRequest(req)
}

// ============================================================
//  WiFi 信号扫描
// ============================================================

// ScanWiFiSignal 扫描 WiFi 信号
func (a *App) ScanWiFiSignal() []network.SignalInfo {
	return network.ScanWiFiSignal()
}

// CheckIPConflict 检测 IP 冲突
func (a *App) CheckIPConflict(localIP string) []string {
	return network.CheckIPConflict(localIP)
}

// ============================================================
//  JWT 解码
// ============================================================

// DecodeJWT 解码 JWT
func (a *App) DecodeJWT(token string) devtools.JWTResult {
	return devtools.DecodeJWT(token)
}

// ============================================================
//  YAML/TOML 格式化
// ============================================================

// FormatYAML 格式化 YAML
func (a *App) FormatYAML(input string) devtools.FormatResult {
	return devtools.FormatYAML(input)
}

// FormatTOML 格式化 TOML
func (a *App) FormatTOML(input string) devtools.FormatResult {
	return devtools.FormatTOML(input)
}

// ============================================================
//  UUID 生成器
// ============================================================

// GenerateUUIDs 批量生成 UUID
func (a *App) GenerateUUIDs(count int, version int) devtools.UUIDGenResult {
	return devtools.GenerateUUIDs(count, version)
}

// ============================================================
//  自动更新
// ============================================================

// CheckUpdate 检查更新
func (a *App) CheckUpdate(currentVersion string) devtools.ReleaseInfo {
	return devtools.CheckUpdate(currentVersion)
}

// ============================================================
//  文件内容替换
// ============================================================
func (a *App) SearchAndReplace(dir string, search string, replace string, fileTypes string) []filetools.ReplaceResult {
	return filetools.SearchAndReplace(dir, search, replace, fileTypes)
}

// ============================================================
//  文件夹大小分析
// ============================================================
func (a *App) AnalyzeFolderSizes(rootPath string, depth int) []filetools.FolderSize {
	return filetools.AnalyzeFolderSizes(rootPath, depth)
}

// ============================================================
//  文件上传
// ============================================================

// UploadFileToServer 将 Base64 文件上传到指定服务器
func (a *App) UploadFileToServer(fileData string, fileName string, serverURL string, fieldName string) upload.UploadResult {
	return upload.UploadFileToServer(fileData, fileName, serverURL, fieldName)
}

func (a *App) GetRecycleBinInfo() filetools.RecycleBinInfo {
	return filetools.GetRecycleBinInfo()
}
func (a *App) EmptyRecycleBin() error {
	return filetools.EmptyRecycleBin()
}

// ============================================================
//  文件差异对比
// ============================================================
func (a *App) DiffFiles(oldPath string, newPath string) ([]filetools.DiffLine, error) {
	return filetools.DiffFiles(oldPath, newPath)
}

// ============================================================
//  配置管理
// ============================================================

// GetConfig 获取当前配置
func (a *App) GetConfig() *common.AppConfig {
	return a.config
}

// SaveConfig 保存配置
func (a *App) SaveConfig(cfg *common.AppConfig) error {
	a.config = cfg
	return common.SaveConfig(common.GetDefaultConfigPath(), cfg)
}

// ============================================================
//  电源方案切换
// ============================================================

// GetPowerPlans 获取电源方案列表
func (a *App) GetPowerPlans() []system.PowerPlan {
	return system.GetPowerPlans()
}

// SetPowerPlan 切换电源方案
func (a *App) SetPowerPlan(guid string) error {
	return system.SetPowerPlan(guid)
}

// ============================================================
//  图片格式转换
// ============================================================

// ConvertImage 转换图片格式
func (a *App) ConvertImage(inputPath string, targetFormat string) filetools.ConvertResult {
	return filetools.ConvertImage(inputPath, targetFormat)
}

// ============================================================
//  文本编码转换
// ============================================================

// ConvertEncoding 转换文本编码
func (a *App) ConvertEncoding(text string, fromCharset string, toCharset string) filetools.EncodingResult {
	return filetools.ConvertEncoding(text, fromCharset, toCharset)
}

// ============================================================
//  代码美化
// ============================================================

// BeautifyHTML 美化 HTML
func (a *App) BeautifyHTML(input string) devtools.CodeBeautifyResult {
	return devtools.BeautifyHTML(input)
}

// BeautifyCSS 美化 CSS
func (a *App) BeautifyCSS(input string) devtools.CodeBeautifyResult {
	return devtools.BeautifyCSS(input)
}

// BeautifySQL 美化 SQL
func (a *App) BeautifySQL(input string) devtools.CodeBeautifyResult {
	return devtools.BeautifySQL(input)
}

// ============================================================
//  二维码生成
// ============================================================

// GenerateQRCode 生成二维码
func (a *App) GenerateQRCode(content string, size int) devtools.QRResult {
	return devtools.GenerateQRCode(content, size)
}

// ============================================================
//  文件对话框（前端辅助）
// ============================================================

// SelectDirectory 选择目录对话框
func (a *App) SelectDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择目录",
	})
	if err != nil {
		return ""
	}
	return dir
}

// SelectFile 选择文件对话框
func (a *App) SelectFile() string {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
	})
	if err != nil {
		return ""
	}
	return file
}
