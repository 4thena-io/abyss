<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 px-2 py-1 rounded text-xs font-medium',
      statusClasses,
    ]"
  >
    <span v-if="showDot" :class="['w-1.5 h-1.5 rounded-full', dotClasses]" />
    {{ status }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
  defineProps<{
    status: 'success' | 'running' | 'pending' | 'failure' | string;
    showDot?: boolean;
  }>(),
  {
    showDot: true,
  },
);

const statusClasses = computed(() => {
  switch (props.status) {
    case 'success':
      return 'bg-ok-s text-ok';
    case 'running':
      return 'bg-run-s text-run';
    case 'pending':
      return 'bg-warn-s text-warn';
    case 'failure':
      return 'bg-fail-s text-fail';
    default:
      return 'bg-surface text-t2';
  }
});

const dotClasses = computed(() => {
  switch (props.status) {
    case 'success':
      return 'bg-ok';
    case 'running':
      return 'bg-run animate-pulse';
    case 'pending':
      return 'bg-warn';
    case 'failure':
      return 'bg-fail';
    default:
      return 'bg-t3';
  }
});
</script>
