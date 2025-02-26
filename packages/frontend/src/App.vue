<template>
  <div class="w-4/5 mx-auto p-6 pt-[25vh]">
    <Title />
    <ActionButtons
      :isDeleteEnabled="isDeleteEnabled"
      @openNewAnimal="openNewAnimalModal"
    />
    <AnimalTable
      :animals="animals"
      @selectionChanged="updateSelection"
    />
    <NewAnimalModal 
    :openModal="openModal"
    @closeNewAnimal="closeNewAnimalModal"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Title from './components/Title.vue'
import ActionButtons from './components/ActionButtons.vue'
import AnimalTable from './components/AnimalTable.vue'
import NewAnimalModal from './components/NewAnimalModal.vue'
import Animal from './types/Animal'
import { getAllAnimals } from './api/animals';

const animals = ref<Animal[]>([])
const isDeleteEnabled = ref(false)
const openModal = ref(false)

const updateSelection = (selectedCount: number) => {
  isDeleteEnabled.value = selectedCount > 0
};

const openNewAnimalModal = () => {
  openModal.value = true
};

const closeNewAnimalModal = () => {
  openModal.value = false
};

function formatDate(date: string): string {
  const d = new Date(date)

  return d.toLocaleDateString(undefined, {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  })
}

async function loadAnimals() {
  try {
    const allAnimals = await getAllAnimals()

    const formattedAnimals = allAnimals.map((a) => {
      return {
        ...a,
        birthday: formatDate(a.birthday),
      }
    })

    animals.value = formattedAnimals
  } catch (error) {
    console.error('Could not load the animals:', error)
  }
}

onMounted(() => {
  loadAnimals()
})
</script>
