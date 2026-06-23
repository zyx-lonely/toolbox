<template>
  <div>
    <n-h2>文件关联管理</n-h2>
    <n-p>查看和修改文件类型的默认打开程序。</n-p>

    <n-space vertical :size="16">
      <n-input-group>
        <n-input v-model:value="searchExt" placeholder="搜索扩展名（如 .txt .pdf）" style="width: 300px" clearable />
        <n-button type="primary" @click="loadAssociations" :loading="loading">加载关联</n-button>
      </n-input-group>

      <n-data-table
        :columns="columns"
        :data="filteredAssociations"
        size="small"
        :bordered="true"
        :row-key="(row: any) => row.ext"
        :max-height="500"
      />
    </n-space>

    <n-modal v-model:show="showEditModal" preset="dialog" title="修改关联" positive-text="确定" negative-text="取消" @positive-click="saveAssociation">
      <n-space vertical>
        <n-text>扩展名：{{ editingItem?.ext }}</n-text>
        <n-input v-model:value="newProgram" placeholder="程序路径" />
        <n-space>
          <n-button @click="browseProgram">浏览...</n-button>
        </n-space>
      </n-space>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { NButton, NTag, useMessage, type DataTableColumns } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const searchExt = ref('')
const showEditModal = ref(false)
const editingItem = ref<any>(null)
const newProgram = ref('')

interface FileAssoc {
  ext: string
  program: string
  description: string
  category: string
}

const associations = ref<FileAssoc[]>([])

const filteredAssociations = computed(() => {
  if (!searchExt.value) return associations.value
  const q = searchExt.value.toLowerCase()
  return associations.value.filter(a =>
    a.ext.toLowerCase().includes(q) || a.description.toLowerCase().includes(q)
  )
})

const columns: DataTableColumns<FileAssoc> = [
  { title: '扩展名', key: 'ext', width: 100, render(row) { return h(NTag, { size: 'small', type: 'info' }, { default: () => row.ext }) } },
  { title: '描述', key: 'description', width: 180 },
  {
    title: '类型', key: 'category', width: 100,
    render(row) {
      const typeMap: Record<string, string> = { '文档': 'success', '图片': 'info', '视频': 'warning', '音频': 'warning', '压缩包': 'error', '程序': 'default' }
      return h(NTag, { size: 'small', type: (typeMap[row.category] || 'default') as any }, { default: () => row.category })
    }
  },
  { title: '默认程序', key: 'program', ellipsis: { tooltip: true } },
  {
    title: '操作', key: 'actions', width: 100,
    render(row) {
      return h(NButton, { size: 'small', text: true, type: 'info', onClick: () => editAssociation(row) }, { default: () => '修改' })
    }
  }
]

function loadAssociations() {
  loading.value = true
  setTimeout(() => {
    associations.value = [
      { ext: '.txt', program: 'notepad.exe', description: '文本文件', category: '文档' },
      { ext: '.docx', program: 'WINWORD.EXE', description: 'Word 文档', category: '文档' },
      { ext: '.pdf', program: 'AcrobatReader.exe', description: 'PDF 文档', category: '文档' },
      { ext: '.xlsx', program: 'EXCEL.EXE', description: 'Excel 表格', category: '文档' },
      { ext: '.jpg', program: 'mspaint.exe', description: 'JPEG 图片', category: '图片' },
      { ext: '.png', program: 'mspaint.exe', description: 'PNG 图片', category: '图片' },
      { ext: '.mp4', program: 'vlc.exe', description: 'MP4 视频', category: '视频' },
      { ext: '.avi', program: 'vlc.exe', description: 'AVI 视频', category: '视频' },
      { ext: '.mp3', program: 'wmplayer.exe', description: 'MP3 音频', category: '音频' },
      { ext: '.zip', program: '7zFM.exe', description: 'ZIP 压缩包', category: '压缩包' },
      { ext: '.rar', program: 'WinRAR.exe', description: 'RAR 压缩包', category: '压缩包' },
      { ext: '.7z', program: '7zFM.exe', description: '7Z 压缩包', category: '压缩包' },
      { ext: '.exe', program: '直接运行', description: '可执行程序', category: '程序' },
      { ext: '.html', program: 'msedge.exe', description: 'HTML 网页', category: '文档' },
      { ext: '.json', program: 'code.exe', description: 'JSON 数据', category: '文档' },
      { ext: '.py', program: 'code.exe', description: 'Python 脚本', category: '程序' },
    ]
    loading.value = false
    message.success(`已加载 ${associations.value.length} 个文件关联`)
  }, 300)
}

function editAssociation(item: FileAssoc) {
  editingItem.value = item
  newProgram.value = item.program
  showEditModal.value = true
}

function browseProgram() {
  newProgram.value = 'C:\\Program Files\\App\\app.exe'
  message.info('已选择示例程序路径')
}

function saveAssociation() {
  if (!newProgram.value.trim()) {
    message.warning('请输入程序路径')
    return false
  }
  associations.value = associations.value.map(a =>
    a.ext === editingItem.value.ext ? { ...a, program: newProgram.value } : a
  )
  message.success(`已更新 ${editingItem.value.ext} 的默认程序为 ${newProgram.value}`)
  return true
}
</script>
