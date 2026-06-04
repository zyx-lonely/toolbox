<template>
  <div>
    <n-h2>网络工具</n-h2>

    <n-tabs type="line" default-value="ping">
      <n-tab-pane name="ping" tab="Ping">
        <n-space vertical>
          <n-input-group>
            <n-input v-model:value="pingHost" placeholder="输入主机名或 IP" />
            <n-input-number v-model:value="pingCount" :min="1" :max="100" style="width: 80px" />
            <n-button type="primary" @click="doPing" :loading="pinging">Ping</n-button>
          </n-input-group>

          <n-empty v-if="!pinging && !pingSummary" description="输入目标后开始 Ping" />

          <n-data-table v-if="pingSummary?.results?.length" :columns="pingColumns" :data="pingSummary.results" size="small" :bordered="true" />

          <n-card v-if="pingSummary" size="small" class="mt-2">
            <n-space justify="space-around">
              <n-statistic label="已发送" :value="pingSummary.sent || 0" />
              <n-statistic label="已接收" :value="pingSummary.received || 0" />
              <n-statistic label="丢包率">
                <span :style="{ color: lossColor }">{{ pingSummary.lossPercent || 0 }}%</span>
              </n-statistic>
              <n-statistic label="最小时延" :value="pingSummary.minLatency || '-'" />
              <n-statistic label="最大时延" :value="pingSummary.maxLatency || '-'" />
              <n-statistic label="平均时延" :value="pingSummary.avgLatency || '-'" />
            </n-space>
          </n-card>
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="port" tab="端口扫描">
        <n-space vertical>
          <n-input-group>
            <n-input v-model:value="portHost" placeholder="输入主机名或 IP" style="width: 200px" />
            <n-input v-model:value="portRange" placeholder="端口范围 (如 1-1000 或 common)" style="width: 200px" />
            <n-button type="primary" @click="doPortScan" :loading="scanning">扫描</n-button>
          </n-input-group>

          <n-empty v-if="!scanning && portResults.length === 0" description="输入目标后开始扫描" />

          <n-data-table v-if="portResults.length" :columns="portColumns" :data="portResults" size="small" :bordered="true" />
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="dns" tab="DNS 查询">
        <n-space vertical>
          <n-input-group>
            <n-input v-model:value="dnsHost" placeholder="输入域名" />
            <n-button type="primary" @click="doDNS" :loading="dnsLoading">查询</n-button>
          </n-input-group>

          <n-card v-if="dnsResult" size="small">
            <n-description-list label-placement="left" :column="1">
              <n-description-item label="域名">{{ dnsResult.hostname }}</n-description-item>
              <n-description-item label="IP 地址">
                <div v-for="ip in dnsResult.answers" :key="ip" style="font-family: monospace;">{{ ip }}</div>
              </n-description-item>
            </n-description-list>
          </n-card>
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="fix" tab="网络修复">
        <n-space vertical>
          <n-alert type="warning" :bordered="false">
            <template #header>⚠️ 需要管理员权限</template>
            部分修复操作（如重置 Winsock、TCP/IP）需要以管理员身份运行本程序才能生效。
          </n-alert>

          <n-button type="primary" size="large" @click="fixAll" :loading="fixing" class="mb-2">
            🔧 一键网络修复
          </n-button>

          <n-grid :cols="2" :x-gap="12" :y-gap="12">
            <n-gi v-for="fix in fixActions" :key="fix.name">
              <n-button quaternary block @click="runFix(fix.name)" :loading="fix.loading">
                {{ fix.label }}
              </n-button>
            </n-gi>
          </n-grid>

          <n-space vertical size="small" v-if="fixResults.length">
            <n-h3>修复结果</n-h3>
            <n-alert v-for="(r, i) in fixResults" :key="i"
              :type="r.success ? 'success' : 'warning'" :bordered="false" closable>
              <template #header>{{ r.action }}</template>
              {{ r.success ? '✅ 执行成功' : '❌ ' + (r.error || '执行失败') }}
              <div v-if="r.output" style="font-family: monospace; font-size: 12px; margin-top: 4px; white-space: pre-wrap;">{{ r.output }}</div>
            </n-alert>
          </n-space>
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, h, reactive } from 'vue'
import { NTag, useMessage } from 'naive-ui'
import { Ping, PortScan, DNSLookup } from '@wails/go/main/App'
import { FlushDNSCache, ResetWinsock, ResetTCPIP, DiagnoseNetwork, FixAllNetwork, ReleaseIP, RenewIP, ResetProxy, ResetArpCache } from '@wails/go/main/App'

