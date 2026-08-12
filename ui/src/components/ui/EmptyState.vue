<template>
  <div class="bg-panel border border-b1 rounded-lg py-16 text-center">
    <component :is="icon" class="w-12 h-12 text-t3 mx-auto" />
    <h3 class="mt-4 text-lg font-medium text-t1">{{ title }}</h3>
    <p v-if="message" class="mt-2 text-t2 text-sm">{{ message }}</p>
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
