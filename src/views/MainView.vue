<template>
  <div class="container" :style="{ gridTemplateColumns: `50px ${sidebarWidth}px 6px 1fr` }">
    <div class="activity-bar">
      <div class="activity-icon-group">
        <button @click="switchSidebar('explorer')" class="activity-icon-item" :class="{ active: sidebarState.explorer }" title="エクスプローラー">
          <svg viewBox="0 0 16 16" fill="currentColor" class="icon">
            <path d="M4 1h8a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2zM4 3v10h8V3H4zm2 2h4v2H6V5zm0 4h4v2H6V9z"/>
          </svg>
        </button>
        <button @click="switchSidebar('search')" class="activity-icon-item" :class="{ active: sidebarState.search }" title="検索">
          <svg viewBox="0 0 24 24" fill="currentColor" class="icon">
            <path d="M15.5 14h-.79l-.28-.27A6.471 6.471 0 0 0 16 9.5 6.5 6.5 0 1 0 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 5 1.5-1.5-5-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
          </svg>
        </button>
        <button @click="switchSidebar('git')" class="activity-icon-item" :class="{ active: sidebarState.git }" title="ソース管理">
          <svg viewBox="0 0 16 16" fill="currentColor" class="icon">
            <path d="M13 10c0 1.66-1.34 3-3 3s-3-1.34-3-3 1.34-3 3-3 3 1.34 3 3zM6 8a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zm4.5 0a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zM14 6H2a1 1 0 0 0-1 1v2a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V7a1 1 0 0 0-1-1z"/>
          </svg>
        </button>
        <button @click="switchSidebar('runcode')" class="activity-icon-item" :class="{ active: sidebarState.runcode }" title="実行とデバッグ">
          <svg viewBox="0 0 24 24" fill="currentColor" class="icon">
            <path d="M8 5v14l11-7z"/>
          </svg>
        </button>
      </div>
    </div>

    <div class="side-bar">
      <div v-if="sidebarState.explorer" class="file-explorer">
        <div class="file-header">
          <h3>エクスプローラー</h3>
        </div>
        <div class="files-title">
          <p>∨ JS_PlayGround</p>
          <button @click="addNewFile" class="add-file-btn">＋</button>
        </div>
        <ul class="file-list">
          <li
            v-for="file in files"
            :key="file.id"
            :class="{ active: file.id === activeFileId }"
            @click="switchFile(file.id)"
          >
            <span
              v-if="editingFileId !== file.id"
              @dblclick.stop="startRename(file.id)"
            >
              {{ file.name }}
            </span>
            <input
              v-else
              type="text"
              v-model="file.name"
              @blur="finishRename"
              @keyup.enter="finishRename"
              v-rename-focus
            />
            <button @click.stop="deleteFile(file.id)" class="delete-btn">×</button>
          </li>
        </ul>
      </div>

      <div v-else-if="sidebarState.search" class="sidebar-content">
        <h3>検索</h3>
        <p>この機能は現在開発中です。</p>
      </div>

      <div v-else-if="sidebarState.runcode" class="sidebar-content">
        <h3>実行とデバッグ</h3>
        <button @click="runCode" class="action-btn">Run Code</button>
        <button @click="downloadCode" class="action-btn">Download Code</button>
      </div>

      <div v-else-if="sidebarState.git" class="sidebar-content">
        <h3>ソース管理</h3>
        <p>この機能は現在開発中です。</p>
      </div>
    </div>

    <div class="divider" @mousedown.prevent="startResize"></div>

    <div class="main-editor-area">
      <div class="editor-header">
        <div class="tab-container">
          <div v-for="file in files" :key="file.id" :class="{ 'editor-tab': true, active: file.id === activeFileId }" @click="switchFile(file.id)">
            {{ file.name }}
            <span class="close-tab-btn" @click.stop="deleteFile(file.id)">×</span>
          </div>
        </div>
        <div class="editor-actions">
          <button @click="runCode" class="action-icon-btn">▷</button>
        </div>
      </div>
      <div class="editor-content-wrapper">
        <div id="editor-container"></div>
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

    <div v-if="showConfirmPopup" class="confirm-overlay">
      <div class="confirm-popup">
        <p>{{ popupMessage }}</p>
        <div class="popup-actions">
          <button @click="confirmAction(true)" class="confirm-ok-btn">OK</button>
          <button @click="confirmAction(false)" class="confirm-cancel-btn">キャンセル</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue';
import * as monaco from 'monaco-editor';
import { v4 as uuidv4 } from 'uuid';

