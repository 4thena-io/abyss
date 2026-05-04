<template>
  <div class="flex items-start justify-between">
    <div class="flex items-center gap-4">
      <div class="w-14 h-14 rounded-lg bg-panel border border-b1 flex items-center justify-center">
        <component :is="icon" class="w-7 h-7 text-t3" />
      </div>
      <div>
        <h1 class="text-2xl font-semibold text-t1">{{ title }}</h1>
        <p class="text-t3 text-sm mt-1">{{ subtitle || 'No description' }}</p>
      </div>
    </div>
    <div v-if="$slots.actions" class="flex items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue';

defineProps<{
  icon: Component;
  title: string;
  subtitle?: string;
}>();
</script>
