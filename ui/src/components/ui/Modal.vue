<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div class="absolute inset-0 bg-black/60" @click="$emit('close')" />
    <div
      class="relative bg-panel border border-b1 rounded-lg w-full flex flex-col max-h-[calc(100vh-2rem)]"
      :class="sizeClass"
    >
      <div class="flex items-center justify-between px-6 pt-6 pb-4 shrink-0">
        <h2 class="text-xl font-semibold text-t1">{{ title }}</h2>
        <button @click="$emit('close')" class="text-t3 hover:text-t1 transition-colors">
          <XMarkIcon class="w-5 h-5" />
        </button>
      </div>
      <div class="overflow-y-auto px-6 pb-6">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { XMarkIcon } from '@heroicons/vue/24/outline';

const props = defineProps<{
  open: boolean;
  title: string;
  size?: 'sm' | 'md' | 'lg';
}>();

defineEmits(['close']);

const sizeClass = computed(() => ({
  'max-w-sm': props.size === 'sm',
  'max-w-md': !props.size || props.size === 'md',
  'max-w-lg': props.size === 'lg',
}));
</script>
