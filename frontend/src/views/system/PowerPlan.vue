<template>
  <div>
    <n-h2>电源方案切换</n-h2>
    <n-card>
      <n-table :bordered="true">
        <thead>
          <tr>
            <th>方案名称</th>
            <th>GUID</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="plan in plans" :key="plan.guid">
            <td>{{ plan.name }}</td>
            <td><n-text depth="3" class="guid-text">{{ plan.guid }}</n-text></td>
            <td>
              <n-tag v-if="plan.active" type="success" round>当前使用</n-tag>
              <n-tag v-else type="default" round>未激活</n-tag>
            </td>
            <td>
              <n-button size="small" :disabled="plan.active" :loading="switching === plan.guid" @click="switchPlan(plan.guid)">
                切换
              </n-button>
            </td>
          </tr>
        </tbody>
      </n-table>
      <n-button type="primary" class="mt-4" @click="loadPlans" :loading="loading">刷新</n-button>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetPowerPlans, SetPowerPlan } from '@wails/go/main/App'

interface PowerPlan { guid: string; name: string; active: boolean }

const plans = ref<PowerPlan[]>([])
const loading = ref(false)
const switching = ref('')

async function loadPlans() {
  loading.value = true
  try { plans.value = await GetPowerPlans() } catch (e: any) { console.error(e) }
  loading.value = false
}

async function switchPlan(guid: string) {
  switching.value = guid
  try {
    await SetPowerPlan(guid)
    await loadPlans()
  } catch (e: any) { console.error(e) }
  switching.value = ''
}

onMounted(loadPlans)
</script>

<style scoped>
.guid-text { font-family: monospace; font-size: 12px; }
.mt-4 { margin-top: 16px; }
</style>
