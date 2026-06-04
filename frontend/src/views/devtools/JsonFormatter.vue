<template>
  <div>
    <n-h2>JSON 格式化</n-h2>
    <n-p>JSON 格式化、校验与压缩。</n-p>

    <n-space vertical>
      <n-input
        v-model:value="inputJSON"
        type="textarea"
        :rows="10"
        placeholder="在此粘贴 JSON 文本..."
      />

      <n-space>
        <n-button type="primary" @click="formatJSON">格式化</n-button>
        <n-button @click="compressJSON">压缩</n-button>
        <n-button @click="validateJSON">校验</n-button>
      </n-space>

      <n-input
        v-model:value="outputJSON"
        type="textarea"
        :rows="10"
        placeholder="结果将显示在这里..."
        readonly
      />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { FormatJSON, MinifyJSON } from '@wails/go/main/App'

const inputJSON = ref('')
const outputJSON = ref('')
const message = useMessage()

async function formatJSON() {
  try {
    const r = await FormatJSON(inputJSON.value)
    if (r?.valid) {
      outputJSON.value = r.formatted
      message.success(`格式化成功，大小: ${r.size} 字符`)
    } else {
      message.error(r?.error || 'JSON 格式无效')
    }
  } catch (e: any) {
    message.error(`错误: ${e}`)
  }
}

async function compressJSON() {
  try {
    const r = await MinifyJSON(inputJSON.value)
    if (r?.valid) {
      outputJSON.value = r.formatted
      message.success(`压缩成功，大小: ${r.size} 字符`)
    } else {
      message.error(r?.error || 'JSON 格式无效')
    }
  } catch (e: any) {
    message.error(`错误: ${e}`)
  }
}

async function validateJSON() {
  try {
    const r = await FormatJSON(inputJSON.value)
    if (r?.valid) message.success('JSON 格式有效 ✓')
    else message.error(r?.error || 'JSON 格式无效')
  } catch (e: any) {
    message.error(`错误: ${e}`)
  }
}
</script>
