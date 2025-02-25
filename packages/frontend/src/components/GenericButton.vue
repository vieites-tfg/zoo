<template>
  <button :type="info.type || 'button'" @click="$emit('click')" :class="[
    'border-2 font-semibold px-4 py-2 rounded transition',
    info.disabled ? 'opacity-50 cursor-not-allowed' : 'opacity-100',
    bgColorClass]">
    <span>{{ info.text }}</span>
  </button>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, computed } from 'vue'
import Button from '../types/Button'
const props = defineProps<{
  info: Button
}>()

console.log(props.info.color)

const bgColorClass = computed(() => {
  if (props.info.color === 'green') {
    return 'border-green-500 text-green-500 hover:bg-green-500 hover:text-white'
  }
  if (props.info.color === 'red') {
    return 'border-red-500 text-red-500 hover:bg-red-500 hover:text-white'
  }
  return 'border-black text-black hover:bg-black hover:text-white'
})

const emit = defineEmits(['click'])
</script>
