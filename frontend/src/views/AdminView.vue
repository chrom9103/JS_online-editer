<template>
  <div class="admin-container">
    <header class="admin-header">
      <h1>Admin - Runs</h1>
      <button class="refresh-btn" @click="fetchRuns">Refresh</button>
    </header>

    <div class="runs-list">
      <div v-if="loading" class="loading">Loading...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="runs.length === 0" class="empty">No runs found</div>
      <table v-else>
        <thead>
          <tr>
            <th>File Name</th>
            <th>Size</th>
            <th>Modified</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="run in runs" :key="run.name" @click="selectRun(run)">
            <td>{{ run.name }}</td>
            <td>{{ formatSize(run.size) }}</td>
            <td>{{ formatDate(run.modTime) }}</td>
            <td>
              <button class="view-btn" @click.stop="viewRun(run)">View</button>
              <a class="download-link" :href="`/api/runs/${run.name}`" target="_blank" @click.stop
                >Download</a
              >
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Code Viewer Modal -->
    <div v-if="selectedRun" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>{{ selectedRun.name }}</h2>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <pre><code>{{ selectedRunContent }}</code></pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

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

const selectRun = (run: RunFile) => {
  viewRun(run)
}

const viewRun = async (run: RunFile) => {
  try {
    const res = await fetch(`/api/runs/${run.name}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    selectedRunContent.value = await res.text()
    selectedRun.value = run
  } catch (e: any) {
    error.value = `Failed to load file: ${e.message}`
  }
}

const closeModal = () => {
  selectedRun.value = null
  selectedRunContent.value = ''
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
})
</script>

<style scoped>
.admin-container {
  min-height: 100vh;
  background-color: #1e1e1e;
  color: #c5c5c5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 30px;
  background-color: #252526;
  border-bottom: 1px solid #000;
}

.admin-header h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 500;
}

.refresh-btn {
  background-color: #007acc;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

.refresh-btn:hover {
  background-color: #0066a3;
}

.runs-list {
  padding: 20px 30px;
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
  font-size: 13px;
}

thead {
  background-color: #252526;
}

th,
td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #333;
}

th {
  font-weight: 600;
  color: #888;
  text-transform: uppercase;
  font-size: 11px;
  letter-spacing: 0.5px;
}

tbody tr {
  cursor: pointer;
  transition: background-color 0.15s;
}

tbody tr:hover {
  background-color: #2a2d2e;
}

.view-btn {
  background-color: #3c3c3c;
  color: white;
  border: none;
  padding: 4px 10px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  margin-right: 8px;
}

.view-btn:hover {
  background-color: #555;
}

.download-link {
  color: #007acc;
  text-decoration: none;
  font-size: 12px;
}

.download-link:hover {
  text-decoration: underline;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: #252526;
  border: 1px solid #000;
  border-radius: 6px;
  width: 80%;
  max-width: 900px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #333;
}

.modal-header h2 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  color: #c5c5c5;
}

.close-btn {
  background: none;
  border: none;
  color: #888;
  font-size: 24px;
  cursor: pointer;
  line-height: 1;
}

.close-btn:hover {
  color: #fff;
}

.modal-body {
  padding: 20px;
  overflow: auto;
  flex: 1;
}

.modal-body pre {
  margin: 0;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.modal-body code {
  color: #d4d4d4;
}
</style>
