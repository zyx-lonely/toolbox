<template>
  <div>
    <n-h2>IP 子网计算器</n-h2>
    <n-p>输入 IP 地址和子网掩码，自动计算网段信息。</n-p>

    <n-space vertical :size="16">
      <n-space>
        <n-input v-model:value="ipAddress" placeholder="IP 地址 (如 192.168.1.100)" style="width: 240px" />
        <n-select v-model:value="cidrValue" :options="cidrOptions" style="width: 140px" placeholder="CIDR" />
        <n-button type="primary" @click="calculate">计算</n-button>
      </n-space>

      <n-grid :cols="3" :x-gap="12" :y-gap="12" v-if="result">
        <n-gi><n-card size="small"><n-statistic label="网络地址" :value="result.network" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="广播地址" :value="result.broadcast" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="子网掩码" :value="result.subnetMask" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="通配符掩码" :value="result.wildcard" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="可用主机数" :value="result.usableHosts" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="总地址数" :value="result.totalHosts" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="地址范围" :value="result.range" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="主机段" :value="result.hostRange" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="CIDR 表示" :value="result.cidr" /></n-card></n-gi>
      </n-grid>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const message = useMessage()
const ipAddress = ref('192.168.1.100')
const cidrValue = ref(24)

const cidrOptions = Array.from({ length: 32 }, (_, i) => ({ label: `/${i + 1}`, value: i + 1 }))

const result = ref<any>(null)

function ipToInt(ip: string): number {
  return ip.split('.').reduce((acc, octet) => (acc << 8) + parseInt(octet), 0) >>> 0
}

function intToIp(int: number): string {
  return [(int >>> 24) & 255, (int >>> 16) & 255, (int >>> 8) & 255, int & 255].join('.')
}

function calculate() {
  const ip = ipAddress.value.trim()
  if (!/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/.test(ip)) {
    message.error('请输入有效的 IP 地址')
    return
  }

  const ipInt = ipToInt(ip)
  const mask = cidrValue.value === 0 ? 0 : (~0 << (32 - cidrValue.value)) >>> 0
  const network = (ipInt & mask) >>> 0
  const broadcast = (network | ~mask) >>> 0
  const wildcard = (~mask) >>> 0
  const totalHosts = Math.pow(2, 32 - cidrValue.value)
  const usableHosts = cidrValue.value >= 31 ? (cidrValue.value === 31 ? 2 : 1) : totalHosts - 2

  result.value = {
    network: intToIp(network),
    broadcast: intToIp(broadcast),
    subnetMask: intToIp(mask),
    wildcard: intToIp(wildcard),
    totalHosts,
    usableHosts,
    range: `${intToIp(network)} - ${intToIp(broadcast)}`,
    hostRange: cidrValue.value >= 31
      ? (cidrValue.value === 31 ? `${intToIp(network)} - ${intToIp(broadcast)}` : intToIp(network))
      : `${intToIp(network + 1)} - ${intToIp(broadcast - 1)}`,
    cidr: `${ip}/${cidrValue.value}`
  }

  message.success('计算完成')
}
</script>
