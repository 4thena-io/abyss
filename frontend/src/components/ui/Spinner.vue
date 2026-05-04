<template>
  <div :class="[
    'border-2 border-b2 border-t-acc rounded-full animate-spin',
    sizeClasses
  ]" />
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(defineProps<{
  size?: 'sm' | 'md' | 'lg';
}>(), {
  size: 'md',
});

const sizeClasses = computed(() => ({
  'w-4 h-4': props.size === 'sm',
  'w-6 h-6': props.size === 'md',
  'w-8 h-8': props.size === 'lg',
}));
</script>
