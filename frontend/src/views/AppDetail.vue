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
      <DetailHeader
        :icon="CubeIcon"
        :title="app.name"
        :subtitle="app.description"
      >
        <template #actions>
          <a :href="formatUrl(app.repoUrl)" target="_blank" rel="noopener noreferrer" class="flex items-center gap-2 px-4 py-2 bg-gray-800 border border-gray-700
              rounded-lg text-gray-300 hover:text-white hover:border-gray-600 transition-colors">
            <CodeBracketIcon class="w-4 h-4" />
            <span>Source Code</span>
          </a>
        </template>
      </DetailHeader>

      <!-- Tabs -->
      <Tabs v-model="activeTab" :tabs="tabs" />

      <!-- Tab Content -->
      <div class="mt-6">
        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'" class="space-y-6">
          <DefinitionList
            title="Details"
            :items="detailItems"
          />
        </div>

        <!-- Builds Tab -->
        <div v-if="activeTab === 'builds'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-medium text-white">Builds</h2>
            <RefreshButton :loading="loadingBuilds" @click="fetchBuilds" />
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
          <DataTable
            v-else
            :columns="buildColumns"
            :rows="paginatedBuilds"
            row-key="id"
          >
            <template #cell-number="{ row }">
              <a :href="row.link" target="_blank" rel="noopener noreferrer"
                class="text-blue-400 hover:text-blue-300 font-medium transition-colors">
                #{{ row.number }}
              </a>
            </template>
            <template #cell-status="{ row }">
              <StatusBadge :status="row.status" />
            </template>
            <template #cell-branch="{ row }">
              <span class="text-gray-300 text-sm">{{ row.branch }}</span>
            </template>
            <template #cell-commit="{ row }">
              <span class="text-gray-400 text-sm font-mono">{{ row.commit?.slice(0, 7) }}</span>
            </template>
            <template #cell-duration="{ row }">
              <span class="text-gray-400 text-sm">{{ formatDuration(row.duration) }}</span>
            </template>
            <template #cell-startedAt="{ row }">
              <span class="text-gray-500 text-sm">{{ formatTime(row.startedAt) }}</span>
            </template>
            <template #footer>
              <Pagination
                v-model="currentPage"
                :total="builds.length"
                :page-size="pageSize"
              />
            </template>
          </DataTable>
        </div>

        <!-- Deployments Tab -->
        <div v-if="activeTab === 'deployments'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-medium text-white">Deployments</h2>
            <RefreshButton :loading="loadingDeployments" @click="fetchDeployments" />
          </div>

          <!-- Loading Deployments -->
          <div v-if="loadingDeployments" class="flex items-center justify-center py-8">
            <Spinner size="md" />
          </div>

          <!-- No Deployments -->
          <EmptyState
            v-else-if="deployments.length === 0"
            :icon="ArrowUpTrayIcon"
            title="No deployments"
            message="No deployments yet"
          />

          <template v-else>
            <!-- Environment Status Cards -->
            <div class="grid grid-cols-3 gap-3 mb-6">
              <div
                v-for="{ env, latest } in envStatuses"
                :key="env"
                class="bg-gray-800 border border-gray-700 rounded-lg p-4"
              >
                <div class="flex justify-between items-center mb-2">
                  <span class="text-xs text-gray-400 font-medium uppercase tracking-wide">
                    {{ envLabel(env) }}
                  </span>
                  <StatusBadge v-if="latest" :status="latest.status" />
                  <span v-else class="px-2 py-0.5 rounded text-xs font-medium bg-gray-700 text-gray-400">
                    none
                  </span>
                </div>
                <p :class="['text-sm font-mono', latest ? 'text-white' : 'text-gray-500']">
                  {{ latest ? latest.commit.slice(0, 7) : '—' }}
                </p>
                <p class="text-xs text-gray-500 mt-1">
                  {{ latest ? `Deployed ${formatTime(latest.deployedAt)}` : 'Not deployed' }}
                </p>
              </div>
            </div>

            <!-- Deployment History -->
            <DataTable
              :columns="deploymentColumns"
              :rows="deployments"
              row-key="id"
            >
              <template #cell-environment="{ row }">
                <span :class="['px-2 py-0.5 rounded text-xs font-medium', envBadgeClass(row.environment)]">
                  {{ row.environment }}
                </span>
              </template>
              <template #cell-status="{ row }">
                <StatusBadge :status="row.status" />
              </template>
              <template #cell-commit="{ row }">
                <span class="text-gray-400 text-sm font-mono">{{ row.commit.slice(0, 7) }}</span>
              </template>
              <template #cell-triggeredBy="{ row }">
                <span class="text-gray-400 text-sm">{{ row.triggeredBy }}</span>
              </template>
              <template #cell-duration="{ row }">
                <span class="text-gray-400 text-sm">{{ formatDuration(row.duration) }}</span>
              </template>
              <template #cell-deployedAt="{ row }">
                <span class="text-gray-500 text-sm">{{ formatTime(row.deployedAt) }}</span>
              </template>
            </DataTable>
          </template>
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
  ArrowUpTrayIcon,
} from '@heroicons/vue/24/outline';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import RefreshButton from '../components/ui/RefreshButton.vue';
import DetailHeader from '../components/ui/DetailHeader.vue';
import DefinitionList from '../components/ui/DefinitionList.vue';
import DataTable from '../components/ui/DataTable.vue';
import Pagination from '../components/ui/Pagination.vue';
import { useFormatters } from '../composables/useFormatters';
import StatusBadge from '../components/ui/StatusBadge.vue';
import type { App } from '../types/App';
import type { Build } from '../types/Build';
import type { Deployment } from '../types/Deployment';

