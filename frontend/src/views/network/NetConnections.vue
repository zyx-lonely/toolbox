<template>
  <div>
    <n-h2>网络连接查看器</n-h2>
    <p>查看当前所有网络连接及对应进程。</p>
    <n-space class="mb-4">
      <n-button type="primary" @click="loadConnections" :loading="loading">刷新</n-button>
      <n-input v-model:value="searchKeyword" placeholder="搜索..." style="width: 200px" clearable />
    </n-space>
    <n-grid :cols="5" :x-gap="12" class="mb-4">
      <n-gi><n-card size="small"><n-statistic label="总数" :value="stats.total" /></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="TCP" :value="stats.tcp" /></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="UDP" :value="stats.udp" /></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="已建立" :value="stats.established" /></n-card></n-gi>
      <n-gi><n-card size="small"><n-statistic label="监听中" :value="stats.listening" /></n-card></n-gi>
    </n-grid>
    <n-empty v-if="!loading && filteredConnections.length === 0" description="暂无数据" />
    <n-data-table v-else :columns="columns" :data="filteredConnections" :bordered="true" :loading="loading" size="small" :pagination="{ pageSize: 20 }" />
  </div>
</template>

<script setup lang="ts">
import { h, ref, computed, onMounted } from 'vue'
import { NTag, NButton, useMessage } from 'naive-ui'
import { GetAllNetConnections, KillNetConnection } from '@wails/go/main/App'

interface NetConn { protocol: string; localAddr: string; localPort: number; remoteAddr: string; remotePort: number; state: string; pid: number; processName: string }

const connections = ref<NetConn[]>([])
const loading = ref(false)
const searchKeyword = ref('')
const message = useMessage()

const stats = computed(() => {
  const total = connections.value.length
  const tcp = connections.value.filter(c => c.protocol?.startsWith('TCP')).length
  const udp = total - tcp
  const established = connections.value.filter(c => c.state === 'ESTABLISHED').length
  const listening = connections.value.filter(c => c.state === 'LISTENING').length
  return { total, tcp, udp, established, listening }
})

const filteredConnections = computed(() => {
  if (!searchKeyword.value) return connections.value
  const kw = searchKeyword.value.toLowerCase()
  return connections.value.filter(c =>
    (c.localAddr || '').toLowerCase().includes(kw) ||
    (c.remoteAddr || '').toLowerCase().includes(kw) ||
    (c.processName || '').toLowerCase().includes(kw) ||
    String(c.pid).includes(kw)
  )
})

const columns = [
  { title: '协议', key: 'protocol', width: 70 },
  { title: '本地地址', key: 'localAddr', width: 150 },
  { title: '本地端口', key: 'localPort', width: 90 },
  { title: '远程地址', key: 'remoteAddr', width: 150 },
  { title: '远程端口', key: 'remotePort', width: 90 },
  { title: '状态', key: 'state', width: 120, render: (row: NetConn) => h(NTag, { type: (row.state === 'ESTABLISHED' ? 'success' : row.state === 'LISTENING' ? 'info' : 'default') as any, size: 'small' }, { default: () => row.state || '-' }) },
  { title: '进程', key: 'processName', width: 120 },
  { title: 'PID', key: 'pid', width: 70 },
  { title: '操作', width: 80, render: (row: NetConn) => h(NButton, { size: 'small', type: 'error', quaternary: true, onClick: () => killConn(row) }, { default: () => '终止' }) }
]

async function loadConnections() {
  loading.value = true
  try {
    const result = await GetAllNetConnections()
    // Wails 多返回值: result 可能是 [connections, group] 或者直接是 connections
    if (Array.isArray(result)) {
      // result[0] 可能是数组（连接列表）
      const first = result[0]
      if (Array.isArray(first)) {
        connections.value = first as NetConn[]
      } else {
        connections.value = result as unknown as NetConn[]
      }
    } else if (result && typeof result === 'object') {
      // 可能直接返回了连接数组
      connections.value = [result] as unknown as NetConn[]
    }
  } catch (e: any) {
    message.error('加载失败: ' + String(e))
  }
  loading.value = false
}

async function killConn(conn: NetConn) {
  try { await KillNetConnection(conn.pid); message.success('已终止'); loadConnections() } catch (e: any) { message.error('失败: ' + String(e)) }
}

onMounted(loadConnections)
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
