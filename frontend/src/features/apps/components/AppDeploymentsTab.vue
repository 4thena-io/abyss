<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium text-t1">Deployments</h2>
      <RefreshButton :loading="loading" @click="fetchDeployments" />
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <Spinner size="md" />
    </div>

    <EmptyState
      v-else-if="deployments.length === 0"
      :icon="ArrowUpTrayIcon"
      title="No deployments"
      message="No deployments yet"
    />

    <template v-else>
      <div class="grid grid-cols-3 gap-3 mb-6">
        <div
          v-for="{ env, latest } in envStatuses"
          :key="env"
          class="bg-panel border border-b1 rounded-lg p-4"
        >
          <div class="flex justify-between items-center mb-2">
            <span class="text-xs text-t3 font-medium uppercase tracking-wide">
              {{ envLabel(env) }}
            </span>
            <StatusBadge v-if="latest" :status="latest.status" />
            <span v-else class="px-2 py-0.5 rounded text-xs font-medium bg-surface text-t3">
              none
            </span>
          </div>
          <p :class="['text-sm font-mono', latest ? 'text-t1' : 'text-t3']">
            {{ latest ? latest.commit.slice(0, 7) : '—' }}
          </p>
          <p class="text-xs text-t3 mt-1">
            {{ latest ? `Deployed ${formatTime(latest.deployedAt)}` : 'Not deployed' }}
          </p>
        </div>
      </div>

      <DataTable :columns="columns" :rows="deployments" row-key="id">
        <template #cell-environment="{ row }">
          <span :class="['px-2 py-0.5 rounded text-xs font-medium', envBadgeClass(row.environment)]">
            {{ row.environment }}
          </span>
        </template>
        <template #cell-status="{ row }">
          <StatusBadge :status="row.status" />
        </template>
        <template #cell-commit="{ row }">
          <span class="text-t2 text-sm font-mono">{{ row.commit.slice(0, 7) }}</span>
        </template>
        <template #cell-triggeredBy="{ row }">
          <span class="text-t2 text-sm">{{ row.triggeredBy }}</span>
        </template>
        <template #cell-duration="{ row }">
          <span class="text-t2 text-sm">{{ formatDuration(row.duration) }}</span>
        </template>
        <template #cell-deployedAt="{ row }">
          <span class="text-t3 text-sm">{{ formatTime(row.deployedAt) }}</span>
        </template>
      </DataTable>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { ArrowUpTrayIcon } from '@heroicons/vue/24/outline';
import type { Deployment } from '../types';
import { useFormatters } from '../../../composables/useFormatters';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';
import RefreshButton from '../../../components/ui/RefreshButton.vue';
import DataTable from '../../../components/ui/DataTable.vue';
import StatusBadge from '../../../components/ui/StatusBadge.vue';

const props = defineProps<{ appId: number }>();
const { formatDuration, formatRelativeTime: formatTime } = useFormatters();

const columns = [
  { key: 'environment', label: 'Environment' },
  { key: 'status', label: 'Status' },
  { key: 'commit', label: 'Commit' },
  { key: 'triggeredBy', label: 'Triggered by' },
  { key: 'duration', label: 'Duration' },
  { key: 'deployedAt', label: 'Deployed' },
];

const deployments = ref<Deployment[]>([]);
const loading = ref(false);

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
  return classes[env] ?? 'bg-surface text-t3';
};

const fetchDeployments = async () => {
  try {
    loading.value = true;
    const response = await fetch(`/api/apps/${props.appId}/deployments`);
    if (!response.ok) throw new Error('Failed to fetch deployments');
    deployments.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch deployments:', error);
    deployments.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(fetchDeployments);
</script>