const route = useRoute();
const { formatDuration, formatRelativeTime: formatTime, formatUrl } = useFormatters();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'builds', label: 'CI/CD' },
  { id: 'deployments', label: 'Deployments' },
];

const buildColumns = [
  { key: 'number', label: 'Build' },
  { key: 'status', label: 'Status' },
  { key: 'branch', label: 'Branch' },
  { key: 'commit', label: 'Commit' },
  { key: 'duration', label: 'Duration' },
  { key: 'startedAt', label: 'Started' },
];

const deploymentColumns = [
  { key: 'environment', label: 'Environment' },
  { key: 'status', label: 'Status' },
  { key: 'commit', label: 'Commit' },
  { key: 'triggeredBy', label: 'Triggered by' },
  { key: 'duration', label: 'Duration' },
  { key: 'deployedAt', label: 'Deployed' },
];

const app = ref<App | null>(null);
const builds = ref<Build[]>([]);
const deployments = ref<Deployment[]>([]);
const loading = ref(true);
const loadingBuilds = ref(false);
const loadingDeployments = ref(false);
const activeTab = ref('overview');

const currentPage = ref(1);
const pageSize = 5;

const paginatedBuilds = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  return builds.value.slice(start, start + pageSize);
});

const envStatuses = computed(() => {
  const envs = ['production', 'staging', 'dev'] as const;
  return envs.map(env => {
    const latest = deployments.value
      .filter(d => d.environment === env)
      .sort((a, b) => new Date(b.deployedAt).getTime() - new Date(a.deployedAt).getTime())[0] ?? null;
    return { env, latest };
  });
});

const envLabel = (env: string) => {
  const labels: Record<string, string> = { production: 'Production', staging: 'Staging', dev: 'Dev' };
  return labels[env] ?? env;
};

const envBadgeClass = (env: string) => {
  const classes: Record<string, string> = {
    production: 'bg-purple-900/30 text-purple-400',
    staging: 'bg-orange-900/30 text-orange-400',
    dev: 'bg-indigo-900/30 text-indigo-400',
  };
  return classes[env] ?? 'bg-gray-700 text-gray-400';
};

const detailItems = computed(() => [
  { label: 'Description', value: app.value?.description, fallback: 'No description' },
  { label: 'Project', value: `Project #${app.value?.projectId}`, to: `/projects/${app.value?.projectId}` },
  { label: 'Repository URL', value: app.value?.repoUrl, to: formatUrl(app.value?.repoUrl || ''), external: true },
  { label: 'CI URL', value: app.value?.ciUrl, to: app.value?.ciUrl ? formatUrl(app.value.ciUrl) : undefined, external: true },
  { label: 'Template', value: `Template #${app.value?.templateId}`, to: `/templates/${app.value?.templateId}` },
]);

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
    currentPage.value = 1;
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

const fetchDeployments = async () => {
  if (!app.value) return;
  try {
    loadingDeployments.value = true;
    const response = await fetch(`/api/apps/${app.value.id}/deployments`);
    if (!response.ok) throw new Error('Failed to fetch deployments');
    deployments.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch deployments:', error);
    deployments.value = [];
  } finally {
    loadingDeployments.value = false;
  }
};

onMounted(async () => {
  await fetchApp();
  if (app.value) {
    fetchBuilds();
    fetchDeployments();
  }
});
</script>
