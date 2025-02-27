<template>
  <div class="flex justify-start pl-4">
    <Button @click="handleNew" :info="greenButton" />
    <Button @click="handleDelete" :info="redButton" class="ml-2" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Button from './GenericButton.vue';

const props = defineProps<{
  isDeleteEnabled: bool
}>()

const emit = defineEmits(['newAnimal'])

const greenButton: Button = {
  text: "New",
  color: "green"
}

const redButton = computed(() => {
  return {
    text: "Delete",
    color: "red",
    disabled: !props.isDeleteEnabled
  }
})

const handleNew = () => {
  emit('newAnimal')
};

const handleDelete = () => {
  if (!props.isDeleteEnabled) return;
  console.log("Botón x Delete presionado");
};
</script>
