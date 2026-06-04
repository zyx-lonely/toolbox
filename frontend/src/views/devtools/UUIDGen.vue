<template>
  <div class="uuid-gen">
    <n-space vertical>
      <n-space align="center">
        <n-input-number v-model:value="count" :min="1" :max="100" style="width: 100px" />
        <n-button type="primary" @click="handleGenerate" :loading="generating">
          生成
        </n-button>
      </n-space>

      <n-list v-if="uuids.length">
        <n-list-item v-for="(uuid, idx) in uuids" :key="idx">
          <n-space align="center" justify="space-between">
            <span>{{ uuid }}</span>
            <n-button size="small" @click="handleCopy(uuid)">复制</n-button>
          </n-space>
        </n-list-item>
      </n-list>

      <n-empty v-else-if="!generating" description="点击生成获取 UUID" />
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const count = ref(5)
const uuids = ref<string[]>([])
const generating = ref(false)

async function handleGenerate() {
  generating.value = true
  try {
    const result = await window.go.main.App.GenerateUUIDs(count.value)
    uuids.value = result || []
  } catch (e: any) {
    uuids.value = []
  } finally {
    generating.value = false
  }
}

async function handleCopy(uuid: string) {
  try {
    await navigator.clipboard.writeText(uuid)
    window.$message?.success('已复制: ' + uuid)
  } catch {
    window.$message?.error('复制失败')
  }
}
</script>

<style scoped>
.uuid-gen {
  padding: 20px;
}
</style>
