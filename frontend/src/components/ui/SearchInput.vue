<template>
  <div :class="['relative flex-1', maxWidth]">
    <MagnifyingGlassIcon class="w-4 h-4 text-t3 absolute left-3 top-1/2 -translate-y-1/2" />
    <input
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      type="text"
      :placeholder="placeholder"
      class="w-full bg-panel border border-b1 rounded-lg pl-10 pr-4 py-2
             text-sm text-t1 placeholder:text-t3 focus:outline-none focus:border-b2"
    />
  </div>
</template>

<script setup lang="ts">
import { MagnifyingGlassIcon } from '@heroicons/vue/24/outline';

withDefaults(defineProps<{
  modelValue: string;
  placeholder?: string;
  maxWidth?: string;
}>(), {
  placeholder: 'Search...',
  maxWidth: 'max-w-xs',
});

defineEmits<{
  'update:modelValue': [value: string];
}>();
</script>
