<template>
  <button
    @click="$emit('click')"
    :disabled="disabled"
    :class="[
      'border-2 border-red-600 font-semibold text-red-500 px-4 py-2 rounded transition hover:bg-red-600 hover:text-white',
      disabled ? 'opacity-50 cursor-not-allowed' : 'opacity-100'
    ]"
  >
    <span>x Delete</span>
  </button>
</template>

<script setup lang="ts">
const props = defineProps<{
  disabled?: boolean;
}>();

const emit = defineEmits(['click'])
</script>

<style scoped>
</style>
