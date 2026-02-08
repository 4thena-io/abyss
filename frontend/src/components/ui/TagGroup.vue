<template>
  <div class="flex items-center gap-2 flex-wrap">
    <span 
      v-for="tag in tags" 
      :key="tag"
      class="px-2 py-1 bg-gray-700 rounded text-xs text-gray-300"
    >
      {{ tag }}
    </span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  tags: string[];
}>();
</script>
