<template>
  <div>
    <n-h2>截屏工具</n-h2>
    <n-p>截取屏幕并保存到本地。</n-p>

    <n-space>
      <n-button type="primary" @click="captureAll" :loading="capturing">
        截取全屏
      </n-button>
      <n-button @click="captureScreen1" :loading="capturing" :disabled="!hasMultiDisplay">
        显示器 1
      </n-button>
      <n-button @click="captureScreen2" :loading="capturing" :disabled="!hasMultiDisplay">
        显示器 2
      </n-button>
    </n-space>

    <n-card v-if="result.success" class="mt-4">
      <template #header>
        <n-space justify="space-between">
          <span>截屏成功</span>
          <n-space>
            <n-button size="small" @click="openFolder">打开文件夹</n-button>
            <n-button size="small" @click="copyPath">复制路径</n-button>
          </n-space>
        </n-space>
      </template>

      <n-image
        :src="imageSrc"
        :width="Math.min(result.width, 600)"
        :height="Math.min(result.height, 400)"
        object-fit="contain"
      />

      <n-description-list label-placement="left" :column="2" class="mt-2">
        <n-description-item label="文件路径">{{ result.path }}</n-description-item>
        <n-description-item label="分辨率">{{ result.width }} × {{ result.height }}</n-description-item>
        <n-description-item label="文件大小">{{ formatBytes(result.size) }}</n-description-item>
      </n-description-list>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { CaptureAllScreens, CaptureScreen, ReadFileAsBase64 } from '@wails/go/main/App'

interface CaptureResult { path: string; width: number; height: number; size: number; success: boolean; error?: string; dataUri?: string }

const result = ref<CaptureResult>({ path: '', width: 0, height: 0, size: 0, success: false })
const capturing = ref(false)
const message = useMessage()

const hasMultiDisplay = computed(() => true) // 简化处理

const imageSrc = computed(() => {
  if (!result.value.path) return ''
  return result.value.dataUri || ''
})

async function captureAll() {
  capturing.value = true
  try {
    const r = await CaptureAllScreens()
    if (r) {
      result.value = r as CaptureResult
      if (r.path) result.value.dataUri = await ReadFileAsBase64(r.path) as string
    }
    if (r?.success) message.success('截屏成功')
    else message.error(r?.error || '截屏失败')
  } catch (e: any) {
    message.error(`截屏失败: ${e}`)
  }
  capturing.value = false
}

async function captureScreen1() {
  capturing.value = true
  try {
    const r = await CaptureScreen(0)
    if (r) {
      result.value = r as CaptureResult
      if (r.path) result.value.dataUri = await ReadFileAsBase64(r.path) as string
    }
    if (r?.success) message.success('截屏成功')
    else message.error(r?.error || '截屏失败')
  } catch (e: any) {
    message.error(`截屏失败: ${e}`)
  }
  capturing.value = false
}

async function captureScreen2() {
  capturing.value = true
  try {
    const r = await CaptureScreen(1)
    if (r) {
      result.value = r as CaptureResult
      if (r.path) result.value.dataUri = await ReadFileAsBase64(r.path) as string
    }
    if (r?.success) message.success('截屏成功')
    else message.error(r?.error || '截屏失败')
  } catch (e: any) {
    message.error(`截屏失败: ${e}`)
  }
  capturing.value = false
}

async function openFolder() {
  try {
    const { OpenInExplorer } = await import('@wails/go/main/App')
    await OpenInExplorer(result.value.path)
  } catch (e) { console.error(e) }
}

async function copyPath() {
  try {
    await navigator.clipboard.writeText(result.value.path)
    message.success('路径已复制')
  } catch { message.warning('复制失败') }
}

function formatBytes(bytes: number): string {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0; let val = bytes
  while (val >= 1024 && i < units.length - 1) { val /= 1024; i++ }
  return `${val.toFixed(i === 0 ? 0 : 2)} ${units[i]}`
}
</script>

<style scoped>
.mt-4 { margin-top: 16px; }
.mt-2 { margin-top: 8px; }
</style>
