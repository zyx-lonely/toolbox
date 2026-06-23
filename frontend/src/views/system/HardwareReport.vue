<template>
  <div>
    <n-h2>硬件检测报告</n-h2>
    <n-p>生成详细的系统硬件信息报告，可导出为 HTML 文件在浏览器中查看。</n-p>

    <n-space>
      <n-button type="primary" @click="generateReport" :loading="generating">
        生成报告
      </n-button>
    </n-space>

    <n-card v-if="report" class="mt-4" size="small">
      <template #header>
        <n-space justify="space-between">
          <span>系统信息摘要</span>
          <n-space>
            <n-button size="small" @click="openReport">在浏览器中打开</n-button>
            <n-button size="small" @click="copyPath">复制路径</n-button>
          </n-space>
        </n-space>
      </template>

      <n-grid :cols="3" :x-gap="16">
        <n-gi>
          <n-statistic label="操作系统" :value="report.os?.name || '-'" />
        </n-gi>
        <n-gi>
          <n-statistic label="处理器" :value="report.cpu?.cores + '核 / ' + report.cpu?.logicalCores + '线程'" />
        </n-gi>
        <n-gi>
          <n-statistic label="内存">
            <span>{{ formatBytes(report.memory?.total) }}</span>
          </n-statistic>
        </n-gi>
      </n-grid>

      <n-description-list label-placement="left" :column="2" class="mt-2">
        <n-description-item label="主机名">{{ report.os?.hostname }}</n-description-item>
        <n-description-item label="用户名">{{ report.os?.userName }}</n-description-item>
        <n-description-item label="系统版本">{{ report.os?.version }}</n-description-item>
        <n-description-item label="构建版本">{{ report.os?.buildNumber }}</n-description-item>
        <n-description-item label="已运行">{{ report.os?.uptime }}</n-description-item>
      </n-description-list>

      <n-h4>磁盘信息</n-h4>
      <n-data-table :columns="diskColumns" :data="report.disks || []" size="small" :bordered="true" />
    </n-card>

    <n-alert v-if="reportPath" type="success" :bordered="false" class="mt-4">
      <template #header>报告已生成</template>
      {{ reportPath }}
    </n-alert>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { GenerateSystemReport } from '@wails/go/main/App'
import { formatBytes } from '@/api/bridge'

const report = ref<any>(null)
const reportPath = ref('')
const generating = ref(false)
const message = useMessage()

const diskColumns = [
  { title: '卷标', key: 'label', width: 80 },
  { title: '文件系统', key: 'fileSystem', width: 90 },
  {
    title: '总容量', key: 'total', width: 100,
    render: (row: any) => formatBytes(row.total)
  },
  {
    title: '可用', key: 'free', width: 100,
    render: (row: any) => formatBytes(row.free)
  },
  {
    title: '使用率', key: 'usage', width: 150,
    render: (row: any) => {
      const pct = Math.round(row.usage || 0)
      return `${pct}%`
    }
  }
]

async function generateReport() {
  generating.value = true
  try {
    const r = await GenerateSystemReport("system")
    if (r) {
      reportPath.value = r as string
      message.success('报告已生成')
      // 自动打开报告
      window.open('file://' + r)
    }
  } catch (e: any) { 
    message.error(`生成失败: ${e}`) 
  }
  generating.value = false
}

</script>

<style scoped>
.mt-4 { margin-top: 16px; }
.mt-2 { margin-top: 8px; }
</style>
