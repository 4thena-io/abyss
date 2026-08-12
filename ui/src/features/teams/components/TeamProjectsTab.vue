<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium text-t1">Projects</h2>
      <RefreshButton :loading="loading" @click="fetchProjects" />
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <Spinner size="md" />
    </div>

    <EmptyState
      v-else-if="projects.length === 0"
      :icon="FolderIcon"
      title="No projects"
      message="This team has no projects yet"
    />

    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <EntityCard
        v-for="project in projects"
        :key="project.id"
        :icon="FolderIcon"
        :title="project.name"
        :description="project.description"
        :to="`/projects/${project.id}`"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { FolderIcon } from '@heroicons/vue/24/outline';
import type { Project } from '../../projects/types';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';
import RefreshButton from '../../../components/ui/RefreshButton.vue';
import EntityCard from '../../../components/ui/EntityCard.vue';

const props = defineProps<{ teamId: number }>();

const projects = ref<Project[]>([]);
const loading = ref(false);

const fetchProjects = async () => {
  try {
    loading.value = true;
    const response = await fetch(`/api/teams/${props.teamId}/projects`);
    if (!response.ok) throw new Error('Failed to fetch projects');
    projects.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch projects:', error);
    projects.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(fetchProjects);
</script>
