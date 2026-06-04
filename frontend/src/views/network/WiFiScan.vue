<template>
  <div class="wifi-scan">
    <n-space vertical>
      <n-space align="center">
        <n-button type="primary" @click="handleScan" :loading="scanning">
          扫描 WiFi
        </n-button>
        <n-tag v-if="scanning" type="info">扫描中...</n-tag>
      </n-space>

      <n-data-table
        :columns="columns"
        :data="signals"
        :bordered="true"
        :loading="scanning"
        :pagination="{ pageSize: 10 }"
      />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { NTag, useMessage, type DataTableColumns } from 'naive-ui'
import { ScanWiFiSignal } from '@wails/go/main/App'

const scanning = ref(false)
const signals = ref<any[]>([])
const message = useMessage()

interface SignalInfo {
  ssid: string
  bssid: string
  signal: number
  channel: number
  auth: string
}

const columns: DataTableColumns<SignalInfo> = [
  { title: 'SSID', key: 'ssid' },
  { title: 'BSSID', key: 'bssid' },
  { title: '信号强度', key: 'signal', render(row) { const s = row.signal; return h(NTag, { type: s > 70 ? 'success' : s > 40 ? 'warning' : 'error' }, { default: () => `${s}%` }) } },
  { title: '信道', key: 'channel' },
  { title: '认证方式', key: 'auth' },
]

async function handleScan() {
  scanning.value = true
  try {
    const result = await ScanWiFiSignal()
    signals.value = (result as any[]) || []
    if (signals.value.length === 0) {
      message.warning('未扫描到 WiFi 信号，请确保 WiFi 已开启')
    }
  } catch (e: any) {
    message.error(String(e))
    signals.value = []
  } finally {
    scanning.value = false
  }
}
</script>

<style scoped>
.wifi-scan {
  padding: 20px;
}
</style>
