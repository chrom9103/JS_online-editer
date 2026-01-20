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
import { PropType } from 'vue';

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

const runCode = () => {
  const file = getActiveFile();
  if (!file) return;
  const originalLog = console.log;
  terminalOutput.value.push({ text: `> Running ${new Date().toLocaleTimeString()}`, type: 'info' });

  console.log = (...args: any[]) => {
    const message = args.map(arg => typeof arg === 'object' ? JSON.stringify(arg) : String(arg)).join(' ');
    terminalOutput.value.push({ text: message, type: 'log' });
  };

  try {
    const result = eval(file.content);
    if (result !== undefined) {
      terminalOutput.value.push({ text: `=> ${typeof result === 'object' ? JSON.stringify(result) : String(result)}`, type: 'log' });
    }
  } catch (e: any) {
    terminalOutput.value.push({ text: `Error: ${e.toString()}`, type: 'error' });
  } finally {
    console.log = originalLog;
  }
};

</script>
