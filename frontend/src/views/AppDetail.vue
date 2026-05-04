<template>
  <div class="space-y-6">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <EmptyState
      v-else-if="!app"
      :icon="CubeIcon"
      title="Application not found"
    >
      <router-link to="/apps" class="text-acc hover:opacity-80">
        Back to applications
      </router-link>
    </EmptyState>

    <template v-else>
      <DetailHeader :icon="CubeIcon" :title="app.name" :subtitle="app.description">
        <template #actions>
          <a :href="formatUrl(app.repoUrl)" target="_blank" rel="noopener noreferrer"
            class="flex items-center gap-2 px-4 py-2 bg-panel border border-b1
                rounded-lg text-t2 hover:text-t1 hover:border-b2 transition-colors">
            <CodeBracketIcon class="w-4 h-4" />
            <span>Source Code</span>
          </a>
        </template>
      </DetailHeader>

      <Tabs v-model="activeTab" :tabs="tabs" />

      <div class="mt-6">
        <AppOverviewTab v-if="activeTab === 'overview'" :app="app" />
        <AppBuildsTab v-else-if="activeTab === 'builds'" :app-id="app.id" />
        <AppDeploymentsTab v-else-if="activeTab === 'deployments'" :app-id="app.id" />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { CubeIcon, CodeBracketIcon } from '@heroicons/vue/24/outline';
import type { App } from '../features/apps/types';
import { useFormatters } from '../composables/useFormatters';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import DetailHeader from '../components/ui/DetailHeader.vue';
import AppOverviewTab from '../features/apps/components/AppOverviewTab.vue';
import AppBuildsTab from '../features/apps/components/AppBuildsTab.vue';
import AppDeploymentsTab from '../features/apps/components/AppDeploymentsTab.vue';

const route = useRoute();
const { formatUrl } = useFormatters();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'builds', label: 'Builds' },
  { id: 'deployments', label: 'Deployments' },
];

const app = ref<App | null>(null);
const loading = ref(true);
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

onMounted(fetchApp);
</script>
