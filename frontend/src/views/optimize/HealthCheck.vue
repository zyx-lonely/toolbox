<template>
  <div>
    <n-h2>系统健康体检</n-h2>
    <n-p>综合检测磁盘/内存/CPU/启动项/临时文件等，给出健康评分和优化建议。</n-p>

    <n-space class="mb-4">
      <n-button type="primary" size="large" @click="runCheck" :loading="checking">
        🔍 开始体检
      </n-button>
    </n-space>

    <n-empty v-if="!report && !checking" description="点击「开始体检」进行系统检测" />

    <n-result v-if="report" :status="report.score >= 60 ? 'success' : 'error'"
      :title="`健康评分: ${report.score}/100`"
      :description="report.score >= 80 ? '系统状态良好' : report.score >= 60 ? '需要优化' : '建议立即处理'">
    </n-result>

    <n-space vertical size="medium" class="mt-4">
      <n-card v-for="item in report?.items || []" :key="item.name" size="small">
        <n-space justify="space-between" align="center">
          <n-space>
            <n-tag :type="item.status === 'good' ? 'success' : item.status === 'warning' ? 'warning' : 'error'">
              {{ item.status === 'good' ? '良好' : item.status === 'warning' ? '警告' : '危险' }}
            </n-tag>
            <strong>{{ item.name }}</strong>
            <span style="color:#888">{{ item.value }}</span>
          </n-space>
        </n-space>
      </n-card>

      <n-card v-if="report?.suggestions?.length" title="优化建议" size="small">
        <n-space vertical>
          <n-alert v-for="(s, i) in report.suggestions" :key="i" type="warning" :bordered="false" closable>
            {{ s }}
          </n-alert>
        </n-space>
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { RunHealthCheck } from '@wails/go/main/App'

const report = ref<any>(null)
const checking = ref(false)
const message = useMessage()

async function runCheck() {
  checking.value = true
  report.value = null
  try {
    const r = await RunHealthCheck()
    if (r) report.value = r
    message.success('体检完成')
  } catch (e: any) { message.error(`体检失败: ${e}`) }
  checking.value = false
}
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
.mt-4 { margin-top: 16px; }
</style>
