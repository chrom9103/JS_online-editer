<template>
  <div class="file-explorer">
    <div class="file-header">
      <h3>エクスプローラー</h3>
    </div>
    <div class="files-title">
      <p>∨ JS_PlayGround</p>
      <button @click="$emit('add-new')" class="add-file-btn">＋</button>
    </div>
    <ul class="file-list">
      <li
        v-for="file in files"
        :key="file.id"
        :class="{ active: file.id === activeFileId }"
        @click="$emit('open-file', file.id)"
      >
        <span v-if="editingFileId !== file.id" @dblclick.stop="$emit('start-rename', file.id)">
          {{ file.name }}
        </span>
        <input
          v-else
          type="text"
          v-model="file.name"
          @blur="$emit('finish-rename')"
          @keyup.enter="$emit('finish-rename')"
          autofocus
        />
        <button @click.stop="$emit('request-delete', file.id)" class="delete-btn">×</button>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  files: Array<{ id: string; name: string; content: string }>
  activeFileId: string
  editingFileId: string | null
}>()
</script>
