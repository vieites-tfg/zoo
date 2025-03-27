<template>
  <AnimalModal
    :open-modal="openModal"
    :info="currentData"
    :clear="clear"
    :reset="reset"
    :provide-form-data="provideFormData"
    @cleared="clearedForm"
    @reseted="resetedForm"
    @providing-form-data="providingFormData"
    @close-animal-modal="$emit('closeAnimalModal')"
  >
    <template v-if="currentData && currentData.action === 'Create'">
      <NewAnimalButtons
        @clear-form="clearForm"
        @create-new-animal="createNewAnimal"
      />
    </template>
    <template v-else>
      <UpdateAnimalButtons
        @reset-form="resetForm"
        @update-animal="updateAnimal"
      />
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

const emptyForm: IAnimal = {
  name: '',
  species: '',
  birthday: '',
  genre: '',
  diet: '',
  condition: '',
  notes: ''
}

const clear = ref<boolean>(false)
const reset = ref<boolean>(false)
const provideFormData = ref<boolean>(false)

const currentData = computed<IAnimalModal>(() => ({
  ...props.animalModalData,
  data: props.animalModalData.data || { ...emptyForm }
}));

// Functions

const providingFormData = (animal: IAnimal) => {
  provideFormData.value = false
  console.log({action: props.animalModalData.action})
  console.log({initial: currentData.value})
  console.log(animal)
  emits('closeAnimalModal')
}

const clearedForm = () => {
  clear.value = false
}

const clearForm = () => {
  clear.value = true
}

const resetForm = () => {
  reset.value = true
}

const resetedForm = () => {
  reset.value = false
}

const createNewAnimal = () => {
  provideFormData.value = true
}

const updateAnimal = () => {
  provideFormData.value = true
}
</script>
