<template>
  <div>
    <n-h2>文件夹大小分析</n-h2>
    <n-p>分析指定目录下各子文件夹的大小占比，快速定位空间占用。</n-p>
    <n-space class="mb-4">
      <n-input-group>
        <n-input v-model:value="scanDir" placeholder="选择要分析的目录" readonly style="width: 300px" />
        <n-button @click="selectDir">选择目录</n-button>
      </n-input-group>
      <n-input-number v-model:value="depth" :min="1" :max="3" style="width: 80px"><template #suffix>级</template></n-input-number>
      <n-button type="primary" @click="analyze" :loading="analyzing">开始分析</n-button>
    </n-space>
    <n-empty v-if="!analyzing && !folders.length && analyzed" description="未找到子文件夹" />
    <n-data-table v-if="folders.length" :columns="columns" :data="folders" size="small" :bordered="true" :max-height="500" />
  </div>
</template>
<script setup lang="ts">
import { ref, h } from 'vue'
import { useMessage } from 'naive-ui'
import { AnalyzeFolderSizes, SelectDirectory } from '@wails/go/main/App'
const scanDir = ref(''); const depth = ref(1); const folders = ref<any[]>([]); const analyzing = ref(false); const analyzed = ref(false)
const message = useMessage()
const columns = [
  { title: '文件夹名', key: 'name', width: 150 },
  { title: '路径', key: 'path', ellipsis: { tooltip: true } },
  { title: '大小', key: 'size', width: 100, render: (r: any) => formatBytes(r.size) },
  { title: '文件数', key: 'files', width: 80 },
  { title: '占比', key: 'pct', width: 150, render: (r: any) => h('span', `${r.pct.toFixed(1)}%`) },
]
async function selectDir() { const d = await SelectDirectory(); if (d) scanDir.value = d as string }
async function analyze() {
  if (!scanDir.value) return; analyzing.value = true; analyzed.value = true
  try { const r = await AnalyzeFolderSizes(scanDir.value, depth.value); if (r) folders.value = r as any[] }
  catch(e:any) { message.error(String(e)) }; analyzing.value = false
}
function formatBytes(b: number) { const u=['B','KB','MB','GB','TB']; let i=0; let v=b; while(v>=1024&&i<4){v/=1024;i++} return `${v.toFixed(i<=1?0:2)} ${u[i]}` }
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
