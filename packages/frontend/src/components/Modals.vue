<template>
  <AnimalModal
    :openModal="openModal"
    :info="currentData"
    :clear="clearModal"
    :provideFormData="provideFormData"
    @cleared="clearedModal"
    @providingFormData="providingFormData"
    @closeAnimalModal="$emit('closeAnimalModal')">
    <template v-if="currentData && currentData.action === 'Create'">
      <NewAnimalButtons @clearForm="clearForm" @createNewAnimal="createNewAnimal" />
    </template>
    <template v-else>
      <UpdateAnimalButtons @resetForm="resetForm" @updateAnimal="updateAnimal" />
    </template>
  </AnimalModal>
</template>

<script setup lang="ts">
import { ref, computed, defineProps } from 'vue'
import IAnimalModal from '../types/AnimalModal'
import AnimalModal from './AnimalModal.vue'
import NewAnimalButtons from './NewAnimalButtons.vue'
import UpdateAnimalButtons from './UpdateAnimalButtons.vue'

const props = defineProps<{
  openModal: boolean,
  animalModalData: IAnimalModal,
}>()

const emits = defineEmits(['closeAnimalModal'])

const clearedForm: IAnimal = {
  name: '',
  species: '',
  birthday: '',
  genre: '',
  diet: '',
  condition: '',
  notes: ''
}

const clearModal = ref<boolean>(false)
const provideFormData = ref<boolean>(false)

const initialData = ref<IAnimalModal>(!props.animalModalData.data ? { ...clearedForm } : props.animalModalData.data)

const currentData = computed<IAnimalModal>(() => {
  props.animalModalData.data = !props.animalModalData.data ? { ...clearedForm } : props.animalModalData.data
  return { ...props.animalModalData }
})

// Functions

const providingFormData = (animal: IAnimal) => {
  provideFormData.value = false
  console.log({action: props.animalModalData.action})
  console.log({initial: currentData.value})
  console.log(animal)
  emits('closeAnimalModal')
}

const clearedModal = () => {
  clearModal.value = false
}

const clearForm = () => {
  clearModal.value = true
}

const resetForm = () => {
  currentData.value.data = { ...clearedForm }
}

const createNewAnimal = () => {
  provideFormData.value = true
}

const updateAnimal = () => {
  provideFormData.value = true
}
</script>
