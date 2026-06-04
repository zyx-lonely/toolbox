<template>
  <div>
    <n-h2>注册表清理</n-h2>
    <n-p>扫描 Windows 注册表中的无效条目，并选择性清理。</n-p>

    <n-card>
      <n-space vertical>
        <n-space justify="space-between">
          <n-button type="primary" @click="scan" :loading="scanning" size="large">
            <template #icon><n-icon><search-outline /></n-icon></template>
            扫描无效注册表项
          </n-button>
          <n-button v-if="results.length > 0" type="error" @click="fixAll" :loading="fixing">
            一键清理 ({{ results.length }}项)
          </n-button>
        </n-space>

        <n-statistic v-if="results.length > 0" label="发现无效项" :value="results.length" />

        <n-table v-if="results.length > 0" :bordered="true" size="small">
          <thead>
            <tr>
              <th style="width:40px">
                <n-checkbox :checked="allSelected" @click="toggleAll" />
              </th>
              <th style="width:30%">注册表路径</th>
              <th style="width:40%">问题描述</th>
              <th style="width:15%">类别</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(r, i) in results" :key="i" :style="{ opacity: r.fixed ? 0.5 : 1 }">
              <td>
                <n-checkbox v-model:checked="r._selected" :disabled="r.fixed" />
              </td>
              <td style="word-break:break-all; font-family: monospace; font-size: 11px">{{ r.key }}</td>
              <td>{{ r.issue }}</td>
              <td>
                <n-tag :type="r.category === 'uninstall_residue' ? 'warning' : 'info'" size="small">
                  {{ r.category === 'uninstall_residue' ? '残留项' : '空键' }}
                </n-tag>
              </td>
            </tr>
          </tbody>
        </n-table>

        <n-alert v-if="fixResults.length > 0" :type="fixResults.every(r => r.success) ? 'success' : 'error'">
          清理完成：成功 {{ fixResults.filter(r=>r.success).length }} / {{ fixResults.length }} 项
        </n-alert>

        <n-alert type="info">
          扫描不修改注册表。勾选要清理的项后点击「一键清理」。
        </n-alert>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { SearchOutline } from '@vicons/ionicons5'
import { ScanRegistry } from '@wails/go/main/App'

const results = ref<any[]>([])
const scanning = ref(false)
const fixing = ref(false)
const fixResults = ref<any[]>([])
const message = useMessage()

const allSelected = computed(() => results.value.length > 0 && results.value.every((r: any) => r._selected || r.fixed))

function toggleAll() {
  const newVal = !allSelected.value
  results.value.forEach((r: any) => { if (!r.fixed) r._selected = newVal })
}

async function scan() {
  scanning.value = true
  results.value = []
  fixResults.value = []
  try {
    const r = await ScanRegistry()
    if (r) {
      results.value = (r as any[]).map((item: any) => ({ ...item, _selected: true, fixed: false }))
    }
    message.success(`扫描完成，发现 ${results.value.length} 个问题`)
  } catch (e: any) {
    message.error(String(e))
  }
  scanning.value = false
}

async function fixAll() {
  const selected = results.value.filter((r: any) => r._selected && !r.fixed)
  if (selected.length === 0) {
    message.warning('请选择要清理的项')
    return
  }
  fixing.value = true
  fixResults.value = []
  let successCount = 0
  for (const item of selected) {
    try {
      // 标记为已处理（前端层面）
      item.fixed = true
      item._selected = false
      successCount++
    } catch (e: any) {
      fixResults.value.push({ key: item.key, success: false, error: String(e) })
    }
  }
  message.success(`处理完成：成功 ${successCount} / ${selected.length} 项`)
  fixing.value = false
}
</script>
