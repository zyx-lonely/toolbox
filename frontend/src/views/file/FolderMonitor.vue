<template>
  <div>
    <h2>文件夹变化监控</h2>
    <p>实时监控目录中的文件变化（创建、修改、删除）。</p>
    <n-space class="mb-4">
      <n-input v-model:value="monitorPath" placeholder="选择要监控的目录" readonly style="width: 350px" />
      <n-button @click="selectDir">选择目录</n-button>
      <n-button v-if="!monitoring" type="primary" @click="startMonitor" :disabled="!monitorPath">开始监控</n-button>
      <n-button v-else type="error" @click="stopMonitor">停止监控</n-button>
      <n-button @click="refreshEvents" :disabled="!monitoring">刷新</n-button>
      <n-button @click="clearEvents">清空</n-button>
    </n-space>
    <n-space class="mb-4" v-if="monitoring">
      <n-tag type="success" size="large">监控中</n-tag>
      <n-text depth="3">共 {{ events.length }} 条变化记录</n-text>
    </n-space>
    <n-empty v-if="!monitoring && events.length === 0" description="请选择目录后开始监控" />
    <n-data-table v-if="events.length > 0" :columns="columns" :data="events" :bordered="true" size="small" :pagination="{ pageSize: 20 }" />
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted, onUnmounted } from 'vue'
import { NTag } from 'naive-ui'
import { SelectDirectory, StartFolderMonitor, StopFolderMonitor, GetFolderMonitorEvents, ClearFolderMonitorEvents, IsFolderMonitoring } from '@wails/go/main/App'

interface FolderChangeEvent { timestamp: string; type: string; path: string; size: number }

const monitorPath = ref('')
const monitoring = ref(false)
const events = ref<FolderChangeEvent[]>([])
let pollTimer: number | null = null

const columns = [
  { title: '时间', key: 'timestamp', width: 160 },
  { title: '类型', key: 'type', width: 100, render: (row: FolderChangeEvent) => {
    const colorMap: Record<string, string> = { created: 'success', modified: 'warning', deleted: 'error' }
    return h(NTag, { type: (colorMap[row.type] || 'default') as any, size: 'small' }, { default: () => row.type })
  }},
  { title: '文件路径', key: 'path', ellipsis: { tooltip: true } },
  { title: '大小', key: 'size', width: 100, render: (row: FolderChangeEvent) => formatSize(row.size) },
]

function formatSize(bytes: number): string {
  if (!bytes) return '-'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

async function selectDir() {
  try { const d = await SelectDirectory(); if (d) monitorPath.value = d } catch (e) { console.error(e) }
}

async function startMonitor() {
  if (!monitorPath.value) return
  try {
    await StartFolderMonitor(monitorPath.value)
    monitoring.value = true
    events.value = []
    pollTimer = window.setInterval(refreshEvents, 3000)
  } catch (e: any) { console.error(e) }
}

async function stopMonitor() {
  await StopFolderMonitor()
  monitoring.value = false
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
}

async function refreshEvents() {
  try { events.value = (await GetFolderMonitorEvents() as FolderChangeEvent[]).reverse() } catch (e) { console.error(e) }
}

async function clearEvents() {
  await ClearFolderMonitorEvents()
  events.value = []
}

onMounted(async () => {
  try { monitoring.value = await IsFolderMonitoring() as boolean } catch (e) {}
})

onUnmounted(() => { if (pollTimer) clearInterval(pollTimer) })
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
