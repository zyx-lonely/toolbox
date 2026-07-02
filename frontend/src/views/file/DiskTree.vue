<template>
  <div>
    <h2>磁盘空间树状图</h2>
    <p>可视化展示各文件夹空间占用情况，点击目录可打开资源管理器。</p>
    <n-space class="mb-4">
      <n-input v-model:value="scanPath" placeholder="选择目录" readonly style="width: 300px" />
      <n-button @click="selectDir">选择目录</n-button>
      <n-button type="primary" @click="scan" :loading="loading" :disabled="!scanPath">扫描</n-button>
      <n-button @click="scan" :loading="loading" :disabled="!scanPath" secondary>刷新</n-button>
    </n-space>
    <n-space v-if="tree" class="mb-4">
      <n-input v-model:value="filterText" placeholder="搜索过滤..." clearable style="width: 250px" />
      <n-button @click="expandAll" size="small">全部展开</n-button>
      <n-button @click="collapseAll" size="small">全部折叠</n-button>
      <n-text depth="3" style="line-height: 32px;">共 {{ filteredCount }} 项，{{ formatSize(totalSize) }}</n-text>
    </n-space>
    <n-empty v-if="!loading && !tree" description="请选择目录后点击扫描" />
    <n-spin :show="loading">
      <n-card v-if="tree" size="small">
        <n-tree
          :data="filteredTreeData"
          :default-expand-all="expandAllFlag"
          :render-label="renderLabel"
          :key="treeKey"
          @update:selected-keys="onSelect"
        />
      </n-card>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { h, ref, computed, watch } from 'vue'
import { NTag } from 'naive-ui'
import { ScanDiskTree, SelectDirectory, OpenDirectory } from '@wails/go/main/App'

interface DiskTreeNode { name: string; path: string; size: number; isDir: boolean; children?: DiskTreeNode[] }

const scanPath = ref('')
const tree = ref<DiskTreeNode | null>(null)
const loading = ref(false)
const filterText = ref('')
const expandAllFlag = ref(true)
const treeKey = ref(0)

const totalSize = computed(() => tree.value?.size || 0)

function countFiles(node: DiskTreeNode): number {
  if (!node.isDir) return 1
  let count = 0
  node.children?.forEach(c => { count += countFiles(c) })
  return count
}

const filteredTreeData = computed(() => {
  if (!tree.value) return []
  if (!filterText.value) return buildTreeData(tree.value)
  return buildTreeData(filterNode(tree.value, filterText.value.toLowerCase()))
})

const filteredCount = computed(() => {
  if (!tree.value) return 0
  if (!filterText.value) return countFiles(tree.value)
  const filtered = filterNode(tree.value, filterText.value.toLowerCase())
  return countFiles(filtered)
})

function filterNode(node: DiskTreeNode, keyword: string): DiskTreeNode | null {
  if (!node.children) {
    if (node.name.toLowerCase().includes(keyword)) return node
    return null
  }
  const matchedChildren: DiskTreeNode[] = []
  for (const child of node.children) {
    const filtered = filterNode(child, keyword)
    if (filtered) matchedChildren.push(filtered)
  }
  if (matchedChildren.length > 0 || node.name.toLowerCase().includes(keyword)) {
    return { ...node, children: matchedChildren.length > 0 ? matchedChildren : node.children }
  }
  return null
}

function formatSize(bytes: number): string {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}

function buildTreeData(node: DiskTreeNode): any[] {
  if (!node.children) return []
  return node.children.map(child => ({
    key: child.path,
    label: child.name,
    children: child.isDir ? buildTreeData(child) : undefined,
    node: child,
  }))
}

function renderLabel({ option }: any) {
  const node = option.node as DiskTreeNode
  const size = formatSize(node.size)
  const color = node.isDir ? 'info' : 'default'
  const style = node.isDir ? 'cursor: pointer; font-weight: 500; color: #18a058; text-decoration: underline;' : ''
  return h('span', [
    h('span', { style }, option.label),
    h(NTag, { size: 'small', type: color as any, style: 'margin-left: 8px' }, { default: () => size })
  ])
}

async function selectDir() {
  try {
    const dir = await SelectDirectory()
    if (dir) scanPath.value = dir
  } catch (e) { console.error(e) }
}

async function scan() {
  if (!scanPath.value) return
  loading.value = true
  tree.value = null
  filterText.value = ''
  try {
    const result = await ScanDiskTree(scanPath.value, 3) as DiskTreeNode
    tree.value = result
  } catch (e: any) { console.error(e) }
  loading.value = false
}

function expandAll() {
  expandAllFlag.value = true
  treeKey.value++
}

function collapseAll() {
  expandAllFlag.value = false
  treeKey.value++
}

function onSelect(keys: string[]) {
  if (keys.length > 0) {
    const path = keys[0]
    if (path && path.includes('\\')) {
      OpenDirectory(path)
    }
  }
}

watch(filterText, () => { treeKey.value++ })
</script>

<style scoped>
.mb-4 { margin-bottom: 16px; }
</style>
