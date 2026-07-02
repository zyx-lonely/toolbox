<template>
  <div>
    <n-h2>自动清理调度</n-h2>
    <p>设置定时自动清理临时文件、缓存等。</p>
    <n-space class="mb-4">
      <n-button type="primary" @click="loadTasks" :loading="loading">刷新</n-button>
    </n-space>
    <n-grid :cols="3" :x-gap="12" class="mb-4">
      <n-gi><n-card size="small"><n-statistic label="任务总数" :value="stats.totalTasks" /></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="已启用" :value="stats.enabledTasks" /></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="已释放空间" :value="formatBytes(stats.totalFreed)" /></n-card></n-gi>
    </n-grid>
    <n-card title="清理任务" size="small" class="mb-4">
      <n-data-table :columns="taskColumns" :data="tasks" :bordered="true" :loading="loading" size="small" />
    </n-card>
    <n-card title="执行日志" size="small">
      <n-data-table :columns="logColumns" :data="logs" :bordered="true" size="small" :pagination="{ pageSize: 10 }" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { NSwitch, NButton, NTag, useMessage } from 'naive-ui'
import { GetCleanupTasks, RunCleanupTask, ToggleCleanupTask, GetCleanupLogs, GetCleanupStats } from '@wails/go/main/App'

interface CleanupTask { id: string; name: string; enabled: boolean; schedule: string; hour: number; minute: number; lastRun: string; nextRun: string; runCount: number }
interface CleanupLog { taskId: string; taskName: string; runTime: string; freedSize: number; fileCount: number; status: string }

const tasks = ref<CleanupTask[]>([])
const logs = ref<CleanupLog[]>([])
const stats = ref<{ totalTasks: number; enabledTasks: number; totalFreed: number }>({ totalTasks: 0, enabledTasks: 0, totalFreed: 0 })
const loading = ref(false)
const message = useMessage()

function formatBytes(bytes: number): string {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + units[i]
}

function formatSchedule(task: CleanupTask): string {
  if (task.schedule === 'daily') return `每天 ${String(task.hour).padStart(2, '0')}:${String(task.minute).padStart(2, '0')}`
  if (task.schedule === 'weekly') return `每周 ${task.hour}:${String(task.minute).padStart(2, '0')}`
  return task.schedule
}

const taskColumns = [
  { title: '名称', key: 'name', width: 150 },
  { title: '计划', width: 150, render: (row: CleanupTask) => formatSchedule(row) },
  { title: '上次执行', key: 'lastRun', width: 150 },
  { title: '执行次数', key: 'runCount', width: 90 },
  { title: '状态', key: 'enabled', width: 80, render: (row: CleanupTask) => h(NSwitch, { value: row.enabled, onUpdateValue: (val: boolean) => toggleTask(row, val) }) },
  { title: '操作', width: 100, render: (row: CleanupTask) => h(NButton, { size: 'small', type: 'primary', quaternary: true, onClick: () => runTask(row) }, { default: () => '立即执行' }) }
]

const logColumns = [
  { title: '任务', key: 'taskName', width: 120 },
  { title: '执行时间', key: 'runTime', width: 150 },
  { title: '释放空间', width: 100, render: (row: CleanupLog) => formatBytes(row.freedSize) },
  { title: '文件数', key: 'fileCount', width: 80 },
  { title: '状态', key: 'status', width: 80, render: (row: CleanupLog) => h(NTag, { type: (row.status === 'success' ? 'success' : 'error') as any, size: 'small' }, { default: () => row.status }) }
]

async function loadTasks() {
  loading.value = true
  try {
    const [t, l, s] = await Promise.all([GetCleanupTasks(), GetCleanupLogs(20), GetCleanupStats()])
    tasks.value = (t || []) as CleanupTask[]
    logs.value = (l || []) as CleanupLog[]
    stats.value = (s || {}) as any
  } catch (e) { console.error(e) }
  loading.value = false
}

async function toggleTask(task: CleanupTask, enabled: boolean) {
  try { await ToggleCleanupTask(task.id, enabled); task.enabled = enabled; message.success(enabled ? '已启用' : '已禁用'); loadTasks() } catch (e: any) { message.error(`失败: ${e}`) }
}

async function runTask(task: CleanupTask) {
  try { const r = await RunCleanupTask(task.id) as CleanupLog; if (r) { message.success(`释放 ${formatBytes(r.freedSize)}`); loadTasks() } } catch (e: any) { message.error(`失败: ${e}`) }
}

onMounted(loadTasks)
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
