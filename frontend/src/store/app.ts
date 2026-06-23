import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useAppStore = defineStore('app', () => {
  const darkMode = ref(localStorage.getItem('darkMode') === 'true')
  const collapsed = ref(localStorage.getItem('collapsed') === 'true')

  function toggleDarkMode() {
    darkMode.value = !darkMode.value
    localStorage.setItem('darkMode', String(darkMode.value))
  }

  function toggleCollapsed() {
    collapsed.value = !collapsed.value
    localStorage.setItem('collapsed', String(collapsed.value))
  }

  return { darkMode, collapsed, toggleDarkMode, toggleCollapsed }
})
