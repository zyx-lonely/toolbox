<template>
  <div>
    <h2>剪贴板翻译</h2>
    <p>读取剪贴板内容并自动翻译。</p>
    <n-space class="mb-4" vertical>
      <n-space>
        <n-button type="primary" @click="translateClipboard" :loading="loading">读取剪贴板并翻译</n-button>
        <n-select v-model:value="targetLang" :options="langOptions" style="width: 150px" />
      </n-space>
      <n-input v-model:value="inputText" type="textarea" placeholder="或直接输入要翻译的文本..." :rows="3" />
      <n-space>
        <n-button type="primary" @click="translateInput" :loading="loading" :disabled="!inputText">翻译</n-button>
        <n-button @click="copyResult" :disabled="!result">复制译文</n-button>
      </n-space>
    </n-space>
    <n-card v-if="result" title="翻译结果" size="small">
      <n-descriptions bordered :column="1">
        <n-descriptions-item label="原文">{{ result.source }}</n-descriptions-item>
        <n-descriptions-item label="译文">{{ result.target }}</n-descriptions-item>
        <n-descriptions-item label="方向">{{ result.sourceLang }} → {{ result.targetLang }}</n-descriptions-item>
        <n-descriptions-item label="引擎">{{ result.provider }}</n-descriptions-item>
      </n-descriptions>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { GetClipboardAndTranslate, Translate } from '@wails/go/main/App'

interface TranslateResult { source: string; target: string; sourceLang: string; targetLang: string; provider: string }

const inputText = ref('')
const result = ref<TranslateResult | null>(null)
const loading = ref(false)
const targetLang = ref('en')
const message = useMessage()

const langOptions = [
  { label: '翻译为英文', value: 'en' },
  { label: '翻译为中文', value: 'zh' },
  { label: '翻译为日文', value: 'ja' },
  { label: '翻译为韩文', value: 'ko' },
  { label: '翻译为法文', value: 'fr' },
  { label: '翻译为德文', value: 'de' },
  { label: '翻译为西班牙文', value: 'es' },
]

async function translateClipboard() {
  loading.value = true
  try {
    const r = await GetClipboardAndTranslate(targetLang.value) as TranslateResult
    result.value = r
    inputText.value = r.source
    message.success('翻译完成')
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

async function translateInput() {
  if (!inputText.value) return
  loading.value = true
  try {
    const r = await Translate(inputText.value, '', targetLang.value) as TranslateResult
    result.value = r
    message.success('翻译完成')
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

function copyResult() {
  if (result.value?.target) {
    navigator.clipboard.writeText(result.value.target)
    message.success('已复制到剪贴板')
  }
}
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
