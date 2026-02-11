<template>
  <div class="admin-container">
    <header class="admin-header">
      <h1>Admin - Runs</h1>
      <button class="refresh-btn" @click="fetchRuns">Refresh</button>
    </header>

    <div class="admin-content">
      <!-- Left: File List -->
      <div class="file-list-panel">
        <div v-if="loading" class="loading">Loading...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="runs.length === 0" class="empty">No runs found</div>
        <table v-else>
          <thead>
            <tr>
              <th>File Name</th>
              <th>Size</th>
              <th>Modified</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="run in runs"
              :key="run.name"
              :class="{ active: selectedRun?.name === run.name }"
              @click="selectRun(run)"
            >
              <td>{{ run.name }}</td>
              <td>{{ formatSize(run.size) }}</td>
              <td>{{ formatDate(run.modTime) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Right: Editor -->
      <div class="editor-panel">
        <div v-if="selectedRun" class="editor-header-bar">
          <span class="file-name">{{ selectedRun.name }}</span>
          <a class="download-link" :href="`/api/runs/${selectedRun.name}`" target="_blank"
            >Download</a
          >
        </div>
        <div v-else class="editor-header-bar">
          <span class="file-name">Select a file to view</span>
        </div>
        <div class="editor-wrapper">
          <MonacoEditor :content="selectedRunContent" :height="editorHeight" :read-only="true" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import MonacoEditor from '../components/main/editorArea/MonacoEditor.vue'

interface RunFile {
  name: string
  size: number
  modTime: string
}

const runs = ref<RunFile[]>([])
const loading = ref(false)
const error = ref('')
const selectedRun = ref<RunFile | null>(null)
const selectedRunContent = ref('')
const editorHeight = ref(600)

const fetchRuns = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await fetch('/api/runs')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    runs.value = data || []
  } catch (e: any) {
    error.value = `Failed to fetch runs: ${e.message}`
  } finally {
    loading.value = false
  }
}

const selectRun = async (run: RunFile) => {
  try {
    const res = await fetch(`/api/runs/${run.name}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    selectedRunContent.value = await res.text()
    selectedRun.value = run
  } catch (e: any) {
    error.value = `Failed to load file: ${e.message}`
  }
}

const updateEditorHeight = () => {
  // ヘッダー(60px) + エディタヘッダー(40px) + 余白を引く
  editorHeight.value = window.innerHeight - 140
}

const formatSize = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

const formatDate = (isoDate: string): string => {
  const d = new Date(isoDate)
  return d.toLocaleString('ja-JP')
}

onMounted(() => {
  fetchRuns()
  updateEditorHeight()
  window.addEventListener('resize', updateEditorHeight)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateEditorHeight)
})
</script>

<style scoped>
.admin-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
  color: #c5c5c5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  overflow: hidden;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: #252526;
  border-bottom: 1px solid #000;
  flex-shrink: 0;
}

.admin-header h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
}

.refresh-btn {
  background-color: #007acc;
  color: white;
  border: none;
  padding: 6px 14px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

.refresh-btn:hover {
  background-color: #0066a3;
}

.admin-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* Left Panel - File List */
.file-list-panel {
  width: 50%;
  border-right: 1px solid #333;
  overflow-y: auto;
  background-color: #1e1e1e;
}

.loading,
.error,
.empty {
  text-align: center;
  padding: 40px;
  font-size: 14px;
}

.error {
  color: #e44f50;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 12px;
}

thead {
  background-color: #252526;
  position: sticky;
  top: 0;
  z-index: 1;
}

th,
td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid #333;
}

th {
  font-weight: 600;
  color: #888;
  text-transform: uppercase;
  font-size: 10px;
  letter-spacing: 0.5px;
}

tbody tr {
  cursor: pointer;
  transition: background-color 0.15s;
}

tbody tr:hover {
  background-color: #2a2d2e;
}

tbody tr.active {
  background-color: #094771;
}

/* Right Panel - Editor */
.editor-panel {
  width: 50%;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
}

.editor-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background-color: #252526;
  border-bottom: 1px solid #333;
  flex-shrink: 0;
}

.file-name {
  font-size: 13px;
  color: #c5c5c5;
}

.download-link {
  color: #007acc;
  text-decoration: none;
  font-size: 12px;
}

.download-link:hover {
  text-decoration: underline;
}

.editor-wrapper {
  flex: 1;
  overflow: hidden;
}
</style>
