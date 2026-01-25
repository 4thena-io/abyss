<template>
  <div class="border-b border-gray-700">
    <nav class="flex gap-6">
      <button 
        v-for="tab in tabs" 
        :key="tab.id" 
        @click="$emit('update:modelValue', tab.id)"
        :class="[
          'pb-3 text-sm font-medium transition-colors relative',
          modelValue === tab.id
            ? 'text-white'
            : 'text-gray-400 hover:text-gray-300'
        ]"
      >
        {{ tab.label }}
        <div 
          v-if="modelValue === tab.id" 
          class="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-500" 
        />
      </button>
    </nav>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string;
  tabs: Array<{ id: string; label: string }>;
}>();

defineEmits<{
  'update:modelValue': [value: string];
}>();
</script>
