import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useToolsStore = defineStore('tools', () => {
  // 最近使用的工具（最多 10 个）
  const recentTools = ref<any[]>([])
  // 收藏的工具
  const favoriteTools = ref<any[]>([])

  // 初始化：从 localStorage 读取
  function init() {
    try {
      const recent = localStorage.getItem('pc-toolbox-recent')
      if (recent) recentTools.value = JSON.parse(recent)
      const fav = localStorage.getItem('pc-toolbox-favorites')
      if (fav) favoriteTools.value = JSON.parse(fav)
    } catch (e) {
      console.warn('读取工具历史失败', e)
    }
  }

  // 添加最近使用
  function addRecent(tool: { title: string; path: string; icon?: string }) {
    const existing = recentTools.value.findIndex(t => t.path === tool.path)
    if (existing >= 0) {
      recentTools.value.splice(existing, 1)
    }
    recentTools.value.unshift(tool)
    if (recentTools.value.length > 10) {
      recentTools.value = recentTools.value.slice(0, 10)
    }
    saveRecent()
  }

  // 切换收藏状态
  function toggleFavorite(tool: { title: string; path: string; icon?: string }) {
    const existing = favoriteTools.value.findIndex(t => t.path === tool.path)
    if (existing >= 0) {
      favoriteTools.value.splice(existing, 1)
    } else {
      favoriteTools.value.push(tool)
    }
    saveFavorites()
  }

  // 检查是否已收藏
  function isFavorite(path: string) {
    return favoriteTools.value.some(t => t.path === path)
  }

  // 移除收藏
  function removeFavorite(path: string) {
    const existing = favoriteTools.value.findIndex(t => t.path === path)
    if (existing >= 0) {
      favoriteTools.value.splice(existing, 1)
      saveFavorites()
    }
  }

  function saveRecent() {
    localStorage.setItem('pc-toolbox-recent', JSON.stringify(recentTools.value))
  }

  function saveFavorites() {
    localStorage.setItem('pc-toolbox-favorites', JSON.stringify(favoriteTools.value))
  }

  return {
    recentTools,
    favoriteTools,
    init,
    addRecent,
    toggleFavorite,
    isFavorite,
    removeFavorite,
  }
})
