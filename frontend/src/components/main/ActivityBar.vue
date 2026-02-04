<template>
  <div class="activity-bar">
    <div class="activity-icon-group">
      <button
        @click="$emit('switch', 'explorer')"
        class="activity-icon-item"
        :class="{ active: active === 'explorer' }"
        title="エクスプローラー"
      >
        <img :src="explorer" class="icon" alt="エクスプローラー" />
      </button>
      <button
        @click="$emit('switch', 'text')"
        class="activity-icon-item"
        :class="{ active: active === 'text' }"
        title="テキスト"
      >
        <img :src="text" class="icon" alt="テキスト" />
      </button>
      <button
        @click="$emit('switch', 'search')"
        class="activity-icon-item"
        :class="{ active: active === 'search' }"
        title="検索"
      >
        <img :src="search" class="icon" alt="検索" />
      </button>
      <button
        @click="$emit('switch', 'git')"
        class="activity-icon-item"
        :class="{ active: active === 'git' }"
        title="ソース管理"
      >
        <img :src="git" class="icon" alt="ソース管理" />
      </button>
      <button
        @click="$emit('switch', 'runcode')"
        class="activity-icon-item"
        :class="{ active: active === 'runcode' }"
        title="実行とデバッグ"
      >
        <img :src="runcode" class="icon" alt="実行とデバッグ" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ active: string }>()

const explorer = new URL('/src/assets/ActivityBar/explorer-icon.svg', import.meta.url).href
const text = new URL('/src/assets/ActivityBar/text-icon.svg', import.meta.url).href
const search = new URL('/src/assets/ActivityBar/search-icon.svg', import.meta.url).href
const git = new URL('/src/assets/ActivityBar/git-icon.svg', import.meta.url).href
const runcode = new URL('/src/assets/ActivityBar/runcode-icon.svg', import.meta.url).href
</script>
