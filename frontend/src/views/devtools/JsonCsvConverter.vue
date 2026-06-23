<template>
  <div>
    <n-h2>JSON / CSV 互转</n-h2>
    <n-p>在 JSON 和 CSV 格式之间互相转换。</n-p>

    <n-space vertical :size="16">
      <n-space>
        <n-button-group>
          <n-button :type="mode === 'json2csv' ? 'primary' : 'default'" @click="mode = 'json2csv'">JSON → CSV</n-button>
          <n-button :type="mode === 'csv2json' ? 'primary' : 'default'" @click="mode = 'csv2json'">CSV → JSON</n-button>
        </n-button-group>
        <n-button @click="loadSample">加载示例</n-button>
      </n-space>

      <n-grid :cols="2" :x-gap="12">
        <n-gi>
          <n-card :title="mode === 'json2csv' ? '输入 JSON' : '输入 CSV'" size="small">
            <n-input
              v-model:value="inputData"
              type="textarea"
              :rows="15"
              :placeholder="mode === 'json2csv' ? '粘贴 JSON 数组...' : '粘贴 CSV 文本...'"
            />
          </n-card>
        </n-gi>
        <n-gi>
          <n-card :title="mode === 'json2csv' ? '输出 CSV' : '输出 JSON'" size="small">
            <n-input
              v-model:value="outputData"
              type="textarea"
              :rows="15"
              placeholder="结果将显示在这里..."
              readonly
            />
          </n-card>
        </n-gi>
      </n-grid>

      <n-space>
        <n-button type="primary" @click="convert">转换</n-button>
        <n-button @click="copyResult">复制结果</n-button>
        <n-button @click="clearAll">清空</n-button>
      </n-space>
    </n-space>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const message = useMessage()
const mode = ref<'json2csv' | 'csv2json'>('json2csv')
const inputData = ref('')
const outputData = ref('')

function loadSample() {
  if (mode.value === 'json2csv') {
    inputData.value = JSON.stringify([
      { name: '张三', age: 25, city: '北京', role: '工程师' },
      { name: '李四', age: 30, city: '上海', role: '设计师' },
      { name: '王五', age: 28, city: '深圳', role: '产品经理' },
    ], null, 2)
  } else {
    inputData.value = 'name,age,city,role\n张三,25,北京,工程师\n李四,30,上海,设计师\n王五,28,深圳,产品经理'
  }
}

function convert() {
  try {
    if (mode.value === 'json2csv') {
      const arr = JSON.parse(inputData.value)
      if (!Array.isArray(arr) || arr.length === 0) {
        message.error('请输入有效的 JSON 数组')
        return
      }
      const headers = Object.keys(arr[0])
      const csvRows = [headers.join(',')]
      for (const row of arr) {
        csvRows.push(headers.map(h => String(row[h] ?? '')).join(','))
      }
      outputData.value = csvRows.join('\n')
      message.success(`转换成功: ${arr.length} 行数据`)
    } else {
      const lines = inputData.value.trim().split('\n')
      if (lines.length < 2) {
        message.error('CSV 至少需要标题行和一行数据')
        return
      }
      const headers = lines[0].split(',').map(h => h.trim())
      const result = []
      for (let i = 1; i < lines.length; i++) {
        const values = lines[i].split(',').map(v => v.trim())
        const obj: any = {}
        headers.forEach((h, idx) => {
          const val = values[idx] ?? ''
          obj[h] = isNaN(Number(val)) || val === '' ? val : Number(val)
        })
        result.push(obj)
      }
      outputData.value = JSON.stringify(result, null, 2)
      message.success(`转换成功: ${result.length} 条记录`)
    }
  } catch (e: any) {
    message.error('转换失败: ' + e.message)
  }
}

function copyResult() {
  if (!outputData.value) {
    message.warning('没有可复制的内容')
    return
  }
  navigator.clipboard.writeText(outputData.value)
  message.success('已复制到剪贴板')
}

function clearAll() {
  inputData.value = ''
  outputData.value = ''
}
</script>
