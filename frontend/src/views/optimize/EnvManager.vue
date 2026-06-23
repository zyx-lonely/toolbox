<template>
  <div>
    <n-h2>&#29992;&#25143;&#21464;&#37327;&#31649;&#29702;</n-h2>
    <n-p>&#30475;&#30475;&#21644;&#32534;&#36753;&#31995;&#32479;/&#29992;&#25143;&#21464;&#37327;&#12290;</n-p>

    <n-space vertical :size="16">
      <n-space>
        <n-radio-group v-model:value="scope" size="medium">
          <n-radio-button value="user">用户变量</n-radio-button>
          <n-radio-button value="system">系统变量</n-radio-button>
        </n-radio-group>
        <n-button type="primary" @click="loadVars" :loading="loading">刷新</n-button>
        <n-button @click="showAddModal = true">新增变量</n-button>
      </n-space>

      <n-data-table
        :columns="columns"
        :data="filteredVars"
        size="small"
        :bordered="true"
        :row-key="(row: any) => row.name"
      />

      <n-card size="small">
        <n-input v-model:value="search" placeholder="搜索变量名或值..." clearable />
      </n-card>
    </n-space>

    <n-modal v-model:show="showAddModal" preset="dialog" title="新增环境变量" positive-text="确定" negative-text="取消" @positive-click="addVariable">
      <n-space vertical>
        <n-input v-model:value="newVar.name" placeholder="变量名" />
        <n-input v-model:value="newVar.value" type="textarea" :rows="3" placeholder="变量值" />
      </n-space>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { NButton, NTag, useMessage, type DataTableColumns } from 'naive-ui'

const message = useMessage()
const loading = ref(false)
const scope = ref('user')
const search = ref('')
const showAddModal = ref(false)
const newVar = ref({ name: '', value: '' })

interface EnvVar { name: string; value: string; type: 'user' | 'system' }
const envVars = ref<EnvVar[]>([])

const filteredVars = computed(() => {
  if (!search.value) return envVars.value
  const q = search.value.toLowerCase()
  return envVars.value.filter(v => v.name.toLowerCase().includes(q) || v.value.toLowerCase().includes(q))
})

const columns: DataTableColumns<EnvVar> = [
  { title: '变量名', key: 'name', width: 250, ellipsis: { tooltip: true } },
  { title: '变量值', key: 'value', ellipsis: { tooltip: true } },
  {
    title: '类型', key: 'type', width: 80,
    render(row) {
      return h(NTag, { type: row.type === 'user' ? 'info' : 'warning', size: 'small' }, { default: () => row.type === 'user' ? '用户' : '系统' })
    }
  },
  {
    title: '操作', key: 'actions', width: 100,
    render(row) {
      return h(NButton, { size: 'small', type: 'error', text: true, onClick: () => deleteVar(row.name) }, { default: () => '删除' })
    }
  }
]

function loadVars() {
  loading.value = true
  setTimeout(() => {
    const isUser = scope.value === 'user'
    envVars.value = (isUser ? [
      { name: 'HOME', value: 'C:\\Users\\Demo', type: 'user' as const },
      { name: 'PATH', value: 'C:\\Users\\Demo\\bin;C:\\Program Files\\Nodejs', type: 'user' as const },
      { name: 'JAVA_HOME', value: 'C:\\Program Files\\Java\\jdk-21', type: 'user' as const },
      { name: 'PYTHONPATH', value: 'C:\\Users\\Demo\\python', type: 'user' as const },
      { name: 'GOPATH', value: 'C:\\Users\\Demo\\go', type: 'user' as const },
    ] : [
      { name: 'PATH', value: 'C:\\Windows\\system32;C:\\Windows;C:\\Program Files', type: 'system' as const },
      { name: 'OS', value: 'Windows_NT', type: 'system' as const },
      { name: 'COMPUTERNAME', value: 'DESKTOP-PC01', type: 'system' as const },
      { name: 'PROCESSOR_ARCHITECTURE', value: 'AMD64', type: 'system' as const },
      { name: 'ProgramFiles', value: 'C:\\Program Files', type: 'system' as const },
      { name: 'SystemRoot', value: 'C:\\Windows', type: 'system' as const },
      { name: 'TEMP', value: 'C:\\Windows\\Temp', type: 'system' as const },
    ])
    loading.value = false
    message.success('已加载 ' + envVars.value.length + ' 个' + (isUser ? '用户' : '系统') + '变量')
  }, 300)
}

function addVariable() {
  if (!newVar.value.name.trim()) {
    message.warning('请输入变量名')
    return false
  }
  envVars.value.push({ name: newVar.value.name, value: newVar.value.value, type: scope.value as 'user' | 'system' })
  message.success('已添加环境变量: ' + newVar.value.name)
  newVar.value = { name: '', value: '' }
  return true
}

function deleteVar(name: string) {
  envVars.value = envVars.value.filter(v => v.name !== name)
  message.success('已删除环境变量: ' + name)
}
</script>
