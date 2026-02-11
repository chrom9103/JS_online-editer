<template>
  <div class="admin-container">
    <!-- 認証画面 -->
    <div v-if="!isAuthenticated" class="auth-overlay">
      <div class="auth-box">
        <h2>Admin Login</h2>
        <form @submit.prevent="handleLogin">
          <input
            v-model="password"
            type="password"
            placeholder="Password"
            class="auth-input"
            :disabled="isLoggingIn"
          />
          <p v-if="authError" class="auth-error">{{ authError }}</p>
          <button type="submit" class="auth-btn" :disabled="isLoggingIn">
            {{ isLoggingIn ? 'Logging in...' : 'Login' }}
          </button>
        </form>
      </div>
    </div>

    <!-- 認証後のメインコンテンツ -->
    <template v-else>
      <header class="admin-header">
        <h1>Admin - Runs</h1>
        <div class="header-actions">
          <button class="refresh-btn" @click="fetchRuns">Refresh</button>
          <button class="logout-btn" @click="handleLogout">Logout</button>
        </div>
      </header>

      <div class="admin-content">
        <!-- Left: File List -->
        <div class="file-list-panel">
          <div v-if="loading" class="loading">Loading...</div>
          <div v-else-if="error" class="error">{{ error }}</div>
          <div v-else-if="runs.length === 0" class="empty">No runs found</div>
          <div class="table-scroll" v-else ref="tableScroll">
            <table ref="tableEl">
              <thead>
                <tr>
                  <th>File Name</th>
                  <th>Size</th>
                  <th>Modified</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(run, idx) in runs"
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
        </div>

        <!-- Right: Editor -->
        <div class="editor-panel">
          <div v-if="selectedRun" class="editor-header-bar">
            <span class="file-name">{{ selectedRun.name }}</span>
            <a
              class="download-link"
              :href="`${getApiBaseUrl()}/runs/${selectedRun.name}`"
              target="_blank"
              @click.prevent="downloadFile"
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
        <div ref="globalScroll" class="global-scroll">
          <div class="global-inner"></div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import MonacoEditor from '../components/main/editorArea/MonacoEditor.vue'

// API base helper (matches EditorArea approach)
const getApiBaseUrl = () => {
  if (import.meta.env.PROD) return '/playground/api'
  return '/api'
}

interface RunFile {
  name: string
  size: number
  modTime: string
}

const AUTH_STORAGE_KEY = 'admin_auth_token'
const AUTH_EXPIRY_KEY = 'admin_auth_expiry'

// 認証状態
const isAuthenticated = ref(false)
const isLoggingIn = ref(false)
const password = ref('')
const authError = ref('')
const authToken = ref('')

// メインコンテンツ
const runs = ref<RunFile[]>([])
const loading = ref(false)
const error = ref('')
const selectedRun = ref<RunFile | null>(null)
const selectedRunContent = ref('')
const editorHeight = ref(600)
const tableScroll = ref<HTMLElement | null>(null)
const tableEl = ref<HTMLElement | null>(null)
const globalScroll = ref<HTMLElement | null>(null)
const tableWidth = ref(0)

let onTableScroll: (() => void) | null = null
let onGlobalScroll: (() => void) | null = null

// 認証トークンをヘッダーに含めてfetch
const authFetch = async (url: string, options: RequestInit = {}) => {
  // map /api prefix to production base if needed
  let resolved = url
  if (typeof url === 'string' && url.startsWith('/api')) {
    resolved = url.replace(/^\/api/, getApiBaseUrl())
  }
  const headers = new Headers(options.headers || {})
  headers.set('Authorization', `Bearer ${authToken.value}`)
  return fetch(resolved, { ...options, headers })
}

// ログイン処理
const handleLogin = async () => {
  if (!password.value.trim()) {
    authError.value = 'Please enter password'
    return
  }

  isLoggingIn.value = true
  authError.value = ''

  try {
    const res = await fetch(`${getApiBaseUrl()}/admin/auth`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password: password.value }),
    })

    const data = await res.json()

    if (data.success && data.token) {
      authToken.value = data.token
      isAuthenticated.value = true
      password.value = ''

      // localStorageに保存（24時間有効）
      const expiry = Date.now() + 24 * 60 * 60 * 1000
      localStorage.setItem(AUTH_STORAGE_KEY, data.token)
      localStorage.setItem(AUTH_EXPIRY_KEY, expiry.toString())

      // ログイン成功後にデータを取得
      fetchRuns()
    } else {
      authError.value = data.error || 'Authentication failed'
    }
  } catch (e: any) {
    authError.value = `Network error: ${e.message}`
  } finally {
    isLoggingIn.value = false
  }
}

// ログアウト処理
const handleLogout = () => {
  authToken.value = ''
  isAuthenticated.value = false
  localStorage.removeItem(AUTH_STORAGE_KEY)
  localStorage.removeItem(AUTH_EXPIRY_KEY)
  runs.value = []
  selectedRun.value = null
  selectedRunContent.value = ''
}

