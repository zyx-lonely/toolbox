import { createRouter, createWebHashHistory, type RouteRecordRaw } from 'vue-router'
import MainLayout from '../layouts/MainLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: MainLayout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '首页', icon: 'grid-outline' }
      },
      {
        path: 'system',
        name: 'System',
        component: () => import('../views/system/SystemInfo.vue'),
        meta: { title: '系统信息', icon: 'hardware-chip-outline' }
      },
      {
        path: 'optimize/disk',
        name: 'DiskClean',
        component: () => import('../views/optimize/DiskClean.vue'),
        meta: { title: '磁盘清理', icon: 'trash-outline' }
      },
      {
        path: 'optimize/startup',
        name: 'StartupManager',
        component: () => import('../views/optimize/StartupManager.vue'),
        meta: { title: '启动项管理', icon: 'power-outline' }
      },
      {
        path: 'optimize/service',
        name: 'ServiceOptimize',
        component: () => import('../views/optimize/ServiceOptimize.vue'),
        meta: { title: '系统服务优化', icon: 'settings-outline' }
      },
      {
        path: 'optimize/browser',
        name: 'BrowserClean',
        component: () => import('../views/optimize/BrowserClean.vue'),
        meta: { title: '浏览器数据管理', icon: 'globe-outline' }
      },
      {
        path: 'optimize/health',
        name: 'HealthCheck',
        component: () => import('../views/optimize/HealthCheck.vue'),
        meta: { title: '系统健康体检', icon: 'heart-outline' }
      },
      {
        path: 'optimize/restore',
        name: 'SystemRestore',
        component: () => import('../views/optimize/SystemRestore.vue'),
        meta: { title: '系统备份与还原', icon: 'shield-outline' }
      },
      {
        path: 'optimize/registry',
        name: 'RegistryClean',
        component: () => import('../views/optimize/RegistryClean.vue'),
        meta: { title: '注册表清理', icon: 'analytics-outline' }
      },
      {
        path: 'system/report',
        name: 'HardwareReport',
        component: () => import('../views/system/HardwareReport.vue'),
        meta: { title: '硬件检测报告', icon: 'document-text-outline' }
      },
      {
        path: 'system/activation',
        name: 'Activation',
        component: () => import('../views/system/Activation.vue'),
        meta: { title: '系统激活', icon: 'key-outline' }
      },
      {
        path: 'system/powerplan',
        name: 'PowerPlan',
        component: () => import('../views/system/PowerPlan.vue'),
        meta: { title: '电源方案', icon: 'flash-outline' }
      },
      {
        path: 'system/process',
        name: 'ProcessManager',
        component: () => import('../views/system/ProcessManager.vue'),
        meta: { title: '进程管理器', icon: 'pulse-outline' }
      },
      {
        path: 'network',
        name: 'NetworkTools',
        component: () => import('../views/network/NetworkTools.vue'),
        meta: { title: '网络工具', icon: 'globe-outline' }
      },
      {
        path: 'network/hosts',
        name: 'HostsEditor',
        component: () => import('../views/network/HostsEditor.vue'),
        meta: { title: 'Hosts 编辑器', icon: 'document-text-outline' }
      },
      {
        path: 'network/geoip',
        name: 'GeoIP',
        component: () => import('../views/network/GeoIP.vue'),
        meta: { title: 'IP 地理位置', icon: 'location-outline' }
      },
      {
        path: 'network/lanscan',
        name: 'LANScan',
        component: () => import('../views/network/LANScan.vue'),
        meta: { title: '局域网扫描', icon: 'wifi-outline' }
      },
      {
        path: 'network/wifi',
        name: 'WiFiPasswords',
        component: () => import('../views/network/WiFiPasswords.vue'),
        meta: { title: 'WiFi 密码', icon: 'lock-closed-outline' }
      },
      { path: 'network/wifiscan', name: 'WiFiScan', component: () => import('../views/network/WiFiScan.vue'), meta: { title: 'WiFi 扫描' } },
      { path: 'network/ipconflict', name: 'IPConflict', component: () => import('../views/network/IPConflict.vue'), meta: { title: 'IP 冲突' } },
      { path: 'network/topology', name: 'LANTopology', component: () => import('../views/network/LANTopology.vue'), meta: { title: '局域网拓扑' } },
      {
        path: 'network/rdp',
        name: 'RemoteDesktop',
        component: () => import('../views/network/RemoteDesktop.vue'),
        meta: { title: '远程桌面', icon: 'desktop-outline' }
      },
      {
        path: 'network/batchping',
        name: 'BatchPing',
        component: () => import('../views/network/BatchPing.vue'),
        meta: { title: '批量 Ping', icon: 'pulse-outline' }
      },
      {
        path: 'network/traffic',
        name: 'TrafficChart',
        component: () => import('../views/network/TrafficChart.vue'),
        meta: { title: '流量图表', icon: 'stats-chart-outline' }
      },
      {
        path: 'file/duplicate',
        name: 'DuplicateFinder',
        component: () => import('../views/file/DuplicateFinder.vue'),
        meta: { title: '重复文件查找', icon: 'copy-outline' }
      },
      {
        path: 'file/organize',
        name: 'FileOrganize',
        component: () => import('../views/file/FileOrganize.vue'),
        meta: { title: '文件批量归类', icon: 'folder-open-outline' }
      },
      {
        path: 'file/largefiles',
        name: 'LargeFiles',
        component: () => import('../views/file/LargeFiles.vue'),
        meta: { title: '大文件查找', icon: 'folder-outline' }
      },
      {
        path: 'file/compress',
        name: 'ImageCompress',
        component: () => import('../views/file/ImageCompress.vue'),
        meta: { title: '图片压缩', icon: 'compress-outline' }
      },
      {
        path: 'file/search',
        name: 'ContentSearch',
        component: () => import('../views/file/ContentSearch.vue'),
        meta: { title: '文件搜索', icon: 'search-outline' }
      },
      {
        path: 'file/replace',
        name: 'ReplaceContent',
        component: () => import('../views/file/ReplaceContent.vue'),
        meta: { title: '文件内容替换', icon: 'swap-horizontal-outline' }
      },
      {
        path: 'file/foldersize',
        name: 'FolderSize',
        component: () => import('../views/file/FolderSize.vue'),
        meta: { title: '文件夹大小分析', icon: 'analytics-outline' }
      },
      {
        path: 'file/filediff',
        name: 'FileDiff',
        component: () => import('../views/file/FileDiff.vue'),
        meta: { title: '文件差异对比', icon: 'document-text-outline' }
      },
      {
        path: 'security/encrypt',
        name: 'FileEncrypt',
        component: () => import('../views/security/FileEncrypt.vue'),
        meta: { title: '文件加密解密', icon: 'lock-closed-outline' }
      },
      { path: 'file/pdf', name: 'PDFTools', component: () => import('../views/file/PDFTools.vue'), meta: { title: 'PDF 工具' } },
      { path: 'file/rename', name: 'BatchRename', component: () => import('../views/file/BatchRename.vue'), meta: { title: '批量重命名' } },
      {
        path: 'file/unlock',
        name: 'FileUnlocker',
        component: () => import('../views/file/FileUnlocker.vue'),
        meta: { title: '文件解锁器', icon: 'lock-open-outline' }
      },
      {
        path: 'file/imageconv',
        name: 'ImageConverter',
        component: () => import('../views/file/ImageConverter.vue'),
        meta: { title: '图片格式转换', icon: 'image-outline' }
      },
      {
        path: 'file/encoding',
        name: 'EncodingConverter',
        component: () => import('../views/file/EncodingConverter.vue'),
        meta: { title: '编码转换', icon: 'text-outline' }
      },
      {
        path: 'file/hash',
        name: 'FileHash',
        component: () => import('../views/file/FileHash.vue'),
        meta: { title: '文件哈希校验', icon: 'finger-print-outline' }
      },
      {
        path: 'security/shred',
        name: 'FileShred',
        component: () => import('../views/security/FileShred.vue'),
        meta: { title: '文件粉碎', icon: 'trash-outline' }
      },
      {
        path: 'security/password',
        name: 'PasswordGen',
        component: () => import('../views/security/PasswordGen.vue'),
        meta: { title: '密码生成器', icon: 'key-outline' }
      },
      {
        path: 'devtools/json',
        name: 'JsonFormatter',
        component: () => import('../views/devtools/JsonFormatter.vue'),
        meta: { title: 'JSON 格式化', icon: 'code-slash-outline' }
      },
      {
        path: 'devtools/codec',
        name: 'CodecTool',
        component: () => import('../views/devtools/CodecTool.vue'),
        meta: { title: '编解码工具', icon: 'swap-horizontal-outline' }
      },
      {
        path: 'devtools/beautify',
        name: 'CodeBeautifier',
        component: () => import('../views/devtools/CodeBeautifier.vue'),
        meta: { title: '代码美化', icon: 'color-wand-outline' }
      },
      {
        path: 'devtools/http',
        name: 'HTTPDebugger',
        component: () => import('../views/devtools/HTTPDebugger.vue'),
        meta: { title: 'HTTP 调试', icon: 'globe-outline' }
      },
      {
        path: 'devtools/regex',
        name: 'RegexTester',
        component: () => import('../views/devtools/RegexTester.vue'),
        meta: { title: '正则测试', icon: 'document-text-outline' }
      },
      { path: 'devtools/jwt', name: 'JWTDecoder', component: () => import('../views/devtools/JWTDecoder.vue'), meta: { title: 'JWT 解码' } },
      { path: 'devtools/yaml', name: 'YAMLFormatter', component: () => import('../views/devtools/YAMLFormatter.vue'), meta: { title: 'YAML/TOML' } },
      { path: 'devtools/uuid', name: 'UUIDGen', component: () => import('../views/devtools/UUIDGen.vue'), meta: { title: 'UUID 生成' } },
      {
        path: 'devtools/cron',
        name: 'CronTool',
        component: () => import('../views/devtools/CronTool.vue'),
        meta: { title: 'Cron 生成器', icon: 'timer-outline' }
      },
      {
        path: 'devtools/qrcode',
        name: 'QRCodeGen',
        component: () => import('../views/devtools/QRCodeGen.vue'),
        meta: { title: '二维码生成', icon: 'qr-code-outline' }
      },
      {
        path: 'devtools/markdown',
        name: 'MarkdownPreview',
        component: () => import('../views/devtools/MarkdownPreview.vue'),
        meta: { title: 'Markdown 预览', icon: 'document-text-outline' }
      },
      {
        path: 'daily/convert',
        name: 'UnitConverter',
        component: () => import('../views/daily/UnitConverter.vue'),
        meta: { title: '单位换算', icon: 'calculator-outline' }
      },
      {
        path: 'daily/screenshot',
        name: 'Screenshot',
        component: () => import('../views/daily/Screenshot.vue'),
        meta: { title: '截屏工具', icon: 'camera-outline' }
      },
      {
        path: 'daily/colorpicker',
        name: 'ColorPicker',
        component: () => import('../views/daily/ColorPicker.vue'),
        meta: { title: '取色器', icon: 'color-palette-outline' }
      },
      {
        path: 'daily/pomodoro',
        name: 'Pomodoro',
        component: () => import('../views/daily/Pomodoro.vue'),
        meta: { title: '番茄时钟', icon: 'timer-outline' }
      },
      {
        path: 'daily/clipboard',
        name: 'Clipboard',
        component: () => import('../views/daily/Clipboard.vue'),
        meta: { title: '剪贴板历史', icon: 'clipboard-outline' }
      },
      {
        path: 'daily/scheduler',
        name: 'Scheduler',
        component: () => import('../views/daily/Scheduler.vue'),
        meta: { title: '定时任务', icon: 'alarm-outline' }
      },
      {
        path: 'daily/reminder',
        name: 'Reminder',
        component: () => import('../views/daily/Reminder.vue'),
        meta: { title: '提醒事项', icon: 'notifications-outline' }
      },
      {
        path: 'daily/camera',
        name: 'CameraCapture',
        component: () => import('../views/daily/CameraCapture.vue'),
        meta: { title: '摄像头拍照', icon: 'camera-outline' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/settings/Settings.vue'),
        meta: { title: '设置', icon: 'cog-outline' }
      },
      {
        path: 'settings/contextmenu',
        name: 'ContextMenu',
        component: () => import('../views/settings/ContextMenu.vue'),
        meta: { title: '右键菜单', icon: 'list-outline' }
      },
      { path: 'settings/update', name: 'AutoUpdate', component: () => import('../views/settings/AutoUpdate.vue'), meta: { title: '检查更新' } },
      { path: 'settings/themes', name: 'ThemeStore', component: () => import('../views/settings/ThemeStore.vue'), meta: { title: '主题商店' } },
      { path: 'settings/logs', name: 'OperationLog', component: () => import('../views/settings/OperationLog.vue'), meta: { title: '操作日志' } }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
