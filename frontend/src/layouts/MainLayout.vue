<template>
  <n-layout position="absolute" has-sider @keydown="handleKeydown">
    <n-layout-sider
      bordered
      :collapsed="appStore.collapsed"
      collapse-mode="width"
      :collapsed-width="60"
      :width="200"
      :native-scrollbar="false"
      show-trigger="arrow-circle"
      @collapse="appStore.toggleCollapsed()"
      @expand="appStore.toggleCollapsed()"
    >
      <div class="sidebar-header">
        <n-icon size="28" color="#18a058">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><rect x="32" y="32" width="192" height="192" rx="24" fill="currentColor"/><rect x="288" y="32" width="192" height="192" rx="24" fill="currentColor"/><rect x="32" y="288" width="192" height="192" rx="24" fill="currentColor"/><rect x="288" y="288" width="192" height="192" rx="24" fill="currentColor"/></svg>
          </n-icon>
          <span v-show="!appStore.collapsed" class="sidebar-title">电脑工具箱</span>
        </div>
      <n-menu
        :collapsed-width="60"
        :collapsed-icon-size="22"
        :options="menuOptions"
        :value="activeKey"
        @update:value="handleMenuSelect"
      />
    </n-layout-sider>

    <n-layout>
      <n-layout-header class="layout-header" bordered>
        <div class="header-left">
          <n-breadcrumb>
            <n-breadcrumb-item>{{ currentTitle }}</n-breadcrumb-item>
          </n-breadcrumb>
        </div>
        <div class="header-right">
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button size="small" quaternary @click="showPalette = true" style="font-size: 12px; font-weight: 500">
                Ctrl+K
              </n-button>
            </template>
            <span>搜索工具（Ctrl+K）</span>
          </n-tooltip>
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button quaternary circle @click="toggleFavorite">
                <template #icon>
                  <n-icon>
                    <StarOutline v-if="!isFavorite" />
                    <Star v-else />
                  </n-icon>
                </template>
              </n-button>
            </template>
            <span>{{ isFavorite ? '取消收藏' : '收藏此工具' }}</span>
          </n-tooltip>
          <n-tooltip trigger="hover">
            <template #trigger">
              <n-button quaternary circle @click="appStore.toggleDarkMode()">
                <template #icon>
                  <n-icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path d="M264 480A232 232 0 0132 248c0-94 54-178.28 137.61-214.67a16 16 0 0118.31 4.85 15.89 15.89 0 011.88 18.91C176.34 86.17 168 116.68 168 152c0 105.87 86.13 192 192 192 35.32 0 65.83-8.34 94.91-23.88a15.89 15.89 0 0118.91 1.88 16 16 0 014.85 18.31C442.28 301.78 358 480 264 480z" fill="currentColor"/></svg></n-icon>
                </template>
              </n-button>
            </template>
            <span>切换主题</span>
          </n-tooltip>
        </div>
      </n-layout-header>

      <n-layout-content class="layout-content" :embedded="true">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
  <CommandPalette :show="showPalette" @close="showPalette = false" />
</template>

<script setup lang="ts">
import { computed, h, watch, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NIcon, type MenuOption, useOsTheme } from 'naive-ui'
import {
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
  StarOutline,
  Star,
  AppsOutline,
  ConstructOutline
} from '@vicons/ionicons5'
import { useAppStore } from '../store/app'
import { useToolsStore } from '../store/tools'
import CommandPalette from '../components/CommandPalette.vue'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const toolsStore = useToolsStore()
const showPalette = ref(false)

// 初始化 tools store
toolsStore.init()

// 应用深色模式
watch(() => appStore.darkMode, (val) => {
  if (val) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}, { immediate: true })

// 当前页面是否已收藏
const isFavorite = computed(() => {
  return toolsStore.isFavorite(route.path)
})

// 切换收藏状态
function toggleFavorite() {
  const title = currentTitle.value
  if (title) {
    toolsStore.toggleFavorite({
      title: title,
      path: route.path
    })
  }
}

// 记录最近使用
function recordRecent() {
  const title = currentTitle.value
  if (title && route.path !== '/') {
    toolsStore.addRecent({
      title: title,
      path: route.path
    })
  }
}

