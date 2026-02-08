<template>
  <component 
    :is="to ? 'router-link' : 'div'" 
    :to="to"
    class="bg-gray-800 rounded-lg p-5 border border-gray-700 hover:border-gray-600 transition-colors"
    :class="{ 'cursor-pointer': to }"
  >
    <div class="flex items-center justify-between">
      <div>
        <p class="text-gray-400 text-sm">{{ label }}</p>
        <p class="text-2xl font-semibold text-white mt-1">
          <span v-if="loading" class="inline-block w-8 h-6 bg-gray-700 rounded animate-pulse"></span>
          <span v-else>{{ value }}</span>
        </p>
      </div>
      <component :is="icon" class="w-8 h-8 text-gray-600" />
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
