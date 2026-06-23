<template>
  <div>
    <n-h2>文本对比工具</n-h2>
    <n-p>逐行对比两段文本，高亮显示差异。</n-p>

    <n-space vertical :size="16">
      <n-grid :cols="2" :x-gap="12">
        <n-gi>
          <n-card title="原始文本" size="small">
            <n-input
              v-model:value="leftText"
              type="textarea"
              :rows="15"
              placeholder="粘贴原始文本..."
            />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card title="修改后文本" size="small">
            <n-input
              v-model:value="rightText"
              type="textarea"
              :rows="15"
              placeholder="粘贴修改后的文本..."
            />
          </n-card>
        </n-gi>
      </n-grid>

      <n-space>
        <n-button type="primary" @click="compare">对比</n-button>
        <n-button @click="clearAll">清空</n-button>
      </n-space>

      <n-card v-if="diffResult.length > 0" title="对比结果" size="small">
        <n-space vertical :size="2">
          <div v-for="(line, idx) in diffResult" :key="idx"
            :style="{
              padding: '4px 8px',
              fontFamily: 'monospace',
              fontSize: '13px',
              borderRadius: '4px',
              backgroundColor: line.type === 'added' ? '#e6ffed' : line.type === 'removed' ? '#ffeef0' : '#f6f8fa',
              color: line.type === 'added' ? '#22863a' : line.type === 'removed' ? '#cb2431' : '#24292e',
              borderLeft: line.type === 'added' ? '3px solid #22863a' : line.type === 'removed' ? '3px solid #cb2431' : '3px solid transparent'
            }"
          >
            <span style="color: #999; user-select: none; margin-right: 8px;">{{ String(line.leftNum || line.rightNum || '').padStart(3, ' ') }}</span>
            {{ line.text }}
          </div>
        </n-space>
        <template #header-extra>
          <n-space :size="16">
            <n-tag type="success" size="small">+{{ stats.added }}</n-tag>
            <n-tag type="error" size="small">-{{ stats.removed }}</n-tag>
            <n-tag type="info" size="small">{{ stats.unchanged }} unchanged</n-tag>
          </n-space>
        </template>
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'

const message = useMessage()
const leftText = ref('')
const rightText = ref('')

interface DiffLine {
  text: string
  type: 'added' | 'removed' | 'unchanged'
  leftNum?: number
  rightNum?: number
}

const diffResult = ref<DiffLine[]>([])
const stats = computed(() => ({
  added: diffResult.value.filter(l => l.type === 'added').length,
  removed: diffResult.value.filter(l => l.type === 'removed').length,
  unchanged: diffResult.value.filter(l => l.type === 'unchanged').length,
}))

function compare() {
  const leftLines = leftText.value.split('\n')
  const rightLines = rightText.value.split('\n')
  const result: DiffLine[] = []

  const maxLen = Math.max(leftLines.length, rightLines.length)
  let li = 0, ri = 0

  while (li < leftLines.length || ri < rightLines.length) {
    if (li >= leftLines.length) {
      result.push({ text: rightLines[ri], type: 'added', rightNum: ri + 1 })
      ri++
    } else if (ri >= rightLines.length) {
      result.push({ text: leftLines[li], type: 'removed', leftNum: li + 1 })
      li++
    } else if (leftLines[li] === rightLines[ri]) {
      result.push({ text: leftLines[li], type: 'unchanged', leftNum: li + 1, rightNum: ri + 1 })
      li++; ri++
    } else {
      result.push({ text: leftLines[li], type: 'removed', leftNum: li + 1 })
      result.push({ text: rightLines[ri], type: 'added', rightNum: ri + 1 })
      li++; ri++
    }
  }

  diffResult.value = result
  message.success(`对比完成: ${stats.value.added} 行新增, ${stats.value.removed} 行删除, ${stats.value.unchanged} 行相同`)
}

function clearAll() {
  leftText.value = ''
  rightText.value = ''
  diffResult.value = []
}
</script>
