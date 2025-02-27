<template>
  <AnimalModal :openModal="openModal" :info="currentData" @closeAnimalModal="$emit('closeAnimalModal')">
    <template v-if="currentData && currentData.action === 'Create'">
      <NewAnimalButtons @clearForm="clearForm" @createNewAnimal="createNewAnimal" />
    </template>
    <template v-else>
      <UpdateAnimalButtons />
    </template>
  </AnimalModal>
</template>

<script setup lang="ts">
import { computed, defineProps } from 'vue'
import IAnimalModal from '../types/AnimalModal'
import AnimalModal from './AnimalModal.vue'
import NewAnimalButtons from './NewAnimalButtons.vue'
import UpdateAnimalButtons from './UpdateAnimalButtons.vue'

const props = defineProps<{
  openModal: boolean,
  animalModalData: IAnimalModal,
}>()

const emits = defineEmits(['closeAnimalModal'])

const currentData = computed<IAnimalModal>(() => {
  props.animalModalData.data = !props.animalModalData.data ? { ...cleanedForm } : props.animalModalData.data
  return props.animalModalData
})

const cleanedForm: IAnimal = {
  name: '',
  species: '',
  birthday: '',
  genre: '',
  diet: '',
  condition: '',
  notes: ''
}

// Functions

const clearForm = () => {
  currentData.value.data = { ...cleanedForm }
}

const createNewAnimal = (animal: IAnimal) => {
  console.log(currentData.value.data)
}
</script>
