<template>
  <div>
    <n-h2>批量 Ping</n-h2>
    <n-p>对指定 IP 段执行批量 Ping 扫描，快速发现在线主机。</n-p>

    <n-card>
      <n-space vertical>
        <n-input-group>
          <n-input v-model:value="ipPrefix" placeholder="IP 前缀, 例如 192.168.1" style="width: 200px" />
          <n-input-number v-model:value="startRange" placeholder="起始" :min="1" :max="254" style="width: 100px" />
          <n-input-number v-model:value="endRange" placeholder="结束" :min="1" :max="254" style="width: 100px" />
          <n-button type="primary" @click="startPing" :loading="pinging">开始扫描</n-button>
          <n-button @click="clearResults">清空</n-button>
        </n-input-group>

        <n-progress v-if="pinging" type="line" :percentage="progress" :indicator-placement="'inside'" />

        <n-empty v-if="!pinging && !results.length" description="输入 IP 段后点击开始扫描" />

        <n-data-table
          v-if="results.length"
          :columns="columns"
          :data="results"
          size="small"
          :bordered="true"
          :max-height="400"
        />
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { BatchPing } from '@wails/go/main/App'

const ipPrefix = ref('192.168.1')
const startRange = ref(1)
const endRange = ref(254)
const pinging = ref(false)
const progress = ref(0)
const results = ref<any[]>([])
const message = useMessage()

const columns = [
  { title: 'IP 地址', key: 'ip', width: 140 },
  { title: '状态', key: 'alive', width: 80, render: (r: any) => r.alive ? '在线' : '离线' },
  { title: '延迟 (ms)', key: 'latency', width: 100, render: (r: any) => r.alive ? r.latency : '-' },
  { title: 'TTL', key: 'ttl', width: 60, render: (r: any) => r.ttl ?? '-' },
  { title: '主机名', key: 'hostname', ellipsis: { tooltip: true } },
]

const totalIps = computed(() => endRange.value - startRange.value + 1)

async function startPing() {
  if (!ipPrefix.value) {
    message.warning('请输入 IP 前缀')
    return
  }
  if (startRange.value > endRange.value) {
    message.warning('起始值不能大于结束值')
    return
  }

  pinging.value = true
  progress.value = 0
  results.value = []

  try {
    const r = await BatchPing(ipPrefix.value, startRange.value, endRange.value)
    if (r) {
      results.value = r as any[]
      message.success(`扫描完成，发现 ${results.value.filter((x: any) => x.alive).length} 台在线主机`)
    }
  } catch (e: any) {
    message.error(String(e))
  }
  pinging.value = false
  progress.value = 100
}

function clearResults() {
  results.value = []
  progress.value = 0
}
</script>

<style scoped>
</style>
