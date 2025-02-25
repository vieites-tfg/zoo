<template>
  <div class="flex justify-start pl-4">
    <NewButton @click="handleNew" />
    <DeleteButton @click="handleDelete" :disabled="!props.isDeleteEnabled" class="ml-2" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import NewButton from './NewButton.vue';
import DeleteButton from './DeleteButton.vue';

const props = defineProps<{
  isDeleteEnabled: bool
}>()

const handleNew = () => {
  console.log("Botón + New presionado");
};

const handleDelete = () => {
  if (!props.isDeleteEnabled) return;
  console.log("Botón x Delete presionado");
};
</script>
