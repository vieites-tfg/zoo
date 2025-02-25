<template>
  <div class="p-4">
    <table class="min-w-full divide-y divide-gray-200">
      <thead>
        <tr>
          <th class="px-4 py-2">
            <div class="flex justify-center items-center">
              <input type="checkbox" :checked="allSelected" @change="toggleAll" />
            </div>
          </th>
          <th v-for="column in columns" :key="column.field" scope="col"
            class="font-head px-2 py-2 text-left text-xs font-semibold text-black uppercase tracking-wide cursor-pointer"
            @click="sortBy(column.field)">
            {{ column.field }}
            <span v-if="sortKey === column.field">
              {{ sortOrder === 'asc' ? '▲' : '▼' }}
            </span>
          </th>
        </tr>
      </thead>
      <tbody class="font-body text-4xl bg-white">
        <tr v-for="(animal, index) in sortedAnimals" :key="animal._id" class="hover:bg-gray-100">
          <td>
            <div class="flex justify-center items-center">
              <input type="checkbox" :checked="animal.selected" @change="toggleRow(index)" />
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
import Animal from './types/Animal'

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

const emit = defineEmits(['selectionChanged']);

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

const toggleRow = (index: number) => {
  props.animals[index].selected = !props.animals[index].selected
  emit('selectionChanged', props.animals.filter((a) => a.selected).length)
};

// Sort animal

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