// 保存されたトークンを検証
const checkStoredAuth = async () => {
  const storedToken = localStorage.getItem(AUTH_STORAGE_KEY)
  const storedExpiry = localStorage.getItem(AUTH_EXPIRY_KEY)

  if (!storedToken || !storedExpiry) return false

  // 有効期限チェック（クライアント側）
  if (Date.now() > parseInt(storedExpiry, 10)) {
    localStorage.removeItem(AUTH_STORAGE_KEY)
    localStorage.removeItem(AUTH_EXPIRY_KEY)
    return false
  }

  // サーバー側でトークン検証
  try {
    const res = await fetch(`${getApiBaseUrl()}/admin/verify`, {
      headers: { Authorization: `Bearer ${storedToken}` },
    })
    const data = await res.json()

    if (data.valid) {
      authToken.value = storedToken
      isAuthenticated.value = true
      return true
    }
  } catch {
    // 検証失敗
  }

  localStorage.removeItem(AUTH_STORAGE_KEY)
  localStorage.removeItem(AUTH_EXPIRY_KEY)
  return false
}

const fetchRuns = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await authFetch('/api/runs')
    if (res.status === 401) {
      handleLogout()
      return
    }
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
    const res = await authFetch(`/api/runs/${run.name}`)
    if (res.status === 401) {
      handleLogout()
      return
    }
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    selectedRunContent.value = await res.text()
    selectedRun.value = run
  } catch (e: any) {
    error.value = `Failed to load file: ${e.message}`
  }
}

const downloadFile = async () => {
  if (!selectedRun.value) return
  try {
    const res = await authFetch(`/api/runs/${selectedRun.value.name}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = selectedRun.value.name
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  } catch (e: any) {
    error.value = `Failed to download: ${e.message}`
  }
}

const updateEditorHeight = () => {
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

onMounted(async () => {
  updateEditorHeight()
  window.addEventListener('resize', updateEditorHeight)

  // set up horizontal scrollbar sync
  const updateTableWidth = () => {
    tableWidth.value = tableEl.value ? tableEl.value.scrollWidth : 0
    if (globalScroll.value && tableWidth.value) {
      const inner = globalScroll.value.querySelector('.global-inner') as HTMLElement | null
      if (inner) inner.style.width = tableWidth.value + 'px'
    }
  }

  onTableScroll = () => {
    if (!tableScroll.value || !globalScroll.value) return
    globalScroll.value.scrollLeft = tableScroll.value.scrollLeft
  }

  onGlobalScroll = () => {
    if (!tableScroll.value || !globalScroll.value) return
    tableScroll.value.scrollLeft = globalScroll.value.scrollLeft
  }

  // attach when elements available
  setTimeout(() => {
    updateTableWidth()
    if (tableScroll.value) tableScroll.value.addEventListener('scroll', onTableScroll!)
    if (globalScroll.value) globalScroll.value.addEventListener('scroll', onGlobalScroll!)
    window.addEventListener('resize', updateTableWidth)
  }, 0)

  // 保存されたトークンをチェック
  const isValid = await checkStoredAuth()
  if (isValid) {
    fetchRuns()
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', updateEditorHeight)
  if (tableScroll.value && onTableScroll)
    tableScroll.value.removeEventListener('scroll', onTableScroll)
  if (globalScroll.value && onGlobalScroll)
    globalScroll.value.removeEventListener('scroll', onGlobalScroll)
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

/* 認証画面 */
.auth-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #1e1e1e;
}

.auth-box {
  background-color: #252526;
  padding: 40px;
  border-radius: 8px;
  border: 1px solid #333;
  width: 320px;
  text-align: center;
}

.auth-box h2 {
  margin: 0 0 24px 0;
  font-size: 20px;
  font-weight: 500;
  color: #c5c5c5;
}

.auth-input {
  width: 100%;
  padding: 12px 16px;
  font-size: 14px;
  background-color: #3c3c3c;
  border: 1px solid #555;
  border-radius: 4px;
  color: #c5c5c5;
  box-sizing: border-box;
  margin-bottom: 12px;
}

.auth-input:focus {
  outline: none;
  border-color: #007acc;
}

.auth-error {
  color: #e44f50;
  font-size: 13px;
  margin: 0 0 12px 0;
}

.auth-btn {
  width: 100%;
  padding: 12px;
  font-size: 14px;
  background-color: #007acc;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.auth-btn:hover:not(:disabled) {
  background-color: #0066a3;
}

.auth-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* ヘッダー */
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

.header-actions {
  display: flex;
  gap: 10px;
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

.logout-btn {
  background-color: #3c3c3c;
  color: #c5c5c5;
  border: none;
  padding: 6px 14px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
}

.logout-btn:hover {
  background-color: #555;
}

.admin-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* Left Panel - File List */
.file-list-panel {
  width: 30%;
  border-right: 1px solid #333;
  overflow-y: auto;
  overflow-x: hidden;
  background-color: #1e1e1e;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, 'Roboto Mono', 'Courier New', monospace;
  padding-bottom: 22px; /* space for global scrollbar */
}

.table-scroll {
  overflow-x: auto;
}

.global-scroll {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 30%;
  height: 18px;
  overflow-x: auto;
  background: transparent;
  z-index: 50;
}

.global-scroll .global-inner {
  height: 1px;
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
  width: max-content;
  border-collapse: collapse;
  font-size: 12px;
  table-layout: auto;
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
  white-space: nowrap;
}

.col-num,
.line-num {
  width: 40px;
  text-align: right;
  padding-right: 16px;
  color: #888;
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
  width: 70%;
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
