<template>
  <div v-if="totalPages > 1" class="flex items-center justify-between px-5 py-3 border-t border-b1">
    <span class="text-sm text-t3"> Showing {{ startItem }} to {{ endItem }} of {{ total }} </span>

    <div class="flex items-center gap-1">
      <button
        @click="emit('update:modelValue', modelValue - 1)"
        :disabled="modelValue === 1"
        class="p-1.5 rounded text-t3 hover:text-t1 hover:bg-surface disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        <ChevronLeftIcon class="w-4 h-4" />
      </button>

      <template v-for="page in totalPages" :key="page">
        <button
          v-if="
            page === 1 || page === totalPages || (page >= modelValue - 1 && page <= modelValue + 1)
          "
          @click="emit('update:modelValue', page)"
          :class="[
            'px-3 py-1 rounded text-sm transition-colors',
            page === modelValue ? 'bg-acc text-white' : 'text-t3 hover:text-t1 hover:bg-surface',
          ]"
        >
          {{ page }}
        </button>
        <span v-else-if="page === modelValue - 2 || page === modelValue + 2" class="px-2 text-t3">
          ...
        </span>
      </template>

      <button
        @click="emit('update:modelValue', modelValue + 1)"
        :disabled="modelValue === totalPages"
        class="p-1.5 rounded text-t3 hover:text-t1 hover:bg-surface disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        <ChevronRightIcon class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline';

const props = defineProps<{
  modelValue: number;
  total: number;
  pageSize: number;
}>();

const emit = defineEmits<{
  'update:modelValue': [page: number];
}>();

const totalPages = computed(() => Math.ceil(props.total / props.pageSize));
const startItem = computed(() => (props.modelValue - 1) * props.pageSize + 1);
const endItem = computed(() => Math.min(props.modelValue * props.pageSize, props.total));
</script>
