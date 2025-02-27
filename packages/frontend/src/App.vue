<template>
  <div class="w-4/5 mx-auto p-6 pt-[25vh]">
    <Title />
    <ActionButtons :isDeleteEnabled="isDeleteEnabled" @newAnimal="newAnimal" />
    <AnimalTable :animals="animals" @selectionChanged="updateSelection" @editAnimal="editAnimal" />
    <Modals :openModal="openModal" :animalModalData="animalModalData" @closeAnimalModal="closeAnimalModal"/>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Title from './components/Title.vue'
import Modals from './components/Modals.vue'
import ActionButtons from './components/ActionButtons.vue'
import AnimalTable from './components/AnimalTable.vue'
import IAnimal from './types/Animal'
import IAnimalModal from './types/AnimalModal'
import { getAllAnimals } from './api/animals';

const animals = ref<IAnimal[]>([])
const isDeleteEnabled = ref(false)
const openModal = ref(false)
const animalModalData = ref<IAnimalModal>({})

const resetButton: Button = {
  text: 'Reset'
}

const updateButton: Button = {
  text: 'Update',
  color: 'green',
  type: 'submit'
}

const updateSelection = (selectedCount: number) => {
  isDeleteEnabled.value = selectedCount > 0
};

const closeAnimalModal = () => {
  openModal.value = false
};

const newAnimal = () => {
  animalModalData.value = {
    action: 'Create',
    title: "New animal"
  }
  openModal.value = true
}

const editAnimal = (animal: IAnimal) => {
  animalModalData.value = {
    action: 'Update',
    title: "Update animal",
    data: animal
  }
  openModal.value = true
}

const formatDate = (date: string): string => {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = ("0" + (d.getMonth() + 1)).slice(-2);
  const day = ("0" + d.getDate()).slice(-2);
  return `${year}-${month}-${day}`;
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
