<template>
  <div class="batch-rename">
    <n-space vertical>
      <n-space align="center">
        <n-button @click="handleSelectDir">选择目录</n-button>
        <span>{{ selectedDir || '未选择目录' }}</span>
      </n-space>

      <n-data-table
        v-if="files.length"
        :columns="columns"
        :data="files"
        :bordered="true"
        :pagination="{ pageSize: 10 }"
      />

      <n-space v-if="files.length" align="center">
        <n-input v-model:value="renamePattern" placeholder="替换模式: 旧文本 -> 新文本" style="width: 300px" />
        <n-button type="primary" @click="handlePreview" :loading="previewing">预览</n-button>
        <n-button type="warning" @click="handleExecute" :loading="executing">执行重命名</n-button>
      </n-space>

      <n-empty v-if="!files.length && !selectedDir" description="请选择要操作的目录" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { DataTableColumns } from 'naive-ui'

const selectedDir = ref('')
const files = ref<any[]>([])
const renamePattern = ref('')

const columns: DataTableColumns<any> = [
  { title: '文件名', key: 'name' },
  { title: '大小', key: 'size' },
  { title: '修改时间', key: 'modTime' },
  { title: '新文件名', key: 'newName' },
]

async function handleSelectDir() {
  try {
    const dir = await window.go.main.App.SelectDirectory()
    if (dir) {
      selectedDir.value = dir
      const result = await window.go.main.App.BatchRenamePreview(dir, '', '')
      files.value = result || []
    }
  } catch (e: any) {
    files.value = []
  }
}

async function handlePreview() {
  if (!selectedDir.value) return
  previewing.value = true
  try {
    const result = await window.go.main.App.BatchRenamePreview(selectedDir.value, renamePattern.value, '')
    if (result) files.value = result
  } finally {
    previewing.value = false
  }
}

const previewing = ref(false)
const executing = ref(false)

async function handleExecute() {
  if (!selectedDir.value) return
  executing.value = true
  try {
    await window.go.main.App.BatchRename(selectedDir.value, renamePattern.value, '')
    window.$message.success('重命名完成')
  } finally {
    executing.value = false
  }
}
</script>

<style scoped>
.batch-rename {
  padding: 20px;
}
</style>
