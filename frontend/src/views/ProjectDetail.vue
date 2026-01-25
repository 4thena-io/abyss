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
      <div class="flex items-start justify-between">
        <div class="flex items-center gap-4">
          <div class="w-14 h-14 rounded-lg bg-gray-800 border border-gray-700 flex items-center justify-center">
            <FolderIcon class="w-7 h-7 text-gray-400" />
          </div>
          <div>
            <h1 class="text-2xl font-semibold text-white">{{ project.name }}</h1>
            <p class="text-gray-400 text-sm mt-1">{{ project.description || 'No description' }}</p>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <Tabs v-model="activeTab" :tabs="tabs" />

      <!-- Tab Content -->
      <div class="mt-6">
        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'" class="space-y-6">
          <div class="bg-gray-800 border border-gray-700 rounded-lg p-6">
            <h2 class="text-lg font-medium text-white mb-4">Details</h2>
            <dl class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <dt class="text-gray-400 text-sm">Name</dt>
                <dd class="text-white mt-1">{{ project.name }}</dd>
              </div>
              <div>
                <dt class="text-gray-400 text-sm">Description</dt>
                <dd class="text-white mt-1">{{ project.description || 'No description' }}</dd>
              </div>
              <div>
                <dt class="text-gray-400 text-sm">Applications</dt>
                <dd class="text-white mt-1">{{ apps.length }} apps</dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Apps Tab -->
        <div v-if="activeTab === 'apps'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-medium text-white">Applications</h2>
            <button @click="fetchApps" class="flex items-center gap-2 px-3 py-1.5 text-sm text-gray-400 
                hover:text-white transition-colors">
              <ArrowPathIcon :class="['w-4 h-4', loadingApps ? 'animate-spin' : '']" />
              <span>Refresh</span>
            </button>
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
            <router-link v-for="app in apps" :key="app.id" :to="`/apps/${app.id}`" class="bg-gray-800 border border-gray-700 rounded-lg p-5 
                     hover:border-gray-600 hover:bg-gray-750 transition-all group">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-gray-700 flex items-center justify-center">
                  <CubeIcon class="w-5 h-5 text-gray-400" />
                </div>
                <div>
                  <h3 class="text-white font-medium group-hover:text-blue-400 transition-colors">
                    {{ app.name }}
                  </h3>
                  <p class="text-gray-500 text-sm">{{ app.kind }}</p>
                </div>
              </div>

              <p v-if="app.description" class="mt-3 text-gray-400 text-sm line-clamp-2">
                {{ app.description }}
              </p>

              <div class="mt-4 flex items-center gap-2 flex-wrap">
                <span class="px-2 py-1 bg-gray-700 rounded text-xs text-gray-300">
                  {{ app.kind }}
                </span>
                <span v-if="app.language" class="px-2 py-1 bg-gray-700 rounded text-xs text-gray-300">
                  {{ app.language }}
                </span>
              </div>
            </router-link>
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
  CubeIcon,
  ArrowPathIcon
} from '@heroicons/vue/24/outline';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
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
