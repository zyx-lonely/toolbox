<template>
  <div>
    <n-h2>系统备份与还原</n-h2>
    <n-p>创建和管理系统还原点。</n-p>
    <n-space class="mb-4">
      <n-input v-model:value="description" placeholder="还原点描述" style="width: 300px" />
      <n-button type="primary" @click="createPoint" :loading="creating">创建还原点</n-button>
      <n-button @click="loadPoints" :loading="loading">刷新列表</n-button>
    </n-space>
    <n-empty v-if="!points.length" description="暂无还原点" />
    <n-data-table v-if="points.length" :columns="columns" :data="points" size="small" :bordered="true" />
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { CreateRestorePoint, GetRestorePoints } from '@wails/go/main/App'
const description = ref('电脑工具箱备份'); const points = ref<any[]>([]); const creating = ref(false); const loading = ref(false); const message = useMessage()
const columns = [{title:'名称',key:'name'},{title:'描述',key:'description',ellipsis:{tooltip:true}},{title:'创建时间',key:'createdAt',width:180}]
async function createPoint() { creating.value = true; try { await CreateRestorePoint(description.value); message.success('还原点创建成功'); loadPoints() } catch(e:any){message.error(String(e))}; creating.value = false }
async function loadPoints() { loading.value = true; try { const r = await GetRestorePoints(); if (r) points.value = r as any[] } catch(e){console.error(e)}; loading.value = false }
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
