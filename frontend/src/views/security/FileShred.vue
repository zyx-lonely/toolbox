<template>
  <div>
    <n-h2>文件粉碎</n-h2>
    <n-p>安全删除文件或目录中的文件，覆写后不可恢复。</n-p>

    <n-card title="选择目标">
      <n-space vertical>
        <n-tabs v-model:value="mode" type="line">
          <n-tab name="file" tab="单个文件" />
          <n-tab name="dir" tab="整个目录" />
        </n-tabs>

        <n-space v-if="mode === 'file'">
          <n-input-group>
            <n-input v-model:value="targetPath" placeholder="选择要粉碎的文件" readonly style="width: 400px" />
            <n-button @click="selectFile">选择文件</n-button>
          </n-input-group>
        </n-space>

        <n-space v-else>
          <n-input-group>
            <n-input v-model:value="targetPath" placeholder="选择要粉碎的目录" readonly style="width: 400px" />
            <n-button @click="selectDir">选择目录</n-button>
          </n-input-group>
        </n-space>

        <n-space align="center">
          <span>覆写次数:</span>
          <n-input-number v-model:value="passes" :min="1" :max="7" style="width: 80px" />
          <n-tag type="info" size="small">{{ passLabel }}</n-tag>
        </n-space>

        <n-alert type="warning" :show-icon="false">
          <template #header>警告</template>
          粉碎后的文件<strong>无法恢复</strong>！请确认文件不再需要。
        </n-alert>

        <n-button type="error" @click="shred" :loading="shredding" :disabled="!targetPath">
          <template #icon><n-icon><trash-outline /></n-icon></template>
          开始粉碎
        </n-button>
      </n-space>
    </n-card>

    <n-card v-if="results.length > 0" title="粉碎结果" style="margin-top: 16px">
      <n-list>
        <n-list-item v-for="(r, i) in results" :key="i">
          <n-thing :title="r.path">
            <template #description>
              <n-tag v-if="r.success" type="success" size="small">成功 ({{ r.passes }}遍)</n-tag>
              <n-tag v-else type="error" size="small">失败: {{ r.error }}</n-tag>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>
      <n-statistic v-if="successCount > 0" label="成功粉碎" :value="successCount" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { TrashOutline } from '@vicons/ionicons5'
import { ShredFile, ShredDir, SelectFile, SelectDirectory } from '@wails/go/main/App'

const mode = ref<'file' | 'dir'>('file')
const targetPath = ref('')
const passes = ref(3)
const shredding = ref(false)
const results = ref<any[]>([])
const message = useMessage()

const passLabel = computed(() => {
  const labels: Record<number, string> = { 1: '1遍', 3: '3遍 (DoD标准)', 7: '7遍 (德国标准)' }
  return labels[passes.value] || `${passes.value}遍`
})

const successCount = computed(() => results.value.filter(r => r.success).length)

async function selectFile() {
  const f = await SelectFile()
  if (f) targetPath.value = f as string
}

async function selectDir() {
  const d = await SelectDirectory()
  if (d) targetPath.value = d as string
}

async function shred() {
  if (!targetPath.value) return
  shredding.value = true
  results.value = []
  try {
    if (mode.value === 'file') {
      const r = await ShredFile(targetPath.value, passes.value)
      if (r) results.value = [r as any]
    } else {
      const rs = await ShredDir(targetPath.value, passes.value)
      if (rs) results.value = rs as any[]
    }
    message.success('粉碎完成')
  } catch (e: any) {
    message.error(String(e))
  }
  shredding.value = false
}
</script>
