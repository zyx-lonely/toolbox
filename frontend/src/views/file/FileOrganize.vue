<template>
  <div>
    <n-h2>文件批量归类</n-h2>
    <n-p>按扩展名、日期或大小自动整理文件夹中的文件。</n-p>
    <n-space class="mb-4">
      <n-input-group>
        <n-input v-model:value="sourceDir" placeholder="选择要整理的目录" readonly style="width: 300px" />
        <n-button @click="selectDir">选择目录</n-button>
      </n-input-group>
      <n-select v-model:value="organizeMode" :options="modeOptions" style="width: 120px" />
      <n-button type="primary" @click="preview" :loading="previewing">预览</n-button>
      <n-button type="warning" @click="execute" :loading="executing" :disabled="!previews.length">执行归类</n-button>
    </n-space>
    <n-empty v-if="!previews.length && sourceDir" description="未找到可归类的文件" />
    <n-data-table v-if="previews.length" :columns="columns" :data="previews" size="small" :bordered="true" :max-height="400" />
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
const sourceDir = ref('')
const organizeMode = ref('extension')
const previews = ref<any[]>([])
const previewing = ref(false)
const executing = ref(false)
const message = useMessage()
const modeOptions = [
  { label: '按扩展名', value: 'extension' },
  { label: '按年份', value: 'date', target: 'year' },
  { label: '按年月', value: 'date', target: 'year-month' },
  { label: '按文件大小', value: 'size' },
]
const columns = [
  { title: '文件名', key: 'sourceName', ellipsis: { tooltip: true } },
  { title: '归类到', key: 'folderName' },
  { title: '大小', key: 'fileSize', width: 100, render: (r: any) => formatBytes(r.fileSize) },
]
async function selectDir() {
  try {
    const { SelectDirectory } = await import('@wails/go/main/App')
    const dir = await SelectDirectory()
    if (dir) sourceDir.value = dir as string
  } catch (e) { console.error(e) }
}
async function preview() {
  if (!sourceDir.value) return
  previewing.value = true
  try {
    const { PreviewOrganize } = await import('@wails/go/main/App')
    const r = await PreviewOrganize(sourceDir.value, { mode: organizeMode.value, target: '', move: false, sortInto: sourceDir.value })
    if (r) previews.value = r as any[]
  } catch (e: any) { message.error(String(e)) }
  previewing.value = false
}
async function execute() {
  executing.value = true
  try {
    const { ExecuteOrganize } = await import('@wails/go/main/App')
    const r = await ExecuteOrganize(sourceDir.value, { mode: organizeMode.value, target: '', move: true, sortInto: sourceDir.value })
    message.success(`归类完成，共处理 ${r?.length || 0} 个文件`)
  } catch (e: any) { message.error(String(e)) }
  executing.value = false
}
function formatBytes(b: number) { const u=['B','KB','MB','GB']; let i=0; let v=b; while(v>=1024&&i<3){v/=1024;i++} return `${v.toFixed(i===0?0:2)} ${u[i]}` }
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
