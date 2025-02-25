<template>
  <button
    @click="$emit('click')"
    class="border-2 border-green-600 font-semibold text-green-500 px-4 py-2 rounded hover:bg-green-600 hover:text-white transition"
  >
    <span>+ New</span>
  </button>
</template>

<script setup lang="ts">
const emit = defineEmits<{
  (e: 'click'): void;
}>();
</script>

<style scoped>
</style>
