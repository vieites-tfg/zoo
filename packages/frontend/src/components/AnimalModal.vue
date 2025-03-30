<template>
  <div>
    <div
      v-if="openModal"
      class="fixed inset-0 flex items-center justify-center bg-gray-100/50 bg-opacity-50 z-50"
    >
      <div class="bg-white w-2/5 p-6 rounded shadow-lg relative">
        <h2 class="text-xl font-title font-bold mb-4 text-center">
          {{ info.title }}
        </h2>

        <form
          class="grid grid-cols-2 gap-4"
          @submit.prevent="handleAdd"
        >
          <div class="mb-2">
            <label
              class="block font-semibold"
              for="name"
            >Name</label>
            <input
              id="name"
              v-model="formData.name"
              type="text"
              class="w-full border border-gray-300 rounded px-2 py-1"
            >
          </div>

          <div class="mb-2">
            <label
              class="block font-semibold"
              for="species"
            >Species</label>
            <input
              id="species"
              v-model="formData.species"
              type="text"
              class="w-full border border-gray-300 rounded px-2 py-1"
            >
          </div>

          <div class="mb-2">
            <label
              class="block font-semibold"
              for="birthday"
            >Birthday</label>
            <input
              id="birthday"
              v-model="formData.birthday"
              type="date"
              class="w-full border border-gray-300 rounded px-2 py-1"
            >
          </div>

          <div class="mb-2">
            <label
              class="block font-semibold"
              for="genre"
            >Genre</label>
            <input
              id="genre"
              v-model="formData.genre"
              type="text"
              class="w-full border border-gray-300 rounded px-2 py-1"
            >
          </div>

          <div class="mb-2">
            <label
              class="block font-semibold"
              for="diet"
            >Diet</label>
            <input
              id="diet"
              v-model="formData.diet"
              type="text"
              class="w-full border border-gray-300 rounded px-2 py-1"
            >
          </div>

          <div class="mb-2">
            <label
              class="block font-semibold"
              for="condition"
            >Condition</label>
            <input
              id="condition"
              v-model="formData.condition"
              type="text"
              class="w-full border border-gray-300 rounded px-2 py-1"
            >
          </div>

          <div class="mb-2 col-span-2">
            <label
              class="block font-semibold"
              for="notes"
            >Notes</label>
            <textarea
              id="notes"
              v-model="formData.notes"
              class="w-full border border-gray-300 rounded px-2 py-1 min-h-8 max-h-56"
            />
          </div>

          <div class="col-start-2 flex justify-end space-x-2 mt-4">
            <slot />
          </div>
        </form>

        <button
          class="absolute top-2 right-2 text-gray-500 hover:text-gray-700"
          @click="closeModal"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="size-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M6 18 18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import IAnimal from '../types/Animal'
import IAnimalModal from '../types/AnimalModal'

const props = defineProps<{
  clear: boolean,
  reset: boolean,
  openModal: boolean,
  provideFormData: boolean,
  info: IAnimalModal,
}>()

const emit = defineEmits<{
  (e: 'closeAnimalModal'): void
  (e: 'cleared'): void
  (e: 'reseted'): void
  (e: 'providingFormData', animal: IAnimal): void
}>()

const emptyForm: IAnimal = {
  name: '',
  species: '',
  birthday: '',
  genre: '',
  diet: '',
  condition: '',
  notes: ''
}

const formData = ref<IAnimal>({ ...props.info.data })

watch(
  () => props.provideFormData,
  (newVal) => {
    if (newVal) {
      emit('providingFormData', { ...formData.value })
    }
  }
)

watch(
  () => props.openModal,
  (opened) => {
    if (opened) {
      formData.value = { ...props.info.data }
    }
  }
)

watch(
  () => props.clear,
  (newVal) => {
    if (newVal) {
      formData.value = { ...emptyForm }
      emit('cleared')
    }
  }
)

watch(
  () => props.reset,
  (newVal) => {
    if (newVal) {
      formData.value = { ...props.info.data }
      emit('reseted')
    }
  }
)

function handleAdd() {
  console.log('Submitted formData:', formData.value)
}

const closeModal = () => {
  emit('closeAnimalModal')
}
</script>
