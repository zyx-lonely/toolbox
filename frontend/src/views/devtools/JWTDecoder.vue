<template>
  <div class="jwt-decoder">
    <n-space vertical>
      <n-input
        v-model:value="token"
        type="textarea"
        placeholder="在此粘贴 JWT Token..."
        :rows="4"
      />

      <n-button type="primary" @click="handleDecode" :disabled="!token.trim()">
        解码
      </n-button>

      <n-card v-if="header" title="Header" size="small">
        <pre class="jwt-block">{{ header }}</pre>
      </n-card>

      <n-card v-if="payload" title="Payload" size="small">
        <pre class="jwt-block">{{ payload }}</pre>
      </n-card>

      <n-alert v-if="error" type="error" :title="error" closable @close="error = ''" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const token = ref('')
const header = ref('')
const payload = ref('')
const error = ref('')

function handleDecode() {
  error.value = ''
  header.value = ''
  payload.value = ''

  try {
    const parts = token.value.trim().split('.')
    if (parts.length !== 3) {
      error.value = '无效的 JWT 格式'
      return
    }
    header.value = JSON.stringify(JSON.parse(atob(parts[0])), null, 2)
    payload.value = JSON.stringify(JSON.parse(atob(parts[1])), null, 2)
  } catch (e: any) {
    error.value = '解码失败: ' + (e.message || 'Token 格式无效')
  }
}
</script>

<style scoped>
.jwt-decoder {
  padding: 20px;
}
.jwt-block {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
