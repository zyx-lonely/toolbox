<template>
  <div>
    <n-h2>HTTP API 调试</n-h2>
    <n-p>发送 HTTP 请求并查看响应，支持常用请求方法。</n-p>

    <n-card>
      <n-space vertical>
        <n-space>
          <n-select v-model:value="method" :options="methodOptions" style="width: 110px" />
          <n-input v-model:value="url" placeholder="请求 URL" clearable style="width: 400px" />
          <n-button type="primary" @click="sendRequest" :loading="sending">发送</n-button>
          <n-button @click="clearResult">清空</n-button>
        </n-space>

        <n-collapse>
          <n-collapse-item title="请求 Headers (JSON)">
            <n-input v-model:value="headers" type="textarea" :rows="3" placeholder='{"Content-Type": "application/json"}' />
          </n-collapse-item>
          <n-collapse-item title="请求 Body">
            <n-input v-model:value="body" type="textarea" :rows="5" placeholder="请求体内容..." />
          </n-collapse-item>
        </n-collapse>

        <n-divider />

        <n-space v-if="result" vertical>
          <n-space>
            <n-tag :type="statusType">{{ result.status }} {{ result.statusText }}</n-tag>
            <n-tag>耗时: {{ result.duration }}ms</n-tag>
            <n-tag>大小: {{ result.size }}</n-tag>
          </n-space>

          <n-tabs type="line" default-value="response">
            <n-tab-pane name="response" tab="响应体">
              <n-input
                v-model:value="result.body"
                type="textarea"
                :rows="12"
                readonly
                :style="{ fontFamily: 'monospace' }"
              />
            </n-tab-pane>
            <n-tab-pane name="responseHeaders" tab="响应 Headers">
              <n-input
                v-model:value="result.responseHeaders"
                type="textarea"
                :rows="8"
                readonly
                :style="{ fontFamily: 'monospace' }"
              />
            </n-tab-pane>
          </n-tabs>
        </n-space>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { SendHTTPRequest } from '@wails/go/main/App'

const method = ref('GET')
const url = ref('')
const headers = ref('')
const body = ref('')
const sending = ref(false)
const result = ref<any>(null)
const message = useMessage()

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' },
  { label: 'HEAD', value: 'HEAD' },
]

const statusType = computed(() => {
  if (!result.value) return 'default'
  const s = result.value.status
  if (s >= 200 && s < 300) return 'success'
  if (s >= 300 && s < 400) return 'warning'
  if (s >= 400) return 'error'
  return 'default'
})

async function sendRequest() {
  if (!url.value) {
    message.warning('请输入请求 URL')
    return
  }

  sending.value = true
  result.value = null

  try {
    let parsedHeaders: string | null = headers.value || null
    if (parsedHeaders) {
      try { JSON.parse(parsedHeaders) } catch { parsedHeaders = null }
    }
    const r = await SendHTTPRequest(method.value, url.value, parsedHeaders, body.value)
    if (r) {
      result.value = r
      message.success(`请求完成: ${r.status} ${r.statusText}`)
    }
  } catch (e: any) {
    message.error(String(e))
  }
  sending.value = false
}

function clearResult() {
  result.value = null
}
</script>

<style scoped>
</style>
