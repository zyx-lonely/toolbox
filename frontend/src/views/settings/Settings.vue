<template>
  <div>
    <n-h2>设置</n-h2>

    <n-card title="外观" class="mb-4">
      <n-space align="center" justify="space-between">
        <span>深色模式</span>
        <n-switch v-model:value="appStore.darkMode" @update:value="onDarkModeChange" />
      </n-space>
    </n-card>

    <n-card title="版本信息" class="mb-4">
      <n-description-list label-placement="left" :column="1">
        <n-description-item label="应用名称">电脑工具箱</n-description-item>
        <n-description-item label="版本">1.0.0 (Build 20260603)</n-description-item>
        <n-description-item label="技术栈">Go + Wails v2 + Vue 3 + Naive UI</n-description-item>
        <n-description-item label="适用平台">Windows 10/11</n-description-item>
      </n-description-list>
    </n-card>

    <n-card title="免责声明" class="mb-4">
      <n-alert type="warning" :bordered="false">
        <template #header>⚠️ 使用须知</template>
        <ol style="padding-left: 20px; line-height: 1.8;">
          <li>本工具为免费开源软件，仅供学习和个人使用。</li>
          <li>使用本工具中的系统优化、清理等功能前，建议先创建系统还原点。</li>
          <li>部分功能需要管理员权限才能正常工作。</li>
          <li>作者不对因使用本工具导致的任何数据丢失或系统问题承担责任。</li>
          <li>本工具不会收集或上传任何用户数据。</li>
        </ol>
      </n-alert>
    </n-card>

    <n-card title="许可证">
      <n-text depth="3">MIT License - 可自由使用、修改和分发，但需保留原版权声明。</n-text>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { useAppStore } from '../../store/app'
import { onMounted, ref } from 'vue'
const appStore = useAppStore()
const version = ref('1.0.0')
const buildDate = ref('20260603')
onMounted(async () => {
  try {
    const { GetAppVersion, GetBuildDate } = await import('@wails/go/main/App')
    const v = await GetAppVersion(); if (v) version.value = v as string
    const d = await GetBuildDate(); if (d) buildDate.value = d as string
  } catch(e) { console.error(e) }
})
function onDarkModeChange(val: boolean) { appStore.darkMode = val }
</script>

<style scoped>
.mb-4 { margin-top: 16px; }
</style>
