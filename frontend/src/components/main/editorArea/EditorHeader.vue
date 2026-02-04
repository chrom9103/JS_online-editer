<template>
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
      <button @click="$emit('run-code')" class="action-icon-btn">▷</button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  files: Array<{ id: string; name: string; content: string }>
  activeFileId: string
}>()

defineEmits<{
  (e: 'open-file', fileId: string): void
  (e: 'request-delete', fileId: string): void
  (e: 'run-code'): void
}>()
</script>

<style scoped>
.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #181818;
  height: 40px;
  padding: 0 10px;
}

.tab-container {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  flex: 1;
}

.editor-tab {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  background-color: #2d2d2d;
  color: #c5c5c5;
  border-radius: 4px 4px 0 0;
  cursor: pointer;
  font-size: 13px;
  user-select: none;
  white-space: nowrap;
}

.editor-tab.active {
  background-color: #1e1e1e;
  color: white;
}

.editor-tab:hover {
  background-color: #3d3d3d;
}

.close-tab-btn {
  margin-left: 8px;
  font-size: 14px;
  cursor: pointer;
  color: #888;
}

.close-tab-btn:hover {
  color: white;
}

.editor-actions {
  display: flex;
  align-items: center;
}

.action-icon-btn {
  background-color: transparent;
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
</style>
