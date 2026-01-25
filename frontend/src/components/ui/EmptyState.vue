<template>
  <div class="bg-gray-800 border border-gray-700 rounded-lg py-16 text-center">
    <component :is="icon" class="w-12 h-12 text-gray-600 mx-auto" />
    <h3 class="mt-4 text-lg font-medium text-white">{{ title }}</h3>
    <p v-if="message" class="mt-2 text-gray-400 text-sm">{{ message }}</p>
    <div v-if="$slots.default" class="mt-4">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue';

defineProps<{
  icon: Component;
  title: string;
  message?: string;
}>();
</script>
