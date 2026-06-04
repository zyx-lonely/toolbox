<template>
  <div>
    <n-h2>Cron 表达式生成器</n-h2>
    <n-p>图形化配置 Cron 定时表达式。</n-p>
    <n-card>
      <n-space vertical>
        <n-space><span style="width:80px">分钟:</span><n-input-number v-model:value="cronMinute" :min="0" :max="59" style="width:80px"/><n-checkbox v-model:checked="cronEveryMin">每分</n-checkbox></n-space>
        <n-space><span style="width:80px">小时:</span><n-input-number v-model:value="cronHour" :min="0" :max="23" style="width:80px"/><n-checkbox v-model:checked="cronEveryHour">每时</n-checkbox></n-space>
        <n-space><span style="width:80px">日期:</span><n-input-number v-model:value="cronDay" :min="1" :max="31" style="width:80px"/><n-checkbox v-model:checked="cronEveryDay">每日</n-checkbox></n-space>
        <n-space><span style="width:80px">月份:</span><n-input-number v-model:value="cronMonth" :min="1" :max="12" style="width:80px"/><n-checkbox v-model:checked="cronEveryMonth">每月</n-checkbox></n-space>
        <n-space><span style="width:80px">星期:</span><n-input-number v-model:value="cronWeek" :min="0" :max="6" style="width:80px"/><n-checkbox v-model:checked="cronEveryWeek">每日</n-checkbox></n-space>
        <n-divider />
        <n-input :value="cronExpression" readonly size="large" style="font-family:monospace;font-size:16px;text-align:center" />
        <n-button type="primary" @click="copy">复制表达式</n-button>
      </n-space>
    </n-card>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
const message = useMessage()
const cronMinute = ref(0); const cronEveryMin = ref(true)
const cronHour = ref(0); const cronEveryHour = ref(true)
const cronDay = ref(1); const cronEveryDay = ref(true)
const cronMonth = ref(1); const cronEveryMonth = ref(true)
const cronWeek = ref(0); const cronEveryWeek = ref(true)
const cronExpression = computed(() => {
  const m = cronEveryMin.value ? '*' : String(cronMinute.value)
  const h = cronEveryHour.value ? '*' : String(cronHour.value)
  const d = cronEveryDay.value ? '*' : String(cronDay.value)
  const mo = cronEveryMonth.value ? '*' : String(cronMonth.value)
  const w = cronEveryWeek.value ? '*' : String(cronWeek.value)
  return `${m} ${h} ${d} ${mo} ${w}`
})
async function copy() {
  try { await navigator.clipboard.writeText(cronExpression.value); message.success('已复制') }
  catch { message.warning('复制失败') }
}
</script>
