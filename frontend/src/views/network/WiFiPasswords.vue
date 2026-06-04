<template>
  <div>
    <n-h2>WiFi 密码查看</n-h2>
    <n-p>显示已保存的 WiFi 密码。</n-p>

    <n-alert type="warning" :bordered="false" class="mb-4">
      <template #header>⚠️ 需要管理员权限</template>
      查看 WiFi 密码需要以管理员身份运行本程序。
    </n-alert>

    <n-button type="primary" @click="loadPasswords" :loading="loading" class="mb-4">获取 WiFi 密码</n-button>
    <n-empty v-if="!loading && !profiles.length" description="点击按钮获取 WiFi 密码" />
    <n-data-table v-if="profiles.length" :columns="columns" :data="profiles" size="small" :bordered="true" />
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { GetWiFiPasswords } from '@wails/go/main/App'
const profiles = ref<any[]>([]); const loading = ref(false); const message = useMessage()
const columns = [
  {title:'SSID',key:'ssid',width:200},
  {title:'密码',key:'password',ellipsis:{tooltip:true}},
  {title:'认证',key:'auth',width:80,render:(r:any)=>r.auth},
]
async function loadPasswords() {
  loading.value = true
  try { const r = await GetWiFiPasswords(); if(r) profiles.value = r as any[]; message.success(`发现 ${profiles.value.length} 个 WiFi 配置`) }
  catch(e:any) { message.error(String(e)) }
  loading.value = false
}
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
