<template>
  <div class="w-4/5 mx-auto p-6 pt-[25vh]">
    <Title />
    <ActionButtons
      :isDeleteEnabled="isDeleteEnabled"
    />
    <AnimalTable
      :animals="animals"
      @selectionChanged="updateSelection"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Title from './components/Title.vue'
import ActionButtons from './components/ActionButtons.vue'
import AnimalTable from './components/AnimalTable.vue'
import Animal from './types/Animal'

const animals = ref<Animal[]>([])
const isDeleteEnabled = ref(false)

const updateSelection = (selectedCount: number) => {
  isDeleteEnabled.value = selectedCount > 0
};

onMounted(() => {
  animals.value = [
    {
      _id: '1',
      name: 'Bella',
      species: 'Giraffe',
      birthday: '2020-01-15',
      genre: 'Female',
      diet: 'Herbivore',
      condition: 'Healthy',
      notes: 'Very social, interacts well with visitors.',
      selected: false
    },
    {
      _id: '2',
      name: 'Rocky',
      species: 'Lion',
      birthday: '2018-06-10',
      genre: 'Male',
      diet: 'Carnivore',
      condition: 'Injured paw',
      notes: 'Under special care due to paw injury.',
      selected: false
    }
  ]
})
</script>
