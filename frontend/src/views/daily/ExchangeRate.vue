<template>
  <div>
    <h2>汇率查询</h2>
    <p>实时汇率查询和货币转换。</p>
    <n-card size="small" class="mb-4">
      <n-space align="center">
        <n-input-number v-model:value="amount" :min="0" :precision="2" style="width: 150px" placeholder="金额" />
        <n-select v-model:value="fromCurrency" :options="currencyOptions" style="width: 180px" />
        <n-button @click="swapCurrency">⇄</n-button>
        <n-select v-model:value="toCurrency" :options="currencyOptions" style="width: 180px" />
        <n-button type="primary" @click="convert" :loading="loading">转换</n-button>
      </n-space>
    </n-card>
    <n-card v-if="result" title="转换结果" size="small">
      <n-statistic :label="`${result.amount} ${result.from}`" :value="result.result.toFixed(2)" />
      <p style="margin-top: 8px; color: #999; font-size: 12px;">
        汇率: 1 {{ result.from }} = {{ result.rate.toFixed(4) }} {{ result.to }}
        <br/>更新时间: {{ result.updateTime }}
      </p>
    </n-card>
    <n-card title="常用汇率" size="small" class="mb-4">
      <n-data-table :columns="rateColumns" :data="rateTable" :bordered="true" size="small" :pagination="false" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { NButton, useMessage } from 'naive-ui'
import { ConvertCurrency, GetSupportedCurrencies, GetExchangeRate } from '@wails/go/main/App'

interface ExchangeResult { from: string; to: string; rate: number; amount: number; result: number; updateTime: string }
interface CurrencyInfo { code: string; name: string }
interface RateRow { pair: string; rate: number }

const amount = ref(100)
const fromCurrency = ref('CNY')
const toCurrency = ref('USD')
const result = ref<ExchangeResult | null>(null)
const loading = ref(false)
const currencies = ref<CurrencyInfo[]>([])
const rateTable = ref<RateRow[]>([])
const message = useMessage()

const currencyOptions = ref<{ label: string; value: string }[]>([])

const rateColumns = [
  { title: '货币对', key: 'pair', width: 200 },
  { title: '汇率', key: 'rate', width: 150, render: (row: RateRow) => row.rate.toFixed(4) },
]

async function loadCurrencies() {
  try {
    const list = await GetSupportedCurrencies() as CurrencyInfo[]
    currencies.value = list
    currencyOptions.value = list.map(c => ({ label: c.name, value: c.code }))
  } catch (e) { console.error(e) }
}

async function convert() {
  loading.value = true
  try {
    const r = await ConvertCurrency(amount.value, fromCurrency.value, toCurrency.value) as ExchangeResult
    result.value = r
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

function swapCurrency() {
  const tmp = fromCurrency.value
  fromCurrency.value = toCurrency.value
  toCurrency.value = tmp
}

async function loadRateTable() {
  const pairs = [
    ['USD', 'CNY'], ['EUR', 'CNY'], ['GBP', 'CNY'], ['JPY', 'CNY'],
    ['KRW', 'CNY'], ['HKD', 'CNY'], ['USD', 'EUR'], ['USD', 'JPY'],
  ]
  const rows: RateRow[] = []
  for (const [from, to] of pairs) {
    try {
      const rate = await GetExchangeRate(from, to) as number
      rows.push({ pair: `${from} → ${to}`, rate })
    } catch { /* skip */ }
  }
  rateTable.value = rows
}

onMounted(() => {
  loadCurrencies()
  loadRateTable()
})
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
