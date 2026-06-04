<template>
  <div>
    <n-h2>进程管理器</n-h2>
    <n-p>查看和管理系统正在运行的进程。</n-p>

    <n-card>
      <n-space style="margin-bottom: 12px">
        <n-input v-model:value="searchQuery" placeholder="搜索进程..." clearable style="width: 250px">
          <template #prefix><n-icon><search-outline /></n-icon></template>
        </n-input>
        <n-button @click="refresh" :loading="loading" type="primary">
          <template #icon><n-icon><refresh-outline /></n-icon></template>
          刷新
        </n-button>
      </n-space>

      <n-data-table
        :columns="columns"
        :data="filteredProcesses"
        :bordered="true"
        :single-line="false"
        :loading="loading"
        :pagination="pagination"
        size="small"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { SearchOutline, RefreshOutline, TrashOutline } from '@vicons/ionicons5'
import { GetProcessList, KillProcess } from '@wails/go/main/App'
import type { DataTableColumn } from 'naive-ui'

const searchQuery = ref('')
const processes = ref<any[]>([])
const loading = ref(false)
const message = useMessage()
const dialog = useDialog()

const pagination = { pageSize: 50 }

const columns: DataTableColumn[] = [
  { title: 'PID', key: 'pid', width: 80, sorter: (a: any, b: any) => a.pid - b.pid },
  { title: '进程名', key: 'name', ellipsis: true },
  { title: '内存', key: 'memory', width: 120 },
  { title: '命令', key: 'command', ellipsis: true },
  {
    title: '操作', key: 'actions', width: 80,
    render(row: any) {
      return h('button', {
        class: 'n-button n-button--error n-button--tiny',
        style: 'padding: 4px 8px; border: none; border-radius: 4px; cursor: pointer; color: #fff; background: #d03050;',
        onClick: () => confirmKill(row)
      }, { default: () => '结束' })
    }
  }
]

const filteredProcesses = computed(() => {
  if (!searchQuery.value) return processes.value
  const q = searchQuery.value.toLowerCase()
  return processes.value.filter((p: any) =>
    p.name?.toLowerCase().includes(q) || String(p.pid).includes(q)
  )
})

async function refresh() {
  loading.value = true
  try {
    const r = await GetProcessList()
    if (r) processes.value = r as any[]
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

function confirmKill(row: any) {
  dialog.warning({
    title: '结束进程',
    content: `确定要结束进程 ${row.name} (PID: ${row.pid}) 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await KillProcess(row.pid)
        message.success(`已结束进程 ${row.pid}`)
        refresh()
      } catch (e: any) {
        message.error(String(e))
      }
    }
  })
}

refresh()
</script>
