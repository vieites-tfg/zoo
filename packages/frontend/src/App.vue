<template>
  <div class="w-4/5 mx-auto p-6 pt-[25vh]">
    <Title />
    <ActionButtons :isDeleteEnabled="isDeleteEnabled" @newAnimal="newAnimal" />
    <AnimalTable :animals="animals" @selectionChanged="updateSelection" @editAnimal="editAnimal" />
    <AnimalModal :openModal="openModal" :info="animalModalData" @closeAnimalModal="closeAnimalModal" >
      <template v-if="animalModalData.action === 'Create'">
        <NewAnimalButtons @clearForm="clearForm" @createNewAnimal="createNewAnimal"/>
      </template>
      <template v-else>
        <UpdateAnimalButtons/>
      </template>
    </AnimalModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Title from './components/Title.vue'
import ActionButtons from './components/ActionButtons.vue'
import NewAnimalButtons from './components/NewAnimalButtons.vue'
import UpdateAnimalButtons from './components/UpdateAnimalButtons.vue'
import AnimalTable from './components/AnimalTable.vue'
import AnimalModal from './components/AnimalModal.vue'
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

const cleanedForm: IAnimal = {
  name: '',
  species: '',
  birthday: '',
  genre: '',
  diet: '',
  condition: '',
  notes: ''
}

const updateSelection = (selectedCount: number) => {
  isDeleteEnabled.value = selectedCount > 0
};

const openAnimalModal = () => {
  console.log(animalModalData.value)
  openModal.value = true
};

const closeAnimalModal = () => {
  openModal.value = false
};

const createNewAnimal = (animal: IAnimal) => {
  console.log(`Creating the animal: ${animal}`)
}

const clearForm = () => {
  animalModalData.value.data = { ...cleanedForm }
}

const newAnimal = () => {
  animalModalData.value = {
    action: 'Create',
    title: "New animal",
    data: { ...cleanedForm }
  }
  openAnimalModal()
}

const editAnimal = (animal: IAnimal) => {
  animalModalData.value = {
    action: 'Update',
    title: "Update animal",
    data: animal
  }
  openAnimalModal()
}

const formatDate = (date: string): string => {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = ("0" + (d.getMonth() + 1)).slice(-2);
  const day = ("0" + d.getDate()).slice(-2);
  return `${day}/${month}/${year}`;
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
