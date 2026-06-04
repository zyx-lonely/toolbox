<template>
  <div>
    <n-h2>Markdown 预览</n-h2>
    <n-card>
      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="编辑">
            <n-input v-model:value="markdown" type="textarea" rows="20" placeholder="在此输入 Markdown 内容..." />
          </n-form-item>
        </n-gi>
        <n-gi>
          <div class="preview-area">
            <n-h3>预览</n-h3>
            <div class="markdown-body" v-html="renderedHTML"></div>
          </div>
        </n-gi>
      </n-grid>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { marked } from 'marked'

const markdown = ref('# Hello\n\n输入 **Markdown** 内容，右侧实时预览。\n\n- 列表项 1\n- 列表项 2\n- 列表项 3')

const renderedHTML = computed(() => {
  try {
    return marked(markdown.value)
  } catch {
    return '<p style="color:red">解析失败</p>'
  }
})
</script>

<style scoped>
.preview-area {
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 12px;
  min-height: 500px;
  overflow-y: auto;
}
</style>

<style>
.markdown-body h1 { font-size: 1.8em; border-bottom: 1px solid #eee; padding-bottom: 8px; margin: 16px 0 8px; }
.markdown-body h2 { font-size: 1.5em; border-bottom: 1px solid #eee; padding-bottom: 6px; margin: 14px 0 6px; }
.markdown-body h3 { font-size: 1.3em; margin: 12px 0 6px; }
.markdown-body p { margin: 8px 0; line-height: 1.6; }
.markdown-body ul, .markdown-body ol { padding-left: 24px; }
.markdown-body li { margin: 4px 0; }
.markdown-body code { background: #f5f5f5; padding: 2px 6px; border-radius: 3px; font-size: 0.9em; }
.markdown-body pre { background: #f5f5f5; padding: 12px; border-radius: 4px; overflow-x: auto; }
.markdown-body pre code { background: transparent; padding: 0; }
.markdown-body blockquote { border-left: 4px solid #1890ff; padding-left: 12px; color: #666; margin: 8px 0; }
.markdown-body table { border-collapse: collapse; width: 100%; margin: 8px 0; }
.markdown-body th, .markdown-body td { border: 1px solid #e0e0e0; padding: 6px 10px; text-align: left; }
.markdown-body th { background: #fafafa; }
.markdown-body a { color: #1890ff; }
.markdown-body img { max-width: 100%; }
</style>
