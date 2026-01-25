<template>
  <div class="p-3 bg-red-500/10 border border-red-500/50 rounded-lg">
    <p class="text-red-400 text-sm">{{ message }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  message: string;
}>();
</script>
