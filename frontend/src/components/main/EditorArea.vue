<template>
  <div class="main-editor-area">
    <EditorHeader
      :files="files"
      :active-file-id="activeFileId"
      @open-file="(id) => $emit('open-file', id)"
      @request-delete="(id) => $emit('request-delete', id)"
      @run-code="runCode"
    />
    <div ref="editorContentWrapper" class="editor-content-wrapper">
      <MonacoEditor
        ref="monacoEditorRef"
        :content="activeFileContent"
        :height="editorHeight"
        @update-content="onEditorContentChange"
      />
      <ResizeDivider @resize="onResize" @resize-end="onResizeEnd" />
      <TerminalPanel :output="terminalOutput" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick, onUnmounted } from 'vue'
import EditorHeader from './editorArea/EditorHeader.vue'
import MonacoEditor from './editorArea/MonacoEditor.vue'
import ResizeDivider from './editorArea/ResizeDivider.vue'
import TerminalPanel from './editorArea/TerminalPanel.vue'

const props = defineProps<{
  files: Array<{ id: string; name: string; content: string }>
  activeFileId: string
}>()

const emit = defineEmits(['open-file', 'request-delete', 'update-file-content'])

const monacoEditorRef = ref<InstanceType<typeof MonacoEditor> | null>(null)
const editorContentWrapper = ref<HTMLElement | null>(null)
const terminalOutput = ref<{ text: string; type: string }[]>([])
const clientId = ref<string>('')

const editorHeight = ref<number>(0)
const editorRatio = ref<number>(0.65)

const MIN_EDITOR_HEIGHT = 120
const DIVIDER_HEIGHT = 6
const MIN_TERMINAL_HEIGHT = 80

// アクティブファイルのコンテンツを計算
const activeFileContent = computed(() => {
  const file = props.files.find((f) => f.id === props.activeFileId)
  return file?.content ?? ''
})

const getActiveFile = () => props.files.find((f) => f.id === props.activeFileId)

const updateEditorOnResize = () => {
  const wrapper = editorContentWrapper.value
  if (!wrapper) return
  const total = wrapper.clientHeight
  let newHeight = Math.round(total * editorRatio.value)
  const maxPossible = total - MIN_TERMINAL_HEIGHT - DIVIDER_HEIGHT
  if (newHeight < MIN_EDITOR_HEIGHT) newHeight = MIN_EDITOR_HEIGHT
  if (newHeight > maxPossible) newHeight = maxPossible
  editorHeight.value = newHeight
  nextTick(() => {
    monacoEditorRef.value?.layout()
  })
}

onMounted(() => {
  // ClientIDをlocalStorageから読み出し、なければ生成して保存
  try {
    const stored = localStorage.getItem('clientId')
    if (stored) {
      clientId.value = stored
    } else if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
      const id = crypto.randomUUID()
      localStorage.setItem('clientId', id)
      clientId.value = id
    } else {
      // フォールバック: 簡易ランダム文字列
      const id = 'id-' + Math.random().toString(36).slice(2, 10)
      localStorage.setItem('clientId', id)
      clientId.value = id
    }
  } catch (e) {
    // localStorageが使えない場合は空のままにする
    clientId.value = ''
  }
})

onMounted(() => {
  // 初期高さをコンテナの65%に設定
  nextTick(() => {
    const wrapper = editorContentWrapper.value
    if (wrapper) {
      editorHeight.value = Math.max(MIN_EDITOR_HEIGHT, Math.floor(wrapper.clientHeight * 0.65))
      editorRatio.value = editorHeight.value / wrapper.clientHeight
    }
    window.addEventListener('resize', updateEditorOnResize)
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', updateEditorOnResize)
})

// エディタのコンテンツ変更ハンドラ
const onEditorContentChange = (content: string) => {
  const file = getActiveFile()
  if (file) {
    emit('update-file-content', { id: file.id, content })
  }
}

// リサイズハンドラ
const onResize = (deltaY: number) => {
  const wrapper = editorContentWrapper.value
  if (!wrapper) return
  const rect = wrapper.getBoundingClientRect()
  let newHeight = editorHeight.value + deltaY
  const maxPossible = rect.height - MIN_TERMINAL_HEIGHT - DIVIDER_HEIGHT
  if (newHeight < MIN_EDITOR_HEIGHT) newHeight = MIN_EDITOR_HEIGHT
  if (newHeight > maxPossible) newHeight = maxPossible
  editorHeight.value = newHeight
  editorRatio.value = editorHeight.value / rect.height
  monacoEditorRef.value?.layout()
}

const onResizeEnd = () => {
  monacoEditorRef.value?.layout()
}

// API実行中フラグ
const isExecuting = ref(false)

// APIのベースURL（環境変数または相対パス）
const getApiBaseUrl = () => {
  // 本番環境では /playground/api を使用
  if (import.meta.env.PROD) {
    return '/playground/api'
  }
  // 開発環境ではローカルのバックエンド
  return '/api'
}

const runCode = async () => {
  const file = getActiveFile()
  if (!file) return

  if (isExecuting.value) {
    terminalOutput.value.push({ text: 'Already executing...', type: 'warn' })
    return
  }

  isExecuting.value = true
  terminalOutput.value.push({ text: `> Running ${new Date().toLocaleTimeString()}`, type: 'info' })

  try {
    const response = await fetch(`${getApiBaseUrl()}/execute`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        code: file.content,
        language: 'javascript',
        clientId: clientId.value,
      }),
    })

    const result = await response.json()

    if (result.output && Array.isArray(result.output)) {
      result.output.forEach((item: { type: string; text: string }) => {
        terminalOutput.value.push({ text: item.text, type: item.type })
      })
    }

    if (!result.success && result.error) {
      terminalOutput.value.push({ text: `Error: ${result.error}`, type: 'error' })
    }
  } catch (e: any) {
    terminalOutput.value.push({
      text: `Network Error: ${e.message || 'Failed to connect to server'}`,
      type: 'error',
    })
  } finally {
    isExecuting.value = false
  }
}

defineExpose({ runCode })
</script>

<style scoped>
.main-editor-area {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.editor-content-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  min-height: 0;
  overflow: hidden;
}
</style>
