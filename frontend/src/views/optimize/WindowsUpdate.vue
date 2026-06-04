<template>
  <div>
    <n-h2>Windows 更新管理</n-h2>
    <n-p>查看和管理 Windows 更新。</n-p>
    <n-space class="mb-4">
      <n-button @click="loadInstalled" :loading="loadingInstalled">已安装更新</n-button>
      <n-button @click="loadPending" :loading="loadingPending">待安装更新</n-button>
    </n-space>
    <n-empty v-if="!installedUpdates.length && !pendingUpdates.length" description="点击按钮查看更新列表" />
    <n-data-table v-if="installedUpdates.length" :columns="columns" :data="installedUpdates" size="small" :bordered="true" />
    <n-data-table v-if="pendingUpdates.length" :columns="columns" :data="pendingUpdates" size="small" :bordered="true" class="mt-4" />
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { GetWindowsUpdates, GetPendingUpdates } from '@wails/go/main/App'
const installedUpdates = ref<any[]>([])
const pendingUpdates = ref<any[]>([])
const loadingInstalled = ref(false)
const loadingPending = ref(false)
const message = useMessage()
const columns = [
  { title: '名称', key: 'name', ellipsis: { tooltip: true } },
  { title: 'KB 编号', key: 'kb' },
  { title: '状态', key: 'status', width: 80 },
  { title: '安装日期', key: 'installDate', width: 120 },
]
async function loadInstalled() {
  loadingInstalled.value = true
  try { const r = await GetWindowsUpdates(); if (r) installedUpdates.value = r as any[] }
  catch (e: any) { message.error(String(e)) }
  loadingInstalled.value = false
}
async function loadPending() {
  loadingPending.value = true
  try { const r = await GetPendingUpdates(); if (r) pendingUpdates.value = r as any[] }
  catch (e: any) { message.error(String(e)) }
  loadingPending.value = false
}
</script>
<style scoped>.mb-4{margin-bottom:16px}.mt-4{margin-top:16px}</style>
