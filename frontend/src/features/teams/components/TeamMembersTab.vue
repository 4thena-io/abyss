<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium text-white">Members</h2>
      <RefreshButton :loading="loading" @click="fetchMembers" />
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <Spinner size="md" />
    </div>

    <EmptyState
      v-else-if="members.length === 0"
      :icon="UsersIcon"
      title="No members"
      message="This team has no members yet"
    />

    <DataTable v-else :columns="columns" :rows="members" row-key="username">
      <template #cell-username="{ row }">
        <span class="text-white text-sm">{{ row.username }}</span>
      </template>
      <template #cell-role="{ row }">
        <span class="px-2 py-0.5 bg-gray-700 rounded text-xs text-gray-300">{{ row.role }}</span>
      </template>
      <template #cell-joinedAt="{ row }">
        <span class="text-gray-500 text-sm">{{ formatTime(row.joinedAt) }}</span>
      </template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { UsersIcon } from '@heroicons/vue/24/outline';
import type { TeamMember } from '../types';
import { useFormatters } from '../../../composables/useFormatters';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';
import RefreshButton from '../../../components/ui/RefreshButton.vue';
import DataTable from '../../../components/ui/DataTable.vue';

const props = defineProps<{ teamId: number }>();
const { formatRelativeTime: formatTime } = useFormatters();

const columns = [
  { key: 'username', label: 'Member' },
  { key: 'role', label: 'Role' },
  { key: 'joinedAt', label: 'Joined' },
];

const members = ref<TeamMember[]>([]);
const loading = ref(false);

const fetchMembers = async () => {
  try {
    loading.value = true;
    const response = await fetch(`/api/teams/${props.teamId}/members`);
    if (!response.ok) throw new Error('Failed to fetch members');
    members.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch members:', error);
    members.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(fetchMembers);
</script>
