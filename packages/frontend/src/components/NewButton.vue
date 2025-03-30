<template>
  <button
    class="border-2 border-green-600 font-semibold text-green-500 px-4 py-2 rounded hover:bg-green-600 hover:text-white transition"
    @click="$emit('click')"
  >
    <span>+ New</span>
  </button>
</template>

<script setup lang="ts">
defineEmits(['click'])
</script>
