<template>
  <div class="dashboard">
    <n-h2>电脑工具箱</n-h2>
    <n-p>综合系统维护与优化工具</n-p>

    <!-- 收藏工具 -->
    <div v-if="toolsStore.favoriteTools.length > 0" class="section">
      <n-h3>⭐ 收藏工具</n-h3>
      <n-space>
        <n-tag
          v-for="tool in toolsStore.favoriteTools"
          :key="tool.path"
          :bordered="true"
          :clickable="true"
          @click="router.push(tool.path)"
          size="large"
        >
          {{ tool.title }}
        </n-tag>
      </n-space>
    </div>

    <!-- 最近使用 -->
    <div v-if="toolsStore.recentTools.length > 0" class="section">
      <n-h3>⏱️ 最近使用</n-h3>
      <n-space>
        <n-tag
          v-for="tool in toolsStore.recentTools.slice(0, 8)"
          :key="tool.path"
          :bordered="true"
          :clickable="true"
          @click="router.push(tool.path)"
          size="medium"
          type="info"
        >
          {{ tool.title }}
        </n-tag>
      </n-space>
    </div>

    <n-grid :cols="2" :x-gap="16" :y-gap="16" class="mt-4">
      <n-gi v-for="card in cards" :key="card.title">
        <n-card :title="card.title" hoverable @click="router.push(card.path); recordRecent(card)">
          <template #header-extra>
            <n-icon size="24" :color="card.color">
              <component :is="card.icon" />
            </n-icon>
          </template>
          {{ card.desc }}
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useToolsStore } from '../store/tools'
import {
  HardwareChipOutline, TrashOutline, PowerOutline,
  SettingsOutline, GlobeOutline, CopyOutline,
  KeyOutline, CodeSlashOutline, CalculatorOutline
} from '@vicons/ionicons5'

const router = useRouter()
const toolsStore = useToolsStore()

// 初始化 store
toolsStore.init()

const cards = [
  { title: '系统信息', desc: '查看 CPU、内存、磁盘等硬件信息', path: '/system', icon: HardwareChipOutline, color: '#18a058' },
  { title: '磁盘清理', desc: '清除系统临时文件和缓存，释放磁盘空间', path: '/optimize/disk', icon: TrashOutline, color: '#d03050' },
  { title: '启动项管理', desc: '管理开机自启程序，加快系统启动速度', path: '/optimize/startup', icon: PowerOutline, color: '#f0a020' },
  { title: '系统服务优化', desc: '按场景优化 Windows 后台服务', path: '/optimize/service', icon: SettingsOutline, color: '#2080f0' },
  { title: '网络工具', desc: 'Ping、端口扫描、DNS 查询等网络诊断', path: '/network', icon: GlobeOutline, color: '#18a058' },
  { title: '重复文件查找', desc: '扫描并清理磁盘中的重复文件', path: '/file/duplicate', icon: CopyOutline, color: '#d03050' },
  { title: '密码生成器', desc: '生成高强度随机密码', path: '/security/password', icon: KeyOutline, color: '#f0a020' },
  { title: 'JSON 格式化', desc: 'JSON 格式化、校验与压缩', path: '/devtools/json', icon: CodeSlashOutline, color: '#2080f0' },
]

function recordRecent(card: any) {
  toolsStore.addRecent({
    title: card.title,
    path: card.path
  })
}
</script>

<style scoped>
.dashboard { max-width: 900px; margin: 0 auto; }
.mt-4 { margin-top: 20px; }
</style>
