<template>
  <div class="w-4/5 mx-auto p-6 pt-[25vh]">
    <MainTitle />
    <ActionButtons
      :is-delete-enabled="isDeleteEnabled"
      @new-animal="newAnimal"
    />
    <AnimalTable
      :animals="animals"
      @selection-changed="updateSelection"
      @edit-animal="editAnimal"
    />
    <ModalWindows
      :open-modal="openModal"
      :animal-modal-data="animalModalData"
      @close-animal-modal="closeAnimalModal"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { MainTitle, ModalWindows, ActionButtons, AnimalTable } from '../components'
import IAnimal from '../types/Animal'
import IAnimalModal from '../types/AnimalModal'
import { getAllAnimals } from '../api/animals';

const animals = ref<IAnimal[]>([])
const isDeleteEnabled = ref(false)
const openModal = ref(false)
const animalModalData = ref<IAnimalModal>({})

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
