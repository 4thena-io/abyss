<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <!-- Not Found -->
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
      <!-- Header -->
      <DetailHeader 
        :icon="FolderIcon" 
        :title="project.name" 
        :subtitle="project.description" 
      />

      <!-- Tabs -->
      <Tabs v-model="activeTab" :tabs="tabs" />

      <!-- Tab Content -->
      <div class="mt-6">
        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'" class="space-y-6">
          <DefinitionList 
            title="Details"
            :items="[
              { label: 'Name', value: project.name },
              { label: 'Description', value: project.description, fallback: 'No description' },
              { label: 'Applications', value: `${apps.length} apps` },
            ]"
          />
        </div>

        <!-- Apps Tab -->
        <div v-if="activeTab === 'apps'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-medium text-white">Applications</h2>
            <RefreshButton :loading="loadingApps" @click="fetchApps" />
          </div>

          <!-- Loading Apps -->
          <div v-if="loadingApps" class="flex items-center justify-center py-8">
            <Spinner size="md" />
          </div>

          <!-- No Apps -->
          <EmptyState 
            v-else-if="apps.length === 0"
            :icon="CubeIcon"
            title="No applications"
            message="No applications in this project"
          />

          <!-- Apps Grid -->
          <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            <EntityCard 
              v-for="app in apps" 
              :key="app.id"
              :icon="CubeIcon"
              :title="app.name"
              :subtitle="app.kind"
              :description="app.description"
              :show-no-description="false"
              :to="`/apps/${app.id}`"
            >
              <template #badges>
                <TagGroup :tags="[app.kind, app.language].filter(Boolean) as string[]" />
              </template>
            </EntityCard>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import {
  FolderIcon,
  CubeIcon
} from '@heroicons/vue/24/outline';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import EntityCard from '../components/ui/EntityCard.vue';
import TagGroup from '../components/ui/TagGroup.vue';
import RefreshButton from '../components/ui/RefreshButton.vue';
import DetailHeader from '../components/ui/DetailHeader.vue';
import DefinitionList from '../components/ui/DefinitionList.vue';
import type { Project } from '../types/Project';
import type { App } from '../types/App';

const route = useRoute();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'apps', label: 'Applications' },
];

const project = ref<Project | null>(null);
const apps = ref<App[]>([]);
const loading = ref(true);
const loadingApps = ref(false);
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

const fetchApps = async () => {
  if (!project.value) return;
  try {
    loadingApps.value = true;
    const response = await fetch(`/api/projects/${project.value.id}/apps`);
    if (!response.ok) throw new Error('Failed to fetch apps');
    apps.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch apps:', error);
    apps.value = [];
  } finally {
    loadingApps.value = false;
  }
};

onMounted(async () => {
  await fetchProject();
  if (project.value) {
    fetchApps();
  }
});
</script>
