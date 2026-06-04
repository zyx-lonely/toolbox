<template>
  <div>
    <n-h2>文件解锁器</n-h2>
    <n-p>检测并解除被其他程序占用的文件。</n-p>

    <n-space class="mb-4">
      <n-input-group>
        <n-input v-model:value="filePath" placeholder="选择或输入文件路径" readonly style="width: 350px" />
        <n-button @click="selectFile">选择文件</n-button>
      </n-input-group>
      <n-button type="primary" @click="tryUnlock" :loading="unlocking" :disabled="!filePath">
        尝试解锁
      </n-button>
    </n-space>

    <n-empty v-if="!filePath" description="选择一个被占用的文件" />

    <n-card v-if="result" size="small">
      <template #header>
        <n-tag :type="result.success ? 'success' : 'error'" size="medium">
          {{ result.success ? '✅ 解锁成功' : '❌ 解锁失败' }}
        </n-tag>
      </template>
      <n-description-list label-placement="left" :column="1">
        <n-description-item label="文件">{{ result.filePath }}</n-description-item>
        <n-description-item v-if="result.success" label="方式">{{ result.releasedBy?.join(', ') }}</n-description-item>
        <n-description-item v-if="!result.success" label="原因">{{ result.error }}</n-description-item>
      </n-description-list>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { SelectFile } from '@wails/go/main/App'

const filePath = ref('')
const unlocking = ref(false)
const result = ref<any>(null)
const message = useMessage()

async function selectFile() {
  try {
    const path = await SelectFile()
    if (path) filePath.value = path as string
  } catch (e) { console.error(e) }
}

async function tryUnlock() {
  if (!filePath.value) return
  unlocking.value = true
  result.value = null
  try {
    // 动态导入 TryUnlock
    const { TryUnlock } = await import('@wails/go/main/App')
    const r = await TryUnlock(filePath.value)
    if (r) {
      result.value = r
      if (r.success) message.success('文件已解锁')
      else message.error('解锁失败: ' + r.error)
    }
  } catch (e: any) {
    message.error(`操作失败: ${e}`)
  }
  unlocking.value = false
}
</script>

<style scoped>
.mb-4 { margin-bottom: 12px; }
</style>
