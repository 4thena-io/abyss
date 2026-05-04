<template>
  <div>
    <DefinitionList
      title="Details"
      :items="[
        { label: 'Name', value: team.name },
        { label: 'Members', value: String(team.memberCount) },
        { label: 'Projects', value: String(team.projectCount) },
        { label: 'Applications', value: String(team.appCount) },
      ]"
    />
  </div>
</template>

<script setup lang="ts">
import type { Team } from '../types';
import DefinitionList from '../../../components/ui/DefinitionList.vue';

defineProps<{ team: Team }>();
</script>
