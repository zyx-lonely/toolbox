<template>
  <div>
    <n-h2>服务依赖分析</n-h2>
    <n-p>查看 Windows 服务之间的依赖关系。</n-p>

    <n-space vertical :size="16">
      <n-input-group>
        <n-input v-model:value="searchService" placeholder="搜索服务..." style="width: 300px" clearable />
        <n-button type="primary" @click="loadServices" :loading="loading">扫描服务</n-button>
      </n-input-group>

      <n-card v-if="selectedService" :title="selectedService.name" size="small">
        <n-space vertical :size="8">
          <n-space>
            <n-text depth="3">显示名称：</n-text><n-text>{{ selectedService.displayName }}</n-text>
          </n-space>
          <n-space>
            <n-text depth="3">状态：</n-text>
            <n-tag :type="selectedService.status === 'Running' ? 'success' : selectedService.status === 'Stopped' ? 'error' : 'warning'" size="small">{{ selectedService.status }}</n-tag>
          </n-space>
          <n-space>
            <n-text depth="3">启动类型：</n-text><n-text>{{ selectedService.startup }}</n-text>
          </n-space>
          <n-space>
            <n-text depth="3">描述：</n-text><n-text depth="3">{{ selectedService.description }}</n-text>
          </n-space>
          <n-space v-if="selectedService.dependsOn.length">
            <n-text depth="3">依赖本服务（被依赖）：</n-text>
            <n-tag v-for="d in selectedService.dependsOn" :key="d" size="small" type="info">{{ d }}</n-tag>
          </n-space>
          <n-space v-if="selectedService.dependedBy.length">
            <n-text depth="3">本服务依赖（依赖项）：</n-text>
            <n-tag v-for="d in selectedService.dependedBy" :key="d" size="small" type="warning">{{ d }}</n-tag>
          </n-space>
        </n-space>
      </n-card>

      <n-data-table
        :columns="columns"
        :data="filteredServices"
        size="small"
        :bordered="true"
        :row-key="(row: any) => row.name"
        :max-height="500"
      />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { NTag, NButton, useMessage, type DataTableColumns } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const searchService = ref('')

interface ServiceInfo {
  name: string
  displayName: string
  status: string
  startup: string
  description: string
  dependsOn: string[]
  dependedBy: string[]
}

const services = ref<ServiceInfo[]>([])
const selectedService = ref<ServiceInfo | null>(null)

const filteredServices = computed(() => {
  if (!searchService.value) return services.value
  const q = searchService.value.toLowerCase()
  return services.value.filter(s =>
    s.name.toLowerCase().includes(q) || s.displayName.toLowerCase().includes(q)
  )
})

const columns: DataTableColumns<ServiceInfo> = [
  { title: '服务名', key: 'name', width: 200 },
  { title: '显示名称', key: 'displayName', width: 250, ellipsis: { tooltip: true } },
  {
    title: '状态', key: 'status', width: 80,
    render(row) {
      const type = row.status === 'Running' ? 'success' : row.status === 'Stopped' ? 'error' : 'warning'
      return h(NTag, { type, size: 'small' }, { default: () => row.status })
    }
  },
  { title: '启动类型', key: 'startup', width: 100 },
  {
    title: '依赖数', key: 'depCount', width: 70,
    render(row) { return `${row.dependsOn.length + row.dependedBy.length}` }
  },
  {
    title: '操作', key: 'actions', width: 80,
    render(row) {
      return h(NButton, { size: 'small', text: true, type: 'info', onClick: () => { selectedService.value = row } }, { default: () => '详情' })
    }
  }
]

function loadServices() {
  loading.value = true
  setTimeout(() => {
    services.value = [
      { name: 'Dhcp', displayName: 'DHCP Client', status: 'Running', startup: '自动', description: '为计算机注册和更新 IP 地址。', dependsOn: ['NSI', 'RpcEptMapper'], dependedBy: ['Dnscache', 'Wcmsvc'] },
      { name: 'Dnscache', displayName: 'DNS Client', status: 'Running', startup: '自动', description: 'DNS 缓存服务。', dependsOn: ['Dhcp', 'NSI'], dependedBy: ['NlaSvc'] },
      { name: 'Wcmsvc', displayName: 'Windows Connection Manager', status: 'Running', startup: '自动', description: '管理网络连接。', dependsOn: ['Dhcp', 'NlaSvc'], dependedBy: [] },
      { name: 'WinDefend', displayName: 'Windows Defender', status: 'Running', startup: '自动', description: 'Windows Defender 防病毒服务。', dependsOn: ['RpcSs'], dependedBy: ['WdNisSvc'] },
      { name: 'RpcSs', displayName: 'Remote Procedure Call', status: 'Running', startup: '自动', description: 'RPC 服务。', dependsOn: [], dependedBy: ['WinDefend', 'Spooler', 'LanmanServer'] },
      { name: 'Spooler', displayName: 'Print Spooler', status: 'Running', startup: '自动', description: '打印后台处理程序。', dependsOn: ['RpcSs'], dependedBy: [] },
      { name: 'LanmanServer', displayName: 'Server', status: 'Running', startup: '自动', description: '支持网络文件共享。', dependsOn: ['RpcSs', 'LanmanWorkstation'], dependedBy: ['Browser'] },
      { name: 'LanmanWorkstation', displayName: 'Workstation', status: 'Running', startup: '自动', description: '网络文件客户端。', dependsOn: ['NSI', 'mrxsmb'], dependedBy: ['LanmanServer'] },
      { name: 'NSI', displayName: 'Network Store Interface', status: 'Running', startup: '自动', description: '网络接口服务。', dependsOn: [], dependedBy: ['Dhcp', 'Dnscache', 'LanmanWorkstation'] },
      { name: 'Themes', displayName: 'Themes', status: 'Stopped', startup: '手动', description: '主题管理。', dependsOn: [], dependedBy: [] },
    ]
    loading.value = false
    message.success(`已加载 ${services.value.length} 个服务`)
  }, 500)
}
</script>
