<template>
  <div class="p-4">
    <table class="min-w-full">
      <thead>
        <tr>
          <th class="py-2 px-2">
            <div class="flex">
              <input type="checkbox" :checked="allSelected" @change="toggleAll" />
            </div>
          </th>
          <th v-for="column in columns" :key="column.field" scope="col" @click="sortBy(column.field)" :class="[
            'font-head px-2 py-2 text-left text-xs font-semibold text-black uppercase tracking-wide',
            column.field !== 'notes' ? 'cursor-pointer' : ''
          ]">
            {{ column.field }}
            <span v-if="sortKey === column.field && sortKey !== 'notes'">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
        </tr>
      </thead>
      <tbody class="font-body text-4xl bg-white">
        <tr v-for="animal in sortedAnimals" :key="animal._id" class="hover:bg-gray-100">
          <td class="px-2 py-2 w-4">
            <div class="flex items-center h-full">
              <div class="flex">
                <input type="checkbox" :checked="animal.selected" @change="toggleRow(animal)" />
              </div>
              <div class="pl-3" >
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="size-5 text-gray-700" @click="$emit('editAnimal', animal)">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L6.832 19.82a4.5 4.5 0 0 1-1.897 1.13l-2.685.8.8-2.685a4.5 4.5 0 0 1 1.13-1.897L16.863 4.487Zm0 0L19.5 7.125" />
                </svg>
              </div>
            </div>
          </td>
          <td v-for="column in columns" :key="column.field" class="px-2 py-1 text-sm text-gray-700 border border-black">
            {{ animal[column.field] }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import IAnimal from '../types/Animal'

const props = defineProps<{
  animals: IAnimal[]
}>()

interface Column {
  field: keyof Omit<IAnimal, '_id'>;
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

const emit = defineEmits<{
  (e: 'selectionChanged'): void
  (e: 'editAnimal', animal: IAnimal): void
}>()

// Select animals

const allSelected = computed<boolean>({
  get() {
    return props.animals.length > 0 && props.animals.every((a) => a.selected);
  },
  set(value: boolean) {
    props.animals.forEach((a) => a.selected = value)
  }
});

const toggleAll = () => {
  allSelected.value = !allSelected.value
  emit('selectionChanged', allSelected.value ? 1 : 0)
};

const toggleRow = (animal: IAnimal) => {
  animal.selected = !animal.selected
  emit('selectionChanged', props.animals.filter((a) => a.selected).length)
};

// Sort animal

const sortKey = ref<keyof IAnimal | ''>('')
const sortOrder = ref<'asc' | 'desc'>('asc')

function sortBy(column: keyof IAnimal) {
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

  if (sortKey.value === 'notes') {
    return props.animals
  }

  const sorted = [...props.animals]
  sorted.sort((a, b) => {
    const valA = (a as any)[sortKey.value]
    const valB = (b as any)[sortKey.value]

    if (!valA || !valB) {
      return 0
    }

    if (sortKey.value === 'birthday') {
      const dateA = new Date(valA)
      const dateB = new Date(valB)

      return sortOrder.value === 'asc'
        ? dateA - dateB
        : dateB - dateA
    }

    const comparison = valA.toString().localeCompare(valB.toString())
    return sortOrder.value === 'asc' ? comparison : -comparison
  })
  return sorted
})
</script>
