<template>
  <div>
    <n-h2>定时任务</n-h2>
    <n-p>创建定时清理、截屏、关机任务。</n-p>
    <n-card title="新建定时任务" class="mb-4">
      <n-space>
        <n-select v-model:value="taskAction" :options="actionOptions" style="width:150px" />
        <n-time-picker v-model:value="taskTime" format="HH:mm" style="width:120px" />
        <n-button type="primary" @click="createTask" :loading="taskCreating">创建</n-button>
      </n-space>
    </n-card>
    <n-h3>已创建的任务</n-h3>
    <n-empty v-if="!tasks.length" description="暂无定时任务" />
    <n-list v-if="tasks.length">
      <n-list-item v-for="t in tasks" :key="t.name">
        <n-space justify="space-between">
          <span>{{ t.name }} — {{ t.action }} 于 {{ t.time }}</span>
          <n-button size="tiny" type="error" @click="deleteTask(t.name)">删除</n-button>
        </n-space>
      </n-list-item>
    </n-list>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { CreateScheduledTask, ListScheduledTasks, DeleteScheduledTask } from '@wails/go/main/App'
const taskAction = ref('cleanup'); const taskTime = ref(new Date()); const tasks = ref<any[]>([]); const taskCreating = ref(false); const message = useMessage()
const actionOptions = [{label:'临时清理',value:'cleanup'},{label:'截屏',value:'screenshot'},{label:'关机',value:'shutdown'}]
async function createTask() { taskCreating.value = true; const h = taskTime.value.getHours(); const m = taskTime.value.getMinutes(); await CreateScheduledTask(taskAction.value, h, m); message.success('任务已创建'); loadTasks(); taskCreating.value = false }
async function loadTasks() { try { const r = await ListScheduledTasks(); if(r) tasks.value = r as any[] } catch(e){console.error(e)} }
async function deleteTask(name: string) { await DeleteScheduledTask(name); loadTasks(); message.success('已删除') }
loadTasks()
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
