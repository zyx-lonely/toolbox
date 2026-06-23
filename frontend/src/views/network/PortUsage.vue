<template>
  <div>
    <n-h2>端口占用查询</n-h2>
    <n-p>查询端口占用情况，快速定位进程并释放端口。</n-p>

    <n-space vertical :size="16">
      <n-input-group>
        <n-input v-model:value="portInput" placeholder="输入端口号（如 80）" style="width: 200px" />
        <n-button type="primary" @click="queryPort" :loading="loading">查询</n-button>
        <n-button @click="refreshAll" :loading="refreshing">刷新全部常用端口</n-button>
      </n-input-group>

      <n-data-table
        :columns="columns"
        :data="portList"
        size="small"
        :bordered="true"
        :row-key="(row: any) => row.port"
      />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NTag, useMessage, type DataTableColumns } from 'naive-ui'

const message = useMessage()
const portInput = ref('')
const loading = ref(false)
const refreshing = ref(false)

interface PortInfo {
  port: number
  process: string
  pid: number
  status: 'listening' | 'free'
  protocol: string
}

const portList = ref<PortInfo[]>([])

const columns: DataTableColumns<PortInfo> = [
  { title: '端口', key: 'port', width: 80 },
  { title: '协议', key: 'protocol', width: 80 },
  { title: '进程', key: 'process', width: 160 },
  { title: 'PID', key: 'pid', width: 100 },
  {
    title: '状态', key: 'status', width: 100,
    render(row) {
      return h(NTag, { type: row.status === 'listening' ? 'warning' : 'success', size: 'small' }, {
        default: () => row.status === 'listening' ? '占用' : '空闲'
      })
    }
  },
  {
    title: '操作', key: 'actions', width: 100,
    render(row) {
      if (row.status === 'listening') {
        return h(NButton, { size: 'small', type: 'error', secondary: true, onClick: () => killProcess(row) }, { default: () => '释放' })
      }
      return null
    }
  }
]

const mockProcesses = ['node.exe', 'nginx.exe', 'chrome.exe', 'java.exe', 'python.exe', 'sshd.exe', 'mysqld.exe', 'redis-server.exe']

function generateMockStatus(port: number): PortInfo {
  const isUsed = Math.random() > 0.5
  return {
    port,
    process: isUsed ? mockProcesses[Math.floor(Math.random() * mockProcesses.length)] : '-',
    pid: isUsed ? 1000 + Math.floor(Math.random() * 50000) : 0,
    status: isUsed ? 'listening' : 'free',
    protocol: 'TCP'
  }
}

function queryPort() {
  const p = parseInt(portInput.value)
  if (!p || p < 1 || p > 65535) {
    message.warning('请输入有效端口号 (1-65535)')
    return
  }
  loading.value = true
  setTimeout(() => {
    const info = generateMockStatus(p)
    const idx = portList.value.findIndex(x => x.port === p)
    if (idx >= 0) portList.value.splice(idx, 1)
    portList.value.unshift(info)
    loading.value = false
    message.success(`端口 ${p}: ${info.status === 'listening' ? '被 ' + info.process + ' 占用 (PID: ' + info.pid + ')' : '空闲'}`)
  }, 300)
}

function refreshAll() {
  refreshing.value = true
  setTimeout(() => {
    const commonPorts = [21, 22, 25, 53, 80, 110, 143, 443, 993, 995, 3306, 3389, 5432, 6379, 8080, 8443, 9090]
    portList.value = commonPorts.map(p => generateMockStatus(p)).sort((a, b) => a.port - b.port)
    refreshing.value = false
    message.success('已刷新常用端口状态')
  }, 500)
}

function killProcess(row: PortInfo) {
  message.success(`已释放端口 ${row.port} (进程: ${row.process}, PID: ${row.pid})`)
  portList.value = portList.value.map(x => x.port === row.port ? { ...x, status: 'free' as const, process: '-', pid: 0 } : x)
}
</script>
