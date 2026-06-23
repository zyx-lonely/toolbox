<template>
  <div>
    <n-h2>剪贴板历史</n-h2>
    <n-p>保存最近 50 条复制记录，支持搜索和再次复制。</n-p>
    <n-space class="mb-4">
      <n-button @click="refresh" :loading="loading">刷新</n-button>
      <n-button @click="readImage" type="info" ghost>读取图片剪贴板</n-button>
      <n-button @click="clearAll" type="warning">清空全部</n-button>
      <n-input v-model:value="searchText" placeholder="搜索剪贴板内容..." style="width: 250px" />
    </n-space>
    <n-empty v-if="!filteredItems.length" description="暂无剪贴板记录" />
    <n-list v-if="filteredItems.length">
      <n-list-item v-for="item in filteredItems" :key="item.id">
        <n-space justify="space-between" align="center" style="width:100%">
          <div style="flex:1; min-width:0">
            <template v-if="item.type === 'image'">
              <img v-if="item.content && item.content.startsWith('data:')" :src="item.content"
                   style="max-height:60px; max-width:200px; border-radius:4px; border:1px solid #eee;" />
              <n-tag v-else size="small" type="info">图片 ({{ formatSize(item.size) }})</n-tag>
            </template>
            <template v-else>
              <n-ellipsis :line-clamp="2" style="max-width: 400px">{{ item.content }}</n-ellipsis>
            </template>
          </div>
          <n-space>
            <n-tag size="tiny" :type="item.type === 'image' ? 'success' : 'info'">{{ item.type }}</n-tag>
            <n-text depth="3" style="font-size: 12px">{{ item.time }}</n-text>
            <n-button size="tiny" @click="copyItem(item)">复制</n-button>
            <n-button size="tiny" @click="removeItem(item.id)">删除</n-button>
          </n-space>
        </n-space>
      </n-list-item>
    </n-list>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { GetClipboardHistory, ClearClipboardHistory, RemoveClipboardItem, ReadClipboardImage, AddClipboardItem } from '@wails/go/main/App'
import type { ClipItem } from '@/types/api'

const items = ref<ClipItem[]>([])
const loading = ref(false)
const searchText = ref('')
const message = useMessage()

const filteredItems = computed(() =>
  items.value.filter(i => !searchText.value || i.content?.includes(searchText.value))
)

async function refresh() {
  loading.value = true
  try {
    const r = await GetClipboardHistory()
    if (r && r.success) {
      items.value = r.data as ClipItem[]
    } else if (r && r.error) {
      message.error(r.error)
    }
  } catch (e) { console.error(e) }
  loading.value = false
}

async function clearAll() {
  await ClearClipboardHistory()
  items.value = []
  message.success('已清空')
}

async function removeItem(id: number) {
  await RemoveClipboardItem(id)
  items.value = items.value.filter(i => i.id !== id)
}

async function readImage() {
  try {
    const r = await ReadClipboardImage()
    if (r && r.success && r.data) {
      const r2 = await AddClipboardItem(r.data as string, 'image')
      if (r2 && r2.success) {
        message.success('图片已添加到剪贴板历史')
        await refresh()
      }
    } else {
      message.warning(r?.error || '剪贴板中没有图片')
    }
  } catch (e: any) {
    message.error(String(e))
  }
}

async function copyItem(item: ClipItem) {
  try {
    if (item.type === 'image' && item.content?.startsWith('data:')) {
      const resp = await fetch(item.content)
      const blob = await resp.blob()
      await navigator.clipboard.write([
        new ClipboardItem({ [blob.type]: blob })
      ])
    } else {
      await navigator.clipboard.writeText(item.content || '')
    }
    message.success('已复制')
  } catch {
    message.warning('复制失败')
  }
}

function formatSize(bytes: number) {
  if (bytes >= 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + ' MB'
  if (bytes >= 1024) return (bytes / 1024).toFixed(0) + ' KB'
  return bytes + ' B'
}
</script>
<style scoped>.mb-4 { margin-bottom: 16px }</style>
