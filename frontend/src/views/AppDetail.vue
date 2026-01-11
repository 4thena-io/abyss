<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="w-8 h-8 border-2 border-gray-600 border-t-blue-500 rounded-full animate-spin" />
    </div>

    <!-- Not Found -->
    <div v-else-if="!app" class="bg-gray-800 border border-gray-700 rounded-lg py-16 text-center">
      <CubeIcon class="w-12 h-12 text-gray-600 mx-auto" />
      <h3 class="mt-4 text-lg font-medium text-white">Application not found</h3>
      <router-link to="/apps" class="mt-4 inline-block text-blue-400 hover:text-blue-300">
        Back to applications
      </router-link>
    </div>

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
            <span>Repository</span>
          </a>
          <a v-if="app.ciUrl" :href="formatUrl(app.ciUrl)" target="_blank" rel="noopener noreferrer" class="flex items-center gap-2 px-4 py-2 bg-gray-800 border border-gray-700 
              rounded-lg text-gray-300 hover:text-white hover:border-gray-600 transition-colors">
            <WrenchScrewdriverIcon class="w-4 h-4" />
            <span>Pipelines</span>
          </a>
        </div>
      </div>

      <!-- Tabs -->
      <div class="border-b border-gray-700">
        <nav class="flex gap-6">
          <button v-for="tab in tabs" :key="tab.id" @click="activeTab = tab.id" :class="[
            'pb-3 text-sm font-medium transition-colors relative',
            activeTab === tab.id
              ? 'text-white'
              : 'text-gray-400 hover:text-gray-300'
          ]">
            {{ tab.label }}
            <div v-if="activeTab === tab.id" class="absolute bottom-0 left-0 right-0 h-0.5 bg-blue-500" />
          </button>
        </nav>
      </div>

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
            <h2 class="text-lg font-medium text-white">CI Builds</h2>
            <button @click="fetchBuilds" class="flex items-center gap-2 px-3 py-1.5 text-sm text-gray-400 
                hover:text-white transition-colors">
              <ArrowPathIcon :class="['w-4 h-4', loadingBuilds ? 'animate-spin' : '']" />
              <span>Refresh</span>
            </button>
          </div>

          <!-- Loading Builds -->
          <div v-if="loadingBuilds" class="flex items-center justify-center py-8">
            <div class="w-6 h-6 border-2 border-gray-600 border-t-blue-500 rounded-full animate-spin" />
          </div>

          <!-- No Builds -->
          <div v-else-if="builds.length === 0" class="bg-gray-800 border border-gray-700 rounded-lg py-12 text-center">
            <WrenchScrewdriverIcon class="w-10 h-10 text-gray-600 mx-auto" />
            <p class="mt-3 text-gray-400">No builds yet</p>
          </div>

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
                <tr v-for="build in builds" :key="build.id" class="hover:bg-gray-750 transition-colors">
                  <td class="px-5 py-4">
                    <a :href="build.link" target="_blank" rel="noopener noreferrer"
                      class="text-blue-400 hover:text-blue-300 font-medium transition-colors">
                      #{{ build.number }}
                    </a>
                  </td>
                  <td class="px-5 py-4">
                    <span :class="[
                      'inline-flex items-center gap-1.5 px-2 py-1 rounded text-xs font-medium',
                      build.status === 'success' ? 'bg-green-500/20 text-green-400' :
                        build.status === 'running' ? 'bg-blue-500/20 text-blue-400' :
                          build.status === 'pending' ? 'bg-yellow-500/20 text-yellow-400' :
                            build.status === 'failure' ? 'bg-red-500/20 text-red-400' :
                              'bg-gray-500/20 text-gray-400'
                    ]">
                      <span :class="[
                        'w-1.5 h-1.5 rounded-full',
                        build.status === 'success' ? 'bg-green-400' :
                          build.status === 'running' ? 'bg-blue-400 animate-pulse' :
                            build.status === 'pending' ? 'bg-yellow-400' :
                              build.status === 'failure' ? 'bg-red-400' : 'bg-gray-400'
                      ]" />
                      {{ build.status }}
                    </span>
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
  CubeIcon,
  CodeBracketIcon,
  WrenchScrewdriverIcon,
  ArrowPathIcon
} from '@heroicons/vue/24/outline';
import type { App } from '../types/App';
import type { Build } from '../types/Build';

const route = useRoute();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'builds', label: 'Builds' },
];

const app = ref<App | null>(null);
const builds = ref<Build[]>([]);
const loading = ref(true);
const loadingBuilds = ref(false);
const activeTab = ref('overview');

const formatUrl = (url: string) => {
  if (!url) return '#';
  if (url.startsWith('http://') || url.startsWith('https://')) return url;
  return `https://${url}`;
};

const formatDuration = (seconds: number) => {
  if (!seconds) return '-';
  if (seconds < 60) return `${seconds}s`;
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}m ${secs}s`;
};

const formatTime = (dateString: string) => {
  if (!dateString) return '-';
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const mins = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (mins < 1) return 'just now';
  if (mins < 60) return `${mins}m ago`;
  if (hours < 24) return `${hours}h ago`;
  if (days < 7) return `${days}d ago`;
  return date.toLocaleDateString();
};

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

const fetchBuilds = async () => {
  if (!app.value) return;
  try {
    loadingBuilds.value = true;
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
