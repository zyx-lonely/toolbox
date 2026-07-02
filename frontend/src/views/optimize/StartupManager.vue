<template>
  <div>
    <n-h2>启动项管理</n-h2>
    <n-p>查看和管理系统开机自启程序。</n-p>

    <n-space class="mb-4">
      <n-button type="primary" @click="loadItems" :loading="loading">
        刷新启动项
      </n-button>
      <n-input v-model:value="searchKeyword" placeholder="搜索启动项..." style="width: 200px" clearable />
    </n-space>

    <n-space class="mb-4" v-if="group">
      <n-statistic label="总数" :value="group.total" />
      <n-statistic label="已启用" :value="group.enabled" />
      <n-statistic label="已禁用" :value="group.disabled" />
    </n-space>

    <n-empty v-if="!loading && filteredItems.length === 0" description="暂无启动项数据" />

    <n-data-table
      v-if="filteredItems.length > 0"
      :columns="columns"
      :data="filteredItems"
      :bordered="true"
      :loading="loading"
      size="small"
      :pagination="{ pageSize: 20 }"
    />
  </div>
</template>

<script setup lang="ts">
import { h, ref, computed, onMounted } from 'vue'
import { NSwitch, NTag, NButton, useMessage } from 'naive-ui'
import { GetAllStartupPrograms, ToggleStartupProgram, DeleteStartupProgram } from '@wails/go/main/App'

interface StartupProgram {
  name: string
  path: string
  command: string
  location: string
  enabled: boolean
  publisher: string
  startTime: string
  size: number
  lastModified: string
}

interface StartupGroup {
  total: number
  enabled: number
  disabled: number
  hkcuCount: number
  hklmCount: number
  folderCount: number
}

const items = ref<StartupProgram[]>([])
const group = ref<StartupGroup | null>(null)
const loading = ref(false)
const searchKeyword = ref('')
const message = useMessage()

const filteredItems = computed(() => {
  if (!searchKeyword.value) return items.value
  const kw = searchKeyword.value.toLowerCase()
  return items.value.filter(item =>
    item.name.toLowerCase().includes(kw) ||
    item.command.toLowerCase().includes(kw) ||
    item.location.toLowerCase().includes(kw)
  )
})

const columns = [
  { title: '名称', key: 'name', width: 180, ellipsis: { tooltip: true } },
  { title: '路径', key: 'path', ellipsis: { tooltip: true } },
  {
    title: '位置', key: 'location', width: 120,
    render: (row: StartupProgram) => {
      const colorMap: Record<string, string> = {
        'HKCU-Run': 'success', 'HKCU-RunOnce': 'success',
        'HKLM-Run': 'warning', 'HKLM-RunOnce': 'warning',
        'StartupFolder': 'info'
      }
      return h(NTag, { type: (colorMap[row.location] || 'default') as any, size: 'small' }, { default: () => row.location })
    }
  },
  {
    title: '状态', key: 'enabled', width: 80,
    render: (row: StartupProgram) => h(NSwitch, {
      value: row.enabled,
      onUpdateValue: (val: boolean) => toggleItem(row, val)
    })
  },
  {
    title: '操作', width: 80,
    render: (row: StartupProgram) => h(NButton, {
      size: 'small', type: 'error', quaternary: true,
      onClick: () => deleteItem(row)
    }, { default: () => '删除' })
  }
]

async function loadItems() {
  loading.value = true
  try {
    const result = await GetAllStartupPrograms()
    if (result && Array.isArray(result)) {
      items.value = result[0] as StartupProgram[]
      group.value = result[1] as StartupGroup
    }
  } catch (e) { console.error(e) }
  loading.value = false
}

async function toggleItem(item: StartupProgram, enable: boolean) {
  try {
    await ToggleStartupProgram(item.name, item.location, enable)
    item.enabled = enable
    message.success(enable ? '已启用' : '已禁用')
    loadItems()
  } catch (e: any) {
    message.error(`操作失败: ${e}`)
  }
}

async function deleteItem(item: StartupProgram) {
  try {
    await DeleteStartupProgram(item.name, item.location)
    message.success('已删除')
    loadItems()
  } catch (e: any) {
    message.error(`删除失败: ${e}`)
  }
}

onMounted(loadItems)
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
