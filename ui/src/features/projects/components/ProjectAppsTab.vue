<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium text-t1">Applications</h2>
      <RefreshButton :loading="loading" @click="fetchApps" />
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <Spinner size="md" />
    </div>

    <EmptyState
      v-else-if="apps.length === 0"
      :icon="CubeIcon"
      title="No applications"
      message="No applications in this project"
    />

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
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { CubeIcon } from '@heroicons/vue/24/outline';
import type { App } from '../../apps/types';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';
import RefreshButton from '../../../components/ui/RefreshButton.vue';
import EntityCard from '../../../components/ui/EntityCard.vue';
import TagGroup from '../../../components/ui/TagGroup.vue';

const props = defineProps<{ projectId: number }>();

const apps = ref<App[]>([]);
const loading = ref(false);

const fetchApps = async () => {
  try {
    loading.value = true;
    const response = await fetch(`/api/projects/${props.projectId}/apps`);
    if (!response.ok) throw new Error('Failed to fetch apps');
    apps.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch apps:', error);
    apps.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(fetchApps);
</script>
