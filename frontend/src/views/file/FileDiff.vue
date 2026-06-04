<template>
  <div>
    <n-h2>文件差异对比</n-h2>
    <n-p>对比两个文本文件的差异。</n-p>
    <n-space class="mb-4">
      <n-input-group>
        <n-input v-model:value="oldFilePath" placeholder="选择旧文件" readonly style="width: 250px" />
        <n-button @click="selectOldFile">选择文件</n-button>
      </n-input-group>
      <n-input-group>
        <n-input v-model:value="newFilePath" placeholder="选择新文件" readonly style="width: 250px" />
        <n-button @click="selectNewFile">选择文件</n-button>
      </n-input-group>
      <n-button type="primary" @click="compare" :loading="comparing">开始对比</n-button>
    </n-space>
    <n-empty v-if="!comparing && !diffs.length && compared" description="文件内容完全相同" />
    <div v-if="diffs.length" class="diff-container">
      <div v-for="(d, i) in diffs" :key="i" :class="'diff-line diff-' + d.type">
        <span class="line-num">{{ d.oldLine || '' }}</span>
        <span class="line-num">{{ d.newLine || '' }}</span>
        <span class="line-content">{{ d.content }}</span>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { DiffFiles, SelectFile } from '@wails/go/main/App'
const oldFilePath = ref(''); const newFilePath = ref(''); const diffs = ref<any[]>([]); const comparing = ref(false); const compared = ref(false)
const message = useMessage()
async function selectOldFile() { const f = await SelectFile(); if (f) oldFilePath.value = f as string }
async function selectNewFile() { const f = await SelectFile(); if (f) newFilePath.value = f as string }
async function compare() {
  if (!oldFilePath.value || !newFilePath.value) { message.warning('请选择两个文件'); return }
  comparing.value = true; compared.value = true
  try { const r = await DiffFiles(oldFilePath.value, newFilePath.value); if (r) diffs.value = r as any[] }
  catch(e:any) { message.error(String(e)) }; comparing.value = false
}
</script>
<style scoped>
.mb-4{margin-bottom:16px}
.diff-container{font-family:monospace;font-size:13px;border:1px solid #e0e0e0;border-radius:4px;max-height:600px;overflow:auto}
.diff-line{display:flex;padding:2px 4px;border-bottom:1px solid #f0f0f0}
.diff-line:last-child{border-bottom:none}
.line-num{width:40px;text-align:right;color:#999;padding-right:8px;user-select:none}
.line-content{flex:1;white-space:pre;overflow:hidden;text-overflow:ellipsis}
.diff-same{background:#fff}
.diff-added{background:#e6ffed}
.diff-removed{background:#ffeef0}
.diff-modified{background:#fff8e1}
</style>
