<template>
  <div class="container" :style="{ gridTemplateColumns: `50px ${sidebarWidth}px 6px 1fr` }">
    <ActivityBar :active="activeSidebar" @switch="switchSidebar" />

    <SideBar
      :sidebarState="sidebarState"
      :files="files"
      :activeFileId="activeFileId"
      :editingFileId="editingFileId"
      @add-new="addNewFile"
      @open-file="switchFile"
      @start-rename="startRename"
      @finish-rename="finishRename"
      @request-delete="confirmDelete"
      @run-code="runCode"
      @download-code="downloadCode"
    />

    <div class="divider" @mousedown.prevent="startResize"></div>

    <EditorArea
      :files="files"
      :activeFileId="activeFileId"
      @open-file="switchFile"
      @request-delete="confirmDelete"
      @update-file-content="handleUpdateFileContent"
    />

    <ConfirmPopup v-if="showConfirmPopup" :message="popupMessage" @confirm="confirmAction" />
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import ActivityBar from '../components/main/ActivityBar.vue';
import SideBar from '../components/main/SideBar.vue';
import EditorArea from '../components/main/EditorArea.vue';
import ConfirmPopup from '../components/main/ConfirmPopup.vue';

const files = ref([
  {
    id: uuidv4(),
    name: 'main.js',
    content: `console.log('Hello world!');`,
  },
]);
const activeFileId = ref(files.value[0].id);
const editingFileId = ref<string | null>(null);

const showConfirmPopup = ref(false);
const popupMessage = ref('');
let popupResolve: ((value: boolean | PromiseLike<boolean>) => void) | null = null;

// サイドバーの状態
const sidebarState = ref({ explorer: true, search: false, runcode: false, git: false });
const activeSidebar = ref('explorer');

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

const switchSidebar = (view: 'explorer' | 'search' | 'runcode' | 'git') => {
  sidebarState.value = { explorer: false, search: false, runcode: false, git: false };
  (sidebarState.value as any)[view] = true;
  activeSidebar.value = view;
};

const addNewFile = () => {
  const newFile = { id: uuidv4(), name: `new-file-${files.value.length + 1}.js`, content: `console.log('Hello world!');` };
  files.value.push(newFile);
  activeFileId.value = newFile.id;
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
  if (file && file.name.trim() === '') file.name = 'untitled.js';
  editingFileId.value = null;
};

const confirmDelete = (fileId: string) => {
  const fileToDelete = files.value.find(f => f.id === fileId);
  if (!fileToDelete) return;
  if (files.value.length <= 1) {
    alert('少なくとも1つのファイルが必要です。');
    return;
  }
  popupMessage.value = `本当にファイル「${fileToDelete.name}」を削除してもよろしいですか？`;
  showConfirmPopup.value = true;
  return new Promise<boolean>((resolve) => { popupResolve = resolve; (popupResolve as any)._target = fileId; });
};

const confirmAction = (result: boolean) => {
  showConfirmPopup.value = false;
  if (popupResolve) {
    const targetFileId = (popupResolve as any)._target as string | undefined;
    popupResolve(result);
    if (result && targetFileId) {
      const idx = files.value.findIndex(f => f.id === targetFileId);
      if (idx > -1) {
        files.value.splice(idx, 1);
        if (activeFileId.value === targetFileId) activeFileId.value = files.value[0].id;
      }
    }
    popupResolve = null;
  }
};

const handleUpdateFileContent = ({ id, content }: { id: string; content: string }) => {
  const f = files.value.find(x => x.id === id);
  if (f) f.content = content;
};

// コード実行（サイドバーからの呼び出し用）
const runCode = async () => {
  const file = files.value.find(f => f.id === activeFileId.value);
  if (!file) return;
  
  // EditorAreaのrunCode機能と重複するが、サイドバーからも実行できるようにする
  alert('コードを実行するには、エディタ上部の実行ボタン（▷）を使用してください。');
};

// コードダウンロード
const downloadCode = () => {
  const file = files.value.find(f => f.id === activeFileId.value);
  if (!file) return;
  
  const blob = new Blob([file.content], { type: 'text/javascript' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = file.name;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
};
</script>

<style>
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
