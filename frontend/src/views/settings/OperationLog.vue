<template>
  <div class="operation-log">
    <n-space vertical>
      <n-space align="center">
        <n-button @click="handleRefresh" :loading="loading" size="small">刷新</n-button>
        <n-button @click="handleClear" size="small" type="warning" ghost>清空日志</n-button>
      </n-space>

      <n-list v-if="logs.length">
        <n-list-item v-for="(log, idx) in logs" :key="idx">
          <n-space align="center">
            <n-tag :type="log.type || 'default'" size="small">{{ log.action }}</n-tag>
            <span class="log-time">{{ log.time }}</span>
          </n-space>
          <template #suffix>
            <span class="log-detail">{{ log.detail }}</span>
          </template>
        </n-list-item>
      </n-list>

      <n-empty v-else description="暂无操作记录">
        <template #extra>
          <n-icon size="48" color="#ccc">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path d="M256 48C141.13 48 48 141.13 48 256s93.13 208 208 208 208-93.13 208-208S370.87 48 256 48zm0 96a40 40 0 11-40 40 40 40 0 0140-40zm64 224H192v-32h32v-96h-32v-32h96v128h32z" fill="currentColor"/></svg>
          </n-icon>
        </template>
      </n-empty>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const loading = ref(false)
const logs = ref<any[]>([])

async function handleRefresh() {
  loading.value = true
  try {
    const result = await window.go.main.App.GetOperationLogs()
    logs.value = result || []
  } catch {
    logs.value = []
  } finally {
    loading.value = false
  }
}

async function handleClear() {
  try {
    await window.go.main.App.ClearOperationLogs()
    logs.value = []
    window.$message?.success('日志已清空')
  } catch (e: any) {
    window.$message?.error('清空失败')
  }
}
</script>

<style scoped>
.operation-log {
  padding: 20px;
}
.log-time {
  font-size: 12px;
  color: #888;
  margin-left: 8px;
}
.log-detail {
  font-size: 13px;
  color: #666;
}
</style>
