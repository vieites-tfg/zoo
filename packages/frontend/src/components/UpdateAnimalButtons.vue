<template>
  <Button
    :info="resetButton"
    @click="$emit('resetForm')"
  />
  <Button
    :info="updateButton"
    class="ml-2"
    @click="$emit('updateAnimal')"
  />
</template>

<script setup lang="ts">
import Button from './GenericButton.vue';

defineEmits([ 'resetForm', 'updateAnimal' ])

const resetButton: Button = {
  text: 'Reset'
}

const updateButton: Button = {
  text: 'Update',
  color: 'green',
  type: 'submit'
}
</script>
