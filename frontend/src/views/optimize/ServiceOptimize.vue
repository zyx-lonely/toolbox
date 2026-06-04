<template>
  <div>
    <n-h2>系统服务优化</n-h2>
    <n-p>按场景一键优化 Windows 系统服务配置。</n-p>

    <n-alert type="warning" :bordered="false" class="mb-4">
      <template #header>⚠️ 需要管理员权限</template>
      修改系统服务需要以管理员身份运行。建议先创建系统还原点。
    </n-alert>

    <n-grid :cols="2" :x-gap="16" :y-gap="16" class="mb-4">
      <n-gi v-for="profile in profiles" :key="profile.name">
        <n-card :title="profile.name" hoverable @click="applyProfile(profile.name)" size="small">
          {{ profile.description }}
        </n-card>
      </n-gi>
    </n-grid>

    <n-h3>服务列表</n-h3>
    <n-button type="primary" @click="loadServices" :loading="loadingServices" size="small" class="mb-4">
      刷新服务列表
    </n-button>

    <n-empty v-if="!loadingServices && services.length === 0" description="暂无服务数据" />

    <n-data-table
      v-if="services.length > 0"
      :columns="svcColumns"
      :data="services"
      :bordered="true"
      :loading="loadingServices"
      size="small"
      :max-height="400"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { NTag, useMessage } from 'naive-ui'
import { GetServices, GetOptimizationProfiles, ApplyOptimizationProfile } from '@wails/go/main/App'
import { startTypeLabel } from '@/api/bridge'

interface ServiceInfo { name: string; displayName: string; description: string; status: string; startType: string; recommended: string }
interface Profile { name: string; description: string; services: { name: string; action: string }[] }

const services = ref<ServiceInfo[]>([])
const profiles = ref<Profile[]>([])
const loadingServices = ref(false)
const message = useMessage()

const svcColumns = [
  { title: '名称', key: 'name', width: 160 },
  { title: '显示名称', key: 'displayName', ellipsis: { tooltip: true } },
  {
    title: '状态', key: 'status', width: 80,
    render: (row: ServiceInfo) => h(NTag, {
      type: row.status === 'running' ? 'success' as const : 'default' as const,
      size: 'small'
    }, { default: () => row.status === 'running' ? '运行中' : '已停止' })
  },
  {
    title: '启动类型', key: 'startType', width: 90,
    render: (row: ServiceInfo) => h(NTag, { size: 'small' }, { default: () => startTypeLabel(row.startType) })
  },
  { title: '建议', key: 'recommended', width: 80 }
]

async function loadServices() {
  loadingServices.value = true
  try {
    const r = await GetServices()
    if (r) services.value = r as ServiceInfo[]
  } catch (e) { console.error(e) }
  loadingServices.value = false
}

async function loadProfiles() {
  try {
    const r = await GetOptimizationProfiles()
    if (r) profiles.value = r as Profile[]
  } catch (e) { console.error(e) }
}

async function applyProfile(name: string) {
  try {
    const r = await ApplyOptimizationProfile(name)
    if (r) {
      message.success(`${name} 应用完成，共 ${r.length} 项变更`)
      loadServices()
    }
  } catch (e: any) {
    message.error(`应用失败: ${e}`)
  }
}

loadServices()
loadProfiles()
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
