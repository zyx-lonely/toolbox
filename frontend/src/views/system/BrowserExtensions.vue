<template>
  <div>
    <n-h2>浏览器扩展管理</n-h2>
    <n-p>查看和管理 Chrome、Edge、Firefox 的扩展。</n-p>

    <n-card>
      <n-space style="margin-bottom: 12px" align="center">
        <n-select
          v-model:value="filterBrowser"
          :options="browserOptions"
          style="width: 150px"
          placeholder="筛选浏览器"
          clearable
        />
        <n-input v-model:value="searchQuery" placeholder="搜索扩展..." clearable style="width: 250px">
          <template #prefix><n-icon><search-outline /></n-icon></template>
        </n-input>
        <n-button @click="refresh" :loading="loading" type="primary">
          <template #icon><n-icon><refresh-outline /></n-icon></template>
          刷新
        </n-button>
      </n-space>

      <n-data-table
        :columns="columns"
        :data="filteredExtensions"
        :bordered="true"
        :single-line="false"
        :loading="loading"
        :pagination="{ pageSize: 20 }"
        size="small"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { SearchOutline, RefreshOutline, TrashOutline, EyeOffOutline, EyeOutline } from '@vicons/ionicons5'
import { GetBrowserExtensions, DisableBrowserExtension, EnableBrowserExtension, RemoveBrowserExtension } from '@wails/go/main/App'
import type { DataTableColumn } from 'naive-ui'

interface BrowserExtension {
  name: string
  version: string
  description: string
  path: string
  browser: string
  enabled: boolean
  id: string
}

const searchQuery = ref('')
const extensions = ref<BrowserExtension[]>([])
const loading = ref(false)
const message = useMessage()
const filterBrowser = ref<string | null>(null)

const browserOptions = [
  { label: 'Chrome', value: 'chrome' },
  { label: 'Edge', value: 'edge' },
  { label: 'Firefox', value: 'firefox' },
]

const browserLabel: Record<string, string> = {
  chrome: 'Chrome',
  edge: 'Edge',
  firefox: 'Firefox',
}

const browserColor: Record<string, string> = {
  chrome: '#4285F4',
  edge: '#0078D7',
  firefox: '#FF7139',
}

const filteredExtensions = computed(() => {
  let list = extensions.value
  if (filterBrowser.value) {
    list = list.filter(e => e.browser === filterBrowser.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(e =>
      e.name.toLowerCase().includes(q) ||
      e.description.toLowerCase().includes(q)
    )
  }
  return list
})

const columns: DataTableColumn<BrowserExtension>[] = [
  {
    title: '名称',
    key: 'name',
    width: 200,
    render: (row: BrowserExtension) => h('span', { style: 'font-weight:500' }, row.name)
  },
  {
    title: '浏览器',
    key: 'browser',
    width: 90,
    render: (row: BrowserExtension) => h('n-tag', { size: 'small', color: { color: browserColor[row.browser] + '22', textColor: browserColor[row.browser] } }, { default: () => browserLabel[row.browser] || row.browser })
  },
  { title: '版本', key: 'version', width: 90 },
  {
    title: '描述',
    key: 'description',
    ellipsis: { tooltip: true }
  },
  {
    title: '状态',
    key: 'enabled',
    width: 80,
    render: (row: BrowserExtension) => h('n-tag', { size: 'small', type: row.enabled ? 'success' : 'error' }, { default: () => row.enabled ? '启用' : '禁用' })
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (row: BrowserExtension) => [
      h('button', {
        class: 'n-button n-button--warning n-button--tiny',
        style: 'padding: 4px 8px; border: none; border-radius: 4px; cursor: pointer; margin-right: 6px; color: #fff; background: #f0a020;',
        onClick: () => toggleExtension(row)
      }, { default: () => row.enabled ? '禁用' : '启用' }),
      h('button', {
        class: 'n-button n-button--error n-button--tiny',
        style: 'padding: 4px 8px; border: none; border-radius: 4px; cursor: pointer; color: #fff; background: #d03050;',
        onClick: () => confirmRemove(row)
      }, { default: () => '删除' })
    ]
  }
]

async function refresh() {
  loading.value = true
  try {
    const r = await GetBrowserExtensions()
    if (r) extensions.value = r as BrowserExtension[]
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

async function toggleExtension(row: BrowserExtension) {
  try {
    const fn = row.enabled ? DisableBrowserExtension : EnableBrowserExtension
    const r = await fn(row.browser, row.id)
    if (r && (r as any).success) {
      row.enabled = !row.enabled
      message.success(row.enabled ? '已启用' : '已禁用')
    } else if (r && (r as any).error) {
      message.error((r as any).error)
    }
  } catch (e: any) {
    message.error(String(e))
  }
}

function confirmRemove(row: BrowserExtension) {
  // 简单确认
  if (confirm(`确定要删除扩展 "${row.name}" 吗？`)) {
    removeExtension(row)
  }
}

async function removeExtension(row: BrowserExtension) {
  try {
    const r = await RemoveBrowserExtension(row.browser, row.id)
    if (r && (r as any).success) {
      extensions.value = extensions.value.filter(e => e.id !== row.id || e.browser !== row.browser)
      message.success('已删除')
    } else if (r && (r as any).error) {
      message.error((r as any).error)
    }
  } catch (e: any) {
    message.error(String(e))
  }
}

onMounted(refresh)
</script>
