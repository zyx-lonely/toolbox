<template>
  <div>
    <n-h2>右键菜单扩展</n-h2>
    <n-p>在文件资源管理器右键菜单中添加「发送到电脑工具箱」入口。</n-p>
    <n-alert type="warning" :bordered="false" class="mb-4">需要管理员权限才能安装/卸载右键菜单。</n-alert>
    <n-space>
      <n-button type="primary" @click="installFile" :loading="installing">安装文件右键菜单</n-button>
      <n-button type="primary" @click="installDir" :loading="installingDir">安装文件夹右键菜单</n-button>
      <n-button @click="uninstall" :loading="uninstalling">卸载右键菜单</n-button>
      <n-button @click="check" :loading="checking">检查状态</n-button>
    </n-space>
    <n-card v-if="status" size="small" class="mt-4">
      <n-tag :type="status.installed ? 'success' : 'default'">{{ status.installed ? '已安装' : '未安装' }}</n-tag>
      <span v-if="status.path" class="ml-2">{{ status.path }}</span>
      <div v-if="status.error" style="color:#e74c3c;margin-top:8px">{{ status.error }}</div>
    </n-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { InstallContextMenu, InstallContextMenuDir, UninstallContextMenu, CheckContextMenu } from '@wails/go/main/App'
const status = ref<any>(null)
const installing = ref(false)
const installingDir = ref(false)
const uninstalling = ref(false)
const checking = ref(false)
const message = useMessage()
async function installFile() {
  installing.value = true; status.value = await InstallContextMenu(); message.success('安装成功'); installing.value = false
}
async function installDir() {
  installingDir.value = true; status.value = await InstallContextMenuDir(); message.success('安装成功'); installingDir.value = false
}
async function uninstall() {
  uninstalling.value = true; status.value = await UninstallContextMenu(); message.success('已卸载'); uninstalling.value = false
}
async function check() {
  checking.value = true; status.value = await CheckContextMenu(); checking.value = false
}
</script>
<style scoped>.mb-4{margin-bottom:16px}.mt-4{margin-top:16px}.ml-2{margin-left:8px}</style>
