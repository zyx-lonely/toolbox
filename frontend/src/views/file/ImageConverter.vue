<template>
  <div>
    <n-h2>图片格式转换</n-h2>
    <n-card>
      <n-space>
        <n-input-group>
          <n-input v-model:value="inputPath" placeholder="选择要转换的图片" readonly style="width: 300px" />
          <n-button @click="selectFile">选择图片</n-button>
        </n-input-group>
        <n-select v-model:value="targetFormat" :options="formatOptions" style="width: 100px" />
        <n-button type="primary" @click="convert" :loading="converting" :disabled="!inputPath">开始转换</n-button>
      </n-space>

      <n-alert v-if="result" :type="result.success ? 'success' : 'error'" class="mt-4" closable>
        {{ result.success ? '转换成功: ' + result.outputPath : '转换失败: ' + result.error }}
      </n-alert>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { SelectFile, ConvertImage } from '@wails/go/main/App'

const inputPath = ref('')
const targetFormat = ref('png')
const converting = ref(false)
const result = ref<any>(null)

const formatOptions = [
  { label: 'PNG', value: 'png' },
  { label: 'JPEG', value: 'jpg' }
]

async function selectFile() {
  try {
    const path = await SelectFile()
    if (path) inputPath.value = path as string
  } catch (e) { console.error(e) }
}

async function convert() {
  if (!inputPath.value) return
  converting.value = true
  result.value = null
  try {
    result.value = await ConvertImage(inputPath.value, targetFormat.value)
  } catch (e: any) {
    result.value = { success: false, error: e.message || String(e) }
  }
  converting.value = false
}
</script>

<style scoped>
.mt-4 { margin-top: 16px; }
</style>
