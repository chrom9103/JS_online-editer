<template>
  <div class="confirm-overlay">
    <div class="confirm-popup">
      <p>{{ message }}</p>
      <div class="popup-actions">
        <button @click="$emit('confirm', true)" class="confirm-ok-btn">OK</button>
        <button @click="$emit('confirm', false)" class="confirm-cancel-btn">キャンセル</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ message: string }>();
</script>
