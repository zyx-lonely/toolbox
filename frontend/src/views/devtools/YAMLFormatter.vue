<template>
  <div class="yaml-formatter">
    <n-space vertical>
      <n-input
        v-model:value="inputText"
        type="textarea"
        placeholder="输入 YAML 或 TOML 内容..."
        :rows="8"
      />

      <n-space align="center">
        <n-button type="primary" @click="handleFormat" :loading="formatting">
          格式化
        </n-button>
        <n-button @click="handleClear">清空</n-button>
      </n-space>

      <n-input
        v-model:value="outputText"
        type="textarea"
        placeholder="格式化结果..."
        :rows="8"
        readonly
      />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const inputText = ref('')
const outputText = ref('')
const formatting = ref(false)

async function handleFormat() {
  if (!inputText.value.trim()) return
  formatting.value = true
  try {
    const result = await window.go.main.App.FormatYAML(inputText.value)
    outputText.value = result || ''
  } catch (e: any) {
    outputText.value = '格式化失败: ' + (e.message || '未知错误')
  } finally {
    formatting.value = false
  }
}

function handleClear() {
  inputText.value = ''
  outputText.value = ''
}
</script>

<style scoped>
.yaml-formatter {
  padding: 20px;
}
</style>
