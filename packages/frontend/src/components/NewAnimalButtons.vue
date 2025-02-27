<template>
  <Button @click="$emit('clearForm')" :info="clearButton" />
  <Button @click="$emit('createNewAnimal', animal)" :info="addButton" class="ml-2" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Button from './GenericButton.vue';
import IAnimal from '../types/Animal';

const emit = defineEmits<{
  (e: 'clearForm'): void
  (e: 'createNewAnimal', animal: IAnimal): void
}>()

const clearButton: Button = {
  text: 'Clear'
}

const addButton: Button = {
  text: 'Add',
  color: 'green',
  type: 'submit'
}
</script>
