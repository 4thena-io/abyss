<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium text-t1">Builds</h2>
      <RefreshButton :loading="loading" @click="fetchBuilds" />
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <Spinner size="md" />
    </div>

    <EmptyState
      v-else-if="builds.length === 0"
      :icon="WrenchScrewdriverIcon"
      title="No builds"
      message="No builds yet"
    />

    <DataTable
      v-else
      :columns="columns"
      :rows="paginatedBuilds"
      row-key="id"
    >
      <template #cell-number="{ row }">
        <a :href="row.link" target="_blank" rel="noopener noreferrer"
          class="text-acc hover:opacity-80 font-medium transition-opacity">
          #{{ row.number }}
        </a>
      </template>
      <template #cell-status="{ row }">
        <StatusBadge :status="row.status" />
      </template>
      <template #cell-branch="{ row }">
        <span class="text-t2 text-sm">{{ row.branch }}</span>
      </template>
      <template #cell-commit="{ row }">
        <span class="text-t2 text-sm font-mono">{{ row.commit?.slice(0, 7) }}</span>
      </template>
      <template #cell-duration="{ row }">
        <span class="text-t2 text-sm">{{ formatDuration(row.duration) }}</span>
      </template>
      <template #cell-startedAt="{ row }">
        <span class="text-t3 text-sm">{{ formatTime(row.startedAt) }}</span>
      </template>
      <template #footer>
        <Pagination v-model="currentPage" :total="builds.length" :page-size="pageSize" />
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { WrenchScrewdriverIcon } from '@heroicons/vue/24/outline';
import type { Build } from '../types';
import { useFormatters } from '../../../composables/useFormatters';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';
import RefreshButton from '../../../components/ui/RefreshButton.vue';
import DataTable from '../../../components/ui/DataTable.vue';
import Pagination from '../../../components/ui/Pagination.vue';
import StatusBadge from '../../../components/ui/StatusBadge.vue';

const props = defineProps<{ appId: number }>();
const { formatDuration, formatRelativeTime: formatTime } = useFormatters();

const columns = [
  { key: 'number', label: 'Build' },
  { key: 'status', label: 'Status' },
  { key: 'branch', label: 'Branch' },
  { key: 'commit', label: 'Commit' },
  { key: 'duration', label: 'Duration' },
  { key: 'startedAt', label: 'Started' },
];

const builds = ref<Build[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = 5;

const paginatedBuilds = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  return builds.value.slice(start, start + pageSize);
});

const fetchBuilds = async () => {
  try {
    loading.value = true;
    currentPage.value = 1;
    const response = await fetch(`/api/apps/${props.appId}/builds`);
    if (!response.ok) throw new Error('Failed to fetch builds');
    builds.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch builds:', error);
    builds.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(fetchBuilds);
</script>
