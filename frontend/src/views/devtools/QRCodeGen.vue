<template>
  <div>
    <n-h2>二维码生成</n-h2>
    <n-card>
      <n-form-item label="内容">
        <n-input v-model:value="content" placeholder="输入要生成二维码的内容（文本、URL 等）" />
      </n-form-item>
      <n-form-item label="尺寸">
        <n-input-number v-model:value="size" :min="128" :max="1024" :step="64" />
      </n-form-item>
      <n-button type="primary" @click="generate" :loading="loading" :disabled="!content">生成二维码</n-button>
      <div v-if="qrResult" class="mt-4 qr-result">
        <img :src="qrResult.dataUri" alt="QR Code" />
        <div class="mt-2">
          <n-button size="small" @click="download">下载</n-button>
          <n-button size="small" class="ml-2" @click="copyDataUri">复制 Data URI</n-button>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { GenerateQRCode } from '@wails/go/main/App'

const message = useMessage()
const content = ref('')
const size = ref(256)
const loading = ref(false)
const qrResult = ref<any>(null)

async function generate() {
  loading.value = true
  try {
    qrResult.value = await GenerateQRCode(content.value, size.value)
    if (!qrResult.value.success) message.error(qrResult.value.error)
  } catch (e: any) {
    message.error('生成失败: ' + (e.message || e))
  }
  loading.value = false
}

function download() {
  if (!qrResult.value?.dataUri) return
  const a = document.createElement('a')
  a.href = qrResult.value.dataUri
  a.download = `qrcode_${Date.now()}.png`
  a.click()
}

function copyDataUri() {
  navigator.clipboard.writeText(qrResult.value.dataUri).then(() => message.success('已复制'))
}
</script>

<style scoped>
.qr-result { text-align: center; }
.qr-result img { max-width: 300px; border: 1px solid #eee; border-radius: 8px; padding: 8px; }
.mt-2 { margin-top: 8px; }
.mt-4 { margin-top: 16px; }
.ml-2 { margin-left: 8px; }
</style>
