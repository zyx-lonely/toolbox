<template>
  <div>
    <n-h2>回收站管理</n-h2>
    <n-p>查看和清空回收站。</n-p>
    <n-card>
      <n-space vertical>
        <n-space>
          <n-statistic title="回收站项目数" :value="binInfo.itemCount" />
          <n-statistic title="回收站大小" :value="binInfo.size" />
        </n-space>
        <n-space>
          <n-button @click="loadInfo" :loading="loading">刷新信息</n-button>
          <n-button type="warning" @click="emptyBin" :loading="emptying" :disabled="binInfo.itemCount === 0">
            <template #icon><n-icon><trash-outline /></n-icon></template>
            清空回收站
          </n-button>
        </n-space>
        <n-alert v-if="lastResult" :type="lastResult.success ? 'success' : 'error'" closable @close="lastResult = null">
          {{ lastResult.message }}
        </n-alert>
      </n-space>
    </n-card>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { TrashOutline } from '@vicons/ionicons5'
import { GetRecycleBinInfo, EmptyRecycleBin } from '@wails/go/main/App'
const binInfo = ref({ itemCount: 0, size: '未知', path: '' })
const loading = ref(false); const emptying = ref(false); const lastResult = ref<any>(null)
const message = useMessage(); const dialog = useDialog()

onMounted(() => loadInfo())

async function loadInfo() {
  loading.value = true
  try { const r = await GetRecycleBinInfo(); if (r) binInfo.value = r as any }
  catch(e:any) { message.error(String(e)) }; loading.value = false
}
function emptyBin() {
  dialog.warning({
    title: '确认清空',
    content: '确定要清空回收站吗？删除的文件将无法恢复。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      emptying.value = true
      try {
        await EmptyRecycleBin()
        lastResult.value = { success: true, message: '回收站已清空' }
        message.success('回收站已清空')
        await loadInfo()
      } catch(e:any) {
        lastResult.value = { success: false, message: String(e) }
        message.error(String(e))
      }; emptying.value = false
    }
  })
}
</script>
