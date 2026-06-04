<template>
  <div>
    <n-h2>图片批量压缩</n-h2>
    <n-p>选择目录批量压缩图片文件，支持质量调节。</n-p>

    <n-card>
      <n-space vertical>
        <n-space>
          <n-input-group>
            <n-input v-model:value="inputDir" placeholder="选择要压缩的目录" readonly style="width: 300px" />
            <n-button @click="selectDir">选择目录</n-button>
          </n-input-group>
          <n-select v-model:value="quality" :options="qualityOptions" style="width: 100px" />
          <n-button type="primary" @click="startCompress" :loading="compressing" :disabled="!inputDir">开始压缩</n-button>
        </n-space>

        <n-progress v-if="compressing" type="line" :percentage="progress" :indicator-placement="'inside'" />

        <n-alert v-if="result" :type="result.success ? 'success' : 'error'" closable>
          {{ result.success
            ? `压缩完成: ${result.processed} 张, 节省 ${result.savedSize}`
            : '压缩失败: ' + (result.error || '未知错误') }}
        </n-alert>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { BatchCompressImages, SelectDirectory } from '@wails/go/main/App'

const inputDir = ref('')
const quality = ref(80)
const compressing = ref(false)
const progress = ref(0)
const result = ref<any>(null)
const message = useMessage()

const qualityOptions = [
  { label: '高清 (90%)', value: 90 },
  { label: '普通 (80%)', value: 80 },
  { label: '均衡 (60%)', value: 60 },
  { label: '压缩优先 (40%)', value: 40 },
]

async function selectDir() {
  try {
    const path = await SelectDirectory()
    if (path) inputDir.value = path as string
  } catch (e) {
    console.error(e)
  }
}

async function startCompress() {
  if (!inputDir.value) return
  compressing.value = true
  progress.value = 0
  result.value = null

  try {
    const r = await BatchCompressImages(inputDir.value, quality.value)
    result.value = r
    if (r?.success) {
      message.success(`压缩完成，共处理 ${r.processed} 张图片`)
    } else {
      message.error(r?.error || '压缩失败')
    }
  } catch (e: any) {
    result.value = { success: false, error: e.message || String(e) }
    message.error(String(e))
  }
  compressing.value = false
  progress.value = 100
}
</script>

<style scoped>
</style>
