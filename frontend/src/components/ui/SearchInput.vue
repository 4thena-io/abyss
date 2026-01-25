<template>
  <div :class="['relative flex-1', maxWidth]">
    <MagnifyingGlassIcon class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
    <input 
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      type="text" 
      :placeholder="placeholder"
      class="w-full bg-gray-800 border border-gray-700 rounded-lg pl-10 pr-4 py-2 
             text-sm text-white placeholder-gray-500 focus:outline-none focus:border-gray-600" 
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
