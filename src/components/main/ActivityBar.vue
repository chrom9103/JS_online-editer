<template>
  <div class="activity-bar">
    <div class="activity-icon-group">
      <button @click="$emit('switch', 'explorer')" class="activity-icon-item" :class="{ active: active === 'explorer' }" title="エクスプローラー">
        <svg viewBox="0 0 16 16" fill="currentColor" class="icon">
          <path d="M4 1h8a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2zM4 3v10h8V3H4zm2 2h4v2H6V5zm0 4h4v2H6V9z"/>
        </svg>
      </button>
      <button @click="$emit('switch', 'search')" class="activity-icon-item" :class="{ active: active === 'search' }" title="検索">
        <svg viewBox="0 0 24 24" fill="currentColor" class="icon">
          <path d="M15.5 14h-.79l-.28-.27A6.471 6.471 0 0 0 16 9.5 6.5 6.5 0 1 0 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 5 1.5-1.5-5-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
        </svg>
      </button>
      <button @click="$emit('switch', 'git')" class="activity-icon-item" :class="{ active: active === 'git' }" title="ソース管理">
        <svg viewBox="0 0 16 16" fill="currentColor" class="icon">
          <path d="M13 10c0 1.66-1.34 3-3 3s-3-1.34-3-3 1.34-3 3-3 3 1.34 3 3zM6 8a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zm4.5 0a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zM14 6H2a1 1 0 0 0-1 1v2a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V7a1 1 0 0 0-1-1z"/>
        </svg>
      </button>
      <button @click="$emit('switch', 'runcode')" class="activity-icon-item" :class="{ active: active === 'runcode' }" title="実行とデバッグ">
        <svg viewBox="0 0 24 24" fill="currentColor" class="icon">
          <path d="M8 5v14l11-7z"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ active: string }>();
</script>
