<template>
  <div>
    <n-h2>文件预览</n-h2>
    <n-p>快速预览文本、图片和代码文件内容。</n-p>

    <n-card>
      <n-space vertical>
        <n-space align="center">
          <n-button @click="selectFile" type="primary">
            <template #icon><n-icon><document-outline /></n-icon></template>
            选择文件
          </n-button>
          <n-button v-if="filePath" @click="refresh" :loading="loading">
            <template #icon><n-icon><refresh-outline /></n-icon></template>
            刷新
          </n-button>
          <n-button v-if="previewData && previewData.isText" @click="copyText" type="success" ghost>
            <template #icon><n-icon><copy-outline /></n-icon></template>
            复制内容
          </n-button>
          <n-button v-if="previewData && previewData.isImage" @click="saveImage" type="info" ghost>
            <template #icon><n-icon><download-outline /></n-icon></template>
            保存图片
          </n-button>
        </n-space>

        <n-alert v-if="filePath" type="info">
          <n-tag type="info">{{ fileName }}</n-tag>
          <n-tag type="default" style="margin-left:8px">{{ formatSize(previewData?.size || 0) }}</n-tag>
          <n-tag v-if="previewData?.isText" type="success" style="margin-left:8px">{{ previewData.totalLines }} 行</n-tag>
        </n-alert>

        <n-spin :show="loading">
          <n-card v-if="previewData" :bordered="true" size="small">
            <template v-if="previewData.isImage">
              <div style="text-align:center; padding:20px;">
                <img :src="previewData.imageBase64" style="max-width:100%; max-height:600px; border:1px solid #eee; border-radius:4px;" />
              </div>
            </template>
            <template v-else-if="previewData.isText">
              <n-code :code="previewData.textContent" language="text" show-line-numbers style="max-height:600px; overflow:auto; font-size:12px;" />
              <n-alert v-if="previewData.truncated" type="warning" style="margin-top:8px;">
                内容已截断显示
              </n-alert>
            </template>
            <template v-else-if="previewData.error">
              <n-alert type="error">{{ previewData.error }}</n-alert>
            </template>
          </n-card>
          <n-empty v-else description="请选择要预览的文件" />
        </n-spin>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { DocumentOutline, RefreshOutline, CopyOutline, DownloadOutline } from '@vicons/ionicons5'
import { SelectFile, PreviewFile, SaveClipboardImage } from '@wails/go/main/App'

const filePath = ref('')
const fileName = ref('')
const loading = ref(false)
const previewData = ref<any>(null)
const message = useMessage()

async function selectFile() {
  const path = await SelectFile()
  if (path) {
    filePath.value = path
    fileName.value = path.split(/[/\\]/).pop() || ''
    await refresh()
  }
}

async function refresh() {
  if (!filePath.value) return
  loading.value = true
  try {
    previewData.value = await PreviewFile(filePath.value)
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

function copyText() {
  if (previewData.value?.textContent) {
    navigator.clipboard.writeText(previewData.value.textContent)
    message.success('内容已复制到剪贴板')
  }
}

async function saveImage() {
  if (!previewData.value?.imageBase64) return
  const path = await SelectFile()
  if (path) {
    const r = await SaveClipboardImage(path)
    if (r && r.success) {
      message.success('图片已保存')
    } else {
      message.error(r?.error || '保存失败')
    }
  }
}

function formatSize(bytes: number) {
  if (bytes >= 1024 * 1024) return (bytes / 1024 / 1024).toFixed(2) + ' MB'
  if (bytes >= 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return bytes + ' B'
}
</script>
