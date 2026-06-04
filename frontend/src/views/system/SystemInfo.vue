<template>
  <div>
    <n-h2>系统信息</n-h2>

    <n-grid :cols="2" :x-gap="16" :y-gap="16">
      <n-gi>
        <n-card title="操作系统" :bordered="true">
          <n-description-list label-placement="left" :column="1">
            <n-description-item label="名称">{{ sysInfo?.os?.name || '加载中...' }}</n-description-item>
            <n-description-item label="版本">{{ sysInfo?.os?.version || '-' }}</n-description-item>
            <n-description-item label="构建版本">{{ sysInfo?.os?.buildNumber || '-' }}</n-description-item>
            <n-description-item label="架构">{{ sysInfo?.os?.architecture || '-' }}</n-description-item>
            <n-description-item label="已运行">{{ sysInfo?.os?.uptime || '-' }}</n-description-item>
          </n-description-list>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="处理器" :bordered="true">
          <n-description-list label-placement="left" :column="1">
            <n-description-item label="型号">{{ sysInfo?.cpu?.name || '加载中...' }}</n-description-item>
            <n-description-item label="物理核心">{{ sysInfo?.cpu?.cores || '-' }}</n-description-item>
            <n-description-item label="逻辑核心">{{ sysInfo?.cpu?.logicalCores || '-' }}</n-description-item>
            <n-description-item label="主频">{{ sysInfo?.cpu?.baseClock || '-' }}</n-description-item>
            <n-description-item label="使用率">
              <n-progress type="line" :percentage="Math.round(sysInfo?.cpu?.usage || 0)"
                :status="cpuStatus" :height="16" />
            </n-description-item>
          </n-description-list>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="内存" :bordered="true">
          <n-description-list label-placement="left" :column="1">
            <n-description-item label="总计">{{ formatBytes(sysInfo?.memory?.total) }}</n-description-item>
            <n-description-item label="已用">{{ formatBytes(sysInfo?.memory?.used) }}</n-description-item>
            <n-description-item label="可用">{{ formatBytes(sysInfo?.memory?.available) }}</n-description-item>
            <n-description-item label="使用率">
              <n-progress type="line" :percentage="Math.round(sysInfo?.memory?.usage || 0)"
                :status="memStatus" :height="16" />
            </n-description-item>
          </n-description-list>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="磁盘" :bordered="true">
          <div v-for="disk in sysInfo?.disks" :key="disk.label" class="disk-item">
            <div class="disk-label">
              <strong>{{ disk.label }}</strong>
              <n-tag size="small">{{ disk.fileSystem }}</n-tag>
            </div>
            <div class="disk-size">{{ formatBytes(disk.used) }} / {{ formatBytes(disk.total) }}</div>
            <n-progress type="line" :percentage="Math.round(disk.usage)"
              :status="disk.usage > 90 ? 'error' : disk.usage > 75 ? 'warning' : 'success'" :height="12" />
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <n-button type="primary" class="mt-4" @click="refresh" :loading="loading">
      刷新信息
    </n-button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { GetSystemInfo } from '@wails/go/main/App'
import { formatBytes } from '@/api/bridge'

interface OsInfo { name: string; version: string; buildNumber: string; architecture: string; uptime: string }
interface CpuInfo { name: string; cores: number; logicalCores: number; baseClock: string; usage: number }
interface MemInfo { total: number; used: number; available: number; usage: number }
interface DiskInfo { label: string; fileSystem: string; total: number; used: number; free: number; usage: number }
interface SystemInfo { os: OsInfo; cpu: CpuInfo; memory: MemInfo; disks: DiskInfo[]; network: any }

const sysInfo = ref<SystemInfo | null>(null)
const loading = ref(false)

const cpuStatus = computed(() => {
  const u = sysInfo.value?.cpu?.usage || 0
  return u > 90 ? 'error' : u > 75 ? 'warning' : 'success'
})

const memStatus = computed(() => {
  const u = sysInfo.value?.memory?.usage || 0
  return u > 90 ? 'error' : u > 80 ? 'warning' : 'success'
})

async function refresh() {
  loading.value = true
  try {
    sysInfo.value = await GetSystemInfo()
  } catch (e: any) {
    console.error('获取系统信息失败:', e)
  }
  loading.value = false
}

onMounted(refresh)
</script>

<style scoped>
.disk-item { margin-bottom: 12px; }
.disk-label { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.disk-size { font-size: 13px; color: #888; margin-bottom: 4px; }
.mt-4 { margin-top: 16px; }
</style>
