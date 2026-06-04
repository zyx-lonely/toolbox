<template>
  <div>
    <n-h2>剪贴板历史</n-h2>
    <n-p>保存最近 50 条复制记录，支持搜索和再次复制。</n-p>
    <n-space class="mb-4">
      <n-button @click="refresh" :loading="loading">刷新</n-button>
      <n-button @click="clearAll" type="warning">清空全部</n-button>
      <n-input v-model:value="searchText" placeholder="搜索剪贴板内容..." style="width: 250px" />
    </n-space>
    <n-empty v-if="!filteredItems.length" description="暂无剪贴板记录" />
    <n-list v-if="filteredItems.length">
      <n-list-item v-for="item in filteredItems" :key="item.id">
        <n-space justify="space-between" align="center">
          <n-ellipsis :line-clamp="2" style="max-width: 400px">{{ item.content }}</n-ellipsis>
          <n-space>
            <n-tag size="tiny" type="info">{{ item.type }}</n-tag>
            <n-text depth="3" style="font-size: 12px">{{ item.time }}</n-text>
            <n-button size="tiny" @click="copyItem(item.content)">复制</n-button>
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
import { GetClipboardHistory, ClearClipboardHistory, RemoveClipboardItem } from '@wails/go/main/App'
const items = ref<any[]>([]); const loading = ref(false); const searchText = ref(''); const message = useMessage()
const filteredItems = computed(() => items.value.filter(i => !searchText.value || i.content?.includes(searchText.value)))
async function refresh() { loading.value = true; try { const r = await GetClipboardHistory(); if (r) items.value = r as any[] } catch(e) {console.error(e)}; loading.value = false }
async function clearAll() { await ClearClipboardHistory(); items.value = []; message.success('已清空') }
async function removeItem(id: number) { await RemoveClipboardItem(id); items.value = items.value.filter(i => i.id !== id) }
async function copyItem(content: string) { try { await navigator.clipboard.writeText(content); message.success('已复制') } catch { message.warning('复制失败') } }
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
