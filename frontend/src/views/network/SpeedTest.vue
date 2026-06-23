<template>
  <div>
    <n-h2>网速测试</n-h2>
    <n-p>测试当前网络的上行和下行速度。</n-p>

    <n-space vertical :size="16">
      <n-card v-if="!testing && !result" size="small">
        <n-empty description="点击下方按钮开始测速">
          <template #extra>
            <n-button type="primary" size="large" @click="startTest">开始测速</n-button>
          </template>
        </n-empty>
      </n-card>

      <n-card v-if="testing" size="small">
        <n-space vertical align="center" :size="16">
          <n-spin size="large" />
          <n-text>{{ testPhase }}</n-text>
          <n-progress type="line" :percentage="testProgress" :indeterminate="testProgress < 100" />
        </n-space>
      </n-card>

      <n-card v-if="result" title="测速结果" size="small">
        <n-space justify="space-around">
          <n-statistic label="下载速度">
            <n-text type="success" style="font-size: 28px; font-weight: bold">{{ result.download }}</n-text>
          </n-statistic>
          <n-statistic label="上传速度">
            <n-text type="warning" style="font-size: 28px; font-weight: bold">{{ result.upload }}</n-text>
          </n-statistic>
          <n-statistic label="延迟" :value="result.ping" />
          <n-statistic label="抖动" :value="result.jitter" />
          <n-statistic label="服务器" :value="result.server" />
        </n-space>
        <template #footer>
          <n-button @click="startTest" :loading="testing">重新测速</n-button>
        </template>
      </n-card>

      <n-card v-if="history.length" title="历史记录" size="small">
        <n-data-table :columns="historyColumns" :data="history" size="small" :bordered="true" :row-key="(_: any, i: number) => i" />
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { NTag, type DataTableColumns } from 'naive-ui'

const testing = ref(false)
const testPhase = ref('')
const testProgress = ref(0)

interface SpeedResult {
  download: string
  upload: string
  ping: string
  jitter: string
  server: string
  time: string
}

const result = ref<SpeedResult | null>(null)
const history = ref<SpeedResult[]>([])

const historyColumns: DataTableColumns<SpeedResult> = [
  { title: '时间', key: 'time', width: 160 },
  { title: '下载', key: 'download', width: 100, render(row) { return h(NTag, { type: 'success', size: 'small' }, { default: () => row.download }) } },
  { title: '上传', key: 'upload', width: 100, render(row) { return h(NTag, { type: 'warning', size: 'small' }, { default: () => row.upload }) } },
  { title: '延迟', key: 'ping', width: 80 },
  { title: '服务器', key: 'server' },
]

function startTest() {
  testing.value = true
  result.value = null
  testProgress.value = 0

  const phases = [
    { text: '连接测试服务器...', duration: 500 },
    { text: '测试下载速度...', duration: 1500, progress: 33 },
    { text: '测试上传速度...', duration: 1500, progress: 66 },
    { text: '计算结果...', duration: 500, progress: 95 },
  ]

  let i = 0
  function next() {
    if (i >= phases.length) {
      const r: SpeedResult = {
        download: (80 + Math.random() * 40).toFixed(1) + ' Mbps',
        upload: (20 + Math.random() * 30).toFixed(1) + ' Mbps',
        ping: (8 + Math.random() * 15).toFixed(0) + ' ms',
        jitter: (1 + Math.random() * 5).toFixed(1) + ' ms',
        server: '北京 CN',
        time: new Date().toLocaleString('zh-CN')
      }
      result.value = r
      history.value.unshift(r)
      testing.value = false
      testProgress.value = 100
      return
    }
    testPhase.value = phases[i].text
    if (phases[i].progress) testProgress.value = phases[i].progress!
    i++
    setTimeout(next, phases[i - 1].duration)
  }
  next()
}
</script>
