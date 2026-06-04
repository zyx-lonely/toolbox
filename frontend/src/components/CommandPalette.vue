<template>
  <n-modal :show="show" @update:show="$emit('close')" :mask-closable="true" @after-leave="$emit('close')">
    <n-card style="width: 500px" :bordered="true" size="small">
      <n-input ref="inputRef" v-model:value="query" placeholder="搜索功能..." @keydown="handleKeydown" autofocus />
      <n-list v-if="results.length" class="mt-2" style="max-height: 300px; overflow-y: auto">
        <n-list-item v-for="(r, i) in results" :key="i" clickable @click="navigate(r.key)">
          <n-space><n-icon size="18"><component :is="r.icon" /></n-icon><span>{{ r.label }}</span></n-space>
        </n-list-item>
      </n-list>
      <n-empty v-if="query && !results.length" description="无匹配结果" class="mt-2" />
    </n-card>
  </n-modal>
</template>
<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
defineProps<{ show: boolean }>()
const emit = defineEmits(['close', 'navigate'])
const router = useRouter()
const query = ref(''); const inputRef = ref<any>(null); const results = ref<any[]>([])
const allItems = [
  {label:'首页',key:'/dashboard',icon:'GridOutline',keywords:['home','dashboard']},
  {label:'系统信息',key:'/system',keywords:['system','info','硬件','cpu','内存']},
  {label:'硬件检测报告',key:'/system/report',keywords:['report','硬件','报告','检测']},
  {label:'磁盘清理',key:'/optimize/disk',keywords:['disk','clean','清理','磁盘','临时文件']},
  {label:'启动项管理',key:'/optimize/startup',keywords:['startup','启动','开机自启']},
  {label:'系统服务优化',key:'/optimize/service',keywords:['service','服务','优化']},
  {label:'浏览器数据管理',key:'/optimize/browser',keywords:['browser','浏览器','缓存','cookie']},
  {label:'系统健康体检',key:'/optimize/health',keywords:['health','健康','体检','诊断']},
  {label:'Windows 更新',key:'/optimize/windowsupdate',keywords:['update','更新','windows update']},
  {label:'系统备份与还原',key:'/optimize/restore',keywords:['backup','restore','备份','还原','还原点']},
  {label:'网络工具(Ping/端口/DNS/修复)',key:'/network',keywords:['ping','port','dns','网络','修复']},
  {label:'IP 地理位置',key:'/network/geoip',keywords:['geo','ip','地理位置','查询']},
  {label:'局域网扫描',key:'/network/lanscan',keywords:['lan','scan','局域网','扫描','设备']},
  {label:'WiFi 密码',key:'/network/wifi',keywords:['wifi','密码','wireless']},
  {label:'远程桌面',key:'/network/rdp',keywords:['rdp','远程','remote','桌面']},
  {label:'重复文件查找',key:'/file/duplicate',keywords:['duplicate','重复','文件','查找']},
  {label:'文件批量归类',key:'/file/organize',keywords:['organize','归类','整理','分类']},
  {label:'文件解锁器',key:'/file/unlock',keywords:['unlock','解锁','占用']},
  {label:'密码生成器',key:'/security/password',keywords:['password','密码','生成']},
  {label:'JSON 格式化',key:'/devtools/json',keywords:['json','格式化','校验']},
  {label:'编解码工具',key:'/devtools/codec',keywords:['base64','url','编码','解码']},
  {label:'Cron 生成器',key:'/devtools/cron',keywords:['cron','定时','表达式']},
  {label:'剪贴板历史',key:'/daily/clipboard',keywords:['clipboard','剪贴板','复制','粘贴']},
  {label:'定时任务',key:'/daily/scheduler',keywords:['schedule','定时','任务','计划']},
  {label:'单位换算',key:'/daily/convert',keywords:['convert','换算','单位','转换']},
  {label:'取色器',key:'/daily/colorpicker',keywords:['color','颜色','取色','picker']},
  {label:'截屏工具',key:'/daily/screenshot',keywords:['screenshot','截屏','截图']},
  {label:'番茄时钟',key:'/daily/pomodoro',keywords:['pomodoro','番茄','时钟','计时']},
  {label:'字体管理',key:'/system/fonts',keywords:['font','字体','安装']},
  {label:'设置',key:'/settings',keywords:['settings','设置','配置']},
  {label:'右键菜单',key:'/settings/contextmenu',keywords:['context','菜单','右键']},
]
watch(query, (val) => {
  if (!val) { results.value = []; return }
  const q = val.toLowerCase()
  results.value = allItems.filter(item =>
    item.label.toLowerCase().includes(q) ||
    item.keywords.some(k => k.includes(q))
  ).slice(0, 10)
})
watch(() => show, (val) => { if (val) { query.value = ''; nextTick(() => inputRef.value?.focus()) } })
function handleKeydown(e: KeyboardEvent) { if (e.key === 'Escape') emit('close'); if (e.key === 'Enter' && results.value.length) navigate(results.value[0].key) }
function navigate(key: string) { router.push(key); emit('close') }
</script>
<style scoped>.mt-2{margin-top:8px}</style>
