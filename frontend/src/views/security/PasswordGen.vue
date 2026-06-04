<template>
  <div>
    <n-h2>密码生成器</n-h2>
    <n-p>生成高强度随机密码。</n-p>

    <n-card>
      <n-space vertical>
        <n-space align="center">
          <n-input v-model:value="generatedPassword" readonly size="large" style="width: 350px; font-family: monospace; font-size: 18px;" />
          <n-button @click="copyPassword" type="primary">复制</n-button>
        </n-space>

        <n-divider />

        <n-space vertical>
          <n-space align="center">
            <span>密码长度:</span>
            <n-slider v-model:value="pwdLength" :min="4" :max="64" style="width: 200px" />
            <n-input-number v-model:value="pwdLength" :min="4" :max="64" style="width: 80px" />
          </n-space>

          <n-space>
            <n-checkbox v-model:checked="useUpper">大写字母 (A-Z)</n-checkbox>
            <n-checkbox v-model:checked="useLower">小写字母 (a-z)</n-checkbox>
            <n-checkbox v-model:checked="useDigits">数字 (0-9)</n-checkbox>
            <n-checkbox v-model:checked="useSpecial">特殊字符 (!@#$%)</n-checkbox>
          </n-space>

          <n-tag v-if="strength" :type="strengthColor as any">
            密码强度: {{ strengthLabel }}
          </n-tag>

          <n-button type="primary" @click="generate" size="large">生成密码</n-button>
        </n-space>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { GeneratePassword } from '@wails/go/main/App'

const generatedPassword = ref('')
const pwdLength = ref(16)
const useUpper = ref(true)
const useLower = ref(true)
const useDigits = ref(true)
const useSpecial = ref(true)
const strength = ref('')
const message = useMessage()

const strengthLabel = computed(() => {
  const map: Record<string, string> = { strong: '强', medium: '中', weak: '弱' }
  return map[strength.value] || ''
})

const strengthColor = computed(() => {
  const map: Record<string, string> = { strong: 'success', medium: 'warning', weak: 'error' }
  return map[strength.value] || 'default'
})

async function generate() {
  try {
    const r = await GeneratePassword(pwdLength.value, useUpper.value, useLower.value, useDigits.value, useSpecial.value)
    if (r) {
      generatedPassword.value = r.password
      strength.value = r.strength
    }
  } catch (e) { console.error(e) }
}

async function copyPassword() {
  if (!generatedPassword.value) return
  try {
    await navigator.clipboard.writeText(generatedPassword.value)
    message.success('已复制到剪贴板')
  } catch {
    message.warning('复制失败，请手动复制')
  }
}

generate()
</script>
