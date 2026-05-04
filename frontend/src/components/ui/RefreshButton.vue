<template>
  <button
    @click="$emit('click')"
    :class="[
      'flex items-center gap-2 text-t2 hover:text-t1 transition-colors',
      size === 'sm' ? 'px-3 py-1.5 text-sm' : 'px-4 py-2 text-base'
    ]"
  >
    <ArrowPathIcon :class="[
      size === 'sm' ? 'w-4 h-4' : 'w-5 h-5',
      loading && 'animate-spin'
    ]" />
    <span>{{ label }}</span>
  </button>
</template>

<script setup lang="ts">
import { ArrowPathIcon } from '@heroicons/vue/24/outline';

withDefaults(defineProps<{
  loading?: boolean;
  label?: string;
  size?: 'sm' | 'md';
}>(), {
  loading: false,
  label: 'Refresh',
  size: 'sm',
});

defineEmits<{
  click: [];
}>();
</script>
