<template>
  <div>
    <n-h2>局域网设备扫描</n-h2>
    <n-p>扫描局域网中活跃的网络设备。</n-p>
    <n-button type="primary" @click="scan" :loading="scanning" class="mb-4">开始扫描</n-button>
    <n-empty v-if="!scanning && !devices.length" description="点击开始扫描" />
    <n-data-table v-if="devices.length" :columns="columns" :data="devices" size="small" :bordered="true" />
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { ScanLAN } from '@wails/go/main/App'
const devices = ref<any[]>([])
const scanning = ref(false)
const message = useMessage()
const columns = [
  { title: 'IP 地址', key: 'ip', width: 140 },
  { title: 'MAC 地址', key: 'mac', width: 160 },
  { title: '主机名', key: 'hostname', ellipsis: { tooltip: true } },
  { title: '状态', key: 'alive', width: 80, render: (r: any) => r.alive ? '活跃' : '离线' },
]
async function scan() {
  scanning.value = true
  try {
    const r = await ScanLAN()
    if (r) devices.value = r as any[]
    message.success(`扫描完成，发现 ${devices.value.length} 台设备`)
  } catch (e: any) { message.error(String(e)) }
  scanning.value = false
}
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