const files = ref([
  {
    id: uuidv4(),
    name: 'main.js',
    content: `console.log('Hello world!');`,
  },
]);
const activeFileId = ref(files.value[0].id);
const code = ref(files.value[0].content);
const terminalOutput = ref<{ text: string; type: string }[]>([]);
const editingFileId = ref<string | null>(null);
let monacoEditor: monaco.editor.IStandaloneCodeEditor;

const showConfirmPopup = ref(false);
const popupMessage = ref('');
let popupResolve: ((value: boolean | PromiseLike<boolean>) => void) | null = null;
const terminalContainer = ref<HTMLElement | null>(null);

// サイドバーの状態を管理するオブジェクト
const sidebarState = ref({
  explorer: true,
  search: false,
  runcode: false,
  git: false,
});

const sidebarWidth = ref(250);
const minSidebarWidth = 158;
const maxSidebarWidth = 600;
let isResizing = false;

const startResize = (event: MouseEvent) => {
  if (event.button !== 0) return;
  isResizing = true;
  document.body.style.userSelect = 'none';
  document.body.style.cursor = 'col-resize';
  window.addEventListener('mousemove', onDragging);
  window.addEventListener('mouseup', stopResize);
};

const onDragging = (event: MouseEvent) => {
  if (!isResizing) return;
  const pageX = event.pageX;
  let newWidth = pageX - 0 - 50;
  if (newWidth < minSidebarWidth) newWidth = minSidebarWidth;
  if (newWidth > maxSidebarWidth) newWidth = maxSidebarWidth;
  sidebarWidth.value = newWidth;
};

const stopResize = () => {
  if (!isResizing) return;
  isResizing = false;
  document.body.style.userSelect = '';
  document.body.style.cursor = '';
  window.removeEventListener('mousemove', onDragging);
  window.removeEventListener('mouseup', stopResize);
};

onUnmounted(() => {
  window.removeEventListener('mousemove', onDragging);
  window.removeEventListener('mouseup', stopResize);
});


// サイドバーを切り替える関数
const switchSidebar = (view: 'explorer' | 'search' | 'runcode' | 'git') => {
  // すべてのビューを非アクティブにする
  sidebarState.value.explorer = false;
  sidebarState.value.search = false;
  sidebarState.value.runcode = false;
  sidebarState.value.git = false;

  // 選択されたビューをアクティブにする
  if (sidebarState.value[view] === false) {
    sidebarState.value[view] = true;
  }
};

watch(activeFileId, (newFileId, oldFileId) => {
  if (monacoEditor && oldFileId) {
    const oldFile = files.value.find(f => f.id === oldFileId);
    if (oldFile) {
      oldFile.content = monacoEditor.getValue();
    }
  }
  const newFile = files.value.find(f => f.id === newFileId);
  if (newFile && monacoEditor) {
    monacoEditor.setValue(newFile.content);
  }
});

onMounted(() => {
  monacoEditor = monaco.editor.create(document.getElementById('editor-container')!, {
    value: code.value,
    language: 'javascript',
    theme: 'vs-dark',
    automaticLayout: true,
  });

  monacoEditor.onDidChangeModelContent(() => {
    code.value = monacoEditor.getValue();
    const activeFile = files.value.find(f => f.id === activeFileId.value);
    if (activeFile) {
      activeFile.content = code.value;
    }
  });
});

watch(terminalOutput, () => {
  nextTick(() => {
    if (terminalContainer.value) {
      terminalContainer.value.scrollTop = terminalContainer.value.scrollHeight;
    }
  });
}, { deep: true });

const runCode = () => {
  const originalLog = console.log;

  terminalOutput.value.push({ text: `> Running ${new Date().toLocaleTimeString()}`, type: 'info' });

  console.log = (...args) => {
    const message = args.map(arg => typeof arg === 'object' ? JSON.stringify(arg) : String(arg)).join(' ');
    terminalOutput.value.push({ text: message, type: 'log' });
  };

  try {
    const result = eval(code.value);
    if (result !== undefined) {
      terminalOutput.value.push({ text: `=> ${typeof result === 'object' ? JSON.stringify(result) : String(result)}`, type: 'log' });
    }
  } catch (e: any) {
    terminalOutput.value.push({ text: `Error: ${e.toString()}`, type: 'error' });
  } finally {
    console.log = originalLog;
  }
};

