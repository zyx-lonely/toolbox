<template>
  <div>
    <n-h2>屏幕标尺</n-h2>
    <n-p>测量屏幕上两点间的像素距离。</n-p>

    <n-card>
      <n-space vertical>
        <n-space align="center">
          <n-button @click="toggleMeasure" :type="isMeasuring ? 'error' : 'primary'" :loading="loading">
            <template #icon><n-icon><resize-outline /></n-icon></template>
            {{ isMeasuring ? '停止测量' : '开始测量' }}
          </n-button>
          <n-button @click="clearMeasurement" v-if="result" type="warning" ghost>
            <template #icon><n-icon><trash-outline /></n-icon></template>
            清除
          </n-button>
          <n-button @click="copyResult" v-if="result" type="success" ghost>
            <template #icon><n-icon><copy-outline /></n-icon></template>
            复制结果
          </n-button>
        </n-space>

        <n-alert v-if="isMeasuring" type="info">
          请在屏幕上点击两个点来测量距离。第一个点：{{ point1 ? `(${point1.x}, ${point1.y})` : '等待点击...' }}
        </n-alert>

        <n-card v-if="result" :bordered="true" size="small">
          <n-descriptions :column="2" bordered size="small">
            <n-descriptions-item label="起点">
              ({{ result.x1 }}, {{ result.y1 }})
            </n-descriptions-item>
            <n-descriptions-item label="终点">
              ({{ result.x2 }}, {{ result.y2 }})
            </n-descriptions-item>
            <n-descriptions-item label="水平距离">
              {{ result.dx }} px
            </n-descriptions-item>
            <n-descriptions-item label="垂直距离">
              {{ result.dy }} px
            </n-descriptions-item>
            <n-descriptions-item label="直线距离">
              {{ result.distance }} px
            </n-descriptions-item>
            <n-descriptions-item label="角度">
              {{ result.angle }}°
            </n-descriptions-item>
          </n-descriptions>
        </n-card>

        <n-empty v-if="!result && !isMeasuring" description="点击'开始测量'后在屏幕上选取两个点" />

        <div ref="canvasRef" v-show="isMeasuring"
             style="position:fixed; top:0; left:0; width:100vw; height:100vh; z-index:99999; cursor:crosshair; background:rgba(0,0,0,0.1);"
             @click="handleClick"
             @mousemove="handleMouseMove">
          <svg style="position:absolute; top:0; left:0; width:100%; height:100%; pointer-events:none;">
            <line v-if="point1" :x1="point1.x" :y1="point1.y" :x2="mousePos.x" :y2="mousePos.y"
                  stroke="#18a058" stroke-width="2" stroke-dasharray="5,5" />
            <circle v-if="point1" :cx="point1.x" :cy="point1.y" r="5" fill="#18a058" />
          </svg>
        </div>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import { ResizeOutline, TrashOutline, CopyOutline } from '@vicons/ionicons5'

const message = useMessage()
const loading = ref(false)
const isMeasuring = ref(false)
const point1 = ref<{x:number, y:number} | null>(null)
const mousePos = ref({x:0, y:0})
const result = ref<any>(null)
const canvasRef = ref<any>(null)

function toggleMeasure() {
  if (isMeasuring.value) {
    isMeasuring.value = false
    point1.value = null
  } else {
    isMeasuring.value = true
    point1.value = null
    result.value = null
  }
}

function handleClick(e: MouseEvent) {
  const x = e.clientX
  const y = e.clientY

  if (!point1.value) {
    point1.value = { x, y }
  } else {
    const dx = Math.abs(x - point1.value.x)
    const dy = Math.abs(y - point1.value.y)
    const distance = Math.sqrt(dx*dx + dy*dy)
    const angle = Math.round(Math.atan2(y - point1.value.y, x - point1.value.x) * 180 / Math.PI)

    result.value = {
      x1: point1.value.x,
      y1: point1.value.y,
      x2: x,
      y2: y,
      dx,
      dy,
      distance: Math.round(distance),
      angle
    }

    isMeasuring.value = false
    point1.value = null
    message.success('测量完成')
  }
}

function handleMouseMove(e: MouseEvent) {
  mousePos.value = { x: e.clientX, y: e.clientY }
}

function clearMeasurement() {
  result.value = null
  point1.value = null
}

function copyResult() {
  if (result.value) {
    const text = `起点: (${result.value.x1}, ${result.value.y1}) → 终点: (${result.value.x2}, ${result.value.y2})\n水平: ${result.value.dx}px | 垂直: ${result.value.dy}px | 距离: ${result.value.distance}px | 角度: ${result.value.angle}°`
    navigator.clipboard.writeText(text)
    message.success('结果已复制')
  }
}

onUnmounted(() => {
  isMeasuring.value = false
})
</script>
