<template>
  <div class="auto-update">
    <n-space vertical>
      <n-space align="center">
        <n-button type="primary" @click="handleCheckUpdate" :loading="checking">
          检查更新
        </n-button>
        <n-tag v-if="currentVersion" type="info">{{ currentVersion }}</n-tag>
      </n-space>

      <n-card v-if="updateInfo" title="更新信息" size="small">
        <n-descriptions :column="1">
          <n-descriptions-item label="最新版本">{{ updateInfo.version }}</n-descriptions-item>
          <n-descriptions-item label="发布日期">{{ updateInfo.releaseDate }}</n-descriptions-item>
        </n-descriptions>

        <n-divider />
        <p style="font-weight: 600">更新说明：</p>
        <pre class="release-notes">{{ updateInfo.releaseNotes }}</pre>

        <n-space v-if="updateInfo.hasUpdate" style="margin-top: 12px">
          <n-button type="primary" @click="handleUpdate">立即更新</n-button>
        </n-space>
      </n-card>

      <n-empty v-else-if="!checking && checked" description="当前已是最新版本" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const checking = ref(false)
const checked = ref(false)
const currentVersion = ref('')
const updateInfo = ref<any>(null)

async function handleCheckUpdate() {
  checking.value = true
  checked.value = false
  try {
    const [ver, info] = await Promise.all([
      window.go.main.App.GetCurrentVersion(),
      window.go.main.App.CheckUpdate()
    ])
    currentVersion.value = ver || ''
    updateInfo.value = info || null
    checked.value = true
  } catch (e: any) {
    updateInfo.value = null
    checked.value = true
  } finally {
    checking.value = false
  }
}

async function handleUpdate() {
  try {
    await window.go.main.App.DoUpdate()
    window.$message?.success('更新已启动')
  } catch (e: any) {
    window.$message?.error('更新失败: ' + (e.message || '未知错误'))
  }
}
</script>

<style scoped>
.auto-update {
  padding: 20px;
}
.release-notes {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  font-size: 13px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
