<template>
  <div>
    <label class="block text-gray-400 text-sm mb-2">
      {{ label }}{{ required ? ' *' : '' }}
    </label>
    
    <!-- Input -->
    <input 
      v-if="type === 'text'"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      type="text"
      :placeholder="placeholder"
      :disabled="disabled"
      :class="inputClasses"
    />
    
    <!-- Textarea -->
    <textarea 
      v-else-if="type === 'textarea'"
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
      :rows="rows"
      :placeholder="placeholder"
      :disabled="disabled"
      :class="[inputClasses, 'resize-none']"
    />
    
    <!-- Select -->
    <select
      v-else-if="type === 'select'"
      :value="modelValue"
      @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
      :disabled="disabled"
      :class="inputClasses"
    >
      <option v-if="placeholder" value="">{{ placeholder }}</option>
      <option 
        v-for="option in options" 
        :key="option.value" 
        :value="option.value"
      >
        {{ option.label }}
      </option>
    </select>
    
    <!-- Helper text -->
    <p v-if="hint" class="text-xs text-gray-500 mt-1">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(defineProps<{
  modelValue: string | number;
  label: string;
  type?: 'text' | 'textarea' | 'select';
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  rows?: number;
  hint?: string;
  options?: Array<{ value: string | number; label: string }>;
}>(), {
  type: 'text',
  required: false,
  disabled: false,
  rows: 3,
});

defineEmits<{
  'update:modelValue': [value: string | number];
}>();

const inputClasses = computed(() => [
  'w-full border border-gray-700 rounded-lg px-4 py-2.5 text-white',
  'placeholder-gray-500 focus:outline-none focus:border-gray-600',
  props.disabled 
    ? 'bg-gray-900/50 text-gray-400 cursor-not-allowed' 
    : 'bg-gray-900',
]);
</script>
