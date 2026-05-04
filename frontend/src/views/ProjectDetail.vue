<template>
  <div class="space-y-6">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <EmptyState
      v-else-if="!project"
      :icon="FolderIcon"
      title="Project not found"
    >
      <router-link to="/projects" class="text-blue-400 hover:text-blue-300">
        Back to projects
      </router-link>
    </EmptyState>

    <template v-else>
      <DetailHeader :icon="FolderIcon" :title="project.name" :subtitle="project.description" />

      <Tabs v-model="activeTab" :tabs="tabs" />

      <div class="mt-6">
        <ProjectOverviewTab v-if="activeTab === 'overview'" :project="project" />
        <ProjectAppsTab v-else-if="activeTab === 'apps'" :project-id="project.id" />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { FolderIcon } from '@heroicons/vue/24/outline';
import type { Project } from '../features/projects/types';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import DetailHeader from '../components/ui/DetailHeader.vue';
import ProjectOverviewTab from '../features/projects/components/ProjectOverviewTab.vue';
import ProjectAppsTab from '../features/projects/components/ProjectAppsTab.vue';

const route = useRoute();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'apps', label: 'Applications' },
];

const project = ref<Project | null>(null);
const loading = ref(true);
const activeTab = ref('overview');

const fetchProject = async () => {
  try {
    loading.value = true;
    const id = route.params.id;
    const response = await fetch(`/api/projects/${id}`);
    if (!response.ok) throw new Error('Project not found');
    project.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch project:', error);
    project.value = null;
  } finally {
    loading.value = false;
  }
};

onMounted(fetchProject);
</script>