// 监听路由变化，记录最近使用
watch(() => route.path, () => {
  // 延迟执行，确保页面已加载
  setTimeout(() => recordRecent(), 500)
})

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions: MenuOption[] = [
  { label: '首页', key: '/dashboard', icon: renderIcon(GridOutline) },
  {
    label: '系统工具',
    key: 'system-group',
    icon: renderIcon(HardwareChipOutline),
    children: [
      { label: '系统信息', key: '/system' },
      { label: '硬件检测报告', key: '/system/report' },
      { label: '系统激活', key: '/system/activation' },
      { label: '电源方案', key: '/system/powerplan' },
      { label: '进程管理器', key: '/system/process' },
      { label: '温度监控', key: '/system/temperature' },
      { label: '电池信息', key: '/system/battery' },
      { label: '浏览器扩展', key: '/system/browser-ext' },
      { label: '导出报告', key: '/system/export' },
      { label: '快捷键管理', key: '/system/shortcuts' },
      { label: '实时性能监控', key: '/system/perfmon' },
    ]
  },
  {
    label: 'Windows 优化',
    key: 'optimize',
    icon: renderIcon(SettingsOutline),
    children: [
      { label: '磁盘清理', key: '/optimize/disk' },
      { label: '启动项管理', key: '/optimize/startup' },
      { label: '系统服务优化', key: '/optimize/service' },
      { label: '系统健康体检', key: '/optimize/health' },
      { label: '系统备份与还原', key: '/optimize/restore' },
      { label: '浏览器数据管理', key: '/optimize/browser' },
      { label: '注册表清理', key: '/optimize/registry' },
      { label: 'Windows 更新管理', key: '/optimize/winupdate' },
      { label: '环境变量管理', key: '/optimize/env' },
      { label: '磁盘健康监控', key: '/optimize/disk-health' },
      { label: '服务依赖分析', key: '/optimize/service-dep' },
      { label: '文件关联管理', key: '/optimize/fileassoc' },
      { label: '软件管理', key: '/optimize/software' },
      { label: '自动清理调度', key: '/optimize/cleanup' }
    ]
  },
  {
    label: '网络工具',
    key: 'network-group',
    icon: renderIcon(GlobeOutline),
    children: [
      { label: '网络诊断', key: '/network' },
      { label: 'Hosts 编辑器', key: '/network/hosts' },
      { label: 'IP 地理位置', key: '/network/geoip' },
      { label: '局域网扫描', key: '/network/lanscan' },
      { label: 'WiFi 密码', key: '/network/wifi' },
      { label: 'WiFi 扫描', key: '/network/wifiscan' },
      { label: 'IP 冲突检测', key: '/network/ipconflict' },
      { label: '批量 Ping', key: '/network/batchping' },
      { label: '流量图表', key: '/network/traffic' },
      { label: '远程桌面', key: '/network/rdp' },
      { label: '端口占用查询', key: '/network/portusage' },
      { label: 'DNS 切换', key: '/network/dns' },
      { label: 'IP 子网计算器', key: '/network/subnet' },
      { label: '网速测试', key: '/network/speedtest' },
      { label: '网络连接查看器', key: '/network/connections' },
    ]
  },
  {
    label: '文件工具',
    key: 'file',
    icon: renderIcon(CopyOutline),
    children: [
      { label: '重复文件查找', key: '/file/duplicate' },
      { label: '大文件查找', key: '/file/largefiles' },
      { label: '图片格式转换', key: '/file/imageconv' },
      { label: '图片批量压缩', key: '/file/compress' },
      { label: '文本编码转换', key: '/file/encoding' },
      { label: '文件哈希校验', key: '/file/hash' },
      { label: '文件内容搜索', key: '/file/search' },
      { label: '文件内容替换', key: '/file/replace' },
      { label: '文件夹大小分析', key: '/file/foldersize' },
      { label: '文件差异对比', key: '/file/filediff' },
      { label: 'PDF 工具', key: '/file/pdf' },
      { label: '批量重命名', key: '/file/rename' },
      { label: '文件批量归类', key: '/file/organize' },
      { label: '文件解锁器', key: '/file/unlock' },
      { label: '文件时间戳修改', key: '/file/timestamp' },
      { label: '批量正则重命名', key: '/file/batchregex' },
      { label: '文件预览', key: '/file/preview' },
      { label: '磁盘空间树状图', key: '/file/disktree' },
      { label: '图片批量加水印', key: '/file/watermark' },
      { label: '文件夹监控', key: '/file/foldermonitor' }
    ]
  },
  {
    label: '外部工具',
    key: 'external-tools',
    icon: renderIcon(ConstructOutline),
    children: [
      { label: '外部工具管理', key: '/tools/external' },
      { label: '系统必备软件', key: '/tools/essential' }
    ]
  },
  {
    label: '开发工具',
    key: 'devtools',
    icon: renderIcon(CodeSlashOutline),
    children: [
      { label: 'JSON 格式化', key: '/devtools/json' },
      { label: '编解码工具', key: '/devtools/codec' },
      { label: '代码美化', key: '/devtools/beautify' },
      { label: 'HTTP 调试', key: '/devtools/http' },
      { label: '正则测试', key: '/devtools/regex' },
      { label: 'JWT 解码', key: '/devtools/jwt' },
      { label: 'YAML/TOML', key: '/devtools/yaml' },
      { label: 'UUID 生成', key: '/devtools/uuid' },
      { label: 'Markdown 预览', key: '/devtools/markdown' },
      { label: '二维码生成', key: '/devtools/qrcode' },
      { label: 'Base32/Hex 互转', key: '/devtools/codec-advanced' },
      { label: '时间戳转换', key: '/devtools/timestamp' },
      { label: '颜色代码转换', key: '/devtools/color' },
      { label: '密码生成器', key: '/security/password' },
      { label: '文件加密解密', key: '/security/encrypt' },
      { label: '文件粉碎', key: '/security/shred' },
      { label: '防火墙管理', key: '/security/firewall' },
      { label: 'Cron 生成器', key: '/devtools/cron' },
      { label: '文本对比', key: '/devtools/textdiff' },
      { label: '文本字数统计', key: '/devtools/textcount' },
      { label: 'JSON/CSV 互转', key: '/devtools/jsoncsv' }
    ]
  },
  {
    label: '日常工具',
    key: 'daily-group',
    icon: renderIcon(CalculatorOutline),
    children: [
      { label: '单位换算', key: '/daily/convert' },
      { label: '取色器', key: '/daily/colorpicker' },
      { label: '截屏工具', key: '/daily/screenshot' },
      { label: '番茄时钟', key: '/daily/pomodoro' },
      { label: '定时提醒', key: '/daily/reminder' },
      { label: '剪贴板历史', key: '/daily/clipboard' },
      { label: '定时任务', key: '/daily/scheduler' },
      { label: '便签', key: '/daily/notepad' },
      { label: '屏幕标尺', key: '/daily/ruler' },
      { label: '剪贴板翻译', key: '/daily/translate' },
      { label: '汇率查询', key: '/daily/exchange' },
      { label: '科学计算器', key: '/daily/calculator' },
      { label: '屏幕录制', key: '/daily/screenrec' },
    ]
  },
  {
    label: '设置',
    key: 'settings-group',
    icon: renderIcon(CogOutline),
    children: [
      { label: '应用设置', key: '/settings' },
      { label: '右键菜单', key: '/settings/contextmenu' },
      { label: '检查更新', key: '/settings/update' },
      { label: '主题商店', key: '/settings/themes' },
      { label: '操作日志', key: '/settings/logs' }
    ]
  }
]

const activeKey = computed(() => route.path)

const currentTitle = computed(() => {
  const matched = menuOptions.find(o => o.key === route.path) as any
  if (matched) return matched.label
  for (const opt of menuOptions) {
    if ((opt as any).children) {
      const child = (opt as any).children.find((c: any) => c.key === route.path)
      if (child) return child.label
    }
  }
  return route.meta?.title as string || '电脑工具箱'
})

function handleMenuSelect(key: string) {
  router.push(key)
}

function handleKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    showPalette.value = !showPalette.value
  }
}

onMounted(() => document.addEventListener('keydown', handleKeydown))
onUnmounted(() => document.removeEventListener('keydown', handleKeydown))
</script>

<style scoped>
.sidebar-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 12px;
  border-bottom: 1px solid var(--n-border-color);
}
.sidebar-title {
  font-size: 16px;
  font-weight: 600;
  white-space: nowrap;
}
.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 20px;
  height: 52px;
}
.header-left {
  display: flex;
  align-items: center;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}
.layout-content {
  padding: 20px;
  min-height: calc(100vh - 52px);
  overflow-y: auto;
}
</style>
