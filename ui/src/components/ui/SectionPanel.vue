<template>
  <div class="bg-panel rounded-lg border border-b1">
    <div class="px-5 py-4 border-b border-b1 flex items-center justify-between">
      <h2 class="text-lg font-medium text-t1">{{ title }}</h2>
      <slot name="header-actions" />
    </div>
    <div :class="noPadding ? '' : 'p-5'">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string;
    noPadding?: boolean;
  }>(),
  {
    noPadding: false,
  },
);
</script>
