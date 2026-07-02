<template>
  <div>
    <h2>时间戳转换</h2>
    <p>Unix 时间戳与日期格式互转。</p>
    <n-space class="mb-4">
      <n-button type="primary" @click="getNow">获取当前时间戳</n-button>
    </n-space>
    <n-grid :cols="2" :x-gap="16">
      <n-gi>
        <n-card title="时间戳 → 日期" size="small">
          <n-space vertical>
            <n-input-number v-model:value="timestamp" placeholder="输入时间戳" style="width: 100%" />
            <n-button @click="toDate">转换</n-button>
            <n-descriptions v-if="tsResult" bordered :column="1" size="small">
              <n-descriptions-item label="日期时间">{{ tsResult.dateTime }}</n-descriptions-item>
              <n-descriptions-item label="10位时间戳">{{ tsResult.unix10 }}</n-descriptions-item>
              <n-descriptions-item label="13位时间戳">{{ tsResult.unix13 }}</n-descriptions-item>
              <n-descriptions-item label="ISO 8601">{{ tsResult.iso8601 }}</n-descriptions-item>
            </n-descriptions>
          </n-space>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="日期 → 时间戳" size="small">
          <n-space vertical>
            <n-input v-model:value="dateStr" placeholder="如 2024-01-01 12:00:00" />
            <n-button @click="fromDate">转换</n-button>
            <n-descriptions v-if="dateResult" bordered :column="1" size="small">
              <n-descriptions-item label="10位时间戳">{{ dateResult.unix10 }}</n-descriptions-item>
              <n-descriptions-item label="13位时间戳">{{ dateResult.unix13 }}</n-descriptions-item>
              <n-descriptions-item label="ISO 8601">{{ dateResult.iso8601 }}</n-descriptions-item>
            </n-descriptions>
          </n-space>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { TimestampToDate, DateToTimestamp, GetNowTimestamp } from '@wails/go/main/App'

interface TsResult { timestamp: number; dateTime: string; unix10: number; unix13: number; iso8601: string }

const timestamp = ref<number | null>(null)
const dateStr = ref('')
const tsResult = ref<TsResult | null>(null)
const dateResult = ref<TsResult | null>(null)
const message = useMessage()

async function getNow() {
  try {
    const r = await GetNowTimestamp() as TsResult
    timestamp.value = r.unix10
    tsResult.value = r
  } catch (e: any) { message.error(String(e)) }
}

async function toDate() {
  if (!timestamp.value) return
  try { tsResult.value = await TimestampToDate(timestamp.value) as TsResult } catch (e: any) { message.error(String(e)) }
}

async function fromDate() {
  if (!dateStr.value) return
  try { dateResult.value = await DateToTimestamp(dateStr.value) as TsResult } catch (e: any) { message.error(String(e)) }
}
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
