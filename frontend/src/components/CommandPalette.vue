<template>
  <n-modal :show="show" @update:show="$emit('close')" :mask-closable="true" @after-leave="$emit('close')">
    <n-card style="width: 560px" :bordered="true" size="small" title="搜索工具">
      <n-input
        ref="inputRef"
        v-model:value="query"
        placeholder="输入关键词搜索工具（支持中文/英文）..."
        clearable
        size="large"
        autofocus
        @keydown="handleKeydown"
      >
        <template #prefix>
          <n-icon :component="SearchOutline" />
        </template>
      </n-input>

      <div v-if="results.length" style="margin-top: 12px; max-height: 360px; overflow-y: auto">
        <n-list hoverable clickable>
          <n-list-item
            v-for="(item, idx) in results"
            :key="item.key"
            :style="idx === activeIndex ? 'background: var(--n-color-embedded-hover, #e8e8e8)' : ''"
            @click="navigate(item.key)"
            @mouseenter="activeIndex = idx"
          >
            <n-space align="center">
              <n-icon size="20">
                <component :is="item.iconComp" />
              </n-icon>
              <div>
                <div style="font-weight: 500">{{ item.label }}</div>
                <div style="font-size: 12px; color: #999">{{ item.desc }}</div>
              </div>
            </n-space>
          </n-list-item>
        </n-list>
      </div>
      <n-empty v-if="query.length > 0 && !results.length" description="未找到匹配的工具" style="margin-top: 20px" />
      <n-text v-if="!query.length" depth="3" style="display: block; margin-top: 16px; font-size: 13px">
        输入关键词开始搜索，支持工具名称、描述、关键字
      </n-text>

      <template #footer>
        <n-text depth="3" style="font-size: 12px">↑ ↓ 导航 · Enter 打开 · Esc 关闭</n-text>
      </template>
    </n-card>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'
import { useRouter } from 'vue-router'
import { SearchOutline, GridOutline, HardwareChipOutline, TrashOutline, PowerOutline, SettingsOutline, GlobeOutline, CopyOutline, KeyOutline, CodeSlashOutline, SwapHorizontalOutline, LockOpenOutline, DocumentTextOutline, AnalyticsOutline, LocationOutline, WifiOutline, ListOutline, CloudDownloadOutline, ShieldOutline, LockClosedOutline, DesktopOutline, ClipboardOutline, AlarmOutline, CalculatorOutline, ColorPaletteOutline, TimerOutline, CogOutline, StatsChartOutline } from '@vicons/ionicons5'

defineProps<{ show: boolean }>()
const emit = defineEmits(['close', 'navigate'])
const router = useRouter()
const query = ref('')
const inputRef = ref<any>(null)
const activeIndex = ref(0)