const pingHost = ref('localhost')
const pingCount = ref(4)
const pingSummary = ref<any>(null)
const pinging = ref(false)

const portHost = ref('localhost')
const portRange = ref('common')
const portResults = ref<any[]>([])
const scanning = ref(false)

const dnsHost = ref('')
const dnsResult = ref<any>(null)
const dnsLoading = ref(false)

const fixResults = ref<any[]>([])
const fixing = ref(false)
const message = useMessage()

const fixActions = reactive([
  { name: 'flushdns', label: '刷新 DNS 缓存', loading: false },
  { name: 'winsock', label: '重置 Winsock', loading: false },
  { name: 'tcpip', label: '重置 TCP/IP', loading: false },
  { name: 'arp', label: '清空 ARP 缓存', loading: false },
  { name: 'release', label: '释放 IP', loading: false },
  { name: 'proxy', label: '清除代理设置', loading: false },
])

const pingColumns = [
  { title: '序号', key: 'sequence', width: 60 },
  { title: '主机', key: 'target' },
  {
    title: '状态', key: 'success', width: 80,
    render: (row: any) => h(NTag, { type: row.success ? 'success' as const : 'error' as const, size: 'small' }, { default: () => row.success ? '成功' : '失败' })
  },
  { title: '延迟', key: 'latency', width: 120 },
  { title: 'TTL', key: 'ttl', width: 60 },
]

const portColumns = [
  { title: '端口', key: 'port', width: 80 },
  {
    title: '状态', key: 'state', width: 80,
    render: (row: any) => h(NTag, { type: row.state === 'open' ? 'error' as const : 'default' as const, size: 'small' }, { default: () => row.state === 'open' ? '开放' : '关闭' })
  },
  { title: '服务', key: 'service' },
]

const lossColor = (pingSummary.value?.lossPercent || 0) > 10 ? '#e74c3c' : '#27ae60'

async function doPing() {
  pinging.value = true
  try {
    const r = await Ping(pingHost.value, pingCount.value, 2000)
    if (r) pingSummary.value = r
  } catch (e) { console.error(e) }
  pinging.value = false
}

async function doPortScan() {
  scanning.value = true
  try {
    const r = await PortScan(portHost.value, portRange.value)
    if (r) portResults.value = r as any[]
  } catch (e) { console.error(e) }
  scanning.value = false
}

async function doDNS() {
  dnsLoading.value = true
  try {
    const r = await DNSLookup(dnsHost.value)
    if (r) dnsResult.value = r
  } catch (e) { console.error(e) }
  dnsLoading.value = false
}

async function runFix(name: string) {
  const action = fixActions.find(a => a.name === name)
  if (!action) return
  action.loading = true
  try {
    let r: any
    switch (name) {
      case 'flushdns': r = await FlushDNSCache(); break
      case 'winsock': r = await ResetWinsock(); break
      case 'tcpip': r = await ResetTCPIP(); break
      case 'arp': r = await ResetArpCache(); break
      case 'release': r = await ReleaseIP(); break
      case 'proxy': r = await ResetProxy(); break
    }
    if (r) {
      fixResults.value.unshift(r)
      if (r.success) message.success(`${action.label} 成功`)
      else message.warning(`${action.label} 失败: ${r.error}`)
    }
  } catch (e: any) {
    message.error(`${action.label} 出错: ${e}`)
  }
  action.loading = false
}

async function fixAll() {
  fixing.value = true
  try {
    const results = await FixAllNetwork()
    if (results) {
      fixResults.value = [...results, ...fixResults.value]
      const success = results.filter((r: any) => r.success).length
      message.success(`修复完成 (${success}/${results.length})`)
    }
  } catch (e: any) {
    message.error(`一键修复失败: ${e}`)
  }
  fixing.value = false
}
</script>

<style scoped>
.mt-2 { margin-top: 8px; }
.mb-2 { margin-bottom: 12px; }
</style>
