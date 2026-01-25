<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <!-- Not Found -->
    <EmptyState 
      v-else-if="!app"
      :icon="CubeIcon"
      title="Application not found"
    >
      <router-link to="/apps" class="text-blue-400 hover:text-blue-300">
        Back to applications
      </router-link>
    </EmptyState>

    <template v-else>
      <!-- Header -->
      <div class="flex items-start justify-between">
        <div class="flex items-center gap-4">
          <div class="w-14 h-14 rounded-lg bg-gray-800 border border-gray-700 flex items-center justify-center">
            <CubeIcon class="w-7 h-7 text-gray-400" />
          </div>
          <div>
            <h1 class="text-2xl font-semibold text-white">{{ app.name }}</h1>
            <p class="text-gray-400 text-sm mt-1">{{ app.description || 'No description' }}</p>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <a :href="formatUrl(app.repoUrl)" target="_blank" rel="noopener noreferrer" class="flex items-center gap-2 px-4 py-2 bg-gray-800 border border-gray-700 
              rounded-lg text-gray-300 hover:text-white hover:border-gray-600 transition-colors">
            <CodeBracketIcon class="w-4 h-4" />
            <span>Source Code</span>
          </a>
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
                <dt class="text-gray-400 text-sm">Description</dt>
                <dd class="text-white mt-1">{{ app.description || 'No description' }}</dd>
              </div>
              <div>
                <dt class="text-gray-400 text-sm">Project</dt>
                <dd class="text-white mt-1">
                  <router-link :to="`/projects/${app.projectId}`" class="text-blue-400 hover:text-blue-300">
                    Project #{{ app.projectId }}
                  </router-link>
                </dd>
              </div>
              <div>
                <dt class="text-gray-400 text-sm">Repository URL</dt>
                <dd class="mt-1">
                  <a :href="formatUrl(app.repoUrl)" target="_blank" rel="noopener noreferrer"
                    class="text-blue-400 hover:text-blue-300 break-all">
                    {{ app.repoUrl }}
                  </a>
                </dd>
              </div>
              <div>
                <dt class="text-gray-400 text-sm">CI URL</dt>
                <dd class="mt-1">
                  <a v-if="app.ciUrl" :href="formatUrl(app.ciUrl)" target="_blank" rel="noopener noreferrer"
                    class="text-blue-400 hover:text-blue-300 break-all">
                    {{ app.ciUrl }}
                  </a>
                  <span v-else class="text-gray-500">Not configured</span>
                </dd>
              </div>
              <div>
                <dt class="text-gray-400 text-sm">Template</dt>
                <dd class="text-white mt-1">
                  <router-link :to="`/templates/${app.templateId}`" class="text-blue-400 hover:text-blue-300">
                    Template #{{ app.templateId }}
                  </router-link>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Builds Tab -->
        <div v-if="activeTab === 'builds'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-medium text-white">Builds</h2>
            <button @click="fetchBuilds" class="flex items-center gap-2 px-3 py-1.5 text-sm text-gray-400 
                hover:text-white transition-colors">
              <ArrowPathIcon :class="['w-4 h-4', loadingBuilds ? 'animate-spin' : '']" />
              <span>Refresh</span>
            </button>
          </div>

          <!-- Loading Builds -->
          <div v-if="loadingBuilds" class="flex items-center justify-center py-8">
            <Spinner size="md" />
          </div>

          <!-- No Builds -->
          <EmptyState 
            v-else-if="builds.length === 0"
            :icon="WrenchScrewdriverIcon"
            title="No builds"
            message="No builds yet"
          />

          <!-- Builds List -->
          <div v-else class="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
            <table class="w-full">
              <thead>
                <tr class="text-left text-gray-400 text-sm border-b border-gray-700">
                  <th class="px-5 py-3 font-medium">Build</th>
                  <th class="px-5 py-3 font-medium">Status</th>
                  <th class="px-5 py-3 font-medium">Branch</th>
                  <th class="px-5 py-3 font-medium">Commit</th>
                  <th class="px-5 py-3 font-medium">Duration</th>
                  <th class="px-5 py-3 font-medium">Started</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-700">
                <tr v-for="build in paginatedBuilds" :key="build.id" class="hover:bg-gray-750 transition-colors">
                  <td class="px-5 py-4">
                    <a :href="build.link" target="_blank" rel="noopener noreferrer"
                      class="text-blue-400 hover:text-blue-300 font-medium transition-colors">
                      #{{ build.number }}
                    </a>
                  </td>
                  <td class="px-5 py-4">
                    <StatusBadge :status="build.status" />
                  </td>
                  <td class="px-5 py-4 text-gray-300 text-sm">{{ build.branch }}</td>
                  <td class="px-5 py-4">
                    <span class="text-gray-400 text-sm font-mono">{{ build.commit?.slice(0, 7) }}</span>
                  </td>
                  <td class="px-5 py-4 text-gray-400 text-sm">{{ formatDuration(build.duration) }}</td>
                  <td class="px-5 py-4 text-gray-500 text-sm">{{ formatTime(build.startedAt) }}</td>
                </tr>
              </tbody>
            </table>

            <!-- Pagination -->
            <div v-if="totalPages > 1" class="flex items-center justify-between px-5 py-3 border-t border-gray-700">
              <span class="text-sm text-gray-400">
                Showing {{ (currentPage - 1) * pageSize + 1 }} to {{ Math.min(currentPage * pageSize, builds.length) }}
                of {{ builds.length }}
              </span>

              <div class="flex items-center gap-1">
                <button @click="goToPage(currentPage - 1)" :disabled="currentPage === 1" class="p-1.5 rounded text-gray-400 hover:text-white hover:bg-gray-700 
               disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
                  <ChevronLeftIcon class="w-4 h-4" />
                </button>

                <template v-for="page in totalPages" :key="page">
                  <button
                    v-if="page === 1 || page === totalPages || (page >= currentPage - 1 && page <= currentPage + 1)"
                    @click="goToPage(page)" :class="[
                      'px-3 py-1 rounded text-sm transition-colors',
                      page === currentPage
                        ? 'bg-blue-600 text-white'
                        : 'text-gray-400 hover:text-white hover:bg-gray-700'
                    ]">
                    {{ page }}
                  </button>
                  <span v-else-if="page === currentPage - 2 || page === currentPage + 2" class="px-2 text-gray-500">
                    ...
                  </span>
                </template>

                <button @click="goToPage(currentPage + 1)" :disabled="currentPage === totalPages" class="p-1.5 rounded text-gray-400 hover:text-white hover:bg-gray-700 
               disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
                  <ChevronRightIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import {
  CubeIcon,
  CodeBracketIcon,
  WrenchScrewdriverIcon,
  ArrowPathIcon,
  ChevronLeftIcon,
  ChevronRightIcon
} from '@heroicons/vue/24/outline';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import { useFormatters } from '../composables/useFormatters';
import StatusBadge from '../components/ui/StatusBadge.vue';
import type { App } from '../types/App';
import type { Build } from '../types/Build';

