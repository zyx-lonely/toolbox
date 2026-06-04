<template>
  <div>
    <n-h2>浏览器数据管理</n-h2>
    <n-p>一键清理浏览器缓存、Cookie、历史记录等数据。</n-p>

    <n-alert type="info" :bordered="false" class="mb-4">
      <template #header>提示</template>
      清理后需要重新登录部分网站。浏览器关闭后才能清理部分数据。
    </n-alert>

    <n-button type="primary" @click="scan" :loading="scanning" class="mb-4">
      扫描浏览器数据
    </n-button>

    <n-empty v-if="!scanning && items.length === 0" description="点击扫描按钮开始检测浏览器数据" />

    <n-space vertical size="medium" v-if="items.length">
      <n-card v-for="(group, idx) in groupedItems" :key="idx" :title="group.browserName" size="small">
        <template #header-extra>
          <n-checkbox :checked="group.allChecked" @update:checked="toggleGroup(group)">全选</n-checkbox>
        </template>
        <n-space vertical>
          <n-space v-for="item in group.items" :key="item.browserId + '-' + item.category"
            justify="space-between" align="center" style="padding: 4px 0">
            <n-space>
              <n-checkbox v-model:checked="item.checked" :disabled="item.fileCount === 0" />
              <n-tag size="small" type="info">{{ item.label }}</n-tag>
              <span v-if="item.fileCount > 0" style="color: #888; font-size: 13px">
                {{ item.fileCount }} 个文件 · {{ formatBytes(item.totalSize) }}
              </span>
              <span v-else style="color: #ccc; font-size: 13px">无数据</span>
            </n-space>
          </n-space>
        </n-space>
      </n-card>

      <n-button type="warning" @click="cleanSelected" :loading="cleaning" :disabled="!hasChecked">
        清理选中项 ({{ checkedCount }})
      </n-button>
    </n-space>

    <n-space vertical size="small" v-if="results.length">
      <n-h3>清理结果</n-h3>
      <n-alert v-for="(r, i) in results" :key="i"
        :type="r.success ? 'success' : 'error'" :bordered="false" closable>
        <template #header>{{ r.label }} ({{ r.browserId }})</template>
        {{ r.success ? `已清理 ${r.fileCount} 个文件，释放 ${formatBytes(r.freedBytes)}` : `失败: ${r.error}` }}
      </n-alert>
      <n-space>
        <n-button type="primary" @click="resetPage">完成</n-button>
        <n-button @click="scan">重新扫描</n-button>
      </n-space>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { ScanBrowserData, CleanBrowserData } from '@wails/go/main/App'
import { formatBytes } from '@/api/bridge'

interface CleanItem { browserId: string; browserName: string; category: string; label: string; path: string; fileCount: number; totalSize: number; checked: boolean }
interface CleanResult { browserId: string; category: string; label: string; success: boolean; fileCount: number; freedBytes: number; error?: string }

const items = ref<CleanItem[]>([])
const results = ref<CleanResult[]>([])
const scanning = ref(false)
const cleaning = ref(false)
const message = useMessage()

const groupedItems = computed(() => {
  const groups: Record<string, { browserName: string; items: CleanItem[] }> = {}
  for (const item of items.value) {
    if (!groups[item.browserId]) {
      groups[item.browserId] = { browserName: item.browserName, items: [] }
    }
    groups[item.browserId].items.push(item)
  }
  return Object.values(groups).map(g => ({
    ...g,
    allChecked: g.items.every(i => i.checked || i.fileCount === 0)
  }))
})

const checkedCount = computed(() => items.value.filter(i => i.checked && i.fileCount > 0).length)
const hasChecked = computed(() => checkedCount.value > 0)

async function scan() {
  scanning.value = true
  results.value = []
  try {
    const result = await ScanBrowserData()
    if (result) items.value = result as CleanItem[]
    const count = result?.length || 0
    if (count === 0) message.info('未检测到浏览器数据')
    else message.success(`检测到 ${count} 项可清理数据`)
  } catch (e: any) { message.error(`扫描失败: ${e}`) }
  scanning.value = false
}

function toggleGroup(group: any) {
  const newVal = !group.allChecked
  group.items.forEach((i: CleanItem) => {
    if (i.fileCount > 0) i.checked = newVal
  })
}

async function cleanSelected() {
  const selected = items.value.filter(i => i.checked && i.fileCount > 0)
  if (!selected.length) return
  cleaning.value = true
  results.value = []
  try {
    const result = await CleanBrowserData(selected)
    if (result) results.value = result as CleanResult[]
    message.success(`清理完成`)
  } catch (e: any) { message.error(`清理失败: ${e}`) }
  cleaning.value = false
}

function resetPage() { items.value = []; results.value = [] }
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
