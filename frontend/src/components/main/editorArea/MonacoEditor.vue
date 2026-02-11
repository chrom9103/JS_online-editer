<template>
  <div
    ref="editorContainer"
    class="monaco-editor-container"
    :style="{ height: height + 'px' }"
  ></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from 'vue'
import * as monaco from 'monaco-editor'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'

// Workerのベースパスを明示的に設定
self.MonacoEnvironment = {
  getWorker(_: unknown, label: string) {
    if (label === 'typescript' || label === 'javascript') {
      return new tsWorker()
    }
    return new editorWorker()
  },
  baseUrl: import.meta.env.BASE_URL,
}

const props = withDefaults(
  defineProps<{
    content: string
    height: number
    readOnly?: boolean
  }>(),
  {
    readOnly: false,
  },
)

const emit = defineEmits<{
  (e: 'update-content', content: string): void
}>()

let monacoEditor: monaco.editor.IStandaloneCodeEditor | null = null
const editorContainer = ref<HTMLElement | null>(null)

onMounted(() => {
  if (!editorContainer.value) return

  monacoEditor = monaco.editor.create(editorContainer.value, {
    value: props.content,
    language: 'javascript',
    theme: 'vs-dark',
    automaticLayout: true,
    readOnly: props.readOnly,
  })

  monacoEditor.onDidChangeModelContent(() => {
    const value = monacoEditor!.getValue()
    emit('update-content', value)
  })
})

onUnmounted(() => {
  if (monacoEditor) {
    monacoEditor.dispose()
    monacoEditor = null
  }
})

watch(
  () => props.content,
  (newContent, oldContent) => {
    if (!monacoEditor) return
    // 外部からのコンテンツ変更時のみ反映（自身の変更は無視）
    if (monacoEditor.getValue() !== newContent) {
      monacoEditor.setValue(newContent)
    }
  },
)

const layout = () => {
  if (monacoEditor) monacoEditor.layout()
}

defineExpose({ layout })
</script>

<style scoped>
.monaco-editor-container {
  flex-shrink: 0;
}
</style>
