<template>
  <div>
    <h2>防火墙规则管理</h2>
    <p>查看和管理 Windows 防火墙规则。</p>
    <n-space class="mb-4">
      <n-select v-model:value="direction" :options="directionOptions" style="width: 120px" @update:value="loadRules" />
      <n-button type="primary" @click="loadRules" :loading="loading">刷新</n-button>
      <n-button @click="showBlockPort = true">阻止端口</n-button>
      <n-button @click="showBlockProgram = true">阻止程序</n-button>
    </n-space>
    <n-space class="mb-4" v-if="firewallStatus">
      <n-tag v-for="(enabled, name) in firewallStatus" :key="name" :type="enabled ? 'success' : 'error'" size="large">
        {{ name }}: {{ enabled ? '已启用' : '已关闭' }}
      </n-tag>
    </n-space>
    <n-data-table :columns="columns" :data="rules" :bordered="true" :loading="loading" size="small" :pagination="{ pageSize: 15 }" />
    <n-modal v-model:show="showBlockPort" preset="dialog" title="阻止端口">
      <n-space vertical>
        <n-input-number v-model:value="blockPort" :min="1" :max="65535" placeholder="端口号" />
        <n-select v-model:value="blockProtocol" :options="protocolOptions" />
        <n-button type="primary" @click="doBlockPort">确认阻止</n-button>
      </n-space>
    </n-modal>
    <n-modal v-model:show="showBlockProgram" preset="dialog" title="阻止程序联网">
      <n-space vertical>
        <n-input v-model:value="blockProgramPath" placeholder="程序完整路径" />
        <n-button type="primary" @click="doBlockProgram">确认阻止</n-button>
      </n-space>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { NTag, NButton, useMessage } from 'naive-ui'
import { GetFirewallRules, ToggleFirewallRule, BlockPort, BlockProgram, GetFirewallStatus } from '@wails/go/main/App'

interface FirewallRule { name: string; displayName: string; direction: string; action: string; protocol: string; localPort: string; remotePort: string; program: string; enabled: boolean }

const rules = ref<FirewallRule[]>([])
const direction = ref('Inbound')
const loading = ref(false)
const firewallStatus = ref<Record<string, boolean> | null>(null)
const showBlockPort = ref(false)
const showBlockProgram = ref(false)
const blockPort = ref(8080)
const blockProtocol = ref('TCP')
const blockProgramPath = ref('')
const message = useMessage()

const directionOptions = [
  { label: '入站规则', value: 'Inbound' },
  { label: '出站规则', value: 'Outbound' },
]
const protocolOptions = [
  { label: 'TCP', value: 'TCP' },
  { label: 'UDP', value: 'UDP' },
  { label: 'TCP+UDP', value: 'TCP,UDP' },
]

const columns = [
  { title: '名称', key: 'displayName', width: 200, ellipsis: { tooltip: true } },
  { title: '操作', key: 'action', width: 80, render: (row: FirewallRule) => h(NTag, { type: row.action === 'Allow' ? 'success' : 'error', size: 'small' }, { default: () => row.action }) },
  { title: '协议', key: 'protocol', width: 80 },
  { title: '本地端口', key: 'localPort', width: 100 },
  { title: '程序', key: 'program', width: 150, ellipsis: { tooltip: true } },
  { title: '状态', key: 'enabled', width: 80, render: (row: FirewallRule) => h(NTag, { type: row.enabled ? 'success' : 'default', size: 'small' }, { default: () => row.enabled ? '启用' : '禁用' }) },
  { title: '操作', width: 100, render: (row: FirewallRule) => h(NButton, { size: 'small', quaternary: true, type: row.enabled ? 'warning' : 'success', onClick: () => toggleRule(row) }, { default: () => row.enabled ? '禁用' : '启用' }) },
]

async function loadRules() {
  loading.value = true
  try {
    rules.value = (await GetFirewallRules(direction.value) as FirewallRule[]) || []
    firewallStatus.value = (await GetFirewallStatus() as Record<string, boolean>) || {}
  } catch (e: any) { message.error(String(e)) }
  loading.value = false
}

async function toggleRule(rule: FirewallRule) {
  try {
    await ToggleFirewallRule(rule.name, !rule.enabled)
    rule.enabled = !rule.enabled
    message.success('已更新')
  } catch (e: any) { message.error(String(e)) }
}

async function doBlockPort() {
  try {
    await BlockPort(blockPort.value, blockProtocol.value, direction.value)
    showBlockPort.value = false
    message.success('已阻止端口')
    loadRules()
  } catch (e: any) { message.error(String(e)) }
}

async function doBlockProgram() {
  if (!blockProgramPath.value) return
  try {
    await BlockProgram(blockProgramPath.value, direction.value)
    showBlockProgram.value = false
    message.success('已阻止程序')
    loadRules()
  } catch (e: any) { message.error(String(e)) }
}

onMounted(loadRules)
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
