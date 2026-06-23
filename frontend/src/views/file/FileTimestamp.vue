<template>
  <div>
    <n-h2>文件时间戳修改</n-h2>
    <n-p>批量修改文件的创建时间、修改时间和访问时间。</n-p>

    <n-space vertical :size="16">
      <n-input-group>
        <n-input v-model:value="filePath" placeholder="输入文件路径" style="flex: 1" />
        <n-button @click="selectFile">选择文件</n-button>
        <n-button type="primary" @click="loadTimestamps" :loading="loading">读取</n-button>
      </n-input-group>

      <n-card v-if="timestamps" title="当前时间戳" size="small">
        <n-space vertical :size="8">
          <n-space align="center">
            <n-text depth="3" style="width: 100px">创建时间：</n-text>
            <n-tag type="info">{{ timestamps.created }}</n-tag>
          </n-space>
          <n-space align="center">
            <n-text depth="3" style="width: 100px">修改时间：</n-text>
            <n-tag type="warning">{{ timestamps.modified }}</n-tag>
          </n-space>
          <n-space align="center">
            <n-text depth="3" style="width: 100px">访问时间：</n-text>
            <n-tag type="success">{{ timestamps.accessed }}</n-tag>
          </n-space>
        </n-space>
      </n-card>

      <n-card title="设置新时间戳" size="small" v-if="timestamps">
        <n-space vertical :size="12">
          <n-space align="center">
            <n-text depth="3" style="width: 100px">创建时间：</n-text>
            <n-date-picker v-model:value="newCreated" type="datetime" clearable style="width: 240px" />
            <n-button size="small" @click="newCreated = Date.now()">设为当前</n-button>
          </n-space>
          <n-space align="center">
            <n-text depth="3" style="width: 100px">修改时间：</n-text>
            <n-date-picker v-model:value="newModified" type="datetime" clearable style="width: 240px" />
            <n-button size="small" @click="newModified = Date.now()">设为当前</n-button>
          </n-space>
          <n-space align="center">
            <n-text depth="3" style="width: 100px">访问时间：</n-text>
            <n-date-picker v-model:value="newAccessed" type="datetime" clearable style="width: 240px" />
            <n-button size="small" @click="newAccessed = Date.now()">设为当前</n-button>
          </n-space>
        </n-space>
        <template #footer>
          <n-space>
            <n-button type="primary" @click="applyChanges">应用修改</n-button>
            <n-button @click="setAllCurrent">全部设为当前时间</n-button>
          </n-space>
        </template>
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const message = useMessage()
const filePath = ref('')
const loading = ref(false)

interface Timestamps {
  created: string
  modified: string
  accessed: string
}

const timestamps = ref<Timestamps | null>(null)
const newCreated = ref<number | null>(null)
const newModified = ref<number | null>(null)
const newAccessed = ref<number | null>(null)

function formatDate(d: Date): string {
  return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function selectFile() {
  filePath.value = 'C:\\Users\\Demo\\Documents\\example.txt'
  message.info('已选择示例文件')
}

function loadTimestamps() {
  if (!filePath.value.trim()) {
    message.warning('请输入文件路径')
    return
  }
  loading.value = true
  setTimeout(() => {
    const now = new Date()
    timestamps.value = {
      created: formatDate(new Date(now.getTime() - 86400000 * 30)),
      modified: formatDate(new Date(now.getTime() - 86400000 * 2)),
      accessed: formatDate(new Date(now.getTime() - 3600000))
    }
    const baseTime = Date.now() - 86400000 * 30
    newCreated.value = baseTime
    newModified.value = Date.now() - 86400000 * 2
    newAccessed.value = Date.now() - 3600000
    loading.value = false
    message.success('已读取文件时间戳')
  }, 300)
}

function applyChanges() {
  const now = new Date()
  timestamps.value = {
    created: newCreated.value ? formatDate(new Date(newCreated.value)) : formatDate(now),
    modified: newModified.value ? formatDate(new Date(newModified.value)) : formatDate(now),
    accessed: newAccessed.value ? formatDate(new Date(newAccessed.value)) : formatDate(now)
  }
  message.success('时间戳修改成功')
}

function setAllCurrent() {
  const now = Date.now()
  newCreated.value = now
  newModified.value = now
  newAccessed.value = now
}
</script>
