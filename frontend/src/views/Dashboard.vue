<template>
  <div class="dashboard">
    <n-h2>电脑工具箱</n-h2>
    <n-p>综合系统维护与优化工具</n-p>

    <n-grid :cols="2" :x-gap="16" :y-gap="16" class="mt-4">
      <n-gi v-for="card in cards" :key="card.title">
        <n-card :title="card.title" hoverable @click="router.push(card.path)">
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
import {
  HardwareChipOutline, TrashOutline, PowerOutline,
  SettingsOutline, GlobeOutline, CopyOutline,
  KeyOutline, CodeSlashOutline, CalculatorOutline
} from '@vicons/ionicons5'

const router = useRouter()

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
</script>

<style scoped>
.dashboard { max-width: 900px; margin: 0 auto; }
.mt-4 { margin-top: 20px; }
</style>
