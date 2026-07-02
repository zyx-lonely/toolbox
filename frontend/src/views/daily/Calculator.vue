<template>
  <div>
    <h2>科学计算器</h2>
    <p>支持基本运算和科学函数。</p>
    <n-card size="small" style="max-width: 400px">
      <n-input v-model:value="expression" placeholder="输入表达式，如 sin(30) + 2^8" :rows="2" />
      <n-space class="mt-4">
        <n-button type="primary" @click="calc" :loading="loading">计算</n-button>
        <n-button @click="clear">清空</n-button>
      </n-space>
      <n-alert v-if="result" :type="result.error ? 'error' : 'success'" class="mt-4">
        <template v-if="result.error">
          错误: {{ result.error }}
        </template>
        <template v-else>
          <strong>{{ result.expression }} = {{ result.result }}</strong>
        </template>
      </n-alert>
      <n-divider />
      <n-text depth="3" style="font-size: 12px">
        函数: sin, cos, tan, asin, acos, atan, sqrt, cbrt, abs, log, ln, log2, exp, floor, ceil, round<br/>
        常量: pi, e<br/>
        运算: + - * / % ^(幂) (括号)
      </n-text>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { Calculate } from '@wails/go/main/App'

interface CalcResult { expression: string; result: number; error?: string }

const expression = ref('')
const result = ref<CalcResult | null>(null)
const loading = ref(false)
const message = useMessage()

async function calc() {
  if (!expression.value) return
  loading.value = true
  try {
    result.value = await Calculate(expression.value) as CalcResult
    if (result.value.error) {
      message.error(result.value.error)
    }
  } catch (e: any) { message.error(String(e)) }
  loading.value = false
}

function clear() {
  expression.value = ''
  result.value = null
}
</script>

<style scoped>
.mt-4 { margin-top: 16px; }
</style>
