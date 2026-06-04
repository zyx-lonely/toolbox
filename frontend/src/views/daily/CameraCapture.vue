<template>
  <div>
    <n-h2>摄像头拍照</n-h2>
    <n-p>调用摄像头拍照并上传到指定服务器。</n-p>

    <n-card>
      <n-space vertical>
        <n-alert v-if="!cameraReady && !errorMsg" type="info">
          请先点击「启动摄像头」按钮授权摄像头访问权限。
        </n-alert>

        <div v-if="cameraReady" class="video-wrapper">
          <video ref="videoRef" autoplay playsinline muted></video>
        </div>

        <div v-if="capturedImage" class="video-wrapper">
          <img :src="capturedImage" alt="拍摄的照片" />
        </div>

        <n-space>
          <n-button v-if="!cameraReady" @click="startCamera" :loading="starting" type="primary">
            <template #icon><n-icon><camera-outline /></n-icon></template>
            启动摄像头
          </n-button>
          <n-button v-if="cameraReady && !capturedImage" @click="capture" type="success">
            <template #icon><n-icon><camera-outline /></n-icon></template>
            拍照
          </n-button>
          <n-button v-if="capturedImage" @click="retake" type="info">
            重新拍照
          </n-button>
          <n-button v-if="cameraReady && !capturedImage" @click="stopCamera" type="warning">
            关闭摄像头
          </n-button>
        </n-space>

        <n-divider v-if="capturedImage" />

        <n-card v-if="capturedImage" title="上传设置" size="small">
          <n-space vertical>
            <n-input v-model:value="serverURL" placeholder="服务器上传地址（如 http://example.com/upload）" />
            <n-input v-model:value="fieldName" placeholder="表单字段名（默认 file）" />
            <n-button @click="upload" :loading="uploading" type="primary" :disabled="!serverURL">
              <template #icon><n-icon><cloud-upload-outline /></n-icon></template>
              上传到服务器
            </n-button>

            <n-alert v-if="uploadResult" :type="uploadResult.success ? 'success' : 'error'" closable @close="uploadResult = null">
              <template #header>{{ uploadResult.success ? '上传成功' : '上传失败' }}</template>
              <div>状态码: {{ uploadResult.statusCode }}</div>
              <div v-if="uploadResult.response">响应: {{ uploadResult.response }}</div>
              <div v-if="uploadResult.error">错误: {{ uploadResult.error }}</div>
            </n-alert>
          </n-space>
        </n-card>

        <n-alert v-if="errorMsg" type="error" closable @close="errorMsg = ''">
          {{ errorMsg }}
        </n-alert>
      </n-space>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import { CameraOutline, CloudUploadOutline } from '@vicons/ionicons5'
import { UploadFileToServer } from '@wails/go/main/App'

const videoRef = ref<HTMLVideoElement | null>(null)
const cameraReady = ref(false)
const starting = ref(false)
const capturedImage = ref('')
const serverURL = ref('')
const fieldName = ref('file')
const uploading = ref(false)
const uploadResult = ref<any>(null)
const errorMsg = ref('')
const message = useMessage()

let mediaStream: MediaStream | null = null

async function startCamera() {
  starting.value = true
  errorMsg.value = ''
  try {
    mediaStream = await navigator.mediaDevices.getUserMedia({
      video: { width: 640, height: 480, facingMode: 'user' }
    })
    if (videoRef.value) {
      videoRef.value.srcObject = mediaStream
    }
    cameraReady.value = true
  } catch (e: any) {
    if (e.name === 'NotAllowedError') {
      errorMsg.value = '摄像头权限被拒绝，请在系统设置中允许摄像头访问'
    } else if (e.name === 'NotFoundError') {
      errorMsg.value = '未检测到摄像头设备'
    } else {
      errorMsg.value = `启动摄像头失败: ${e.message || e}`
    }
  }
  starting.value = false
}

function stopCamera() {
  if (mediaStream) {
    mediaStream.getTracks().forEach(t => t.stop())
    mediaStream = null
  }
  cameraReady.value = false
}

function capture() {
  const video = videoRef.value
  if (!video) return

  const canvas = document.createElement('canvas')
  canvas.width = video.videoWidth || 640
  canvas.height = video.videoHeight || 480
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
  capturedImage.value = canvas.toDataURL('image/jpeg', 0.85)
  // 拍照后自动关闭摄像头以释放资源
  stopCamera()
}

function retake() {
  capturedImage.value = ''
  uploadResult.value = null
}

async function upload() {
  if (!capturedImage.value || !serverURL.value) return
  uploading.value = true
  uploadResult.value = null
  try {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-')
    const fileName = `photo_${timestamp}.jpg`
    const r = await UploadFileToServer(capturedImage.value, fileName, serverURL.value, fieldName.value)
    if (r) uploadResult.value = r as any
    if (r?.success) message.success('上传成功')
  } catch (e: any) {
    uploadResult.value = { success: false, error: String(e), statusCode: 0 }
  }
  uploading.value = false
}

onUnmounted(() => {
  stopCamera()
})
</script>

<style scoped>
.video-wrapper {
  background: #000;
  border-radius: 8px;
  overflow: hidden;
  max-width: 640px;
}
.video-wrapper video,
.video-wrapper img {
  display: block;
  width: 100%;
  height: auto;
}
</style>
