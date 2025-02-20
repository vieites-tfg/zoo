<template>
  <div class="p-4">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          <th
            v-for="column in columns"
            :key="column.field"
            scope="col"
            class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer"
            @click="sortBy(column.field)"
          >
            {{ column.field }}
            <span v-if="sortKey === column.field">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr
          v-for="animal in sortedAnimals"
          :key="animal._id"
          class="hover:bg-gray-100"
        >
          <td
            v-for="column in columns"
            :key="column.field"
            class="px-4 py-2 text-sm text-gray-700"
          >
            {{ animal[column.field] }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Animal from '@/types/Animal'

const props = defineProps<{
  animals: Animal[]
}>()

interface Column {
  field: keyof Omit<Animal, '_id'>;
}

const columns: Column[] = [
  { field: 'name' },
  { field: 'species' },
  { field: 'birthday' },
  { field: 'genre' },
  { field: 'diet' },
  { field: 'condition' },
  { field: 'notes' },
];

const sortKey = ref<keyof Animal | ''>('')
const sortOrder = ref<'asc' | 'desc'>('asc')

function sortBy(column: keyof Animal) {
  if (sortKey.value === column) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = column
    sortOrder.value = 'asc'
  }
}

const sortedAnimals = computed(() => {
  if (!sortKey.value) {
    return props.animals
  }

  const sorted = [...props.animals]
  sorted.sort((a, b) => {
    const valA = (a as any)[sortKey.value]
    const valB = (b as any)[sortKey.value]

    if (!valA || !valB) {
      return 0
    }

    const comparison = valA.toString().localeCompare(valB.toString())
    return sortOrder.value === 'asc' ? comparison : -comparison
  })
  return sorted
})
</script>

<style scoped></style>
