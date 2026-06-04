import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      '@wails': path.resolve(__dirname, 'wailsjs')
    }
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true
  }
})
