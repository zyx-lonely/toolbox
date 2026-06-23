<template>
  <div>
    <n-h2>导出系统报告</n-h2>
    <n-p>生成系统健康报告，可打印或保存为 PDF。</n-p>

    <n-card title="报告选项">
      <n-space vertical style="width: 100%">
        <n-divider style="margin: 12px 0">报告类型</n-divider>
        <n-radio-group v-model:value="reportType" name="reportType">
          <n-space>
            <n-radio value="system">系统信息</n-radio>
            <n-radio value="monitor">硬件监控</n-radio>
            <n-radio value="full">完整报告</n-radio>
          </n-space>
        </n-radio-group>

        <n-divider style="margin: 12px 0">包含内容</n-divider>
        <n-space>
          <n-checkbox v-model:checked="includeSystem">系统信息</n-checkbox>
          <n-checkbox v-model:checked="includeMonitor">硬件监控</n-checkbox>
          <n-checkbox v-model:checked="includeTemp">温度信息</n-checkbox>
          <n-checkbox v-model:checked="includeStartup">启动项</n-checkbox>
          <n-checkbox v-model:checked="includeProcess">进程列表</n-checkbox>
        </n-space>

        <n-divider style="margin: 12px 0"></n-divider>
        <n-space>
          <n-button type="primary" @click="generateReport" :loading="generating">
            <template #icon><n-icon><document-text-outline /></n-icon></template>
            生成报告
          </n-button>
          <n-button @click="openReport" :disabled="!reportHTML" v-if="reportHTML">
            <template #icon><n-icon><open-outline /></n-icon></template>
            打开报告
          </n-button>
          <n-button @click="printReport" :disabled="!reportHTML" v-if="reportHTML">
            <template #icon><n-icon><print-outline /></n-icon></template>
            打印 / 保存为 PDF
          </n-button>
        </n-space>
      </n-space>
    </n-card>

    <n-card title="预览" style="margin-top: 16px" v-if="reportHTML">
      <n-iframe :srcdoc="reportHTML" style="width: 100%; height: 600px; border: 1px solid #e8e8e8; border-radius: 6px"></n-iframe>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { DocumentTextOutline, OpenOutline, PrintOutline } from '@vicons/ionicons5'
import { GenerateSystemReport } from '@wails/go/main/App'

const reportType = ref('full')
const includeSystem = ref(true)
const includeMonitor = ref(true)
const includeTemp = ref(true)
const includeStartup = ref(true)
const includeProcess = ref(true)
const generating = ref(false)
const reportHTML = ref('')
const message = useMessage()

async function generateReport() {
  generating.value = true
  try {
    const r = await GenerateSystemReport(reportType.value)
    if (r) {
      reportHTML.value = r as string
      message.success('报告生成成功')
    }
  } catch (e: any) {
    message.error(`生成失败: ${e}`)
  }
  generating.value = false
}

function openReport() {
  // 创建 Blob 并打开
  const blob = new Blob([reportHTML.value], { type: 'text/html' })
  const url = URL.createObjectURL(blob)
  window.open(url, '_blank')
}

function printReport() {
  // 创建新窗口并打印
  const blob = new Blob([reportHTML.value], { type: 'text/html' })
  const url = URL.createObjectURL(blob)
  const win = window.open(url, '_blank')
  if (win) {
    win.onload = () => {
      win.print()
    }
  }
}
</script>
