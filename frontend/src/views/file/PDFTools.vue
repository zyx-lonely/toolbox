<template>
  <div>
    <n-h2>文档转换工具</n-h2>
    <n-p>将 Word/Excel/PPT 文档转换为 PDF 格式。</n-p>

    <n-alert type="info" :bordered="false" style="margin-bottom: 16px">
      <template #header>使用说明</template>
      需要安装 LibreOffice 或 Microsoft Office。转换过程可能需要几秒钟。
    </n-alert>

    <n-card>
      <n-space vertical>
        <n-input-group>
          <n-input v-model:value="filePath" placeholder="选择要转换的 Office 文档" readonly style="width: 400px" />
          <n-button @click="selectFile">选择文件</n-button>
        </n-input-group>

        <n-alert v-if="filePath" type="info">
          当前文件: {{ filePath }} → {{ outputName }}.pdf
        </n-alert>

        <n-button type="primary" @click="convert" :loading="converting" :disabled="!filePath">
          转换为 PDF
        </n-button>

        <n-alert v-if="result" :type="result.success ? 'success' : 'error'" closable @close="result = null">
          <template #header>{{ result.success ? '转换成功' : '转换失败' }}</template>
          {{ result.success ? `输出文件: ${result.outputPath}` : result.error }}
        </n-alert>
      </n-space>
    </n-card>

    <n-card title="支持的功能" size="small" style="margin-top: 16px">
      <n-table :single-line="false" size="small">
        <thead><tr><th>输入格式</th><th>输出格式</th><th>说明</th></tr></thead>
        <tbody>
          <tr><td>.doc / .docx</td><td>.pdf</td><td>Word 文档转 PDF</td></tr>
          <tr><td>.xls / .xlsx</td><td>.pdf</td><td>Excel 表格转 PDF</td></tr>
          <tr><td>.ppt / .pptx</td><td>.pdf</td><td>PPT 演示文稿转 PDF</td></tr>
        </tbody>
      </n-table>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { ConvertToPDF, SelectFile } from '@wails/go/main/App'

const filePath = ref('')
const converting = ref(false)
const result = ref<any>(null)
const message = useMessage()

const outputName = computed(() => filePath.value ? filePath.value.replace(/\.\w+$/, '') : '')

async function selectFile() {
  const f = await SelectFile()
  if (f) filePath.value = f as string
}

async function convert() {
  if (!filePath.value) return
  converting.value = true
  result.value = null
  try {
    const r = await ConvertToPDF(filePath.value)
    if (r) result.value = r as any
    if (r?.success) message.success('转换成功')
  } catch (e: any) {
    message.error(String(e))
  }
  converting.value = false
}
</script>
