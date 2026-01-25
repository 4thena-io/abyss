<template>
  <span :class="[
    'inline-flex items-center gap-1.5 px-2 py-1 rounded text-xs font-medium',
    statusClasses
  ]">
    <span 
      v-if="showDot"
      :class="['w-1.5 h-1.5 rounded-full', dotClasses]" 
    />
    {{ status }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(defineProps<{
  status: 'success' | 'running' | 'pending' | 'failure' | string;
  showDot?: boolean;
}>(), {
  showDot: true,
});

const statusClasses = computed(() => {
  switch (props.status) {
    case 'success':
      return 'bg-green-500/20 text-green-400';
    case 'running':
      return 'bg-blue-500/20 text-blue-400';
    case 'pending':
      return 'bg-yellow-500/20 text-yellow-400';
    case 'failure':
      return 'bg-red-500/20 text-red-400';
    default:
      return 'bg-gray-500/20 text-gray-400';
  }
});

const dotClasses = computed(() => {
  switch (props.status) {
    case 'success':
      return 'bg-green-400';
    case 'running':
      return 'bg-blue-400 animate-pulse';
    case 'pending':
      return 'bg-yellow-400';
    case 'failure':
      return 'bg-red-400';
    default:
      return 'bg-gray-400';
  }
});
</script>
