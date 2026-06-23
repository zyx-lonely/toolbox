<template>
  <div>
    <n-h2>文本字数统计</n-h2>
    <n-p>统计文本的字数、行数、词数、字节等信息。</n-p>

    <n-space vertical :size="16">
      <n-input
        v-model:value="inputText"
        type="textarea"
        :rows="12"
        placeholder="在此粘贴或输入文本..."
        @input="analyze"
      />

      <n-grid :cols="4" :x-gap="12" :y-gap="12">
        <n-gi><n-card size="small"><n-statistic label="总字符数" :value="stats.chars" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="中文字符" :value="stats.chinese" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="英文单词" :value="stats.words" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="总行数" :value="stats.lines" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="非空行" :value="stats.nonEmptyLines" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="空行" :value="stats.emptyLines" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="字节数" :value="stats.bytes" /></n-card></n-gi>
        <n-gi><n-card size="small"><n-statistic label="UTF-8 字节" :value="stats.utf8Bytes" /></n-card></n-gi>
      </n-grid>

      <n-card title="详细信息" size="small" v-if="inputText">
        <n-descriptions bordered :column="2" label-placement="left" size="small">
          <n-descriptions-item label="数字字符">{{ stats.digits }}</n-descriptions-item>
          <n-descriptions-item label="标点符号">{{ stats.punctuation }}</n-descriptions-item>
          <n-descriptions-item label="大写字母">{{ stats.uppercase }}</n-descriptions-item>
          <n-descriptions-item label="小写字母">{{ stats.lowercase }}</n-descriptions-item>
          <n-descriptions-item label="空格数">{{ stats.spaces }}</n-descriptions-item>
          <n-descriptions-item label="Tab 数">{{ stats.tabs }}</n-descriptions-item>
        </n-descriptions>
      </n-card>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'

const inputText = ref('')

const stats = reactive({
  chars: 0,
  chinese: 0,
  words: 0,
  lines: 0,
  nonEmptyLines: 0,
  emptyLines: 0,
  bytes: 0,
  utf8Bytes: 0,
  digits: 0,
  punctuation: 0,
  uppercase: 0,
  lowercase: 0,
  spaces: 0,
  tabs: 0,
})

function analyze() {
  const text = inputText.value
  const lineArr = text.split('\n')

  stats.chars = text.length
  stats.lines = lineArr.length
  stats.nonEmptyLines = lineArr.filter(l => l.trim().length > 0).length
  stats.emptyLines = stats.lines - stats.nonEmptyLines
  stats.bytes = new Blob([text]).size
  stats.utf8Bytes = new TextEncoder().encode(text).byteLength
  stats.chinese = (text.match(/[\u4e00-\u9fff]/g) || []).length
  stats.words = text.trim() ? text.trim().split(/\s+/).length : 0
  stats.digits = (text.match(/\d/g) || []).length
  stats.punctuation = (text.match(/[^\w\s\u4e00-\u9fff]/g) || []).length
  stats.uppercase = (text.match(/[A-Z]/g) || []).length
  stats.lowercase = (text.match(/[a-z]/g) || []).length
  stats.spaces = (text.match(/ /g) || []).length
  stats.tabs = (text.match(/\t/g) || []).length
}
</script>
