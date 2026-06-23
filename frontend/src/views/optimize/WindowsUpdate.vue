<template>
  <div>
    <n-h2>Windows 更新管理</n-h2>
    <n-p>查看和管理系统 Windows 更新。</n-p>

    <n-space vertical :size="16">
      <n-space>
        <n-button type="primary" @click="scanUpdates" :loading="scanning">扫描更新</n-button>
        <n-button @click="checkHistory">查看更新历史</n-button>
      </n-space>

      <n-card v-if="updates.length" title="可用更新" size="small">
        <n-data-table
          :columns="columns"
          :data="updates"
          size="small"
          :bordered="true"
          :row-key="(row: any) => row.id"
          :checked-row-keys="checkedKeys"
          @update:checked-row-keys="onCheck"
        />
        <template #footer>
          <n-space>
            <n-button type="warning" :disabled="!checkedKeys.length" @click="installSelected">安装选中 ({{ checkedKeys.length }})</n-button>
            <n-button type="error" :disabled="!checkedKeys.length" @click="hideSelected">隐藏选中</n-button>
          </n-space>
        </template>
      </n-card>

      <n-card v-if="history.length" title="更新历史" size="small">
        <n-timeline>
          <n-timeline-item
            v-for="item in history"
            :key="item.id"
            :type="item.success ? 'success' : 'error'"
            :title="item.name"
            :content="`${item.date} - ${item.success ? '安装成功' : '安装失败'}`"
          />
        </n-timeline>
      </n-card>

      <n-empty v-if="!updates.length && !history.length && !scanning" description="点击上方按钮扫描可用更新" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { NButton, NTag, NCheckbox, useMessage, type DataTableColumns, type DataTableRowKey } from 'naive-ui'

const message = useMessage()
const scanning = ref(false)
const checkedKeys = ref<DataTableRowKey[]>([])

interface UpdateItem {
  id: string
  name: string
  kb: string
  size: string
  severity: string
  date: string
}

interface HistoryItem {
  id: string
  name: string
  date: string
  success: boolean
}

const updates = ref<UpdateItem[]>([])
const history = ref<HistoryItem[]>([])

const columns: DataTableColumns<UpdateItem> = [
  { type: 'selection' },
  { title: '名称', key: 'name', width: 280 },
  { title: 'KB 编号', key: 'kb', width: 120 },
  { title: '大小', key: 'size', width: 80 },
  {
    title: '严重级别', key: 'severity', width: 100,
    render(row) {
      const type = row.severity === '重要' ? 'error' : row.severity === '推荐' ? 'warning' : 'info'
      return h(NTag, { type, size: 'small' }, { default: () => row.severity })
    }
  },
  { title: '发布日期', key: 'date', width: 120 },
]

function onCheck(keys: DataTableRowKey[]) {
  checkedKeys.value = keys
}

function scanUpdates() {
  scanning.value = true
  setTimeout(() => {
    updates.value = [
      { id: '1', name: '2026-06 累积更新 (KB5034441)', kb: 'KB5034441', size: '1.2 GB', severity: '重要', date: '2026-06-11' },
      { id: '2', name: 'Windows Defender 安全 intelligence', kb: 'KB2267602', size: '156 MB', severity: '重要', date: '2026-06-17' },
      { id: '3', name: '.NET Framework 4.8.1 更新', kb: 'KB5032252', size: '78 MB', severity: '推荐', date: '2026-06-11' },
      { id: '4', name: 'Windows Malicious Software Removal Tool', kb: 'KB890830', size: '68 MB', severity: '推荐', date: '2026-06-11' },
      { id: '5', name: 'Servicing Stack Update', kb: 'KB5035942', size: '16 MB', severity: '推荐', date: '2026-06-05' },
    ]
    checkedKeys.value = []
    scanning.value = false
    message.success(`扫描完成，发现 ${updates.value.length} 个可用更新`)
  }, 1000)
}

function checkHistory() {
  history.value = [
    { id: 'h1', name: '2026-05 累积更新 (KB5032189)', date: '2026-05-15', success: true },
    { id: 'h2', name: '.NET 8.0 运行时更新', date: '2026-05-10', success: true },
    { id: 'h3', name: 'Windows Defender 更新', date: '2026-05-08', success: true },
    { id: 'h4', name: '2026-04 累积更新 (KB5031356)', date: '2026-04-12', success: false },
  ]
  message.info('已加载更新历史')
}

function installSelected() {
  message.success(`正在安装 ${checkedKeys.value.length} 个更新...`)
  checkedKeys.value = []
}

function hideSelected() {
  message.success(`已隐藏 ${checkedKeys.value.length} 个更新`)
  updates.value = updates.value.filter(u => !checkedKeys.value.includes(u.id))
  checkedKeys.value = []
}
</script>
