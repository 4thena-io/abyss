<template>
  <div class="space-y-6">
    <DefinitionList
      title="Details"
      :items="[
        { label: 'Name', value: project.name },
        { label: 'Description', value: project.description, fallback: 'No description' },
      ]"
    />
  </div>
</template>

<script setup lang="ts">
import type { Project } from '../types';
import DefinitionList from '../../../components/ui/DefinitionList.vue';

defineProps<{ project: Project }>();
</script>