// 所有工具的索引（完整列表，带图标和描述）
const allItems = [
  { label: '首页', key: '/dashboard', icon: 'GridOutline', desc: '系统概览与快捷操作', keywords: ['home', 'dashboard', '首页', '概览'] },
  { label: '系统信息', key: '/system', icon: 'HardwareChipOutline', desc: '查看硬件和系统详情', keywords: ['system', 'info', '硬件', 'cpu', '内存', '主板'] },
  { label: '系统监控', key: '/system/monitor', icon: 'StatsChartOutline', desc: '实时监控系统资源（CPU/内存/磁盘/网络）', keywords: ['monitor', '监控', '实时', 'cpu', '内存', '磁盘', '网络', '图表'] },
  { label: '硬件检测报告', key: '/system/report', icon: 'DocumentTextOutline', desc: '生成硬件健康报告', keywords: ['report', '硬件', '报告', '检测', '健康'] },
  { label: '系统激活', key: '/system/activation', icon: 'LockClosedOutline', desc: '查看 Windows 激活状态', keywords: ['activation', '激活', 'windows', '正版'] },
  { label: '电源方案', key: '/system/powerplan', icon: 'PowerOutline', desc: '管理电源计划和节能设置', keywords: ['power', '电源', '节能', '计划'] },
  { label: '进程管理器', key: '/system/process', icon: 'ListOutline', desc: '查看和管理运行中的进程', keywords: ['process', '进程', '任务管理器', '结束进程'] },
  { label: '磁盘清理', key: '/optimize/disk', icon: 'TrashOutline', desc: '清理临时文件、回收站、更新缓存', keywords: ['disk', 'clean', '清理', '磁盘', '临时文件', '垃圾'] },
  { label: '启动项管理', key: '/optimize/startup', icon: 'PowerOutline', desc: '管理开机启动项，加速开机', keywords: ['startup', '启动', '开机自启', '加速'] },
  { label: '系统服务优化', key: '/optimize/service', icon: 'SettingsOutline', desc: '优化系统服务，提升性能', keywords: ['service', '服务', '优化', '禁用'] },
  { label: '系统健康体检', key: '/optimize/health', icon: 'AnalyticsOutline', desc: '全面体检系统健康状态', keywords: ['health', '健康', '体检', '诊断', '检查'] },
  { label: '浏览器数据管理', key: '/optimize/browser', icon: 'LogoChromeOutline', desc: '清理浏览器缓存、Cookie、历史记录', keywords: ['browser', '浏览器', '缓存', 'cookie', 'chrome', 'edge'] },
  { label: 'Windows 更新管理', key: '/optimize/winupdate', icon: 'CloudDownloadOutline', desc: '管理 Windows 更新，暂停或恢复', keywords: ['update', '更新', 'windows update', '暂停'] },
  { label: '环境变量管理', key: '/optimize/env', icon: 'CodeSlashOutline', desc: '查看和编辑系统环境变量', keywords: ['env', '环境变量', 'path', '系统变量'] },
  { label: '磁盘健康监控', key: '/optimize/disk-health', icon: 'HardwareChipOutline', desc: '监控硬盘健康状态和温度', keywords: ['disk', '硬盘', '健康', 'smart', '温度'] },
  { label: '服务依赖分析', key: '/optimize/service-dep', icon: 'AnalyticsOutline', desc: '分析系统服务依赖关系', keywords: ['service', '依赖', '分析', '关系图'] },
  { label: '文件关联管理', key: '/optimize/fileassoc', icon: 'DocumentTextOutline', desc: '管理系统文件关联', keywords: ['file', '文件', '关联', '默认程序'] },
  { label: '网络工具', key: '/network', icon: 'GlobeOutline', desc: 'Ping、端口扫描、DNS 修复、网络诊断', keywords: ['ping', 'port', 'dns', '网络', '修复', '诊断', 'speedtest', '网速'] },
  { label: 'IP 地理位置', key: '/network/geoip', icon: 'LocationOutline', desc: 'IP 地址归属地查询', keywords: ['geo', 'ip', '地理位置', '查询', '归属地'] },
  { label: '局域网扫描', key: '/network/lanscan', icon: 'WifiOutline', desc: '扫描局域网设备', keywords: ['lan', 'scan', '局域网', '扫描', '设备'] },
  { label: 'WiFi 密码查看', key: '/network/wifi', icon: 'WifiOutline', desc: '查看已保存的 WiFi 密码', keywords: ['wifi', '密码', 'wireless', '查看'] },
  { label: 'WiFi 扫描', key: '/network/wifiscan', icon: 'WifiOutline', desc: '扫描附近 WiFi 信号', keywords: ['wifi', '扫描', '信号', '附近'] },
  { label: 'IP 冲突检测', key: '/network/ipconflict', icon: 'GlobeOutline', desc: '检测局域网 IP 冲突', keywords: ['ip', '冲突', '检测', '局域网'] },
  { label: '局域网拓扑', key: '/network/topology', icon: 'AnalyticsOutline', desc: '生成局域网设备拓扑图', keywords: ['topology', '拓扑', '局域网', '地图'] },
  { label: '批量 Ping', key: '/network/batchping', icon: 'GlobeOutline', desc: '批量 Ping 多个 IP 地址', keywords: ['ping', '批量', '多个', 'ip'] },
  { label: '流量图表', key: '/network/traffic', icon: 'StatsChartOutline', desc: '查看网络流量图表', keywords: ['traffic', '流量', '图表', '监控'] },
  { label: '远程桌面', key: '/network/rdp', icon: 'DesktopOutline', desc: '连接到远程桌面', keywords: ['rdp', '远程', 'remote', '桌面', 'mstsc'] },
  { label: '端口占用查询', key: '/network/portusage', icon: 'ListOutline', desc: '查看端口占用情况', keywords: ['port', '端口', '占用', '查看'] },
  { label: 'DNS 切换', key: '/network/dns', icon: 'GlobeOutline', desc: '快速切换 DNS 服务器', keywords: ['dns', '切换', '域名', '解析'] },
  { label: 'IP 子网计算器', key: '/network/subnet', icon: 'CodeSlashOutline', desc: 'IP 子网划分计算', keywords: ['subnet', '子网', '计算', 'ip', 'cidr'] },
  { label: '网速测试', key: '/network/speedtest', icon: 'StatsChartOutline', desc: '测试网络上传下载速度', keywords: ['speedtest', '网速', '测试', '上传', '下载'] },
  { label: 'Hosts 编辑器', key: '/network/hosts', icon: 'DocumentTextOutline', desc: '编辑系统 hosts 文件', keywords: ['hosts', '编辑', '域名', '映射'] },
  { label: '重复文件查找', key: '/file/duplicate', icon: 'CopyOutline', desc: '扫描并清理重复文件，释放空间', keywords: ['duplicate', '重复', '文件', '查找', '清理'] },
  { label: '大文件查找', key: '/file/largefiles', icon: 'TrashOutline', desc: '查找占用空间较大的文件', keywords: ['large', '大文件', '查找', '空间', '清理'] },
  { label: '图片格式转换', key: '/file/imageconv', icon: 'CopyOutline', desc: '批量转换图片格式', keywords: ['image', '图片', '转换', '格式', 'png', 'jpg'] },
  { label: '图片批量压缩', key: '/file/compress', icon: 'CopyOutline', desc: '批量压缩图片大小', keywords: ['compress', '压缩', '图片', '批量'] },
  { label: '文本编码转换', key: '/file/encoding', icon: 'CodeSlashOutline', desc: '转换文本文件编码格式', keywords: ['encoding', '编码', '转换', 'utf8', 'gbk'] },
  { label: '文件哈希校验', key: '/file/hash', icon: 'KeyOutline', desc: '计算文件 MD5/SHA1/SHA256 哈希值', keywords: ['hash', '哈希', 'md5', 'sha1', 'sha256', '校验'] },
  { label: '文件内容搜索', key: '/file/search', icon: 'SearchOutline', desc: '在文件中搜索指定内容', keywords: ['search', '搜索', '内容', '文件内容'] },
  { label: '文件内容替换', key: '/file/replace', icon: 'SwapHorizontalOutline', desc: '批量替换文件中的文本内容', keywords: ['replace', '替换', '批量', '文本'] },
  { label: '文件夹大小分析', key: '/file/foldersize', icon: 'AnalyticsOutline', desc: '分析文件夹大小分布', keywords: ['folder', '文件夹', '大小', '分析', '分布'] },
  { label: '文件差异对比', key: '/file/filediff', icon: 'SwapHorizontalOutline', desc: '对比两个文件的差异', keywords: ['diff', '差异', '对比', '文件对比'] },
  { label: 'PDF 工具', key: '/file/pdf', icon: 'DocumentTextOutline', desc: 'PDF 合并、拆分、压缩、转换', keywords: ['pdf', '合并', '拆分', '压缩', '转换'] },
  { label: '批量重命名', key: '/file/rename', icon: 'CalculatorOutline', desc: '批量重命名文件', keywords: ['rename', '重命名', '批量', '文件'] },
  { label: '文件批量归类', key: '/file/organize', icon: 'FolderOutline', desc: '按规则批量归类文件', keywords: ['organize', '归类', '整理', '分类'] },
  { label: '文件解锁器', key: '/file/unlock', icon: 'LockOpenOutline', desc: '解锁被占用的文件', keywords: ['unlock', '解锁', '占用', '删除'] },
  { label: '文件时间戳修改', key: '/file/timestamp', icon: 'TimerOutline', desc: '批量修改文件创建/修改时间', keywords: ['timestamp', '时间戳', '时间', '修改', '创建'] },
  { label: '批量正则重命名', key: '/file/batchregex', icon: 'CodeSlashOutline', desc: '使用正则表达式批量重命名', keywords: ['regex', '正则', '重命名', '批量'] },
  { label: 'JSON 格式化', key: '/devtools/json', icon: 'CodeSlashOutline', desc: 'JSON 格式化、校验、压缩、对比', keywords: ['json', '格式化', '校验', '压缩', '对比'] },
  { label: '编解码工具', key: '/devtools/codec', icon: 'KeyOutline', desc: 'Base64、URL、Hex 编解码', keywords: ['base64', 'url', '编码', '解码', 'hex'] },
  { label: '代码美化', key: '/devtools/beautify', icon: 'CodeSlashOutline', desc: '格式化 HTML/CSS/JS/SQL 代码', keywords: ['beautify', '代码', '格式化', 'html', 'css', 'js'] },
  { label: 'HTTP 调试', key: '/devtools/http', icon: 'GlobeOutline', desc: '发送 HTTP 请求并查看响应', keywords: ['http', '请求', '调试', 'api', 'rest'] },
  { label: '正则测试', key: '/devtools/regex', icon: 'CodeSlashOutline', desc: '在线测试正则表达式', keywords: ['regex', '正则', '测试', '表达式'] },
  { label: 'JWT 解码', key: '/devtools/jwt', icon: 'KeyOutline', desc: '解码和验证 JWT Token', keywords: ['jwt', 'token', '解码', '验证'] },
  { label: 'YAML/TOML 转换', key: '/devtools/yaml', icon: 'CodeSlashOutline', desc: 'YAML 与 TOML 互转', keywords: ['yaml', 'toml', '转换', '格式'] },
  { label: 'UUID 生成', key: '/devtools/uuid', icon: 'KeyOutline', desc: '批量生成 UUID', keywords: ['uuid', '生成', '批量'] },
  { label: 'Markdown 预览', key: '/devtools/markdown', icon: 'DocumentTextOutline', desc: 'Markdown 实时预览和导出', keywords: ['markdown', '预览', '导出', '编辑器'] },
  { label: '二维码生成', key: '/devtools/qrcode', icon: 'CopyOutline', desc: '生成二维码图片', keywords: ['qrcode', '二维码', '生成'] },
  { label: 'Cron 生成器', key: '/devtools/cron', icon: 'TimerOutline', desc: '生成和解析 Cron 表达式', keywords: ['cron', '定时', '表达式', '生成'] },
  { label: '文本对比', key: '/devtools/textdiff', icon: 'SwapHorizontalOutline', desc: '对比两段文本的差异', keywords: ['text', '文本', '对比', 'diff', '差异'] },
  { label: '文本字数统计', key: '/devtools/textcount', icon: 'CalculatorOutline', desc: '统计文本字数、字符数', keywords: ['text', '文本', '字数', '统计', '字符'] },
  { label: 'JSON/CSV 互转', key: '/devtools/jsoncsv', icon: 'SwapHorizontalOutline', desc: 'JSON 与 CSV 格式互转', keywords: ['json', 'csv', '转换', '格式'] },
  { label: '密码生成器', key: '/security/password', icon: 'KeyOutline', desc: '生成高强度随机密码', keywords: ['password', '密码', '生成', '随机', '强度'] },
  { label: '文件加密解密', key: '/security/encrypt', icon: 'LockClosedOutline', desc: '加密/解密文件', keywords: ['encrypt', '加密', '解密', '文件'] },
  { label: '文件粉碎', key: '/security/shred', icon: 'TrashOutline', desc: '彻底删除文件不可恢复', keywords: ['shred', '粉碎', '删除', '彻底', '不可恢复'] },
  { label: '单位换算', key: '/daily/convert', icon: 'CalculatorOutline', desc: '长度/重量/温度/数据单位换算', keywords: ['convert', '换算', '单位', '转换', '长度', '重量'] },
  { label: '取色器', key: '/daily/colorpicker', icon: 'ColorPaletteOutline', desc: '屏幕取色并复制颜色值', keywords: ['color', '颜色', '取色', 'picker', '十六进制', 'rgb'] },
  { label: '截屏工具', key: '/daily/screenshot', icon: 'CopyOutline', desc: '截图并识别文字', keywords: ['screenshot', '截屏', '截图', '识别', 'ocr'] },
  { label: '番茄时钟', key: '/daily/pomodoro', icon: 'TimerOutline', desc: '番茄工作法计时器', keywords: ['pomodoro', '番茄', '时钟', '计时', '专注'] },
  { label: '定时提醒', key: '/daily/reminder', icon: 'AlarmOutline', desc: '倒计时和闹钟提醒', keywords: ['reminder', '提醒', '倒计时', '闹钟', '定时'] },
  { label: '剪贴板历史', key: '/daily/clipboard', icon: 'ClipboardOutline', desc: '查看和管理剪贴板历史', keywords: ['clipboard', '剪贴板', '复制', '粘贴', '历史'] },
  { label: '定时任务', key: '/daily/scheduler', icon: 'TimerOutline', desc: '创建和管理定时任务', keywords: ['schedule', '定时', '任务', '计划', 'cron'] },
  { label: '便签', key: '/daily/notepad', icon: 'DocumentTextOutline', desc: '创建和管理便签', keywords: ['notepad', '便签', '笔记', '记事'] },
  { label: '应用设置', key: '/settings', icon: 'CogOutline', desc: '应用设置和主题切换', keywords: ['settings', '设置', '配置', '主题', '深色模式'] },
  { label: '右键菜单', key: '/settings/contextmenu', icon: 'ListOutline', desc: '管理右键菜单项', keywords: ['context', '菜单', '右键', '管理'] },
  { label: '检查更新', key: '/settings/update', icon: 'CloudDownloadOutline', desc: '检查应用更新', keywords: ['update', '更新', '检查', '版本'] },
  { label: '主题商店', key: '/settings/themes', icon: 'ColorPaletteOutline', desc: '下载和应用主题', keywords: ['theme', '主题', '商店', '下载'] },
  { label: '操作日志', key: '/settings/logs', icon: 'DocumentTextOutline', desc: '查看应用操作日志', keywords: ['log', '日志', '操作', '查看'] },
]

