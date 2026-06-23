// 统一 API 响应类型
export interface APIResponse<T = unknown> {
  success: boolean
  data?: T
  error?: string
}

// 剪贴板条目
export interface ClipItem {
  id: number
  content: string
  type: 'text' | 'image'
  time: string
  size: number
}

// 文件信息
export interface FileItem {
  path: string
  name: string
  size: number
  modified: string
  isDir: boolean
}

// 计划任务
export interface TaskItem {
  id: string
  name: string
  cron: string
  nextRun: string
  enabled: boolean
}

// 正则匹配结果
export interface MatchItem {
  index: number
  text: string
  groups?: string[]
}

// 文件差异
export interface DiffItem {
  line: number
  type: 'equal' | 'add' | 'remove'
  oldText: string
  newText: string
}

// 文件夹大小信息
export interface FolderItem {
  path: string
  size: number
  fileCount: number
  pct: number
}

// 重复文件组
export interface DuplicateGroup {
  hash: string
  files: string[]
  totalSize: number
  wastedSize: number
}

// 系统信息
export interface SystemInfo {
  cpu?: CpuInfo
  memory?: MemoryInfo
  disk?: DiskInfo[]
  network?: NetworkInfo
}

export interface CpuInfo {
  model: string
  cores: number
  usage: number
}

export interface MemoryInfo {
  total: number
  used: number
  free: number
  usage: number
}

export interface DiskInfo {
  path: string
  total: number
  used: number
  free: number
  usage: number
}

export interface NetworkInfo {
  bytesSent: number
  bytesRecv: number
}

// 进程信息
export interface ProcessItem {
  pid: number
  name: string
  cpu: number
  memory: number
  status: string
}

// 启动项
export interface StartupItem {
  name: string
  command: string
  location: string
  enabled: boolean
}
