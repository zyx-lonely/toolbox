<template>
  <div>
    <n-h2>文件内容替换</n-h2>
    <n-p>在指定目录中搜索并批量替换文件内容。</n-p>
    <n-space vertical class="mb-4">
      <n-space>
        <n-input-group>
          <n-input v-model:value="scanDir" placeholder="选择要扫描的目录" readonly style="width: 300px" />
          <n-button @click="selectDir">选择目录</n-button>
        </n-input-group>
      </n-space>
      <n-input v-model:value="searchText" placeholder="输入搜索文本" clearable style="width: 400px" />
      <n-input v-model:value="replaceText" placeholder="输入替换文本" clearable style="width: 400px" />
      <n-space>
        <n-input v-model:value="fileTypes" placeholder="文件扩展名过滤，如 .txt,.md (留空表示所有)" style="width: 300px" />
        <n-button type="primary" @click="startReplace" :loading="loading">开始替换</n-button>
      </n-space>
    </n-space>
    <n-empty v-if="!loading && !results.length && searched" description="未找到匹配内容" />
    <n-data-table v-if="results.length" :columns="columns" :data="results" size="small" :bordered="true" :max-height="500" />
  </div>
</template>
<script setup lang="ts">
import { ref, h } from 'vue'
import { useMessage } from 'naive-ui'
import { SearchAndReplace, SelectDirectory } from '@wails/go/main/App'
const scanDir = ref(''); const searchText = ref(''); const replaceText = ref(''); const fileTypes = ref('')
const results = ref<any[]>([]); const loading = ref(false); const searched = ref(false)
const message = useMessage()
const columns = [
  { title: '文件路径', key: 'path', ellipsis: { tooltip: true } },
  { title: '匹配数', key: 'matches', width: 80 },
  { title: '替换数', key: 'replaced', width: 80 },
  { title: '错误', key: 'error', width: 200, render: (r: any) => r.error ? h('span', { style: 'color:red' }, r.error) : '' },
]
async function selectDir() { const d = await SelectDirectory(); if (d) scanDir.value = d as string }
async function startReplace() {
  if (!scanDir.value || !searchText.value) { message.warning('请选择目录并输入搜索文本'); return }
  loading.value = true; searched.value = true
  try {
    const r = await SearchAndReplace(scanDir.value, searchText.value, replaceText.value, fileTypes.value)
    if (r) results.value = r as any[]
    message.success(`处理完成，共 ${results.value.length} 个文件`)
  } catch(e:any) { message.error(String(e)) }
  loading.value = false
}
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
