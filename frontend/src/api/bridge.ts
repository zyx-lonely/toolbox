/** Wails 后端 API 桥接 - 工具函数 */

export function formatBytes(bytes: number | undefined | null): string {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  let val = bytes
  while (val >= 1024 && i < units.length - 1) { val /= 1024; i++ }
  return `${val.toFixed(i === 0 ? 0 : 2)} ${units[i]}`
}

export function riskType(risk: string): string {
  return { low: 'success', medium: 'warning', high: 'error' }[risk] || 'default'
}

export function riskLabel(risk: string): string {
  return { low: '低风险', medium: '中风险', high: '高风险' }[risk] || risk
}

export function startTypeLabel(t: string): string {
  const map: Record<string, string> = { auto: '自动', manual: '手动', disabled: '禁用' }
  return map[t] || t
}
