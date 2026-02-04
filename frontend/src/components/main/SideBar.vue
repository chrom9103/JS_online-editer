<template>
  <div class="side-bar">
    <component v-if="sidebarState.explorer" :is="FileExplorer" class="file-explorer"
      :files="files" :activeFileId="activeFileId" :editingFileId="editingFileId"
      @add-new="$emit('add-new')"
      @open-file="$emit('open-file', $event)"
      @start-rename="$emit('start-rename', $event)"
      @finish-rename="$emit('finish-rename')"
      @request-delete="$emit('request-delete', $event)"
    />

    <div v-else-if="sidebarState.text" class="sidebar-content">
      <h3>テキスト</h3>
      <p>この機能は現在開発中です。</p>
    </div>

    <div v-else-if="sidebarState.search" class="sidebar-content">
      <h3>検索</h3>
      <p>この機能は現在開発中です。</p>
    </div>

    <div v-else-if="sidebarState.runcode" class="sidebar-content">
      <h3>実行とデバッグ</h3>
      <button @click="$emit('run-code')" class="action-btn">Run Code</button>
      <button @click="$emit('download-code')" class="action-btn">Download Code</button>
    </div>

    <div v-else-if="sidebarState.git" class="sidebar-content">
      <h3>ソース管理</h3>
      <p>この機能は現在開発中です。</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import FileExplorer from './FileExplorer.vue';

const props = defineProps<{
  sidebarState: { explorer: boolean; text: boolean; search: boolean; runcode: boolean; git: boolean };
  files: Array<{ id: string; name: string; content: string }>;
  activeFileId: string;
  editingFileId: string | null;
}>();
</script>
