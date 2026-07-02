<template>
  <div>
    <h2>电池信息</h2>
    <p>查看笔记本电池健康度和状态。</p>
    <n-empty v-if="!loading && batteries.length === 0" description="未检测到电池（台式机不支持）" />
    <n-spin :show="loading">
      <n-grid :cols="2" :x-gap="12" :y-gap="12">
        <n-gi v-for="b in batteries" :key="b.name">
          <n-card :title="b.name || '电池'" size="small">
            <n-space vertical>
              <n-progress type="circle" :percentage="b.chargeLevel" :color="getColor(b.chargeLevel)" style="width: 100px; height: 100px; display: block; margin: 0 auto;">
                <span style="font-size: 20px; font-weight: bold;">{{ b.chargeLevel }}%</span>
              </n-progress>
              <n-descriptions bordered :column="1" size="small">
                <n-descriptions-item label="状态">{{ b.status }}</n-descriptions-item>
                <n-descriptions-item label="健康度">{{ b.healthPercent.toFixed(1) }}%</n-descriptions-item>
                <n-descriptions-item label="设计容量">{{ b.designCapacity }} mWh</n-descriptions-item>
                <n-descriptions-item label="满充容量">{{ b.fullCapacity }} mWh</n-descriptions-item>
                <n-descriptions-item label="电压">{{ b.voltage.toFixed(2) }} V</n-descriptions-item>
                <n-descriptions-item label="温度" v-if="b.temperature > 0">{{ b.temperature.toFixed(1) }} °C</n-descriptions-item>
              </n-descriptions>
            </n-space>
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetBatteryInfo } from '@wails/go/main/App'

interface BatteryInfo { name: string; status: string; chargeLevel: number; designCapacity: number; fullCapacity: number; currentCapacity: number; voltage: number; temperature: number; cycleCount: number; healthPercent: number; isPresent: boolean }

const batteries = ref<BatteryInfo[]>([])
const loading = ref(false)

function getColor(level: number): string {
  if (level > 60) return '#18a058'
  if (level > 20) return '#f0a020'
  return '#d03050'
}

onMounted(async () => {
  loading.value = true
  try { batteries.value = await GetBatteryInfo() as BatteryInfo[] } catch (e) { console.error(e) }
  loading.value = false
})
</script>
