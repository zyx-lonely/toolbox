<template>
  <div>
    <n-h2>取色器</n-h2>
    <n-p>颜色值转换工具 — 在下方拾取或输入颜色。</n-p>

    <n-grid :cols="2" :x-gap="32">
      <n-gi>
        <n-card title="拾色器">
          <n-space vertical>
            <div class="color-preview" :style="{ backgroundColor: hexColor }" />
            <n-color-picker v-model:value="hexColor" :show-alpha="false" style="width: 100%" />
            <n-input v-model:value="hexColor" placeholder="#FF0000" />
          </n-space>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="颜色值">
          <n-description-list label-placement="left" :column="1">
            <n-description-item label="HEX">
              <n-input v-model:value="hexColor" size="small" />
            </n-description-item>
            <n-description-item label="RGB">
              <n-input :value="rgbValue" readonly size="small" />
            </n-description-item>
            <n-description-item label="HSL">
              <n-input :value="hslValue" readonly size="small" />
            </n-description-item>
            <n-description-item label="HSV">
              <n-input :value="hsvValue" readonly size="small" />
            </n-description-item>
          </n-description-list>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const hexColor = ref('#18a058')

function hexToRgb(hex: string): { r: number; g: number; b: number } | null {
  const clean = hex.replace('#', '')
  if (clean.length !== 6) return null
  const r = parseInt(clean.substring(0, 2), 16)
  const g = parseInt(clean.substring(2, 4), 16)
  const b = parseInt(clean.substring(4, 6), 16)
  if (isNaN(r) || isNaN(g) || isNaN(b)) return null
  return { r, g, b }
}

function rgbToHsl(r: number, g: number, b: number): { h: number; s: number; l: number } {
  r /= 255; g /= 255; b /= 255
  const max = Math.max(r, g, b), min = Math.min(r, g, b)
  let h = 0, s = 0, l = (max + min) / 2
  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    switch (max) {
      case r: h = ((g - b) / d + (g < b ? 6 : 0)) / 6; break
      case g: h = ((b - r) / d + 2) / 6; break
      case b: h = ((r - g) / d + 4) / 6; break
    }
  }
  return { h: Math.round(h * 360), s: Math.round(s * 100), l: Math.round(l * 100) }
}

function rgbToHsv(r: number, g: number, b: number): { h: number; s: number; v: number } {
  r /= 255; g /= 255; b /= 255
  const max = Math.max(r, g, b), min = Math.min(r, g, b)
  let h = 0, s = max === 0 ? 0 : (max - min) / max, v = max
  if (max !== min) {
    const d = max - min
    switch (max) {
      case r: h = ((g - b) / d + (g < b ? 6 : 0)) * 60; break
      case g: h = ((b - r) / d + 2) * 60; break
      case b: h = ((r - g) / d + 4) * 60; break
    }
  }
  return { h: Math.round(h), s: Math.round(s * 100), v: Math.round(v * 100) }
}

const rgbValue = computed(() => {
  const c = hexToRgb(hexColor.value)
  return c ? `rgb(${c.r}, ${c.g}, ${c.b})` : '无效颜色'
})

const hslValue = computed(() => {
  const c = hexToRgb(hexColor.value)
  if (!c) return '无效颜色'
  const h = rgbToHsl(c.r, c.g, c.b)
  return `hsl(${h.h}, ${h.s}%, ${h.l}%)`
})

const hsvValue = computed(() => {
  const c = hexToRgb(hexColor.value)
  if (!c) return '无效颜色'
  const h = rgbToHsv(c.r, c.g, c.b)
  return `hsv(${h.h}, ${h.s}%, ${h.v}%)`
})
</script>

<style scoped>
.color-preview {
  width: 100%;
  height: 80px;
  border-radius: 8px;
  border: 1px solid #d9d9d9;
}
</style>
