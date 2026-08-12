<template>
  <component
    :is="to ? 'router-link' : 'div'"
    :to="to"
    class="bg-panel rounded-lg p-5 border border-b1 hover:border-b2 transition-colors"
    :class="{ 'cursor-pointer': to }"
  >
    <div class="flex items-center justify-between">
      <div>
        <p class="text-t2 text-sm">{{ label }}</p>
        <p class="text-2xl font-semibold text-t1 mt-1">
          <span v-if="loading" class="inline-block w-8 h-6 bg-surface rounded animate-pulse"></span>
          <span v-else>{{ value }}</span>
        </p>
      </div>
      <component :is="icon" class="w-8 h-8 text-t3" />
    </div>
  </component>
</template>

<script setup lang="ts">
import type { Component } from 'vue';

defineProps<{
  label: string;
  value: string | number;
  icon: Component;
  to?: string;
  loading?: boolean;
}>();
</script>
