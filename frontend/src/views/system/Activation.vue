<template>
  <div>
    <n-h2>Windows 激活信息</n-h2>
    <n-p>查看当前系统激活状态和可用的激活方式。</n-p>
    <n-button type="primary" @click="loadInfo" :loading="loading" class="mb-4">获取激活信息</n-button>

    <n-card v-if="info" size="small" class="mb-4">
      <n-description-list label-placement="left" :column="1">
        <n-description-item label="激活状态">
          <n-tag :type="info.activationStatus === '已激活' ? 'success' : 'error'">{{ info.activationStatus }}</n-tag>
        </n-description-item>
        <n-description-item label="系统版本">{{ info.edition || '未知' }}</n-description-item>
      </n-description-list>
    </n-card>

    <n-card title="开源激活工具" size="small" class="mb-4">
      <n-list>
        <n-list-item v-for="tool in tools" :key="tool.name">
          <template #prefix>
            <n-tag :type="tool.type === 'hwid' ? 'success' : 'warning'" size="tiny">{{ tool.type === 'hwid' ? 'HWID' : tool.type === 'kms' ? 'KMS' : tool.type }}</n-tag>
          </template>
          <n-space vertical>
            <strong>{{ tool.name }}</strong>
            <n-text depth="3">{{ tool.description }}</n-text>
            <n-button size="tiny" @click="openURL(tool.url)">访问项目</n-button>
          </n-space>
        </n-list-item>
      </n-list>
    </n-card>

    <n-card title="KMS 激活方法" size="small" class="mb-4">
      <n-space vertical>
        <n-empty v-if="!kmsMethods.length" description="点击获取激活信息后显示" />
        <n-card v-for="(method, i) in kmsMethods" :key="i" size="small" :title="method.title">
          <n-space vertical>
            <n-text v-for="(step, j) in method.steps" :key="j">{{ j+1 }}. {{ step }}</n-text>
            <n-input v-if="method.command" :value="method.command" type="textarea" :rows="3" readonly style="font-family:monospace;font-size:12px" />
          </n-space>
        </n-card>
      </n-space>
    </n-card>

    <n-card title="激活方式说明" size="small" class="mb-4">
      <n-table :single-line="false" size="small">
        <thead><tr><th>方式</th><th>说明</th><th>有效期</th><th>适用版本</th></tr></thead>
        <tbody>
          <tr><td>HWID（数字许可证）</td><td>绑定硬件，重装自动激活</td><td>永久</td><td>Win10/Win11</td></tr>
          <tr><td>KMS38</td><td>本地 KMS 模拟，无需外部服务器</td><td>至 2038 年</td><td>Win10/Win11/Server</td></tr>
          <tr><td>KMS</td><td>连接 KMS 服务器激活</td><td>180 天（可续）</td><td>VL 版本/Office</td></tr>
          <tr><td>Ohook</td><td>劫持 Office 激活检验</td><td>永久</td><td>Office 2016-2024</td></tr>
        </tbody>
      </n-table>
    </n-card>

    <n-card title="免责声明" size="small">
      <n-alert type="warning" :bordered="false">
        以上信息仅用于技术研究和学习。请尊重软件版权，在合法授权范围内使用。
        作者不对使用这些信息产生的任何后果负责。
      </n-alert>
    </n-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { GetActivationInfo, GetActivationTools, GetKMSMethods, OpenExternalURL } from '@wails/go/main/App'
const info = ref<any>(null); const tools = ref<any[]>([]); const kmsMethods = ref<any[]>([]); const loading = ref(false); const message = useMessage()
async function loadInfo() {
  loading.value = true
  try {
    const r = await GetActivationInfo(); if (r) info.value = r
    const t = await GetActivationTools(); if (t) tools.value = t as any[]
    const k = await GetKMSMethods(); if (k) kmsMethods.value = k as any[]
  } catch(e:any) { message.error(String(e)) }; loading.value = false
}
async function openURL(url: string) { try { await OpenExternalURL(url) } catch { window.open(url, '_blank') } }
</script>
<style scoped>.mb-4{margin-bottom:16px}</style>
