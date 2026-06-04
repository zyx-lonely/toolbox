<template>
  <div>
    <n-h2>文件内容搜索</n-h2>
    <n-p>在指定目录中搜索包含关键字的文件。</n-p>

    <n-card>
      <n-space vertical>
        <n-space>
          <n-input-group>
            <n-input v-model:value="searchDir" placeholder="选择搜索目录" readonly style="width: 300px" />
            <n-button @click="selectDir">选择目录</n-button>
          </n-input-group>
          <n-input v-model:value="keyword" placeholder="搜索关键字" clearable style="width: 200px" />
          <n-input v-model:value="filePattern" placeholder="文件匹配, 如 *.txt" clearable style="width: 160px" />
          <n-button type="primary" @click="startSearch" :loading="searching" :disabled="!searchDir || !keyword">
            搜索
          </n-button>
        </n-space>

        <n-progress v-if="searching" type="line" :percentage="progress" :indicator-placement="'inside'" />

        <n-empty v-if="!searching && !results.length && searched" description="未找到匹配内容" />
        <n-empty v-if="!searching && !results.length && !searched" description="选择目录并输入关键字后点击搜索" />

        <n-data-table
          v-if="results.length"
          :columns="columns"
          :data="results"
          size="small"
          :bordered="true"
          :max-height="500"
        />
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { SearchFileContent, SelectDirectory } from '@wails/go/main/App'

const searchDir = ref('')
const keyword = ref('')
const filePattern = ref('*')
const searching = ref(false)
const progress = ref(0)
const searched = ref(false)
const results = ref<any[]>([])
const message = useMessage()

const columns = [
  { title: '文件路径', key: 'path', ellipsis: { tooltip: true }, width: 300 },
  { title: '行号', key: 'line', width: 60 },
  { title: '内容', key: 'content', ellipsis: { tooltip: true } },
]

async function selectDir() {
  try {
    const path = await SelectDirectory()
    if (path) searchDir.value = path as string
  } catch (e) {
    console.error(e)
  }
}

async function startSearch() {
  if (!searchDir.value || !keyword.value) return
  searching.value = true
  progress.value = 0
  searched.value = true
  results.value = []

  try {
    const r = await SearchFileContent(searchDir.value, keyword.value, filePattern.value)
    if (r) {
      results.value = r as any[]
      message.success(`搜索完成，共找到 ${results.value.length} 条结果`)
    }
  } catch (e: any) {
    message.error(String(e))
  }
  searching.value = false
  progress.value = 100
}
</script>

<style scoped>
</style>
