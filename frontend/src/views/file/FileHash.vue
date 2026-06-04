<template>
  <div>
    <n-h2>文件哈希校验</n-h2>
    <n-p>计算文件的 MD5、SHA1、SHA256 哈希值，用于校验文件完整性。</n-p>

    <n-card>
      <n-space vertical>
        <n-space>
          <n-input-group>
            <n-input v-model:value="filePath" placeholder="选择文件" readonly style="width: 400px" />
            <n-button @click="selectFile">选择文件</n-button>
          </n-input-group>
        </n-space>

        <n-radio-group v-model:value="algorithm">
          <n-radio-button value="md5">MD5</n-radio-button>
          <n-radio-button value="sha1">SHA1</n-radio-button>
          <n-radio-button value="sha256">SHA256</n-radio-button>
        </n-radio-group>

        <n-button type="primary" @click="compute" :loading="computing" :disabled="!filePath">
          计算哈希值
        </n-button>

        <n-input v-if="hash" :value="hash" readonly placeholder="哈希值将显示在这里..." type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" style="font-family: monospace" />

        <n-space v-if="hash" align="center">
          <n-button @click="copyHash">复制</n-button>
          <n-button @click="compareMode = !compareMode">对比验证</n-button>
        </n-space>

        <n-space v-if="compareMode && hash" vertical>
          <n-input v-model:value="expectedHash" placeholder="粘贴预期哈希值进行比较..." />
          <n-tag v-if="compareResult !== null" :type="compareResult ? 'success' : 'error'">
            {{ compareResult ? '✓ 哈希值匹配，文件完整' : '✗ 哈希值不匹配，文件已损坏或被篡改' }}
          </n-tag>
        </n-space>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { ComputeFileHash, SelectFile } from '@wails/go/main/App'

const filePath = ref('')
const algorithm = ref('sha256')
const hash = ref('')
const computing = ref(false)
const compareMode = ref(false)
const expectedHash = ref('')
const message = useMessage()

const compareResult = computed(() => {
  if (!hash.value || !expectedHash.value) return null
  return hash.value.toLowerCase() === expectedHash.value.toLowerCase().trim()
})

async function selectFile() {
  const f = await SelectFile()
  if (f) {
    filePath.value = f as string
    hash.value = ''
    compareMode.value = false
    expectedHash.value = ''
  }
}

async function compute() {
  if (!filePath.value) return
  computing.value = true
  hash.value = ''
  try {
    const h = await ComputeFileHash(filePath.value, algorithm.value)
    if (h) hash.value = h as string
  } catch (e: any) {
    message.error(String(e))
  }
  computing.value = false
}

async function copyHash() {
  if (!hash.value) return
  try {
    await navigator.clipboard.writeText(hash.value)
    message.success('已复制')
  } catch {
    message.warning('复制失败')
  }
}
</script>
