<template>
  <div>
    <n-h2>主题商店</n-h2>
    <n-p>选择应用主题风格。</n-p>

    <n-grid :cols="3" :x-gap="16" :y-gap="16">
      <n-gi v-for="t in themes" :key="t.key">
        <n-card :title="t.label" size="small" hoverable
          :style="{ border: isActive(t.key) ? '2px solid #18a058' : '', cursor: 'pointer' }"
          @click="applyTheme(t.key)">
          <div class="preview" :style="{ background: t.bg, color: t.fg, height: '80px', borderRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '24px' }">
            {{ t.icon }}
          </div>
          <template #action>
            <n-button block :type="isActive(t.key) ? 'success' : 'primary'" size="small">
              {{ isActive(t.key) ? '使用中' : '应用' }}
            </n-button>
          </template>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>
<script setup lang="ts">
const themes = [
  { key: 'light', label: '明亮主题', bg: '#ffffff', fg: '#333', icon: '☀️' },
  { key: 'dark', label: '深色主题', bg: '#1e1e2e', fg: '#cdd6f4', icon: '🌙' },
  { key: 'auto', label: '跟随系统', bg: 'linear-gradient(135deg,#fff,#1e1e2e)', fg: '#18a058', icon: '🔄' },
]
const appStore = useAppStore()
import { useAppStore } from '@/store/app'
function isActive(key: string) { return key === 'dark' ? appStore.darkMode : key === 'light' ? !appStore.darkMode : false }
function applyTheme(key: string) {
  if (key === 'dark') appStore.darkMode = true
  else if (key === 'light') appStore.darkMode = false
  else appStore.darkMode = window.matchMedia('(prefers-color-scheme: dark)').matches
}
</script>
<style scoped>.preview{margin-bottom:8px}</style>
