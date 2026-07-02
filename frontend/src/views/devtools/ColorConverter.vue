<template>
  <div>
    <h2>颜色代码转换</h2>
    <p>HEX / RGB / HSL 颜色格式互转。</p>
    <n-grid :cols="3" :x-gap="12">
      <n-gi>
        <n-card title="HEX → RGB/HSL" size="small">
          <n-space vertical>
            <n-input v-model:value="hexInput" placeholder="#FF5733 或 FF5733" />
            <n-button type="primary" @click="fromHex">转换</n-button>
            <div v-if="hexResult" class="color-preview" :style="{ background: hexResult.hex }" />
            <n-descriptions v-if="hexResult" bordered :column="1" size="small">
              <n-descriptions-item label="HEX">{{ hexResult.hex }}</n-descriptions-item>
              <n-descriptions-item label="RGB">{{ hexResult.rgb }}</n-descriptions-item>
              <n-descriptions-item label="HSL">{{ hexResult.hsl }}</n-descriptions-item>
            </n-descriptions>
          </n-space>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="RGB → HEX/HSL" size="small">
          <n-space vertical>
            <n-space>
              <n-input-number v-model:value="r" :min="0" :max="255" style="width: 80px" placeholder="R" />
              <n-input-number v-model:value="g" :min="0" :max="255" style="width: 80px" placeholder="G" />
              <n-input-number v-model:value="b" :min="0" :max="255" style="width: 80px" placeholder="B" />
            </n-space>
            <n-button type="primary" @click="fromRGB">转换</n-button>
            <div v-if="rgbResult" class="color-preview" :style="{ background: rgbResult.hex }" />
            <n-descriptions v-if="rgbResult" bordered :column="1" size="small">
              <n-descriptions-item label="HEX">{{ rgbResult.hex }}</n-descriptions-item>
              <n-descriptions-item label="RGB">{{ rgbResult.rgb }}</n-descriptions-item>
              <n-descriptions-item label="HSL">{{ rgbResult.hsl }}</n-descriptions-item>
            </n-descriptions>
          </n-space>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="HSL → HEX/RGB" size="small">
          <n-space vertical>
            <n-space>
              <n-input-number v-model:value="h" :min="0" :max="360" style="width: 80px" placeholder="H" />
              <n-input-number v-model:value="s" :min="0" :max="100" style="width: 80px" placeholder="S" />
              <n-input-number v-model:value="l" :min="0" :max="100" style="width: 80px" placeholder="L" />
            </n-space>
            <n-button type="primary" @click="fromHSL">转换</n-button>
            <div v-if="hslResult" class="color-preview" :style="{ background: hslResult.hex }" />
            <n-descriptions v-if="hslResult" bordered :column="1" size="small">
              <n-descriptions-item label="HEX">{{ hslResult.hex }}</n-descriptions-item>
              <n-descriptions-item label="RGB">{{ hslResult.rgb }}</n-descriptions-item>
              <n-descriptions-item label="HSL">{{ hslResult.hsl }}</n-descriptions-item>
            </n-descriptions>
          </n-space>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { HexToColor, RGBToColor, HSLToColor } from '@wails/go/main/App'

interface ColorResult { hex: string; rgb: string; hsl: string; r: number; g: number; b: number; hue: number; sat: number; lit: number }

const hexInput = ref('')
const r = ref<number | null>(null)
const g = ref<number | null>(null)
const b = ref<number | null>(null)
const h = ref<number | null>(null)
const s = ref<number | null>(null)
const l = ref<number | null>(null)

const hexResult = ref<ColorResult | null>(null)
const rgbResult = ref<ColorResult | null>(null)
const hslResult = ref<ColorResult | null>(null)
const message = useMessage()

async function fromHex() {
  try { hexResult.value = await HexToColor(hexInput.value) as ColorResult } catch (e: any) { message.error(String(e)) }
}

async function fromRGB() {
  if (r.value === null || g.value === null || b.value === null) return
  try { rgbResult.value = await RGBToColor(r.value, g.value, b.value) as ColorResult } catch (e: any) { message.error(String(e)) }
}

async function fromHSL() {
  if (h.value === null || s.value === null || l.value === null) return
  try { hslResult.value = await HSLToColor(h.value, s.value, l.value) as ColorResult } catch (e: any) { message.error(String(e)) }
}
</script>

<style scoped>
.color-preview { width: 100%; height: 40px; border-radius: 6px; border: 1px solid #ddd; margin: 8px 0; }
</style>
