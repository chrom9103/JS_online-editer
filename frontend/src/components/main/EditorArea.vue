<template>
  <div class="main-editor-area">
    <div class="editor-header">
      <div class="tab-container">
        <div v-for="file in files" :key="file.id" :class="{ 'editor-tab': true, active: file.id === activeFileId }" @click="$emit('open-file', file.id)">
          {{ file.name }}
          <span class="close-tab-btn" @click.stop="$emit('request-delete', file.id)">×</span>
        </div>
      </div>
      <div class="editor-actions">
        <button @click="runCode" class="action-icon-btn">▷</button>
      </div>
    </div>
    <div class="editor-content-wrapper">
      <div ref="editorContainer" id="editor-container"></div>
      <div class="terminal-area">
        <div class="terminal-header">
          <span class="terminal-title">TERMINAL</span>
        </div>
        <div class="output-container" ref="terminalContainer">
          <p v-for="(line, index) in terminalOutput" :key="index" :class="line.type">{{ line.text }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue';
import * as monaco from 'monaco-editor';

const props = defineProps<{
  files: Array<{ id: string; name: string; content: string }>;
  activeFileId: string;
}>();

const emit = defineEmits(['open-file', 'request-delete', 'update-file-content']);

let monacoEditor: monaco.editor.IStandaloneCodeEditor | null = null;
const editorContainer = ref<HTMLElement | null>(null);
const terminalContainer = ref<HTMLElement | null>(null);
const terminalOutput = ref<{ text: string; type: string }[]>([]);

const getActiveFile = () => props.files.find(f => f.id === props.activeFileId);

onMounted(() => {
  (window as any).MonacoEnvironment = {
    getWorkerUrl: function (_moduleId: string, label: string) {
      if (label === 'typescript' || label === 'javascript') {
        return new URL('monaco-editor/esm/vs/language/typescript/ts.worker.js', import.meta.url).toString();
      }
      return new URL('monaco-editor/esm/vs/editor/editor.worker.js', import.meta.url).toString();
    }
  };

  monacoEditor = monaco.editor.create(editorContainer.value!, {
    value: getActiveFile()?.content ?? '',
    language: 'javascript',
    theme: 'vs-dark',
    automaticLayout: true,
  });

  monacoEditor.onDidChangeModelContent(() => {
    const value = monacoEditor!.getValue();
    const file = getActiveFile();
    if (file) {
      emit('update-file-content', { id: file.id, content: value });
    }
  });
});

onUnmounted(() => {
  if (monacoEditor) {
    monacoEditor.dispose();
    monacoEditor = null;
  }
});

watch(() => props.activeFileId, (newId) => {
  if (!monacoEditor) return;
  const file = props.files.find(f => f.id === newId);
  if (file) monacoEditor.setValue(file.content);
});

watch(() => props.files, () => {
  nextTick(() => {
    if (terminalContainer.value) terminalContainer.value.scrollTop = terminalContainer.value.scrollHeight;
  });
}, { deep: true });

// API実行中フラグ
const isExecuting = ref(false);

// APIのベースURL（環境変数または相対パス）
const getApiBaseUrl = () => {
  // 本番環境では同一ドメインのAPIを使用
  if (import.meta.env.PROD) {
    return '/api';
  }
  // 開発環境ではローカルのバックエンドを使用
  return import.meta.env.VITE_API_URL || 'http://localhost:8081/api';
};

const runCode = async () => {
  const file = getActiveFile();
  if (!file) return;
  
  if (isExecuting.value) {
    terminalOutput.value.push({ text: 'Already executing...', type: 'warn' });
    return;
  }

  isExecuting.value = true;
  terminalOutput.value.push({ text: `> Running ${new Date().toLocaleTimeString()}`, type: 'info' });

  try {
    const response = await fetch(`${getApiBaseUrl()}/execute`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        code: file.content,
        language: 'javascript',
      }),
    });

    const result = await response.json();

    if (result.output && Array.isArray(result.output)) {
      result.output.forEach((item: { type: string; text: string }) => {
        terminalOutput.value.push({ text: item.text, type: item.type });
      });
    }

    if (!result.success && result.error) {
      terminalOutput.value.push({ text: `Error: ${result.error}`, type: 'error' });
    }
  } catch (e: any) {
    terminalOutput.value.push({ text: `Network Error: ${e.message || 'Failed to connect to server'}`, type: 'error' });
  } finally {
    isExecuting.value = false;
    nextTick(() => {
      if (terminalContainer.value) terminalContainer.value.scrollTop = terminalContainer.value.scrollHeight;
    });
  }
};

</script>
