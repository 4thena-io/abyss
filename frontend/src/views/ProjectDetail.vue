<template>
  <div class="space-y-6">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <EmptyState v-else-if="!project" :icon="FolderIcon" title="Project not found">
      <router-link to="/projects" class="text-acc hover:opacity-80"> Back to projects </router-link>
    </EmptyState>

    <template v-else>
      <DetailHeader :icon="FolderIcon" :title="project.name" :subtitle="project.description" />

      <Tabs v-model="activeTab" :tabs="tabs" />

      <div class="mt-6">
        <ProjectOverviewTab v-if="activeTab === 'overview'" :project="project" />
        <ProjectAppsTab v-else-if="activeTab === 'apps'" :project-id="project.id" />
        <ProjectSettingsTab
          v-else-if="activeTab === 'settings'"
          :project="project"
          @updated="project = $event"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { FolderIcon } from '@heroicons/vue/24/outline';
import type { Project } from '../features/projects/types';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import DetailHeader from '../components/ui/DetailHeader.vue';
import ProjectOverviewTab from '../features/projects/components/ProjectOverviewTab.vue';
import ProjectAppsTab from '../features/projects/components/ProjectAppsTab.vue';
import ProjectSettingsTab from '../features/projects/components/ProjectSettingsTab.vue';

const route = useRoute();
const router = useRouter();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'apps', label: 'Applications' },
  { id: 'settings', label: 'Settings' },
];

const validTabs = tabs.map((t) => t.id);
const project = ref<Project | null>(null);
const loading = ref(true);
const activeTab = ref(
  validTabs.includes(route.query.tab as string) ? (route.query.tab as string) : 'overview',
);

watch(activeTab, (tab) => router.replace({ query: { tab } }));
watch(
  () => route.query.tab,
  (tab) => {
    if (validTabs.includes(tab as string) && tab !== activeTab.value) {
      activeTab.value = tab as string;
    }
  },
);

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
