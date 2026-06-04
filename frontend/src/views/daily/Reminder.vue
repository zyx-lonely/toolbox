<template>
  <div>
    <n-h2>提醒事项</n-h2>
    <n-card>
      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="标题">
            <n-input v-model:value="newReminder.title" placeholder="提醒标题" />
          </n-form-item>
          <n-form-item label="时间">
            <n-date-picker v-model:value="newReminder.time" type="datetime" style="width: 100%" />
          </n-form-item>
          <n-form-item label="重复">
            <n-select v-model:value="newReminder.repeat" :options="repeatOptions" />
          </n-form-item>
          <n-button type="primary" @click="addReminder" :disabled="!newReminder.title || !newReminder.time">添加提醒</n-button>
        </n-gi>
        <n-gi>
          <n-h3>提醒列表</n-h3>
          <n-list v-if="reminders.length > 0">
            <n-list-item v-for="(item, idx) in reminders" :key="idx">
              <template #suffix>
                <n-button size="tiny" quaternary circle @click="removeReminder(idx)">
                  <template #icon><n-icon><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg></n-icon></template>
                </n-button>
              </template>
              <n-thing :title="item.title" :description="formatTime(item.time)">
                <template #extra>
                  <n-tag v-if="item.repeat" size="tiny">{{ item.repeat }}</n-tag>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
          <n-empty v-else description="暂无提醒" />
        </n-gi>
      </n-grid>
      <n-alert v-if="notification" :type="notification" class="mt-4" closable>
        {{ notificationText }}
      </n-alert>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'

interface Reminder {
  title: string
  time: number | null
  repeat: string
}

const message = useMessage()
const newReminder = ref<Reminder>({ title: '', time: null, repeat: 'none' })
const reminders = ref<Reminder[]>([])
const notification = ref<string | null>(null)
const notificationText = ref('')
let timer: ReturnType<typeof setInterval> | null = null

const repeatOptions = [
  { label: '不重复', value: 'none' },
  { label: '每天', value: 'daily' },
  { label: '每周', value: 'weekly' },
  { label: '每月', value: 'monthly' }
]

function addReminder() {
  if (!newReminder.value.title || !newReminder.value.time) return
  reminders.value.push({ ...newReminder.value })
  newReminder.value = { title: '', time: null, repeat: 'none' }
  saveReminders()
  message.success('提醒已添加')
}

function removeReminder(idx: number) {
  reminders.value.splice(idx, 1)
  saveReminders()
}

function formatTime(ts: number | null) {
  if (!ts) return ''
  return new Date(ts).toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

function checkReminders() {
  const now = Date.now()
  for (const item of reminders.value) {
    if (item.time && now >= item.time && now - item.time < 10000) {
      notificationText.value = `提醒: ${item.title}`
      notification.value = 'success'
      showNotification(item.title)
      if (item.repeat === 'none') {
        item.time = null
      } else {
        const d = new Date(item.time)
        if (item.repeat === 'daily') d.setDate(d.getDate() + 1)
        else if (item.repeat === 'weekly') d.setDate(d.getDate() + 7)
        else if (item.repeat === 'monthly') d.setMonth(d.getMonth() + 1)
        item.time = d.getTime()
      }
      saveReminders()
    }
  }
}

function showNotification(title: string) {
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification('PC Toolbox', { body: title })
  }
}

function saveReminders() {
  localStorage.setItem('pc-toolbox-reminders', JSON.stringify(reminders.value))
}

function loadReminders() {
  try {
    const saved = localStorage.getItem('pc-toolbox-reminders')
    if (saved) reminders.value = JSON.parse(saved)
  } catch {}
}

onMounted(() => {
  loadReminders()
  if ('Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission()
  }
  timer = setInterval(checkReminders, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.mt-4 { margin-top: 16px; }
</style>