const downloadCode = () => {
  const activeFile = files.value.find(f => f.id === activeFileId.value);
  if (!activeFile) return;

  const blob = new Blob([activeFile.content], { type: 'text/javascript' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = activeFile.name;
  a.click();
  URL.revokeObjectURL(url);
};

const addNewFile = () => {
  const newFile = {
    id: uuidv4(),
    name: `new-file-${files.value.length + 1}.js`,
    content: `console.log('Hello world!');`,
  };
  files.value.push(newFile);
  switchFile(newFile.id);
};

const switchFile = (fileId: string) => {
  if (editingFileId.value === fileId) return;
  activeFileId.value = fileId;
};

const startRename = (fileId: string) => {
  editingFileId.value = fileId;
};

const finishRename = () => {
  const file = files.value.find(f => f.id === editingFileId.value);
  if (file && file.name.trim() === '') {
    file.name = 'untitled.js';
  }
  editingFileId.value = null;
};

const deleteFile = async (fileId: string) => {
  if (files.value.length <= 1) {
    alert('少なくとも1つのファイルが必要です。');
    return;
  }
  const fileToDelete = files.value.find(f => f.id === fileId);
  if (!fileToDelete) return;

  popupMessage.value = `本当にファイル「${fileToDelete.name}」を削除してもよろしいですか？`;
  showConfirmPopup.value = true;
  const confirmed = await new Promise<boolean>(resolve => {
    popupResolve = resolve;
  });

  if (confirmed) {
    const fileIndex = files.value.findIndex(f => f.id === fileId);
    if (fileIndex > -1) {
      files.value.splice(fileIndex, 1);
      if (activeFileId.value === fileId) {
        activeFileId.value = files.value[0].id;
      }
    }
  }
};

const confirmAction = (result: boolean) => {
  showConfirmPopup.value = false;
  if (popupResolve) {
    popupResolve(result);
    popupResolve = null;
  }
};

const vRenameFocus = {
  mounted: (el: HTMLElement) => {
    el.focus();
  },
};
</script>

<style scoped>
html, body {
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.container {
  display: grid;
  grid-template-rows: 1fr;
  width: 100vw;
  height: 100vh;
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  color: #c5c5c5;
  background-color: #1e1e1e;
  overflow: hidden;
}

.activity-bar {
  grid-column: 1;
  background-color: #333333;
  padding-top: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}
.activity-icon-group {
  display: flex;
  flex-direction: column;
  gap: 15px;
}
.activity-icon-item {
  width: 50px;
  height: 48px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  color: #858585;
  position: relative;
  background: none;
  border: none;
}
.activity-icon-item:hover, .activity-icon-item.active {
  color: #ffffff;
}
.activity-icon-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: #ffffff;
}
.icon {
  width: 24px;
  height: 24px;
}

.side-bar {
  grid-column: 2;
  background-color: #252526;
  padding: 10px;
  border-right: 1px solid #000000;
  overflow-y: auto;
}
.sidebar-content {
  padding: 10px 0;
}
.sidebar-content h3 {
  color: #c5c5c5;
  font-size: 11px;
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.search-box input {
  width: 100%;
  box-sizing: border-box;
  background-color: #3e3e40;
  border: 1px solid #333333;
  color: #cccccc;
  padding: 5px 10px;
  border-radius: 3px;
}

.file-explorer h3 {
  color: #c5c5c5;
  font-size: 11px;
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.file-explorer .file-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 5px;
  margin-bottom: 5px;
  border-bottom: 1px solid transparent;
}
.add-file-btn {
  background: transparent;
  border: none;
  color: #c5c5c5;
  cursor: pointer;
  font-size: 20px;
  line-height: 1;
  padding: 0 5px;
}
.add-file-btn:hover {
  color: #ffffff;
}

.files-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px 0;
}
.files-title p {
  font-size: 11px;
  color: #888888;
  text-transform: uppercase;
  font-weight: normal;
  margin: 0;
  letter-spacing: 0.5px;
}
.file-explorer ul.file-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.file-explorer ul.file-list li {
  display: flex;
  align-items: center;
  padding: 2px 10px;
  cursor: pointer;
  color: #c5c5c5;
  font-size: 13px;
  transition: background-color 0.1s;
  position: relative;
}
.file-explorer ul.file-list li:hover {
  background-color: #2a2d2e;
}
.file-explorer ul.file-list li.active {
  background-color: #006090;
  color: #ffffff;
}
.file-explorer ul.file-list li span {
  flex-grow: 1;
  user-select: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.file-explorer ul.file-list li input {
  flex-grow: 1;
  background-color: #3e3e40;
  color: white;
  border: 1px solid #007acc;
  padding: 2px 4px;
  outline: none;
  box-sizing: border-box;
  font-size: 13px;
}
.delete-btn {
  background: transparent;
  border: none;
  color: #c5c5c5;
  cursor: pointer;
  font-size: 18px;
  padding: 0;
  line-height: 1;
  margin-left: 10px;
  opacity: 0;
  transition: opacity 0.2s;
}
.file-explorer ul.file-list li:hover .delete-btn {
  opacity: 1;
}
.file-explorer ul.file-list li.active .delete-btn {
  opacity: 1;
  color: #ffffff;
}
.delete-btn:hover {
  color: #e44f50;
}

.divider {
  grid-column: 3;
  width: 6px;
  background: transparent;
  cursor: col-resize;
  -webkit-user-select: none;
  -moz-user-select: none;
  user-select: none;
  transition: background-color 0.08s;
}
.divider:hover {
  background: rgba(112,129,144,0.25);
}

.main-editor-area {
  grid-column: 4;
  background-color: #1e1e1e;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #2d2d2d;
  border-bottom: 1px solid #000000;
  min-height: 35px;
  padding-right: 10px;
}
.tab-container {
  display: flex;
  overflow-x: auto;
  flex-grow: 1;
  white-space: nowrap;
}
.editor-tab {
  padding: 8px 15px;
  cursor: pointer;
  border-right: 1px solid #000000;
  color: #888;
  background-color: #2d2d2d;
  font-size: 13px;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
.editor-tab:hover {
  background-color: #333333;
  color: #ffffff;
}
.editor-tab.active {
  background-color: #1e1e1e;
  color: #ffffff;
  border-bottom: 2px solid #007acc;
  padding-bottom: 6px;
}
.close-tab-btn {
  margin-left: 10px;
  font-weight: normal;
  font-size: 16px;
  line-height: 1;
  color: #888;
  opacity: 0.7;
  transition: opacity 0.2s, color 0.2s;
}
.close-tab-btn:hover {
  opacity: 1;
  color: #ffffff;
}
.editor-tab.active .close-tab-btn {
  color: #ffffff;
}
.editor-actions {
  display: flex;
  gap: 8px;
}
.action-btn {
  background-color: #007acc;
  color: white;
  border: none;
  margin: 6px 0 6px 0;
  padding: 6px 10px;
  border-radius: 3px;
  width: 100%;
  cursor: pointer;
  font-size: 13px;
  transition: background-color 0.2s;
  white-space: nowrap;
}
.action-btn:hover {
  background-color: #0066a3;
}
.action-icon-btn {
  background-color: #2d2d2d;
  color: white;
  border: none;
  margin: 6px 0 6px 0;
  padding: 6px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 13px;
  transition: background-color 0.2s;
  width: 32px;
  height: 32px;
}
.action-icon-btn:hover {
  background-color: #1d1d1d;
}

.editor-content-wrapper {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  overflow: hidden;
}

#editor-container {
  height: 65%;
  flex-shrink: 0;
}
.terminal-area {
  height: 35%;
  background-color: #1e1e1e;
  color: white;
  display: flex;
  flex-direction: column;
  border-top: 1px solid #000000;
}
.terminal-header {
  background-color: #2d2d2d;
  padding: 5px 10px;
  display: flex;
  align-items: center;
  min-height: 28px;
}
.terminal-title {
  font-size: 13px;
  font-weight: 500;
  color: #c5c5c5;
  user-select: none;
}
.output-container {
  padding: 10px;
  overflow-y: auto;
  flex-grow: 1;
  font-size: 13px;
}
.terminal-area p {
  margin: 0;
  font-family: 'Courier New', Courier, monospace;
  white-space: pre-wrap;
  word-wrap: break-word;
  line-height: 1.4;
  color: white;
}
.terminal-area p.log {
  color: #cccccc;
}
.terminal-area p.error {
  color: #e44f50;
}
.terminal-area p.info {
  color: #888888;
}

.confirm-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}
.confirm-popup {
  background-color: #252526;
  color: #c5c5c5;
  padding: 20px;
  border-radius: 5px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.4);
  text-align: center;
  max-width: 400px;
  width: 90%;
  border: 1px solid #000000;
}
.confirm-popup p {
  margin-bottom: 20px;
  font-size: 14px;
  color: #c5c5c5;
}
.popup-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
.confirm-ok-btn, .confirm-cancel-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 13px;
  transition: background-color 0.2s;
}
.confirm-ok-btn {
  background-color: #007acc;
  color: white;
}
.confirm-ok-btn:hover {
  background-color: #0066a3;
}
.confirm-cancel-btn {
  background-color: #3c3c3c;
  color: white;
}
.confirm-cancel-btn:hover {
  background-color: #555555;
}
</style>
