<template>
  <div>
    <n-h2>重复文件查找</n-h2>
    <n-p>扫描目录中的重复文件，支持快速（大小+时间）和精确（MD5 哈希）两种模式。</n-p>

    <n-space>
      <n-input-group>
        <n-input v-model:value="scanPath" placeholder="选择要扫描的目录" readonly style="width: 300px" />
        <n-button @click="selectDir">选择目录</n-button>
      </n-input-group>
      <n-select v-model:value="scanMode" :options="modeOptions" style="width: 120px" />
      <n-button type="primary" @click="doScan" :loading="scanning">开始扫描</n-button>
    </n-space>

    <n-space vertical size="medium" v-if="groups.length" class="mt-4">
      <n-alert type="info" :bordered="false">
        <template #header>扫描结果</template>
        发现 {{ groups.length }} 组重复文件，可释放 {{ formatBytes(totalRedundant) }}
      </n-alert>

      <n-card v-for="(group, i) in groups" :key="i" size="small">
        <template #header>
          <n-space>
            <span>第 {{ i + 1 }} 组 — {{ group.fileCount }} 个文件</span>
            <n-tag size="small">{{ formatBytes(group.totalSize) }}</n-tag>
          </n-space>
        </template>
        <n-space vertical>
          <n-text v-for="f in group.files" :key="f.path" depth="2">
            {{ f.path }} <n-tag size="tiny" type="info">{{ f.matchType }}</n-tag>
          </n-text>
        </n-space>
      </n-card>
    </n-space>

    <n-empty v-if="!scanning && !groups.length && scanPath" description="未找到重复文件" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { SelectDirectory, FindDuplicateFiles } from '@wails/go/main/App'
import { formatBytes } from '@/api/bridge'

const scanPath = ref('')
const scanMode = ref('quick')
const groups = ref<any[]>([])
const scanning = ref(false)

const modeOptions = [
  { label: '快速模式', value: 'quick' },
  { label: '精确模式 (MD5)', value: 'exact' }
]

const totalRedundant = computed(() => groups.value.reduce((s: number, g: any) => s + (g.totalSize || 0), 0))

async function selectDir() {
  try {
    const dir = await SelectDirectory()
    if (dir) scanPath.value = dir as string
  } catch (e) { console.error(e) }
}

async function doScan() {
  if (!scanPath.value) return
  scanning.value = true
  groups.value = []
  try {
    const r = await FindDuplicateFiles(scanPath.value, scanMode.value)
    if (r) groups.value = r as any[]
  } catch (e) { console.error(e) }
  scanning.value = false
}
</script>

<style scoped>
.mt-4 { margin-top: 16px; }
</style>
