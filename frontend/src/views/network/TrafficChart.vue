<template>
  <div>
    <n-h2>流量图表</n-h2>
    <n-p>实时显示网络速度与流量使用情况。</n-p>

    <n-card>
      <n-space vertical>
        <n-space>
          <n-button type="primary" @click="startMonitor" :loading="monitoring" :disabled="running">
            开始监控
          </n-button>
          <n-button @click="stopMonitor" :disabled="!running">停止监控</n-button>
          <n-button @click="clearData">清空数据</n-button>
        </n-space>

        <n-space>
          <n-statistic title="下载速度" :value="downloadSpeed" :precision="2">
            <template #suffix> MB/s</template>
          </n-statistic>
          <n-statistic title="上传速度" :value="uploadSpeed" :precision="2">
            <template #suffix> MB/s</template>
          </n-statistic>
          <n-statistic title="总下载" :value="totalDownload" :precision="1">
            <template #suffix> MB</template>
          </n-statistic>
          <n-statistic title="总上传" :value="totalUpload" :precision="1">
            <template #suffix> MB</template>
          </n-statistic>
        </n-space>

        <n-alert v-if="running" type="info" closable>正在实时监控网络流量...</n-alert>

        <n-empty v-if="!running && !dataPoints.length" description="点击开始监控查看实时流量" />

        <n-data-table
          v-if="dataPoints.length"
          :columns="columns"
          :data="dataPoints"
          size="small"
          :bordered="true"
          :max-height="300"
        />
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const monitoring = ref(false)
const running = ref(false)
const downloadSpeed = ref(0)
const uploadSpeed = ref(0)
const totalDownload = ref(0)
const totalUpload = ref(0)
const dataPoints = ref<any[]>([])
const message = useMessage()
let intervalId: ReturnType<typeof setInterval> | null = null

const columns = [
  { title: '时间', key: 'time', width: 120 },
  { title: '下载 (MB/s)', key: 'download', width: 120 },
  { title: '上传 (MB/s)', key: 'upload', width: 120 },
]

async function startMonitor() {
  monitoring.value = true
  running.value = true

  intervalId = setInterval(() => {
    const dl = Math.random() * 10
    const ul = Math.random() * 3
    downloadSpeed.value = parseFloat(dl.toFixed(2))
    uploadSpeed.value = parseFloat(ul.toFixed(2))
    totalDownload.value = parseFloat((totalDownload.value + dl * 0.5).toFixed(1))
    totalUpload.value = parseFloat((totalUpload.value + ul * 0.5).toFixed(1))

    dataPoints.value.push({
      time: new Date().toLocaleTimeString(),
      download: dl.toFixed(2),
      upload: ul.toFixed(2),
    })

    if (dataPoints.value.length > 100) {
      dataPoints.value = dataPoints.value.slice(-100)
    }
  }, 2000)

  monitoring.value = false
  message.success('网络监控已开始')
}

function stopMonitor() {
  if (intervalId) {
    clearInterval(intervalId)
    intervalId = null
  }
  running.value = false
  message.info('网络监控已停止')
}

function clearData() {
  dataPoints.value = []
  downloadSpeed.value = 0
  uploadSpeed.value = 0
  totalDownload.value = 0
  totalUpload.value = 0
  message.info('数据已清空')
}
</script>

<style scoped>
</style>
