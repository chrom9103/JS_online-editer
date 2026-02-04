<template>
  <div class="terminal-area">
    <div class="terminal-header">
      <span class="terminal-title">TERMINAL</span>
    </div>
    <div class="output-container" ref="outputContainer">
      <p v-for="(line, index) in output" :key="index" :class="line.type">
        {{ line.text }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'

const props = defineProps<{
  output: Array<{ text: string; type: string }>
}>()

const outputContainer = ref<HTMLElement | null>(null)

// 出力が追加されたらスクロールを下に移動
watch(
  () => props.output.length,
  () => {
    nextTick(() => {
      if (outputContainer.value) {
        outputContainer.value.scrollTop = outputContainer.value.scrollHeight
      }
    })
  },
)
</script>

<style scoped>
.terminal-area {
  background-color: #1e1e1e;
  color: white;
  display: flex;
  flex-direction: column;
  border-top: 1px solid #000000;
  flex: 1 1 auto;
  min-height: 0;
  overflow: hidden;
}

.terminal-header {
  background-color: #2d2d2d;
  padding: 5px 10px;
  display: flex;
  align-items: center;
  min-height: 28px;
  flex-shrink: 0;
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
  flex: 1 1 auto;
  font-size: 13px;
  min-height: 0;
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

.terminal-area p.warn {
  color: #f0ad4e;
}
</style>
