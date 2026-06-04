<template>
  <div>
    <n-h2>编解码工具</n-h2>
    <n-p>Base64 / URL 编解码。</n-p>

    <n-tabs type="line" default-value="base64">
      <n-tab-pane name="base64" tab="Base64">
        <n-space vertical>
          <n-input type="textarea" v-model:value="base64Input" :rows="5" placeholder="输入要编码/解码的文本" />
          <n-space>
            <n-button type="primary" @click="base64Encode">编码</n-button>
            <n-button @click="base64Decode">解码</n-button>
          </n-space>
          <n-input type="textarea" v-model:value="base64Output" :rows="5" placeholder="结果" readonly />
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="url" tab="URL">
        <n-space vertical>
          <n-input type="textarea" v-model:value="urlInput" :rows="5" placeholder="输入要编码/解码的 URL" />
          <n-space>
            <n-button type="primary" @click="urlEncode">编码</n-button>
            <n-button @click="urlDecode">解码</n-button>
          </n-space>
          <n-input type="textarea" v-model:value="urlOutput" :rows="5" placeholder="结果" readonly />
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { EncodeBase64, DecodeBase64, EncodeURL, DecodeURL } from '@wails/go/main/App'

const base64Input = ref('')
const base64Output = ref('')
const urlInput = ref('')
const urlOutput = ref('')

async function base64Encode() {
  try {
    const r = await EncodeBase64(base64Input.value)
    if (r) base64Output.value = r.output
  } catch (e) { console.error(e) }
}

async function base64Decode() {
  try {
    const r = await DecodeBase64(base64Input.value)
    if (r) base64Output.value = r.error || r.output
  } catch (e) { console.error(e) }
}

async function urlEncode() {
  try {
    const r = await EncodeURL(urlInput.value)
    if (r) urlOutput.value = r.output
  } catch (e) { console.error(e) }
}

async function urlDecode() {
  try {
    const r = await DecodeURL(urlInput.value)
    if (r) urlOutput.value = r.error || r.output
  } catch (e) { console.error(e) }
}
</script>
