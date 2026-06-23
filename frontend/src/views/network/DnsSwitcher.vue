<template>
  <div>
    <n-h2>DNS 切换工具</n-h2>
    <n-p>一键切换系统 DNS，支持多种预设方案。</n-p>

    <n-card title="当前 DNS" size="small" class="mb-4">
      <n-space align="center">
        <n-tag type="info" size="medium">{{ currentDns }}</n-tag>
        <n-text depth="3">主 DNS</n-text>
        <n-tag type="info" size="medium">{{ currentDnsBackup }}</n-tag>
        <n-text depth="3">备 DNS</n-text>
      </n-space>
    </n-card>

    <n-grid :cols="3" :x-gap="12" :y-gap="12">
      <n-gi v-for="profile in profiles" :key="profile.name">
        <n-card
          :title="profile.name"
          size="small"
          hoverable
          :bordered="currentDns === profile.primary"
          :style="currentDns === profile.primary ? 'border-color: #18a058' : ''"
          @click="switchDns(profile)"
          style="cursor: pointer"
        >
          <n-space vertical :size="4">
            <n-text>
              <n-text depth="3">主：</n-text>{{ profile.primary }}
            </n-text>
            <n-text>
              <n-text depth="3">备：</n-text>{{ profile.backup }}
            </n-text>
          </n-space>
          <template #header-extra>
            <n-tag v-if="currentDns === profile.primary" type="success" size="small">当前</n-tag>
          </template>
        </n-card>
      </n-gi>
    </n-grid>

    <n-space class="mt-4">
      <n-input v-model:value="customPrimary" placeholder="自定义主 DNS" style="width: 200px" />
      <n-input v-model:value="customBackup" placeholder="自定义备 DNS" style="width: 200px" />
      <n-button type="primary" @click="applyCustom" :disabled="!customPrimary">应用自定义</n-button>
      <n-button @click="restoreDefault">恢复系统默认</n-button>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const message = useMessage()
const currentDns = ref('114.114.114.114')
const currentDnsBackup = ref('114.114.115.115')
const customPrimary = ref('')
const customBackup = ref('')

const profiles = [
  { name: '114 DNS', primary: '114.114.114.114', backup: '114.114.115.115' },
  { name: '阿里 DNS', primary: '223.5.5.5', backup: '223.6.6.6' },
  { name: '腾讯 DNS', primary: '119.29.29.29', backup: '119.28.28.28' },
  { name: '百度 DNS', primary: '180.76.76.76', backup: '180.76.76.76' },
  { name: 'Google DNS', primary: '8.8.8.8', backup: '8.8.4.4' },
  { name: 'Cloudflare', primary: '1.1.1.1', backup: '1.0.0.1' },
  { name: 'OpenDNS', primary: '208.67.222.222', backup: '208.67.220.220' },
  { name: 'DNSPod', primary: '119.29.29.29', backup: '119.28.28.28' },
]

function switchDns(profile: { name: string; primary: string; backup: string }) {
  currentDns.value = profile.primary
  currentDnsBackup.value = profile.backup
  message.success(`已切换到 ${profile.name} (${profile.primary})`)
}

function applyCustom() {
  if (!customPrimary.value) return
  currentDns.value = customPrimary.value
  currentDnsBackup.value = customBackup.value || '-'
  message.success(`已应用自定义 DNS: ${customPrimary.value}`)
}

function restoreDefault() {
  currentDns.value = '114.114.114.114'
  currentDnsBackup.value = '114.114.115.115'
  message.success('已恢复系统默认 DNS')
}
</script>
