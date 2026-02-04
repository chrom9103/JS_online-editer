<template>
  <div class="main-editor-area">
    <div class="editor-header">
      <div class="tab-container">
        <div
          v-for="file in files"
          :key="file.id"
          :class="{ 'editor-tab': true, active: file.id === activeFileId }"
          @click="$emit('open-file', file.id)"
        >
          {{ file.name }}
          <span class="close-tab-btn" @click.stop="$emit('request-delete', file.id)">×</span>
        </div>
      </div>
      <div class="editor-actions">
        <button @click="runCode" class="action-icon-btn">▷</button>
      </div>
    </div>
    <div ref="editorContentWrapper" class="editor-content-wrapper">
      <div
        ref="editorContainer"
        id="editor-container"
        :style="{ height: editorHeight + 'px' }"
      ></div>
      <div class="v-divider" @mousedown.prevent="startVerticalResize"></div>
      <div class="terminal-area">
        <div class="terminal-header">
          <span class="terminal-title">TERMINAL</span>
        </div>
        <div class="output-container" ref="terminalContainer">
          <p v-for="(line, index) in terminalOutput" :key="index" :class="line.type">
            {{ line.text }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue';
import * as monaco from 'monaco-editor';
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker';

// Workerのベースパスを明示的に設定
self.MonacoEnvironment = {
  getWorker(_: unknown, label: string) {
    if (label === 'typescript' || label === 'javascript') {
      return new tsWorker();
    }
    return new editorWorker();
  },
  // 本番環境でのworkerパス解決用
  baseUrl: import.meta.env.BASE_URL,
};

const props = defineProps<{
  files: Array<{ id: string; name: string; content: string }>;
  activeFileId: string;
}>();

const emit = defineEmits(['open-file', 'request-delete', 'update-file-content']);

let monacoEditor: monaco.editor.IStandaloneCodeEditor | null = null;
const editorContainer = ref<HTMLElement | null>(null);
const editorContentWrapper = ref<HTMLElement | null>(null);
const terminalContainer = ref<HTMLElement | null>(null);
const terminalOutput = ref<{ text: string; type: string }[]>([]);
const clientId = ref<string>('');

const editorHeight = ref<number>(0);
const editorRatio = ref<number>(0.65);
let isVertResizing = false;

const MIN_EDITOR_HEIGHT = 120;

const updateEditorOnResize = () => {
  const wrapper = editorContentWrapper.value;
  if (!wrapper) return;
  const total = wrapper.clientHeight;
  let newHeight = Math.round(total * editorRatio.value);
  const maxPossible = total - MIN_EDITOR_HEIGHT - 6;
  if (newHeight < MIN_EDITOR_HEIGHT) newHeight = MIN_EDITOR_HEIGHT;
  if (newHeight > maxPossible) newHeight = maxPossible;
  editorHeight.value = newHeight;
  if (monacoEditor) monacoEditor.layout();
}

const getActiveFile = () => props.files.find((f) => f.id === props.activeFileId);

onMounted(() => {
  // ClientIDをlocalStorageから読み出し、なければ生成して保存
  try {
    const stored = localStorage.getItem('clientId');
    if (stored) {
      clientId.value = stored;
    } else if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
      const id = crypto.randomUUID();
      localStorage.setItem('clientId', id);
      clientId.value = id;
    } else {
      // フォールバック: 簡易ランダム文字列
      const id = 'id-' + Math.random().toString(36).slice(2, 10);
      localStorage.setItem('clientId', id);
      clientId.value = id;
    }
  } catch (e) {
    // localStorageが使えない場合は空のままにする
    clientId.value = '';
  }
})

onMounted(() => {
  // 初期高さをコンテナの65%に設定
  nextTick(() => {
    const wrapper = editorContentWrapper.value;
    if (wrapper) {
      editorHeight.value = Math.max(MIN_EDITOR_HEIGHT, Math.floor(wrapper.clientHeight * 0.65));
      editorRatio.value = editorHeight.value / wrapper.clientHeight;
    };
    initMonaco();
    if (monacoEditor) monacoEditor.layout();
    window.addEventListener('mousemove', onVerticalDragging as any);
    window.addEventListener('mouseup', stopVerticalResize as any);
    window.addEventListener('resize', updateEditorOnResize);
  });
});

// 初期化を高さ確定後に行う
const initMonaco = () => {
  if (monacoEditor) return;
  if (!editorContainer.value) return;
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
};

onUnmounted(() => {
  if (monacoEditor) {
    monacoEditor.dispose();
    monacoEditor = null;
  }
  window.removeEventListener('mousemove', onVerticalDragging as any);
  window.removeEventListener('mouseup', stopVerticalResize as any);
  window.removeEventListener('resize', updateEditorOnResize);
});

watch(() => props.activeFileId, (newId) => {
    if (!monacoEditor) return;
    const file = props.files.find((f) => f.id === newId);
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
  // 本番環境では /playground/api を使用
  if (import.meta.env.PROD) {
    return '/playground/api';
  }
  // 開発環境ではローカルのバックエンド
  return '/api';
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
        clientId: clientId.value,
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
}

const startVerticalResize = (event: MouseEvent) => {
  if (event.button !== 0) return;
  isVertResizing = true;
  document.body.style.userSelect = 'none';
  document.body.style.cursor = 'row-resize';
}

const onVerticalDragging = (event: MouseEvent) => {
  if (!isVertResizing) return;
  const wrapper = editorContentWrapper.value;
  if (!wrapper) return;
  const rect = wrapper.getBoundingClientRect();
  let newHeight = event.clientY - rect.top;
  const maxPossible = rect.height - MIN_EDITOR_HEIGHT - 6; // reserve for divider
  if (newHeight < MIN_EDITOR_HEIGHT) newHeight = MIN_EDITOR_HEIGHT;
  if (newHeight > maxPossible) newHeight = maxPossible;
  editorHeight.value = newHeight;
  editorRatio.value = editorHeight.value / rect.height;
  if (monacoEditor) monacoEditor.layout();
}

const stopVerticalResize = () => {
  if (!isVertResizing) return;
  isVertResizing = false;
  document.body.style.userSelect = '';
  document.body.style.cursor = '';
  if (monacoEditor) monacoEditor.layout();
};

defineExpose({ runCode });

</script>
