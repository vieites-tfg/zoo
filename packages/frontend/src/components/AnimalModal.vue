<template>
  <div>
    <div v-if="openModal" class="fixed inset-0 flex items-center justify-center bg-gray-100/50 bg-opacity-50 z-50">
      <div class="bg-white w-2/5 p-6 rounded shadow-lg relative">
        <h2 class="text-xl font-title font-bold mb-4 text-center">{{ info.title }}</h2>

        <form class="grid grid-cols-2 gap-4" @submit.prevent="handleAdd">
          <div class="mb-2">
            <label class="block font-semibold" for="name">Name</label>
            <input id="name" v-model="info.data.name" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="species">Species</label>
            <input id="species" v-model="info.data.species" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="birthday">Birthday</label>
            <input id="birthday" v-model="formattedBirthday" type="date"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="genre">Genre</label>
            <input id="genre" v-model="info.data.genre" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="diet">Diet</label>
            <input id="diet" v-model="info.data.diet" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2">
            <label class="block font-semibold" for="condition">Condition</label>
            <input id="condition" v-model="info.data.condition" type="text"
              class="w-full border border-gray-300 rounded px-2 py-1" />
          </div>

          <div class="mb-2 col-span-2">
            <label class="block font-semibold" for="notes">Notes</label>
            <textarea id="notes" v-model="info.data.notes"
              class="w-full border border-gray-300 rounded px-2 py-1 min-h-8 max-h-56"></textarea>
          </div>

          <div class="col-start-2 flex justify-end space-x-2 mt-4">
            <slot/>
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
import IAnimalModal from '../types/AnimalModal'

const props = defineProps<{
  openModal: boolean,
  info: IAnimalModal,
}>()

const formattedBirthday = computed(() => {
  const [day, month, year] = props.info.data.birthday.split('/');
  return `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
})

const emit = defineEmits(['closeAnimalModal'])

function closeModal() {
  emit('closeAnimalModal')
}
</script>
