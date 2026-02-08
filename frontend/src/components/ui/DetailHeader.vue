<template>
  <div class="flex items-start justify-between">
    <div class="flex items-center gap-4">
      <div class="w-14 h-14 rounded-lg bg-gray-800 border border-gray-700 flex items-center justify-center">
        <component :is="icon" class="w-7 h-7 text-gray-400" />
      </div>
      <div>
        <h1 class="text-2xl font-semibold text-white">{{ title }}</h1>
        <p class="text-gray-400 text-sm mt-1">{{ subtitle || 'No description' }}</p>
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
