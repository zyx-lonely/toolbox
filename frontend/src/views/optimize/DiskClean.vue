<template>
  <div>
    <n-h2>磁盘清理</n-h2>
    <n-p>扫描并清理系统临时文件和缓存，释放磁盘空间。</n-p>

    <n-alert type="warning" :bordered="false" class="mb-4">
      <template #header>⚠️ 需要管理员权限</template>
      部分系统目录（如 WinSxS、SoftwareDistribution）需要管理员权限才能清理。
    </n-alert>

    <n-alert type="info" :bordered="false" class="mb-4">
      <template #header>提示</template>
      低风险项默认勾选，中/高风险项需要手动确认。清理前建议创建系统还原点。
    </n-alert>

    <n-button type="primary" @click="scan" :loading="scanning" class="mb-4">
      扫描可清理空间
    </n-button>

    <n-space vertical size="medium" v-if="targets.length">
      <n-card v-for="t in targets" :key="t.path" :title="t.description" size="small">
        <template #header-extra>
          <n-tag :type="riskType(t.risk) as any" size="small">{{ riskLabel(t.risk) }}</n-tag>
        </template>
        <n-space justify="space-between" align="center">
          <span>{{ t.path }} — {{ formatBytes(t.totalSize) }} · {{ t.fileCount }} 个文件</span>
          <n-checkbox v-model:checked="t.checked">清理</n-checkbox>
        </n-space>
      </n-card>

      <n-space>
        <n-button type="warning" @click="cleanSelected" :loading="cleaning" :disabled="!hasChecked">
          清理选中项 ({{ checkedCount }})
        </n-button>
        <n-button @click="selectAll" size="small">全选</n-button>
        <n-button @click="deselectAll" size="small">取消全选</n-button>
      </n-space>
    </n-space>

    <n-space vertical size="small" v-if="results.length">
      <n-h3>清理结果</n-h3>
      <n-alert v-for="r in results" :key="r.target" :type="r.success ? 'success' : 'error'" :bordered="false">
        <template #header>{{ r.target }}</template>
        {{ r.success ? `成功释放 ${formatBytes(r.freedBytes)}，清理 ${r.fileCount} 个文件` : `失败: ${r.error}` }}
      </n-alert>

      <n-space>
        <n-button type="primary" @click="resetPage">完成</n-button>
        <n-button @click="scan">重新扫描</n-button>
      </n-space>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ScanCleanupPaths, CleanTargets } from '@wails/go/main/App'
import { formatBytes, riskType, riskLabel } from '@/api/bridge'

interface CleanTarget { path: string; description: string; risk: string; browser?: string; fileCount: number; totalSize: number; checked: boolean }
interface CleanResult { target: string; success: boolean; freedBytes: number; fileCount: number; error?: string }

const targets = ref<CleanTarget[]>([])
const results = ref<CleanResult[]>([])
const scanning = ref(false)
const cleaning = ref(false)

const checkedCount = computed(() => targets.value.filter(t => t.checked).length)
const hasChecked = computed(() => checkedCount.value > 0)

async function scan() {
  scanning.value = true
  results.value = []
  try {
    const result = await ScanCleanupPaths()
    if (result) targets.value = result as CleanTarget[]
  } catch (e) { console.error(e) }
  scanning.value = false
}

async function cleanSelected() {
  cleaning.value = true
  results.value = []
  const paths = targets.value.filter(t => t.checked).map(t => t.path)
  try {
    const result = await CleanTargets(paths)
    if (result) results.value = result as CleanResult[]
  } catch (e) { console.error(e) }
  cleaning.value = false
}

function selectAll() { targets.value.forEach(t => t.checked = true) }
function deselectAll() { targets.value.forEach(t => t.checked = false) }
function resetPage() { targets.value = []; results.value = [] }
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
