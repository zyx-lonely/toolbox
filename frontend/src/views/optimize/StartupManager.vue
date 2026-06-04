<template>
  <div>
    <n-h2>启动项管理</n-h2>
    <n-p>查看和管理系统开机自启程序。</n-p>

    <n-button type="primary" @click="loadItems" :loading="loading" class="mb-4">
      刷新启动项
    </n-button>

    <n-empty v-if="!loading && items.length === 0" description="暂无启动项数据" />

    <n-data-table
      v-if="items.length > 0"
      :columns="columns"
      :data="items"
      :bordered="true"
      :loading="loading"
      size="small"
    />
  </div>
</template>

<script setup lang="ts">
import { h, ref, onMounted } from 'vue'
import { NSwitch, NTag, useMessage } from 'naive-ui'
import { GetStartupItems, ToggleStartupItem } from '@wails/go/main/App'

interface StartupItem { name: string; command: string; location: string; publisher: string; enabled: boolean; impact: string }

const items = ref<StartupItem[]>([])
const loading = ref(false)
const message = useMessage()

const columns = [
  { title: '名称', key: 'name', width: 180 },
  { title: '命令', key: 'command', ellipsis: { tooltip: true } },
  { title: '位置', key: 'location', width: 120 },
  {
    title: '影响', key: 'impact', width: 80,
    render: (row: StartupItem) => {
      const color = { high: 'error' as const, medium: 'warning' as const, low: 'success' as const }[row.impact] || 'default'
      return h(NTag, { type: color, size: 'small' }, { default: () => row.impact })
    }
  },
  {
    title: '状态', key: 'enabled', width: 80,
    render: (row: StartupItem) => h(NSwitch, {
      value: row.enabled,
      onUpdateValue: (val: boolean) => toggleItem(row, val)
    })
  }
]

async function loadItems() {
  loading.value = true
  try {
    const result = await GetStartupItems()
    if (result) items.value = result as StartupItem[]
  } catch (e) { console.error(e) }
  loading.value = false
}

async function toggleItem(item: StartupItem, enable: boolean) {
  try {
    await ToggleStartupItem(item.name, enable)
    item.enabled = enable
    message.success(enable ? '已启用' : '已禁用')
  } catch (e: any) {
    message.error(`操作失败: ${e}`)
  }
}

onMounted(loadItems)
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
