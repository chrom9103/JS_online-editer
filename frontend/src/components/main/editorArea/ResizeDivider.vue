<template>
  <div
    class="resize-divider"
    :class="{ resizing: isResizing }"
    @mousedown.prevent="startResize"
  ></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const emit = defineEmits<{
  (e: 'resize', deltaY: number): void
  (e: 'resize-start'): void
  (e: 'resize-end'): void
}>()

const isResizing = ref(false)
let lastClientY = 0

const startResize = (event: MouseEvent) => {
  if (event.button !== 0) return
  isResizing.value = true
  lastClientY = event.clientY
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'row-resize'
  emit('resize-start')
}

const onMouseMove = (event: MouseEvent) => {
  if (!isResizing.value) return
  const deltaY = event.clientY - lastClientY
  lastClientY = event.clientY
  emit('resize', deltaY)
}

const stopResize = () => {
  if (!isResizing.value) return
  isResizing.value = false
  document.body.style.userSelect = ''
  document.body.style.cursor = ''
  emit('resize-end')
}

onMounted(() => {
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', stopResize)
})

onUnmounted(() => {
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', stopResize)
})
</script>

<style scoped>
.resize-divider {
  height: 6px;
  cursor: row-resize;
  background: transparent;
  transition: background-color 0.08s;
  flex-shrink: 0;
  z-index: 10;
  position: relative;
}

.resize-divider:hover,
.resize-divider.resizing {
  background: rgba(112, 129, 144, 0.25);
}
</style>
