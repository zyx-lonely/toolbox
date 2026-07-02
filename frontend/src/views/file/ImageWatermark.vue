<template>
  <div>
    <h2>图片批量加水印</h2>
    <p>选择图片文件批量添加文字水印。</p>
    <n-space class="mb-4" vertical>
      <n-space>
        <n-button type="primary" @click="selectFiles">选择图片</n-button>
        <n-button @click="addWatermark" :loading="loading" :disabled="files.length === 0">添加水印</n-button>
      </n-space>
      <n-space>
        <n-input v-model:value="watermarkText" placeholder="水印文字" style="width: 200px" />
        <n-select v-model:value="position" :options="positionOptions" style="width: 150px" />
        <n-slider v-model:value="opacity" :max="100" style="width: 200px" />
        <span>透明度 {{ opacity }}%</span>
      </n-space>
      <n-text v-if="files.length > 0">已选择 {{ files.length }} 个文件</n-text>
    </n-space>
    <n-data-table v-if="results.length > 0" :columns="columns" :data="results" :bordered="true" size="small" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { SelectFile, AddWatermarkToImages } from '@wails/go/main/App'

interface WatermarkResult { file: string; width: number; height: number; size: number }

const files = ref<string[]>([])
const results = ref<WatermarkResult[]>([])
const watermarkText = ref('')
const position = ref('bottom-right')
const opacity = ref(50)
const loading = ref(false)
const message = useMessage()

const positionOptions = [
  { label: '右下角', value: 'bottom-right' },
  { label: '左下角', value: 'bottom-left' },
  { label: '右上角', value: 'top-right' },
  { label: '左上角', value: 'top-left' },
  { label: '居中', value: 'center' },
]

const columns = [
  { title: '文件', key: 'file', ellipsis: { tooltip: true } },
  { title: '尺寸', width: 120, render: (row: WatermarkResult) => `${row.width}x${row.height}` },
  { title: '大小', width: 100, render: (row: WatermarkResult) => (row.size / 1024).toFixed(0) + ' KB' },
]

async function selectFiles() {
  try {
    const f = await SelectFile()
    if (f) files.value = [f]
  } catch (e) { console.error(e) }
}

async function addWatermark() {
  if (!watermarkText.value) { message.warning('请输入水印文字'); return }
  loading.value = true
  try {
    const r = await AddWatermarkToImages(files.value, watermarkText.value, position.value, opacity.value) as WatermarkResult[]
    results.value = r || []
    message.success(`成功处理 ${r?.length || 0} 个文件`)
  } catch (e: any) { message.error(String(e)) }
  loading.value = false
}
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
