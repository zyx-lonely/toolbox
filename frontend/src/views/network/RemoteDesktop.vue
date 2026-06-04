<template>
  <div>
    <n-h2>远程桌面</n-h2>
    <n-p>通过局域网远程连接其他电脑。</n-p>
    <n-card title="新建连接">
      <n-space vertical>
        <n-input v-model:value="computer" placeholder="计算机名 (可选)" />
        <n-input v-model:value="address" placeholder="IP 地址 *" />
        <n-input-number v-model:value="port" :min="1" :max="65535" placeholder="端口 (默认 3389)" style="width: 150px" />
        <n-button type="primary" @click="connect" :loading="connecting">连接</n-button>
      </n-space>
    </n-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { LaunchMSTSC } from '@wails/go/main/App'
const computer = ref(''); const address = ref(''); const port = ref(3389); const connecting = ref(false); const message = useMessage()
async function connect() {
  if (!address.value) { message.warning('请输入 IP 地址'); return }
  connecting.value = true
  try { await LaunchMSTSC(computer.value, address.value, port.value); message.success('正在启动远程桌面连接') }
  catch(e:any) { message.error(String(e)) }
  connecting.value = false
}
</script>
