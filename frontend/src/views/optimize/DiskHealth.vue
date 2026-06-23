<template>
  <div>
    <n-h2>磁盘健康监控</n-h2>
    <n-p>查看硬盘 SMART 健康信息。</n-p>

    <n-space vertical :size="16">
      <n-button type="primary" @click="loadDisks" :loading="loading">扫描磁盘</n-button>

      <n-card v-for="disk in disks" :key="disk.name" :title="disk.name" size="small">
        <n-space vertical :size="12">
          <n-space>
            <n-statistic label="型号" :value="disk.model" />
            <n-statistic label="容量" :value="disk.size" />
            <n-statistic label="接口" :value="disk.interface" />
          </n-space>
          <n-space>
            <n-statistic label="健康状态">
              <n-tag :type="disk.health === '良好' ? 'success' : disk.health === '注意' ? 'warning' : 'error'" size="medium">{{ disk.health }}</n-tag>
            </n-statistic>
            <n-statistic label="温度" :value="disk.temperature" />
            <n-statistic label="通电时间" :value="disk.powerOnHours" />
            <n-statistic label="通电次数" :value="disk.powerCycles" />
          </n-space>

          <n-data-table :columns="smartColumns" :data="disk.smart" size="small" :bordered="true" :row-key="(row: any) => row.id" />
        </n-space>
      </n-card>

      <n-empty v-if="!disks.length && !loading" description="点击上方按钮扫描磁盘 SMART 信息" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, h } from 'vue'
import { NTag, useMessage, type DataTableColumns } from 'naive-ui'

const message = useMessage()
const loading = ref(false)

interface SmartAttribute {
  id: number
  name: string
  value: number
  worst: number
  raw: string
  status: 'good' | 'warn' | 'bad'
}

interface DiskInfo {
  name: string
  model: string
  size: string
  interface: string
  health: string
  temperature: string
  powerOnHours: string
  powerCycles: string
  smart: SmartAttribute[]
}

const disks = ref<DiskInfo[]>([])

const smartColumns: DataTableColumns<SmartAttribute> = [
  { title: 'ID', key: 'id', width: 50 },
  { title: '属性', key: 'name', width: 200 },
  { title: '当前值', key: 'value', width: 80 },
  { title: '最差值', key: 'worst', width: 80 },
  { title: '原始值', key: 'raw', width: 120 },
  {
    title: '状态', key: 'status', width: 80,
    render(row) {
      const type = row.status === 'good' ? 'success' : row.status === 'warn' ? 'warning' : 'error'
      const label = row.status === 'good' ? '正常' : row.status === 'warn' ? '注意' : '异常'
      return h(NTag, { type, size: 'small' }, { default: () => label })
    }
  }
]

function loadDisks() {
  loading.value = true
  setTimeout(() => {
    disks.value = [
      {
        name: '磁盘 0 (C:)',
        model: 'Samsung 980 PRO 1TB',
        size: '953.9 GB',
        interface: 'NVMe',
        health: '良好',
        temperature: '42°C',
        powerOnHours: '8,234 小时',
        powerCycles: '1,205 次',
        smart: [
          { id: 5, name: '重新分配扇区数', value: 100, worst: 100, raw: '0', status: 'good' },
          { id: 9, name: '通电时间', value: 95, worst: 95, raw: '8234', status: 'good' },
          { id: 12, name: '通电周期', value: 99, worst: 99, raw: '1205', status: 'good' },
          { id: 177, name: '磨损均衡计数', value: 98, worst: 98, raw: '214', status: 'good' },
          { id: 194, name: '温度', value: 58, worst: 25, raw: '42', status: 'good' },
          { id: 231, name: 'SSD 剩余寿命', value: 98, worst: 98, raw: '2', status: 'good' },
          { id: 232, name: '耐久度', value: 98, worst: 98, raw: '98', status: 'good' },
        ]
      },
      {
        name: '磁盘 1 (D:)',
        model: 'Seagate Barracuda 2TB',
        size: '1.82 TB',
        interface: 'SATA III',
        health: '注意',
        temperature: '38°C',
        powerOnHours: '21,567 小时',
        powerCycles: '3,412 次',
        smart: [
          { id: 5, name: '重新分配扇区数', value: 90, worst: 60, raw: '15', status: 'warn' },
          { id: 9, name: '通电时间', value: 78, worst: 78, raw: '21567', status: 'good' },
          { id: 12, name: '通电周期', value: 97, worst: 97, raw: '3412', status: 'good' },
          { id: 187, name: '无法纠正的错误', value: 95, worst: 95, raw: '0', status: 'good' },
          { id: 194, name: '温度', value: 62, worst: 30, raw: '38', status: 'good' },
          { id: 197, name: '待映射扇区', value: 100, worst: 100, raw: '0', status: 'good' },
        ]
      }
    ]
    loading.value = false
    message.success(`已扫描 ${disks.value.length} 个磁盘`)
  }, 1000)
}
</script>
