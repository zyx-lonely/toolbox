<template>
  <div>
    <n-h2>代码美化</n-h2>
    <n-card>
      <n-radio-group v-model:value="lang" class="lang-selector">
        <n-radio-button value="sql">SQL</n-radio-button>
        <n-radio-button value="html">HTML</n-radio-button>
        <n-radio-button value="css">CSS</n-radio-button>
      </n-radio-group>
      <div class="mt-4">
        <n-form-item label="输入代码">
          <n-input v-model:value="inputCode" type="textarea" rows="8" placeholder="粘贴要美化的代码" />
        </n-form-item>
      </div>
      <n-button type="primary" @click="beautify" :loading="loading" :disabled="!inputCode">美化</n-button>
      <div v-if="outputCode" class="mt-4">
        <n-form-item label="美化结果">
          <n-input :value="outputCode" type="textarea" rows="8" readonly />
        </n-form-item>
        <n-button size="small" @click="copy">复制结果</n-button>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const message = useMessage()
const lang = ref('sql')
const inputCode = ref('')
const outputCode = ref('')
const loading = ref(false)

async function beautify() {
  loading.value = true
  try {
    const { BeautifySQL, BeautifyHTML, BeautifyCSS } = await import('@wails/go/main/App')
    let r
    if (lang.value === 'sql') r = await BeautifySQL(inputCode.value)
    else if (lang.value === 'html') r = await BeautifyHTML(inputCode.value)
    else r = await BeautifyCSS(inputCode.value)
    outputCode.value = r.output
  } catch (e: any) {
    message.error('美化失败: ' + (e.message || e))
  }
  loading.value = false
}

function copy() {
  navigator.clipboard.writeText(outputCode.value).then(() => message.success('已复制'))
}
</script>

<style scoped>
.lang-selector { margin-bottom: 8px; }
.mt-4 { margin-top: 16px; }
</style>
