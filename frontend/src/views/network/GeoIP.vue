<template>
  <div>
    <n-h2>IP 地理位置查询</n-h2>
    <n-p>查询 IP 地址所属地区、运营商等信息。</n-p>
    <n-space class="mb-4">
      <n-input-group>
        <n-input v-model:value="ipAddr" placeholder="输入 IP 地址" style="width: 250px" />
        <n-button type="primary" @click="query" :loading="querying">查询</n-button>
      </n-input-group>
    </n-space>
    <n-empty v-if="!result && !querying" description="输入 IP 地址后点击查询" />
    <n-card v-if="result" size="small">
      <n-description-list label-placement="left" :column="1">
        <n-description-item label="IP">{{ result.ip }}</n-description-item>
        <n-description-item label="国家/地区">{{ result.country }}</n-description-item>
        <n-description-item label="区域">{{ result.region }}</n-description-item>
        <n-description-item label="城市">{{ result.city }}</n-description-item>
        <n-description-item label="ISP">{{ result.isp }}</n-description-item>
        <n-description-item label="经纬度">{{ result.lat }}, {{ result.lon }}</n-description-item>
      </n-description-list>
    </n-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { QueryGeoIP } from '@wails/go/main/App'
const ipAddr = ref('8.8.8.8')
const result = ref<any>(null)
const querying = ref(false)
const message = useMessage()
async function query() {
  if (!ipAddr.value) return
  querying.value = true
  try {
    const r = await QueryGeoIP(ipAddr.value)
    if (r) result.value = r
    if (r?.error) message.error(r.error)
  } catch (e: any) { message.error(String(e)) }
  querying.value = false
}
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
