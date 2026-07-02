<template>
  <div>
    <h2>高级编解码</h2>
    <p>Base32/Hex 互转。</p>
    <n-tabs>
      <n-tab-pane name="b2h" tab="Base32 → Hex">
        <n-space vertical>
          <n-input v-model:value="base32Input" type="textarea" placeholder="输入 Base32 字符串" :rows="3" />
          <n-button type="primary" @click="base32ToHex" :loading="loading">转换</n-button>
          <n-input v-if="hexOutput" :value="hexOutput" readonly type="textarea" :rows="2" />
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="h2b" tab="Hex → Base32">
        <n-space vertical>
          <n-input v-model:value="hexInput" type="textarea" placeholder="输入 Hex 字符串" :rows="3" />
          <n-button type="primary" @click="hexToBase32" :loading="loading">转换</n-button>
          <n-input v-if="base32Output" :value="base32Output" readonly type="textarea" :rows="2" />
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { Base32ToHex, HexToBase32 } from '@wails/go/main/App'

const base32Input = ref('')
const hexInput = ref('')
const base32Output = ref('')
const hexOutput = ref('')
const loading = ref(false)
const message = useMessage()

async function base32ToHex() {
  loading.value = true
  try { hexOutput.value = await Base32ToHex(base32Input.value) as string } catch (e: any) { message.error(String(e)) }
  loading.value = false
}

async function hexToBase32() {
  loading.value = true
  try { base32Output.value = await HexToBase32(hexInput.value) as string } catch (e: any) { message.error(String(e)) }
  loading.value = false
}
</script>