const route = useRoute();
const { formatDuration, formatRelativeTime: formatTime, formatUrl } = useFormatters();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'builds', label: 'CI/CD' },
];

const app = ref<App | null>(null);
const builds = ref<Build[]>([]);
const loading = ref(true);
const loadingBuilds = ref(false);
const activeTab = ref('overview');

const fetchApp = async () => {
  try {
    loading.value = true;
    const id = route.params.id;
    const response = await fetch(`/api/apps/${id}`);
    if (!response.ok) throw new Error('App not found');
    app.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch app:', error);
    app.value = null;
  } finally {
    loading.value = false;
  }
};

const currentPage = ref(1);
const pageSize = 5;

const paginatedBuilds = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  return builds.value.slice(start, start + pageSize);
});

const totalPages = computed(() => Math.ceil(builds.value.length / pageSize));

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
  }
};

// Reset to page 1 when builds refresh
const fetchBuilds = async () => {
  if (!app.value) return;
  try {
    loadingBuilds.value = true;
    currentPage.value = 1; // Add this line
    const response = await fetch(`/api/apps/${app.value.id}/builds`);
    if (!response.ok) throw new Error('Failed to fetch builds');
    builds.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch builds:', error);
    builds.value = [];
  } finally {
    loadingBuilds.value = false;
  }
};

onMounted(async () => {
  await fetchApp();
  if (app.value) {
    fetchBuilds();
  }
});
</script>
