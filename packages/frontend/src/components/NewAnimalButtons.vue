<template>
  <Button
    :info="clearButton"
    @click="$emit('clearForm')"
  />
  <Button
    :info="addButton"
    class="ml-2"
    @click="$emit('createNewAnimal')"
  />
</template>

<script setup lang="ts">
import Button from './GenericButton.vue';

defineEmits(['clearForm', 'createNewAnimal'])

const clearButton: Button = {
  text: 'Clear'
}

const addButton: Button = {
  text: 'Add',
  color: 'green',
  type: 'submit'
}
</script>
