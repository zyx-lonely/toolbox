<template>
  <div class="ip-conflict">
    <n-space vertical>
      <n-space align="center">
        <n-input v-model:value="localIP" placeholder="请输入本地 IP 地址" style="width: 240px" />
        <n-button type="primary" @click="handleCheck" :loading="checking">
          检测冲突
        </n-button>
      </n-space>

      <n-empty v-if="!checking && results.length === 0" description="暂无检测结果" />

      <n-list v-else>
        <n-list-item v-for="(item, idx) in results" :key="idx">
          <template #prefix>
            <n-tag :type="item.conflict ? 'error' : 'success'">
              {{ item.conflict ? '冲突' : '正常' }}
            </n-tag>
          </template>
          {{ item.message }}
        </n-list-item>
      </n-list>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const localIP = ref('')
const checking = ref(false)
const results = ref<any[]>([])

async function handleCheck() {
  if (!localIP.value) return
  checking.value = true
  try {
    const result = await window.go.main.App.CheckIPConflict(localIP.value)
    results.value = Array.isArray(result) ? result : [{ conflict: true, message: String(result) }]
  } catch (e: any) {
    results.value = [{ conflict: true, message: e.message || '检测失败' }]
  } finally {
    checking.value = false
  }
}
</script>

<style scoped>
.ip-conflict {
  padding: 20px;
}
</style>
