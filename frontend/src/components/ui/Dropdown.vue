<template>
  <div>
    <label v-if="label" class="block text-gray-400 text-sm mb-2">
      {{ label }}{{ required ? ' *' : '' }}
    </label>
    <div class="relative" ref="dropdownRef">
      <button 
        @click="isOpen = !isOpen"
        type="button"
        :class="buttonClasses"
      >
        <span :class="{ 'text-gray-500': !modelValue && placeholder }">
          {{ displayValue }}
        </span>
        <ChevronDownIcon class="w-4 h-4 ml-auto shrink-0" />
      </button>
      <div 
        v-if="isOpen"
        class="absolute z-10 w-full mt-1 bg-gray-800 border border-gray-700 rounded-lg overflow-hidden max-h-60 overflow-y-auto"
      >
        <button 
          v-if="allLabel"
          @click="selectOption('')"
          type="button"
          :class="[
            'w-full px-3 py-2 text-left text-sm transition-colors',
            !modelValue ? 'bg-gray-700 text-white' : 'text-gray-300 hover:bg-gray-700'
          ]"
        >
          {{ allLabel }}
        </button>
        <button 
          v-for="option in options" 
          :key="option.value"
          type="button"
          @click="selectOption(option.value)"
          :class="[
            'w-full px-3 py-2 text-left text-sm transition-colors',
            modelValue === option.value ? 'bg-gray-700 text-white' : 'text-gray-300 hover:bg-gray-700'
          ]"
        >
          {{ option.label }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { ChevronDownIcon } from '@heroicons/vue/24/outline';

const props = withDefaults(defineProps<{
  modelValue: string | number;
  options: Array<{ value: string | number; label: string }>;
  placeholder?: string;
  allLabel?: string;
  label?: string;
  required?: boolean;
  variant?: 'default' | 'form';
}>(), {
  placeholder: 'Select...',
  required: false,
  variant: 'default',
});

const emit = defineEmits<{
  'update:modelValue': [value: string | number];
}>();

const isOpen = ref(false);
const dropdownRef = ref<HTMLElement | null>(null);

const buttonClasses = computed(() => {
  const base = 'flex items-center gap-2 border border-gray-700 rounded-lg transition-colors w-full text-left';
  if (props.variant === 'form') {
    return `${base} bg-gray-900 px-4 py-2.5 text-white hover:border-gray-600`;
  }
  return `${base} bg-gray-800 px-3 py-2 text-sm text-gray-300 hover:border-gray-600 min-w-32`;
});

const displayValue = computed(() => {
  if (!props.modelValue && props.modelValue !== 0) {
    return props.allLabel || props.placeholder;
  }
  const selected = props.options.find(o => o.value === props.modelValue);
  return selected?.label || props.modelValue;
});

const selectOption = (value: string | number) => {
  emit('update:modelValue', value);
  isOpen.value = false;
};

const handleClickOutside = (event: Event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false;
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>
