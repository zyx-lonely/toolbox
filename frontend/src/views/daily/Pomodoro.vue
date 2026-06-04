<template>
  <div>
    <n-h2>番茄时钟</n-h2>
    <n-p>专注工作 25 分钟，休息 5 分钟，提高工作效率。</n-p>

    <n-card>
      <n-space vertical align="center" size="large">
        <div class="timer-display" :class="{ 'timer-work': isWork, 'timer-break': !isWork }">
          <div class="timer-type">{{ isWork ? '专注中' : '休息中' }}</div>
          <div class="timer-time">{{ displayTime }}</div>
          <div class="timer-rounds">已完成 {{ completedRounds }} 轮</div>
        </div>

        <n-progress
          type="circle"
          :percentage="progressPercentage"
          :color="isWork ? '#e74c3c' : '#27ae60'"
          :rail-color="isWork ? '#fde8e8' : '#e8f8e8'"
          :size="120"
          :stroke-width="8"
        />

        <n-space>
          <n-button
            :type="isRunning ? 'warning' : 'primary'"
            size="large"
            @click="toggleTimer"
          >
            {{ isRunning ? '暂停' : '开始' }}
          </n-button>
          <n-button @click="resetTimer" size="large">重置</n-button>
          <n-button v-if="isWork && isRunning" @click="skipToBreak" size="large">跳过到休息</n-button>
        </n-space>

        <n-space>
          <n-input-number v-model:value="workMinutes" :min="1" :max="120" size="small" style="width: 80px">
            <template #prefix>专注</template>
            <template #suffix>分</template>
          </n-input-number>
          <n-input-number v-model:value="breakMinutes" :min="1" :max="60" size="small" style="width: 80px">
            <template #prefix>休息</template>
            <template #suffix>分</template>
          </n-input-number>
        </n-space>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'

const workMinutes = ref(25)
const breakMinutes = ref(5)
const isWork = ref(true)
const isRunning = ref(false)
const completedRounds = ref(0)

const remainingSeconds = ref(workMinutes.value * 60)
let timer: ReturnType<typeof setInterval> | null = null

const displayTime = computed(() => {
  const m = Math.floor(remainingSeconds.value / 60)
  const s = remainingSeconds.value % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
})

const totalSeconds = computed(() => (isWork.value ? workMinutes.value : breakMinutes.value) * 60)

const progressPercentage = computed(() => {
  return Math.round(((totalSeconds.value - remainingSeconds.value) / totalSeconds.value) * 100)
})

const message = useMessage()

function startTimer() {
  if (timer) return
  isRunning.value = true
  timer = setInterval(() => {
    remainingSeconds.value--
    if (remainingSeconds.value <= 0) {
      switchPhase()
    }
  }, 1000)
}

function pauseTimer() {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
  isRunning.value = false
}

function toggleTimer() {
  if (isRunning.value) pauseTimer()
  else startTimer()
}

function resetTimer() {
  pauseTimer()
  isWork.value = true
  remainingSeconds.value = workMinutes.value * 60
}

function switchPhase() {
  pauseTimer()
  if (isWork.value) {
    message.success('🎉 专注时间结束！休息一下吧')
    isWork.value = false
    remainingSeconds.value = breakMinutes.value * 60
    completedRounds.value++
  } else {
    message.info('休息结束，开始新一轮专注')
    isWork.value = true
    remainingSeconds.value = workMinutes.value * 60
  }
  startTimer()
}

function skipToBreak() {
  if (isWork.value && isRunning.value) {
    switchPhase()
  }
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.timer-display {
  text-align: center;
  padding: 32px;
  border-radius: 16px;
  min-width: 280px;
}
.timer-work { background: linear-gradient(135deg, #fde8e8, #f8d7da); }
.timer-break { background: linear-gradient(135deg, #e8f8e8, #d4edda); }
.timer-type { font-size: 18px; font-weight: 600; margin-bottom: 8px; }
.timer-time { font-size: 64px; font-weight: 700; font-family: 'Courier New', monospace; letter-spacing: 4px; }
.timer-rounds { font-size: 14px; color: #888; margin-top: 8px; }
</style>
