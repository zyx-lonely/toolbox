<template>
  <div>
    <n-h2>软件管理</n-h2>
    <n-p>查看和管理已安装的软件，支持批量卸载。</n-p>

    <n-card>
      <n-space style="margin-bottom: 12px" align="center">
        <n-input v-model:value="searchQuery" placeholder="搜索软件..." clearable style="width: 250px">
          <template #prefix><n-icon><search-outline /></n-icon></template>
        </n-input>
        <n-button @click="refresh" :loading="loading" type="primary">
          <template #icon><n-icon><refresh-outline /></n-icon></template>
          刷新
        </n-button>
        <n-button v-if="checkedRowKeys.length > 0" @click="confirmBatchUninstall" type="error">
          <template #icon><n-icon><trash-outline /></n-icon></template>
          批量卸载 ({{ checkedRowKeys.length }})
        </n-button>
        <n-button @click="exportList" :loading="exporting">
          <template #icon><n-icon><download-outline /></n-icon></template>
          导出列表
        </n-button>
      </n-space>

      <n-data-table
        :columns="columns"
        :data="filteredSoftware"
        :bordered="true"
        :single-line="false"
        :loading="loading"
        :pagination="pagination"
        size="small"
        :row-key="(row: any) => row.name"
        v-model:checked-row-keys="checkedRowKeys"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { SearchOutline, RefreshOutline, TrashOutline, DownloadOutline } from '@vicons/ionicons5'
import { GetInstalledSoftware, UninstallSoftware, BatchUninstallSoftware, ExportSoftwareList } from '@wails/go/main/App'
import type { DataTableColumn } from 'naive-ui'

const searchQuery = ref('')
const softwareList = ref<any[]>([])
const loading = ref(false)
const exporting = ref(false)
const checkedRowKeys = ref<string[]>([])
const message = useMessage()
const dialog = useDialog()

const pagination = { pageSize: 100 }

const columns: DataTableColumn[] = [
  { type: 'selection', width: 40 },
  { title: '软件名称', key: 'name', ellipsis: true, sorter: (a: any, b: any) => a.name.localeCompare(b.name) },
  { title: '版本', key: 'version', width: 100 },
  { title: '发布者', key: 'publisher', ellipsis: true },
  { title: '安装日期', key: 'installDate', width: 120 },
  { title: '大小', key: 'size', width: 100 },
  {
    title: '操作', key: 'actions', width: 100,
    render(row: any) {
      return h('button', {
        class: 'n-button n-button--error n-button--tiny',
        style: 'padding: 4px 8px; border: none; border-radius: 4px; cursor: pointer; color: #fff; background: #d03050;',
        onClick: () => confirmUninstall(row)
      }, { default: () => '卸载' })
    }
  }
]

const filteredSoftware = computed(() => {
  if (!searchQuery.value) return softwareList.value
  const q = searchQuery.value.toLowerCase()
  return softwareList.value.filter((s: any) =>
    s.name?.toLowerCase().includes(q) ||
    s.publisher?.toLowerCase().includes(q) ||
    s.version?.toLowerCase().includes(q)
  )
})

async function refresh() {
  loading.value = true
  try {
    const r = await GetInstalledSoftware()
    if (r) {
      softwareList.value = r as any[]
    }
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

function confirmUninstall(row: any) {
  dialog.warning({
    title: '卸载软件',
    content: `确定要卸载 ${row.name} 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const r = await UninstallSoftware(row.uninstall)
        if (r && r.success) {
          message.success(`已开始卸载 ${row.name}`)
        } else {
          message.error(r?.error || '卸载失败')
        }
      } catch (e: any) {
        message.error(String(e))
      }
    }
  })
}

function confirmBatchUninstall() {
  if (checkedRowKeys.value.length === 0) return
  dialog.warning({
    title: '批量卸载软件',
    content: `确定要卸载选中的 ${checkedRowKeys.value.length} 个软件吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const selectedSoftware = softwareList.value.filter((s: any) => checkedRowKeys.value.includes(s.name))
        const uninstallCmds = selectedSoftware.map((s: any) => s.uninstall)
        const r = await BatchUninstallSoftware(uninstallCmds)
        if (r && r.success) {
          message.success(`已开始批量卸载 ${checkedRowKeys.value.length} 个软件`)
          checkedRowKeys.value = []
        } else {
          message.error(r?.error || '批量卸载失败')
        }
      } catch (e: any) {
        message.error(String(e))
      }
    }
  })
}

async function exportList() {
  exporting.value = true
  try {
    const filePath = await showSaveDialog('software_list.json')
    if (!filePath) {
      exporting.value = false
      return
    }
    const r = await ExportSoftwareList(filePath, JSON.stringify(softwareList.value))
    if (r && r.success) {
      message.success('软件列表已导出')
    } else {
      message.error(r?.error || '导出失败')
    }
  } catch (e: any) {
    message.error(String(e))
  }
  exporting.value = false
}

function showSaveDialog(defaultName: string): Promise<string> {
  return new Promise((resolve) => {
    // 这里需要使用 Wails 的文件对话框
    // 暂时使用 prompt 作为替代
    const path = prompt('请输入保存路径:', `C:\\Users\\${getCurrentUser()}\\Downloads\\${defaultName}`)
    resolve(path || '')
  })
}

function getCurrentUser(): string {
  // 获取当前用户名
  return 'user'
}

refresh()
</script>
