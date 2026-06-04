<template>
  <div>
    <n-h2>局域网拓扑</n-h2>
    <n-p>扫描局域网中的活跃设备。</n-p>
    <n-space class="mb-4">
      <n-button type="primary" @click="scan" :loading="scanning">扫描局域网设备</n-button>
      <n-tag v-if="devices.length" type="info">发现 {{ devices.length }} 台设备</n-tag>
    </n-space>

    <n-empty v-if="!scanning && !devices.length" description="点击按钮开始扫描" />

    <n-grid v-if="devices.length" :cols="4" :x-gap="12" :y-gap="12">
      <n-gi v-for="(d, i) in devices" :key="i">
        <n-card :title="d.hostname || '未知设备'" size="small" hoverable>
          <n-descriptions :column="1" size="small">
            <n-descriptions-item label="IP">{{ d.ip }}</n-descriptions-item>
            <n-descriptions-item label="MAC">{{ d.mac || '-' }}</n-descriptions-item>
            <n-descriptions-item label="状态">
              <n-tag :type="d.alive ? 'success' : 'warning'" size="small">{{ d.alive ? '在线' : '离线' }}</n-tag>
            </n-descriptions-item>
          </n-descriptions>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { ScanLAN } from '@wails/go/main/App'
const devices = ref<any[]>([]); const scanning = ref(false); const message = useMessage()
async function scan() {
  scanning.value = true
  try { const r = await ScanLAN(); if (r) devices.value = r as any[]; message.success(`发现 ${devices.value.length} 台设备`) }
  catch(e:any) { message.error(String(e)) }
  scanning.value = false
}
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
