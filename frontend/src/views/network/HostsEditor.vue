<template>
  <div>
    <n-h2>Hosts 文件编辑器</n-h2>
    <n-p>编辑系统的 hosts 文件，用于域名解析映射。</n-p>

    <n-card>
      <n-space vertical>
        <n-space justify="space-between">
          <n-button @click="loadHosts" :loading="loading" type="primary">加载 hosts</n-button>
          <n-space>
            <n-button @click="saveHosts" :loading="saving" type="warning">保存</n-button>
            <n-button ghost @click="switchMode">{{ rawMode ? '表格模式' : '编辑模式' }}</n-button>
          </n-space>
        </n-space>

        <div v-if="!rawMode">
          <n-space style="margin-bottom: 8px">
            <n-button size="small" @click="addRow" type="success">添加行</n-button>
          </n-space>
          <n-data-table
            :columns="columns"
            :data="parsedEntries"
            :bordered="true"
            :single-line="false"
            size="small"
            :max-height="500"
          />
        </div>

        <n-input
          v-else
          v-model:value="rawContent"
          type="textarea"
          placeholder="等待加载 hosts 文件..."
          :autosize="{ minRows: 15, maxRows: 30 }"
          style="font-family: 'Consolas', 'Courier New', monospace; font-size: 13px"
        />

        <n-alert type="info">
          修改 hosts 文件需要管理员权限。编辑后点击「保存」生效。
        </n-alert>
      </n-space>
    </n-card>

    <n-modal v-model:show="showAddModal" title="添加 Hosts 条目" preset="card" style="width: 480px">
      <n-space vertical>
        <n-input v-model:value="editItem.ip" placeholder="IP 地址（如 127.0.0.1）" />
        <n-input v-model:value="editItem.hostname" placeholder="域名（如 example.com）" />
        <n-space justify="end">
          <n-button @click="showAddModal = false">取消</n-button>
          <n-button type="primary" @click="confirmAdd" :disabled="!editItem.ip || !editItem.hostname">
            确定
          </n-button>
        </n-space>
      </n-space>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { useMessage, useDialog } from 'naive-ui'
import { GetHostsEntries, SaveHostsEntries } from '@wails/go/main/App'
import type { DataTableColumn } from 'naive-ui'

const loading = ref(false)
const saving = ref(false)
const rawMode = ref(true)
const rawContent = ref('')
const entries = ref<any[]>([])
const showAddModal = ref(false)
const editItem = ref({ ip: '', hostname: '' })
const message = useMessage()
const dialog = useDialog()

const parsedEntries = computed(() =>
  entries.value.filter((e: any) => e.ip && e.hostname)
)

const columns: DataTableColumn[] = [
  { title: 'IP 地址', key: 'ip', width: 140 },
  { title: '域名', key: 'hostname', ellipsis: true },
  {
    title: '操作', key: 'actions', width: 80,
    render(_row: any, index: number) {
      return h('a', {
        style: 'cursor: pointer; color: #d03050;',
        onClick: () => {
          const realIdx = entries.value.findIndex((e: any) => e.ip === _row.ip && e.hostname === _row.hostname)
          if (realIdx >= 0) deleteEntry(realIdx)
        }
      }, { default: () => '删除' })
    }
  }
]

function switchMode() {
  rawMode.value = !rawMode.value
  if (rawMode.value && entries.value.length > 0) {
    rebuildRaw()
  }
}

function rebuildRaw() {
  rawContent.value = entries.value.map((e: any) => {
    if (e.comment) return `# ${e.comment}`
    if (!e.ip || !e.hostname) return e.line || ''
    return `${e.ip} ${e.hostname}`
  }).join('\n')
}

async function loadHosts() {
  loading.value = true
  try {
    const r = await GetHostsEntries()
    if (r) {
      entries.value = r as any[]
      rebuildRaw()
    }
  } catch (e: any) {
    message.error(String(e))
  }
  loading.value = false
}

async function saveHosts() {
  saving.value = true
  try {
    if (rawMode.value) {
      // 从文本重建条目再保存
      const lines = rawContent.value.split('\n')
      const newEntries = []
      for (const line of lines) {
        const trimmed = line.trim()
        if (!trimmed) continue
        if (trimmed.startsWith('#')) {
          newEntries.push({ comment: trimmed.replace(/^#\s*/, ''), line: trimmed })
        } else {
          const parts = trimmed.split(/\s+/)
          if (parts.length >= 2) {
            newEntries.push({ ip: parts[0], hostname: parts.slice(1).join(' '), line: trimmed })
          }
        }
      }
      entries.value = newEntries
    }
    await SaveHostsEntries(entries.value)
    message.success('已保存')
  } catch (e: any) {
    message.error(String(e))
  }
  saving.value = false
}

function addRow() {
  editItem.value = { ip: '', hostname: '' }
  showAddModal.value = true
}

function confirmAdd() {
  const line = `${editItem.value.ip} ${editItem.value.hostname}`
  entries.value.push({ ip: editItem.value.ip, hostname: editItem.value.hostname, line })
  if (rawMode.value) {
    rawContent.value += '\n' + line
  }
  showAddModal.value = false
}

function deleteEntry(index: number) {
  dialog.warning({
    title: '删除',
    content: `确定删除 ${entries.value[index].ip || ''} ${entries.value[index].hostname || entries.value[index].comment || ''} 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      entries.value.splice(index, 1)
      if (rawMode.value) rebuildRaw()
    }
  })
}

loadHosts()
</script>
