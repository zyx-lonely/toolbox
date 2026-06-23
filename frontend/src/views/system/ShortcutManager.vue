<template>
  <div>
    <n-h2>快捷键管理</n-h2>
    <n-p>查看系统和应用的快捷键列表，检测可能的冲突。</n-p>

    <n-card>
      <n-space style="margin-bottom: 12px" align="center">
        <n-input v-model:value="searchQuery" placeholder="搜索快捷键..." clearable style="width: 250px">
          <template #prefix>
            <n-icon><search-outline /></n-icon>
          </template>
        </n-input>
        <n-select v-model:value="filterSource" :options="sourceOptions" placeholder="来源筛选" clearable style="width:150px" />
        <n-button @click="refresh" :loading="loading" type="primary">
          <template #icon>
            <n-icon><refresh-outline /></n-icon>
          </template>
          刷新
        </n-button>
      </n-space>

      <n-data-table
        :columns="columns"
        :data="filteredKeys"
        :bordered="true"
        :single-line="false"
        :loading="loading"
        :pagination="pagination"
        size="small"
        :row-key="(row) => row.id"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useMessage } from 'naive-ui'
import { SearchOutline, RefreshOutline } from '@vicons/ionicons5'
import { GetShortcutKeys } from '@wails/go/main/App'
import type { DataTableColumn } from 'naive-ui'

const searchQuery = ref('')
const filterSource = ref<string | null>(null)
const keys = ref<any[]>([])
const loading = ref(false)
const message = useMessage()

const pagination = { pageSize: 50 }

const sourceOptions = [
  { label: '系统', value: 'system' },
  { label: '注册表', value: 'registry' },
  { label: '已注册', value: 'registered' },
]

const columns: DataTableColumn[] = [
  { title: '名称', key: 'name', width: 200, sorter: (a: any, b: any) => a.name.localeCompare(b.name) },
  { title: '修饰键', key: 'modifiers', width: 150 },
  { title: '按键', key: 'key', width: 100 },
  { title: '完整快捷键', key: 'fullPath', width: 180 },
  { title: '功能', key: 'appName', ellipsis: true },
  {
    title: '来源', key: 'source', width: 100,
    render(row: any) {
      const typeMap: Record<string, string> = { system: 'info', registry: 'warning', registered: 'success' }
      return h('n-tag', { type: typeMap[row.source] || 'default', size: 'small' }, { default: () => row.source })
    }
  },
]

const filteredKeys = computed(() => {
  let result = keys.value
  if (filterSource.value) {
    result = result.filter((k: any) => k.source === filterSource.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter((k: any) =>
      k.name?.toLowerCase().includes(q) ||
      k.fullPath?.toLowerCase().includes(q) ||
      k.modifiers?.toLowerCase().includes(q) ||
      k.key?.toLowerCase().includes(q) ||
      k.appName?.toLowerCase().includes(q)
    )
  }
  return result
})

async function refresh() {
  loading.value = true
  try {
    const r = await GetShortcutKeys()
    if (r) {
      keys.value = r as any[]
    }
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

refresh()
</script>
