<template>
  <div>
    <n-h2>大文件查找</n-h2>
    <n-p>扫描指定目录下的大文件，快速定位占用磁盘空间的文件。</n-p>
    <n-space class="mb-4">
      <n-input-group>
        <n-input v-model:value="scanDir" placeholder="选择要扫描的目录" readonly style="width: 300px" />
        <n-button @click="selectDir">选择目录</n-button>
      </n-input-group>
      <n-input-number v-model:value="minSize" :min="1" :max="10000" style="width: 100px"><template #suffix>MB</template></n-input-number>
      <n-button type="primary" @click="scan" :loading="scanning">开始扫描</n-button>
    </n-space>
    <n-empty v-if="!scanning && !files.length && scanDir" description="未找到大文件" />
    <n-data-table v-if="files.length" :columns="columns" :data="files" size="small" :bordered="true" :max-height="500" />
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { FindLargeFiles, SelectDirectory } from '@wails/go/main/App'
const scanDir = ref(''); const minSize = ref(100); const files = ref<any[]>([]); const scanning = ref(false); const message = useMessage()
const columns = [
  { title: '文件名', key: 'name', ellipsis: { tooltip: true } },
  { title: '路径', key: 'path', ellipsis: { tooltip: true } },
  { title: '大小', key: 'size', width: 100, render: (r: any) => formatBytes(r.size) },
  { title: '修改时间', key: 'modified', width: 150 },
]
async function selectDir() { const d = await SelectDirectory(); if (d) scanDir.value = d as string }
async function scan() {
  if (!scanDir.value) return; scanning.value = true
  try { const r = await FindLargeFiles(scanDir.value, minSize.value, 100); if (r) files.value = r as any[]; message.success(`找到 ${files.value.length} 个大文件`) }
  catch(e:any) { message.error(String(e)) }; scanning.value = false
}
function formatBytes(b: number) { const u=['B','KB','MB','GB','TB']; let i=0; let v=b; while(v>=1024&&i<4){v/=1024;i++} return `${v.toFixed(i<=1?0:2)} ${u[i]}` }
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
