<template>
  <Button @click="$emit('resetForm')" :info="resetButton" />
  <Button @click="$emit('updateAnimal')" :info="updateButton" class="ml-2" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Button from './GenericButton.vue';
import IAnimal from '../types/Animal';

const emit = defineEmits<{
  (e: 'resetForm'): void
  (e: 'updateAnimal'): void
}>()

const resetButton: Button = {
  text: 'Reset'
}

const updateButton: Button = {
  text: 'Update',
  color: 'green',
  type: 'submit'
}
</script>
