<template>
  <div>
    <n-h2>进程管理器</n-h2>
    <n-p>查看和管理系统正在运行的进程。</n-p>

    <n-card>
      <n-space style="margin-bottom: 12px" align="center">
        <n-input v-model:value="searchQuery" placeholder="搜索进程..." clearable style="width: 250px">
          <template #prefix><n-icon><search-outline /></n-icon></template>
        </n-input>
        <n-switch v-model:value="useVirtualScroll" size="small">
          <template #checked>虚拟滚动</template>
          <template #unchecked>分页</template>
        </n-switch>
        <n-button @click="refresh" :loading="loading" type="primary">
          <template #icon><n-icon><refresh-outline /></n-icon></template>
          刷新
        </n-button>
        <n-button v-if="checkedRowKeys.length > 0" @click="confirmBatchKill" type="error">
          <template #icon><n-icon><trash-outline /></n-icon></template>
          批量结束 ({{ checkedRowKeys.length }})
        </n-button>
      </n-space>

      <n-data-table
        v-if="!useVirtualScroll"
        :columns="columns"
        :data="filteredProcesses"
        :bordered="true"
        :single-line="false"
        :loading="loading"
        :pagination="pagination"
        size="small"
        :row-key="(row: any) => row.pid"
        v-model:checked-row-keys="checkedRowKeys"
      />

      <n-data-table
        v-else
        :columns="columns"
        :data="filteredProcesses"
        :bordered="true"
        :single-line="false"
        :loading="loading"
        :max-height="600"
        :virtual-scroll="true"
        size="small"
        :row-height="32"
        :scroll-x="800"
        :row-key="(row: any) => row.pid"
        v-model:checked-row-keys="checkedRowKeys"
      />
    </n-card>

    <!-- 优先级设置弹窗 -->
    <n-modal v-model:show="showPriorityModal" preset="dialog" title="设置进程优先级">
      <n-space vertical style="margin-top: 12px">
        <n-text>进程：{{ selectedProcess?.name }} (PID: {{ selectedProcess?.pid }})</n-text>
        <n-select
          v-model:value="selectedPriority"
          :options="priorityOptions"
          placeholder="选择优先级"
        />
        <n-space justify="end">
          <n-button @click="showPriorityModal = false">取消</n-button>
          <n-button type="primary" @click="setPriority" :loading="settingPriority">确定</n-button>
        </n-space>
      </n-space>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { SearchOutline, RefreshOutline, TrashOutline } from '@vicons/ionicons5'
import { GetProcessList, KillProcess, SetProcessPriority, BatchKillProcesses } from '@wails/go/main/App'
import type { DataTableColumn } from 'naive-ui'

const searchQuery = ref('')
const processes = ref<any[]>([])
const loading = ref(false)
const useVirtualScroll = ref(false)
const settingPriority = ref(false)
const showPriorityModal = ref(false)
const selectedProcess = ref<any>(null)
const selectedPriority = ref('normal')
const checkedRowKeys = ref<number[]>([])
const message = useMessage()
const dialog = useDialog()

const pagination = { pageSize: 50 }

const priorityOptions = [
  { label: '低（空闲）', value: 'idle' },
  { label: '低于正常', value: 'below_normal' },
  { label: '正常', value: 'normal' },
  { label: '高于正常', value: 'above_normal' },
  { label: '高', value: 'high' },
  { label: '实时（慎用）', value: 'realtime' }
]

const columns: DataTableColumn[] = [
  { type: 'selection', width: 40 },
  { title: 'PID', key: 'pid', width: 80, sorter: (a: any, b: any) => a.pid - b.pid },
  { title: '进程名', key: 'name', ellipsis: true },
  { title: '内存', key: 'memory', width: 120 },
  { title: '命令', key: 'command', ellipsis: true },
  {
    title: '操作', key: 'actions', width: 180,
    render(row: any) {
      return h('div', { style: 'display: flex; gap: 4px' }, [
        h('button', {
          class: 'n-button n-button--primary n-button--tiny',
          style: 'padding: 4px 8px; border: none; border-radius: 4px; cursor: pointer; color: #fff; background: #18a058;',
          onClick: () => openPriorityModal(row)
        }, { default: () => '优先级' }),
        h('button', {
          class: 'n-button n-button--error n-button--tiny',
          style: 'padding: 4px 8px; border: none; border-radius: 4px; cursor: pointer; color: #fff; background: #d03050;',
          onClick: () => confirmKill(row)
        }, { default: () => '结束' })
      ])
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

function openPriorityModal(row: any) {
  selectedProcess.value = row
  selectedPriority.value = 'normal'
  showPriorityModal.value = true
}

async function setPriority() {
  if (!selectedProcess.value) return
  settingPriority.value = true
  try {
    const r = await SetProcessPriority(selectedProcess.value.pid, selectedPriority.value)
    if (r && r.success) {
      message.success(`已将进程 ${selectedProcess.value.name} 的优先级设置为 ${selectedPriority.value}`)
      showPriorityModal.value = false
    } else {
      message.error(r?.error || '设置失败')
    }
  } catch (e: any) {
    message.error(String(e))
  }
  settingPriority.value = false
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

function confirmBatchKill() {
  if (checkedRowKeys.value.length === 0) return
  dialog.warning({
    title: '批量结束进程',
    content: `确定要结束选中的 ${checkedRowKeys.value.length} 个进程吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const r = await BatchKillProcesses(checkedRowKeys.value)
        if (r && r.success) {
          message.success(`已批量结束 ${checkedRowKeys.value.length} 个进程`)
          checkedRowKeys.value = []
          refresh()
        } else {
          message.error(r?.error || '批量结束失败')
        }
      } catch (e: any) {
        message.error(String(e))
      }
    }
  })
}

refresh()
</script>
