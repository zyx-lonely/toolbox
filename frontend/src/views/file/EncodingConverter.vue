<template>
  <div>
    <n-h2>文本编码转换</n-h2>
    <n-card>
      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="源编码">
            <n-select v-model:value="fromEnc" :options="encOptions" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="目标编码">
            <n-select v-model:value="toEnc" :options="encOptions" />
          </n-form-item>
        </n-gi>
      </n-grid>
      <n-form-item label="输入文本">
        <n-input v-model:value="inputText" type="textarea" rows="6" placeholder="请输入要转换的文本" />
      </n-form-item>
      <n-button type="primary" @click="convert" :loading="loading" :disabled="!inputText">转换</n-button>
      <div v-if="outputText" class="mt-4">
        <n-form-item label="输出结果">
          <n-input :value="outputText" type="textarea" rows="6" readonly />
        </n-form-item>
        <n-button size="small" @click="copy">复制结果</n-button>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { ConvertEncoding } from '@wails/go/main/App'

const message = useMessage()
const inputText = ref('')
const outputText = ref('')
const fromEnc = ref('UTF-8')
const toEnc = ref('GBK')
const loading = ref(false)

const encOptions = [
  { label: 'UTF-8', value: 'UTF-8' },
  { label: 'GBK', value: 'GBK' },
  { label: 'UTF-16 LE', value: 'UTF-16 LE' }
]

async function convert() {
  loading.value = true
  try {
    const r = await ConvertEncoding(inputText.value, fromEnc.value, toEnc.value)
    outputText.value = r.output
  } catch (e: any) {
    message.error('转换失败: ' + (e.message || e))
  }
  loading.value = false
}

function copy() {
  navigator.clipboard.writeText(outputText.value).then(() => message.success('已复制'))
}
</script>

<style scoped>
.mt-4 { margin-top: 16px; }
</style>