// 图标组件映射
const iconMap: Record<string, any> = {
  GridOutline,
  HardwareChipOutline,
  TrashOutline,
  PowerOutline,
  SettingsOutline,
  GlobeOutline,
  CopyOutline,
  KeyOutline,
  CodeSlashOutline,
  SwapHorizontalOutline,
  LockOpenOutline,
  DocumentTextOutline,
  AnalyticsOutline,
  LocationOutline,
  WifiOutline,
  ListOutline,
  CloudDownloadOutline,
  ShieldOutline,
  LockClosedOutline,
  DesktopOutline,
  ClipboardOutline,
  AlarmOutline,
  CalculatorOutline,
  ColorPaletteOutline,
  TimerOutline,
  CogOutline,
  StatsChartOutline,
}

const results = computed(() => {
  if (!query.value.trim()) return []
  const q = query.value.toLowerCase()
  return allItems
    .filter(item =>
      item.label.toLowerCase().includes(q) ||
      item.desc.toLowerCase().includes(q) ||
      item.keywords.some(k => k.includes(q))
    )
    .map(item => ({
      ...item,
      iconComp: iconMap[item.icon] || GridOutline
    }))
    .slice(0, 15)
})

watch(query, () => { activeIndex.value = 0 })
watch(() => query.value, (val) => { if (val) activeIndex.value = 0 })

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('close')
    return
  }
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = Math.min(activeIndex.value + 1, results.value.length - 1)
    return
  }
  if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = Math.max(activeIndex.value - 1, 0)
    return
  }
  if (e.key === 'Enter' && results.value.length > 0) {
    e.preventDefault()
    navigate(results.value[activeIndex.value].key)
  }
}

function navigate(key: string) {
  router.push(key)
  emit('close')
}

watch(() => query.value, (val) => {
  if (val) nextTick(() => inputRef.value?.focus())
})
</script>
