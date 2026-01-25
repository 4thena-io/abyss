<template>
  <div class="bg-gray-800 rounded-lg border border-gray-700">
    <div class="px-5 py-4 border-b border-gray-700 flex items-center justify-between">
      <h2 class="text-lg font-medium text-white">{{ title }}</h2>
      <slot name="header-actions" />
    </div>
    <div :class="noPadding ? '' : 'p-5'">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title: string;
  noPadding?: boolean;
}>(), {
  noPadding: false,
});
</script>
