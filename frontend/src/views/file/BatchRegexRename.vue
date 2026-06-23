<template>
  <div>
    <n-h2>批量正则重命名</n-h2>
    <n-p>使用正则表达式批量重命名文件。</n-p>

    <n-space vertical :size="16">
      <n-input-group>
        <n-input v-model:value="folderPath" placeholder="文件夹路径" style="flex: 1" />
        <n-button @click="selectFolder">选择文件夹</n-button>
        <n-button type="primary" @click="loadFiles" :loading="loading">加载文件</n-button>
      </n-input-group>

      <n-card title="重命名规则" size="small">
        <n-space vertical :size="12">
          <n-space align="center">
            <n-text depth="3" style="width: 80px">查找：</n-text>
            <n-input v-model:value="findPattern" placeholder="正则表达式，如 (\\d+)" style="width: 300px" />
          </n-space>
          <n-space align="center">
            <n-text depth="3" style="width: 80px">替换为：</n-text>
            <n-input v-model:value="replacePattern" placeholder="替换内容，如 file_$1" style="width: 300px" />
          </n-space>
          <n-space>
            <n-button @click="applyRegex">应用预览</n-button>
            <n-button type="primary" @click="executeRename" :disabled="!previewResults.length">执行重命名</n-button>
            <n-button @click="resetAll">重置</n-button>
          </n-space>
        </n-space>
      </n-card>

      <n-space :size="16">
        <n-card title="原始文件名" size="small" style="width: 400px">
          <n-list bordered>
            <n-list-item v-for="(f, i) in files" :key="i">
              <n-thing :title="f" />
            </n-list-item>
            <n-empty v-if="!files.length" description="暂无文件" />
          </n-list>
        </n-card>
        <n-card title="重命名预览" size="small" style="width: 400px">
          <n-list bordered>
            <n-list-item v-for="(item, i) in previewResults" :key="i">
              <n-thing>
                <template #header>
                  <n-text>{{ item.original }}</n-text>
                </template>
                <template #description>
                  <n-text depth="3">→</n-text>
                  <n-text type="success">{{ item.renamed }}</n-text>
                </template>
              </n-thing>
            </n-list-item>
            <n-empty v-if="!previewResults.length" description="应用正则后显示预览" />
          </n-list>
        </n-card>
      </n-space>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const message = useMessage()
const folderPath = ref('')
const findPattern = ref('(\\d+)')
const replacePattern = ref('file_$1')
const loading = ref(false)
const files = ref<string[]>([])
const previewResults = ref<{ original: string; renamed: string }[]>([])

function selectFolder() {
  folderPath.value = 'C:\\Users\\Demo\\Downloads\\photos'
  message.info('已选择示例文件夹')
}

function loadFiles() {
  if (!folderPath.value.trim()) {
    message.warning('请输入文件夹路径')
    return
  }
  loading.value = true
  setTimeout(() => {
    files.value = [
      'IMG_001.jpg', 'IMG_002.jpg', 'IMG_003.jpg', 'IMG_004.jpg',
      'screenshot_20260601.png', 'screenshot_20260615.png',
      'document_v1.docx', 'document_v2.docx', 'document_v3.docx'
    ]
    previewResults.value = []
    loading.value = false
    message.success(`已加载 ${files.value.length} 个文件`)
  }, 300)
}

function applyRegex() {
  try {
    const regex = new RegExp(findPattern.value)
    previewResults.value = files.value.map(f => ({
      original: f,
      renamed: f.replace(regex, replacePattern.value)
    }))
    message.success('预览已更新')
  } catch (e: any) {
    message.error('正则表达式无效: ' + e.message)
  }
}

function executeRename() {
  const count = previewResults.value.filter(r => r.original !== r.renamed).length
  files.value = previewResults.value.map(r => r.renamed)
  previewResults.value = []
  message.success(`已成功重命名 ${count} 个文件`)
}

function resetAll() {
  files.value = []
  previewResults.value = []
  findPattern.value = '(\\d+)'
  replacePattern.value = 'file_$1'
}
</script>
