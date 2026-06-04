<template>
  <div>
    <n-h2>单位换算</n-h2>
    <n-p>文件大小、长度、温度等单位换算。</n-p>

    <n-tabs type="line" default-value="file">
      <n-tab-pane name="file" tab="文件大小">
        <n-space vertical>
          <n-input-group>
            <n-input-number v-model:value="fileValue" :min="0" style="width: 150px" />
            <n-select v-model:value="fileFromUnit" :options="fileUnits" style="width: 100px" />
            <n-text class="eq">=</n-text>
            <n-input :value="formatFileResult" readonly style="width: 200px" />
            <n-select v-model:value="fileToUnit" :options="fileUnits" style="width: 100px" />
          </n-input-group>
          <n-text depth="3">结果: {{ formatFileResult }} {{ fileToUnit }}</n-text>
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="temp" tab="温度">
        <n-space vertical>
          <n-input-group>
            <n-input-number v-model:value="tempValue" :min="-273" style="width: 150px" />
            <n-select v-model:value="tempFromUnit" :options="tempUnits" style="width: 80px" />
            <n-text class="eq">=</n-text>
            <n-input :value="tempResult" readonly style="width: 200px" />
            <n-select v-model:value="tempToUnit" :options="tempUnits" style="width: 80px" />
          </n-input-group>
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="length" tab="长度">
        <n-space vertical>
          <n-input-group>
            <n-input-number v-model:value="lengthValue" :min="0" style="width: 150px" />
            <n-select v-model:value="lengthFromUnit" :options="lengthUnits" style="width: 100px" />
            <n-text class="eq">=</n-text>
            <n-input :value="lengthResult" readonly style="width: 200px" />
            <n-select v-model:value="lengthToUnit" :options="lengthUnits" style="width: 100px" />
          </n-input-group>
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const fileValue = ref(1024)
const fileFromUnit = ref('KB')
const fileToUnit = ref('MB')

const fileUnits = [
  { label: 'B', value: 'B' },
  { label: 'KB', value: 'KB' },
  { label: 'MB', value: 'MB' },
  { label: 'GB', value: 'GB' },
  { label: 'TB', value: 'TB' }
]

const formatFileResult = computed(() => {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const fromIdx = units.indexOf(fileFromUnit.value)
  const toIdx = units.indexOf(fileToUnit.value)
  if (fromIdx === -1 || toIdx === -1) return '-'
  const bytes = fileValue.value * Math.pow(1024, fromIdx)
  const result = bytes / Math.pow(1024, toIdx)
  return result.toFixed(2)
})

const tempValue = ref(25)
const tempFromUnit = ref('°C')
const tempToUnit = ref('°F')

const tempUnits = [
  { label: '°C', value: '°C' },
  { label: '°F', value: '°F' },
  { label: 'K', value: 'K' }
]

const tempResult = computed(() => {
  let celsius: number
  if (tempFromUnit.value === '°C') celsius = tempValue.value
  else if (tempFromUnit.value === '°F') celsius = (tempValue.value - 32) * 5 / 9
  else celsius = tempValue.value - 273.15

  if (tempToUnit.value === '°C') return celsius.toFixed(2)
  if (tempToUnit.value === '°F') return (celsius * 9 / 5 + 32).toFixed(2)
  return (celsius + 273.15).toFixed(2)
})

const lengthValue = ref(1)
const lengthFromUnit = ref('m')
const lengthToUnit = ref('cm')

const lengthUnits = [
  { label: 'mm', value: 'mm' },
  { label: 'cm', value: 'cm' },
  { label: 'm', value: 'm' },
  { label: 'km', value: 'km' },
  { label: '英寸', value: 'in' },
  { label: '英尺', value: 'ft' }
]

const lengthResult = computed(() => {
  const toMM: Record<string, number> = { mm: 1, cm: 10, m: 1000, km: 1000000, in: 25.4, ft: 304.8 }
  const fromMM = toMM[lengthFromUnit.value]
  const toMMVal = toMM[lengthToUnit.value]
  if (!fromMM || !toMMVal) return '-'
  return ((lengthValue.value * fromMM) / toMMVal).toFixed(4)
})
</script>

<style scoped>
.eq { padding: 0 12px; font-size: 18px; font-weight: bold; line-height: 34px; }
</style>
