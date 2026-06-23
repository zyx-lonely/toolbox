<template>
  <div>
    <n-h2>便签 / 记事本</n-h2>
    <n-p>快速记录临时内容，支持自动保存。</n-p>

    <n-space vertical :size="16">
      <n-input-group>
        <n-input v-model:value="newTitle" placeholder="便签标题" style="width: 200px" />
        <n-button type="primary" @click="addNote">新建</n-button>
      </n-input-group>

      <n-grid :cols="3" :x-gap="12" :y-gap="12" v-if="notes.length">
        <n-gi v-for="(note, idx) in notes" :key="idx">
          <n-card size="small" :title="note.title" :style="{ borderTop: `3px solid ${note.color}` }">
            <template #header-extra>
              <n-dropdown :options="colorOptions" @select="(c: any) => changeColor(idx, c)" trigger="click">
                <n-button quaternary circle size="tiny" style="color: #999">●</n-button>
              </n-dropdown>
            </template>
            <n-input
              v-model:value="note.content"
              type="textarea"
              :rows="6"
              placeholder="在此输入内容..."
              @input="autoSave"
            />
            <template #footer>
              <n-space justify="space-between" align="center">
                <n-text depth="3" style="font-size: 12px">{{ note.time }}</n-text>
                <n-button size="tiny" type="error" text @click="deleteNote(idx)">删除</n-button>
              </n-space>
            </template>
          </n-card>
        </n-gi>
      </n-grid>

      <n-empty v-else description="暂无便签，点击上方新建" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const message = useMessage()
const newTitle = ref('')

interface Note {
  title: string
  content: string
  time: string
  color: string
}

const notes = ref<Note[]>([])

const colorOptions = [
  { label: '绿色', key: '#18a058' },
  { label: '蓝色', key: '#2080f0' },
  { label: '橙色', key: '#f0a020' },
  { label: '红色', key: '#d03050' },
  { label: '紫色', key: '#8b5cf6' },
]

function addNote() {
  const title = newTitle.value.trim() || `便签 ${notes.value.length + 1}`
  notes.value.unshift({
    title,
    content: '',
    time: new Date().toLocaleString('zh-CN'),
    color: colorOptions[Math.floor(Math.random() * colorOptions.length)].key
  })
  newTitle.value = ''
  message.success(`已创建便签「${title}」`)
}

function deleteNote(idx: number) {
  const title = notes.value[idx].title
  notes.value.splice(idx, 1)
  message.success(`已删除便签「${title}」`)
}

function changeColor(idx: number, color: string) {
  notes.value[idx].color = color
}

function autoSave() {
  // Auto-save would persist to localStorage in production
}
</script>
