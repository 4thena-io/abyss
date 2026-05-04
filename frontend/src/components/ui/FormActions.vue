<template>
  <div class="flex items-center justify-end gap-3 pt-2">
    <button
      type="button"
      @click="$emit('cancel')"
      class="px-4 py-2 text-t2 hover:text-t1 transition-colors"
    >
      {{ cancelLabel }}
    </button>
    <button
      type="submit"
      :disabled="disabled || saving"
      :class="[
        'px-4 py-2 rounded-lg font-medium transition-colors',
        !disabled && !saving
          ? 'bg-acc hover:opacity-90 text-white'
          : 'bg-surface text-t3 cursor-not-allowed'
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
