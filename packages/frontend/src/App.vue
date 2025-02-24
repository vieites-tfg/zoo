<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AnimalTable from './components/AnimalTable.vue'
import Animal from './types/Animal'

const animals = ref<Animal[]>([])

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
      notes: 'Very social, interacts well with visitors.'
    },
    {
      _id: '2',
      name: 'Rocky',
      species: 'Lion',
      birthday: '2018-06-10',
      genre: 'Male',
      diet: 'Carnivore',
      condition: 'Injured paw',
      notes: 'Under special care due to paw injury.'
    }
  ]
})
</script>

<template>
  <div class="w-4/5 mx-auto p-6">
    <h1 class="font-title text-4xl font-bold mb-4 flex justify-center">Zoo</h1>
    <AnimalTable :animals="animals" />
  </div>
</template>

<style scoped>
</style>
