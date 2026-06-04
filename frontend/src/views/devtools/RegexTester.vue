<template>
  <div>
    <n-h2>正则测试</n-h2>
    <n-p>在线编写与测试正则表达式，支持实时匹配与分组展示。</n-p>

    <n-card>
      <n-space vertical>
        <n-space>
          <n-input v-model:value="pattern" placeholder="输入正则表达式" clearable style="width: 350px" />
          <n-select v-model:value="flags" multiple :options="flagOptions" style="width: 200px" placeholder="修饰符" />
        </n-space>

        <n-space>
          <n-checkbox v-model:checked="caseInsensitive" @update:checked="syncFlags">忽略大小写 (i)</n-checkbox>
          <n-checkbox v-model:checked="global" @update:checked="syncFlags">全局匹配 (g)</n-checkbox>
          <n-checkbox v-model:checked="multiline" @update:checked="syncFlags">多行模式 (m)</n-checkbox>
        </n-space>

        <n-input
          v-model:value="testText"
          type="textarea"
          :rows="8"
          placeholder="输入要测试的文本..."
        />

        <n-button type="primary" @click="testRegex" :disabled="!pattern">测试匹配</n-button>

        <n-divider />

        <n-empty v-if="matches.length === 0 && tested" description="无匹配结果" />
        <n-empty v-if="matches.length === 0 && !tested" description="输入正则和文本后点击测试匹配" />

        <n-data-table
          v-if="matches.length"
          :columns="matchColumns"
          :data="matches"
          size="small"
          :bordered="true"
        />

        <n-card v-if="matches.length" title="匹配统计" size="small">
          共 {{ matches.length }} 个匹配 |
          匹配文本总长度: {{ totalMatchLen }} 字符 |
          修饰符: /{{ pattern }}/{{ flagsStr }}
        </n-card>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'

const pattern = ref('')
const testText = ref('')
const flags = ref<string[]>([])
const tested = ref(false)
const matches = ref<any[]>([])
const message = useMessage()

const caseInsensitive = ref(false)
const global = ref(true)
const multiline = ref(false)

const flagOptions = [
  { label: '忽略大小写 (i)', value: 'i' },
  { label: '全局匹配 (g)', value: 'g' },
  { label: '多行模式 (m)', value: 'm' },
  { label: '点号匹配换行 (s)', value: 's' },
  { label: 'Unicode (u)', value: 'u' },
]

const matchColumns = [
  { title: '#', key: 'index', width: 50 },
  { title: '匹配文本', key: 'text', ellipsis: { tooltip: true } },
  { title: '位置', key: 'pos', width: 100 },
  { title: '分组', key: 'groups', ellipsis: { tooltip: true }, render: (r: any) => r.groups || '-' },
]

const flagsStr = computed(() => flags.value.join(''))
const totalMatchLen = computed(() => matches.value.reduce((sum: number, m: any) => sum + m.text.length, 0))

function syncFlags() {
  const f: string[] = []
  if (caseInsensitive.value) f.push('i')
  if (global.value) f.push('g')
  if (multiline.value) f.push('m')
  flags.value = f
}

function testRegex() {
  if (!pattern.value) {
    message.warning('请输入正则表达式')
    return
  }

  tested.value = true
  matches.value = []

  try {
    const regex = new RegExp(pattern.value, flags.value.join(''))
    let m: RegExpExecArray | null
    let idx = 0

    if (flags.value.includes('g')) {
      while ((m = regex.exec(testText.value)) !== null) {
        const groups = m.length > 1 ? m.slice(1).filter(g => g !== undefined).join(', ') : ''
        matches.value.push({
          index: ++idx,
          text: m[0],
          pos: m.index,
          groups,
        })
        if (m.index === regex.lastIndex) regex.lastIndex++
      }
    } else {
      m = regex.exec(testText.value)
      if (m) {
        const groups = m.length > 1 ? m.slice(1).filter(g => g !== undefined).join(', ') : ''
        matches.value.push({
          index: 1,
          text: m[0],
          pos: m.index,
          groups,
        })
      }
    }

    if (matches.value.length) {
      message.success(`找到 ${matches.value.length} 个匹配`)
    }
  } catch (e: any) {
    message.error('正则表达式错误: ' + (e.message || String(e)))
  }
}
</script>

<style scoped>
</style>
