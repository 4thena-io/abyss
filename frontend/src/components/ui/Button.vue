<template>
  <button :type="type" :disabled="disabled || loading" :class="[
    'inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-colors',
    sizeClasses,
    variantClasses,
    { 'cursor-not-allowed opacity-50': disabled || loading }
  ]">
    <Spinner v-if="loading" size="sm" />
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import Spinner from './Spinner.vue';

const props = withDefaults(defineProps<{
  type?: 'button' | 'submit';
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  loading?: boolean;
}>(), {
  type: 'button',
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
});

const sizeClasses = computed(() => ({
  'px-3 py-1.5 text-sm': props.size === 'sm',
  'px-4 py-2 text-sm': props.size === 'md',
  'px-5 py-2.5 text-base': props.size === 'lg',
}));

const variantClasses = computed(() => ({
  'bg-acc hover:opacity-90 text-white': props.variant === 'primary' && !props.disabled,
  'bg-surface text-t3': props.variant === 'primary' && props.disabled,
  'bg-panel border border-b1 hover:border-b2 text-t2 hover:text-t1': props.variant === 'secondary',
  'text-t2 hover:text-t1': props.variant === 'ghost',
  'bg-fail hover:opacity-90 text-white': props.variant === 'danger',
}));
</script>
