<template>
  <div class="p-3 bg-fail-s border border-fail/30 rounded-lg">
    <p class="text-fail text-sm">{{ message }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  message: string;
}>();
</script>
