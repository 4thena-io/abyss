<template>
  <div class="flex items-center justify-end gap-3 pt-2">
    <button 
      type="button" 
      @click="$emit('cancel')" 
      class="px-4 py-2 text-gray-400 hover:text-white transition-colors"
    >
      {{ cancelLabel }}
    </button>
    <button 
      type="submit" 
      :disabled="disabled || saving"
      :class="[
        'px-4 py-2 rounded-lg font-medium transition-colors',
        !disabled && !saving
          ? 'bg-blue-600 hover:bg-blue-500 text-white'
          : 'bg-gray-700 text-gray-500 cursor-not-allowed'
      ]"
    >
      {{ saving ? submittingLabel : submitLabel }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(defineProps<{
  submitLabel: string;
  submittingLabel?: string;
  cancelLabel?: string;
  disabled?: boolean;
  saving?: boolean;
}>(), {
  cancelLabel: 'Cancel',
  disabled: false,
  saving: false,
});

const submittingLabel = computed(() => 
  props.submittingLabel || props.submitLabel.replace(/^(\w+)/, '$1ing...')
);

defineEmits<{
  cancel: [];
}>();
</script>
