<template>
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-semibold text-white">{{ title }}</h1>
      <p v-if="subtitle" class="text-gray-400 mt-1">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string;
  subtitle?: string;
}>();
</script>
