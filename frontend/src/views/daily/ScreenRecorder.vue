<template>
  <div>
    <h2>屏幕录制</h2>
    <p>录制屏幕视频，支持设置时长和帧率。</p>
    <n-card size="small">
      <n-space vertical>
        <n-space>
          <n-input-number v-model:value="duration" :min="10" :max="3600" style="width: 150px">
            <template #prefix>时长</template>
            <template #suffix>秒</template>
          </n-input-number>
          <n-select v-model:value="fps" :options="fpsOptions" style="width: 120px" />
          <n-switch v-model:value="audio" /> <span>录制音频</span>
        </n-space>
        <n-space>
          <n-button v-if="!recording" type="primary" @click="startRecording" :loading="loading">开始录制</n-button>
          <n-button v-else type="error" @click="stopRecording">停止录制</n-button>
          <n-tag v-if="recording" type="error" size="large">录制中 {{ elapsed }}s</n-tag>
        </n-space>
        <n-alert v-if="result" :type="result.success ? 'success' : 'error'">
          <template v-if="result.success">
            录制完成: {{ result.filePath }}
          </template>
          <template v-else>
            {{ result.error }}
          </template>
        </n-alert>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import { StartScreenRecording, StopScreenRecording } from '@wails/go/main/App'

interface ScreenRecordResult { filePath: string; size: number; duration: string; success: boolean; error?: string }

const duration = ref(60)
const fps = ref(30)
const audio = ref(false)
const recording = ref(false)
const loading = ref(false)
const result = ref<ScreenRecordResult | null>(null)
const elapsed = ref(0)
const message = useMessage()
let timer: number | null = null

const fpsOptions = [
  { label: '15 FPS', value: 15 },
  { label: '24 FPS', value: 24 },
  { label: '30 FPS', value: 30 },
  { label: '60 FPS', value: 60 },
]

async function startRecording() {
  loading.value = true
  try {
    const r = await StartScreenRecording(duration.value, fps.value, audio.value) as ScreenRecordResult
    result.value = r
    if (r.success) {
      recording.value = true
      elapsed.value = 0
      timer = window.setInterval(() => { elapsed.value++ }, 1000)
      message.success('录制开始')
    }
  } catch (e: any) { message.error(String(e)) }
  loading.value = false
}

async function stopRecording() {
  const r = await StopScreenRecording() as ScreenRecordResult
  result.value = r
  recording.value = false
  if (timer) { clearInterval(timer); timer = null }
  message.success('录制已停止')
}

onUnmounted(() => { if (timer) clearInterval(timer) })
</script>
