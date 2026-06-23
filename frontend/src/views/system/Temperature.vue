<template>
  <div>
    <n-h2>温度监控</n-h2>
    <n-p>查看硬件温度信息（CPU、GPU、硬盘等）。</n-p>

    <n-alert v-if="errorMsg" type="warning" style="margin-bottom: 16px">
      {{ errorMsg }}
    </n-alert>

    <n-button type="primary" @click="loadTemperatures" :loading="loading" class="mb-4">
      刷新温度
    </n-button>

    <n-empty v-if="!loading && temperatures.length === 0 && !errorMsg" description="暂无温度数据" />

    <n-grid v-if="temperatures.length > 0" :cols="3" :x-gap="16" :y-gap="16" class="mt-4">
      <n-gi v-for="temp in temperatures" :key="temp.name">
        <n-card :title="getSensorName(temp.name)" hoverable>
          <n-space vertical>
            <n-progress
              type="dashboard"
              :percentage="getTemperaturePercentage(temp.temperature)"
              :color="getTemperatureColor(temp.temperature)"
              :height="120"
            >
              {{ temp.temperature.toFixed(1) }}°C
            </n-progress>
            <n-text depth="3" style="font-size: 12px">{{ temp.name }}</n-text>
          </n-space>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { GetTemperatures } from '@wails/go/main/App'

interface TemperatureInfo {
  name: string
  temperature: number
  unit: string
}

const temperatures = ref<TemperatureInfo[]>([])
const loading = ref(false)
const errorMsg = ref('')
const message = useMessage()

async function loadTemperatures() {
  loading.value = true
  errorMsg.value = ''
  try {
    const result = await GetTemperatures()
    if (result) {
      temperatures.value = result as TemperatureInfo[]
      if (temperatures.value.length === 0) {
        errorMsg.value = '未检测到温度传感器，可能需要管理员权限或安装第三方工具（如 Open Hardware Monitor）'
      }
    }
  } catch (e: any) {
    errorMsg.value = String(e)
    temperatures.value = []
  }
  loading.value = false
}

function getSensorName(name: string): string {
  // 简化传感器名称
  if (name.includes('CPU') || name.includes('Processor')) return 'CPU'
  if (name.includes('GPU') || name.includes('Display')) return 'GPU'
  if (name.includes('Disk') || name.includes('HDD')) return '硬盘'
  if (name.includes('Battery')) return '电池'
  return name.split('\\').pop() || name
}

function getTemperaturePercentage(temp: number): number {
  // 将温度转换为百分比（0-100°C 对应 0-100%）
  return Math.min(Math.max(temp, 0), 100)
}

function getTemperatureColor(temp: number): string {
  if (temp < 40) return '#18a058'  // 绿色 - 正常
  if (temp < 60) return '#f0a020'  // 黄色 - 注意
  if (temp < 80) return '#f56c6c'  // 橙色 - 警告
  return '#d03050'  // 红色 - 危险
}

onMounted(loadTemperatures)
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
.mt-4 { margin-top: 20px; }
</style>
