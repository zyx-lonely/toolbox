<template>
  <div>
    <n-h2>外部工具</n-h2>
    <n-p>运行外部工具，扩展 PC Toolbox 功能。</n-p>

    <n-card>
      <n-space style="margin-bottom: 12px" align="center">
        <n-button @click="loadTools" :loading="loading" type="primary">
          <template #icon><n-icon><refresh-outline /></n-icon></template>
          刷新
        </n-button>
      </n-space>

      <n-alert type="info" style="margin-bottom: 12px">
        请将外部工具的可执行文件放到 <n-tag type="info">{{ toolsDir }}</n-tag> 目录中
      </n-alert>

      <n-data-table
        :columns="columns"
        :data="tools"
        :bordered="true"
        :single-line="false"
        :loading="loading"
        size="small"
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, h, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { RefreshOutline } from '@vicons/ionicons5'
import { GetExternalTools, RunExternalTool, GetToolsDir, CheckToolExists } from '@wails/go/main/App'
import type { DataTableColumn } from 'naive-ui'

const tools = ref<any[]>([])
const loading = ref(false)
const toolsDir = ref('')
const message = useMessage()

const pagination = { pageSize: 50 }

const columns: DataTableColumn[] = [
  { title: '工具名称', key: 'name', width: 200, sorter: (a: any, b: any) => a.name.localeCompare(b.name) },
  { title: '描述', key: 'description', ellipsis: true },
  { title: '分类', key: 'category', width: 120 },
  {
    title: '状态', key: 'status', width: 100,
    render(row: any) {
      return h('n-tag', {
        type: row.exists ? 'success' : 'error'
      }, { default: () => row.exists ? '可用' : '不可用' })
    }
  },
  {
    title: '操作', key: 'actions', width: 100,
    render(row: any) {
      return h('button', {
        class: 'n-button n-button--primary n-button--tiny',
        style: 'padding: 4px 8px; border: none; border-radius: 4px; cursor: pointer; color: #fff; background: #18a058;',
        onClick: () => runTool(row),
        disabled: !row.exists
      }, { default: () => '运行' })
    }
  }
]

async function loadTools() {
  loading.value = true
  try {
    const result = await GetExternalTools() as any
    if (result && result.success && result.data) {
      // 检查工具是否可用
      const toolsList = result.data as any[]
      for (const tool of toolsList) {
        tool.exists = await checkToolExists(tool.executable)
      }
      tools.value = toolsList
    } else if (result && result.error) {
      message.error('加载工具列表失败：' + result.error)
    }
    
    const dir = await GetToolsDir()
    if (dir) {
      toolsDir.value = dir as string
    }
  } catch (e: any) {
    message.error('加载工具列表失败：' + String(e))
  } finally {
    loading.value = false
  }
}

async function checkToolExists(executable: string) {
  try {
    return await CheckToolExists(executable)
  } catch {}
  return false
}

async function runTool(row: any) {
  try {
    // 后端 RunExternalTool 接收 JSON 字符串，需要序列化
    const toolJSON = JSON.stringify({
      name: row.name,
      description: row.description,
      executable: row.executable,
      args: row.args || '',
      workingDir: row.workingDir || '',
      icon: row.icon || '',
      category: row.category || ''
    })
    const result = await RunExternalTool(toolJSON) as any
    if (result && result.success) {
      message.success(`已启动 ${row.name}`)
    } else {
      message.error(`启动工具失败: ${result?.error || '未知错误'}`)
    }
  } catch (e: any) {
    message.error(`启动工具失败: ${e}`)
  }
}

onMounted(() => {
  loadTools()
})
</script>
