<template>
  <div>
    <div v-if="openModal" class="fixed inset-0 flex items-center justify-center bg-gray-100/50 bg-opacity-50 z-50">
      <div class="bg-white w-2/5 p-6 rounded shadow-lg relative">
        <h2 class="text-xl font-title font-bold mb-4 text-center">New animal</h2>

        <form @submit.prevent="handleAdd">
          <div class="mb-2">
            <label class="block font-semibold" for="name">Name</label>
            <input id="name" v-model="formData.name" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="species">Species</label>
            <input id="species" v-model="formData.species" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="birthday">Birthday</label>
            <input id="birthday" v-model="formData.birthday" type="date"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="genre">Genre</label>
            <input id="genre" v-model="formData.genre" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="diet">Diet</label>
            <input id="diet" v-model="formData.diet" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="condition">Condition</label>
            <input id="condition" v-model="formData.condition" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="notes">Notes</label>
            <textarea id="notes" v-model="formData.notes"
              class="w-full border border-gray-300 rounded px-2 py-1"></textarea>
          </div>

          <div class="flex justify-end space-x-2 mt-4">
            <Button :info="clearButton" @click="clearForm" />
            <Button :info="addButton" @click="handleAdd" />
          </div>
        </form>

        <button @click="closeModal" class="absolute top-2 right-2 text-gray-500 hover:text-gray-700">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
            stroke="currentColor" class="size-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import Button from './GenericButton.vue'

const props = defineProps<{
  openModal: boolean
}>()

const emit = defineEmits(['closeNewAnimal'])

const clearButton: Button = {
  text: 'Clear'
}

const addButton: Button = {
  text: 'Add',
  color: 'green',
  type: 'submit'
}

const formData = reactive({
  name: '',
  species: '',
  birthday: '',
  genre: '',
  diet: '',
  condition: '',
  notes: ''
})

function closeModal() {
  emit('closeNewAnimal')
}

function clearForm() {
  formData.name = ''
  formData.species = ''
  formData.birthday = ''
  formData.genre = ''
  formData.diet = ''
  formData.condition = ''
  formData.notes = ''
}

function handleAdd() {
  console.log('Datos del formulario:', { ...formData })

  closeModal()
  clearForm()
}
</script>
