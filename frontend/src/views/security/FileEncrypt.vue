<template>
  <div>
    <n-h2>文件加密/解密</n-h2>
    <n-p>使用 AES-GCM 算法加密或解密文件。</n-p>
    <n-card>
      <n-space vertical>
        <n-space>
          <n-input-group>
            <n-input v-model:value="filePath" placeholder="选择要处理的文件" readonly style="width: 350px" />
            <n-button @click="selectFile">选择文件</n-button>
          </n-input-group>
        </n-space>
        <n-input v-model:value="password" type="password" placeholder="输入密码" show-password-on="click" style="width: 350px" />
        <n-space>
          <n-button type="primary" @click="encrypt" :loading="encrypting">
            <template #icon><n-icon><lock-closed-outline /></n-icon></template>
            加密文件
          </n-button>
          <n-button type="warning" @click="decrypt" :loading="decrypting">
            <template #icon><n-icon><lock-open-outline /></n-icon></template>
            解密文件
          </n-button>
        </n-space>
        <n-alert v-if="result" :type="result.success ? 'success' : 'error'" closable @close="result = null">
          {{ result.success ? (result.outputPath ? `处理成功: ${result.outputPath}` : '操作成功') : `失败: ${result.error}` }}
        </n-alert>
      </n-space>
    </n-card>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { LockClosedOutline, LockOpenOutline } from '@vicons/ionicons5'
import { EncryptFile, DecryptFile, SelectFile } from '@wails/go/main/App'
const filePath = ref(''); const password = ref(''); const result = ref<any>(null)
const encrypting = ref(false); const decrypting = ref(false)
const message = useMessage()

async function selectFile() { const f = await SelectFile(); if (f) filePath.value = f as string }
async function encrypt() {
  if (!filePath.value || !password.value) { message.warning('请选择文件并输入密码'); return }
  encrypting.value = true; result.value = null
  try { const r = await EncryptFile(filePath.value, password.value); if (r) result.value = r as any }
  catch(e:any) { message.error(String(e)) }; encrypting.value = false
}
async function decrypt() {
  if (!filePath.value || !password.value) { message.warning('请选择文件并输入密码'); return }
  decrypting.value = true; result.value = null
  try { const r = await DecryptFile(filePath.value, password.value); if (r) result.value = r as any }
  catch(e:any) { message.error(String(e)) }; decrypting.value = false
}
</script>
