<template>
  <div>
    <n-h2>实时性能监控</n-h2>
    <p>查看系统 CPU、内存等实时使用情况。</p>
    <n-space class="mb-4">
      <n-button type="primary" @click="refresh" :loading="loading">刷新</n-button>
    </n-space>
    <n-grid :cols="4" :x-gap="12" :y-gap="12" class="mb-4">
      <n-gi><n-card size="small"><n-statistic label="CPU 使用率"><template #suffix>%</template>{{ snapshot.cpuUsage.toFixed(1) }}</n-statistic></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="内存使用率"><template #suffix>%</template>{{ snapshot.memoryUsage.toFixed(1) }}</n-statistic></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="内存使用" :value="formatBytes(snapshot.memoryUsed)" /></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="进程数" :value="snapshot.processCount" /></n-card></n-gi>
    </n-grid>
    <n-card title="详细信息" size="small">
      <n-descriptions bordered :column="2">
        <n-descriptions-item label="CPU 使用率">{{ snapshot.cpuUsage.toFixed(1) }}%</n-descriptions-item>
        <n-descriptions-item label="内存使用率">{{ snapshot.memoryUsage.toFixed(1) }}%</n-descriptions-item>
        <n-descriptions-item label="内存总量">{{ formatBytes(snapshot.memoryTotal) }}</n-descriptions-item>
        <n-descriptions-item label="内存已用">{{ formatBytes(snapshot.memoryUsed) }}</n-descriptions-item>
        <n-descriptions-item label="内存可用">{{ formatBytes(snapshot.memoryAvail) }}</n-descriptions-item>
        <n-descriptions-item label="进程数">{{ snapshot.processCount }}</n-descriptions-item>
        <n-descriptions-item label="线程数">{{ snapshot.threadCount }}</n-descriptions-item>
        <n-descriptions-item label="句柄数">{{ snapshot.handleCount }}</n-descriptions-item>
      </n-descriptions>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { GetPerfSnapshot, ClearPerfHistory } from '@wails/go/main/App'

interface PerfSnapshot { timestamp: number; cpuUsage: number; memoryUsage: number; memoryTotal: number; memoryUsed: number; memoryAvail: number; processCount: number; threadCount: number; handleCount: number }

const snapshot = ref<PerfSnapshot>({ timestamp: 0, cpuUsage: 0, memoryUsage: 0, memoryTotal: 0, memoryUsed: 0, memoryAvail: 0, processCount: 0, threadCount: 0, handleCount: 0 })
const loading = ref(false)
let refreshInterval: number | null = null

function formatBytes(bytes: number): string {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + units[i]
}

async function refresh() {
  loading.value = true
  try { const r = await GetPerfSnapshot() as PerfSnapshot; if (r) snapshot.value = r } catch (e) { console.error(e) }
  loading.value = false
}

onMounted(() => { refresh(); refreshInterval = window.setInterval(refresh, 2000) })
onUnmounted(() => { if (refreshInterval) clearInterval(refreshInterval) })
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
