<template>
  <div>
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <EmptyState
      v-else-if="apps.length === 0"
      :icon="CubeIcon"
      title="No applications"
      message="No apps have been created from this template yet"
    />

    <div v-else class="bg-panel border border-b1 rounded-xl overflow-hidden">
      <table class="w-full border-collapse">
        <thead>
          <tr>
            <th
              class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1"
            >
              Application
            </th>
            <th
              class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1"
            >
              Kind
            </th>
            <th
              class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1"
            >
              Language
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="app in apps"
            :key="app.id"
            class="cursor-pointer hover:bg-surface transition-colors"
            @click="router.push(`/apps/${app.id}`)"
          >
            <td class="px-4 py-3 border-b border-b0 text-sm font-medium text-t1">{{ app.name }}</td>
            <td class="px-4 py-3 border-b border-b0">
              <span
                v-if="app.kind"
                class="px-2 py-0.5 rounded-full text-xs font-medium bg-surface border border-b1 text-t2"
              >
                {{ app.kind }}
              </span>
            </td>
            <td class="px-4 py-3 border-b border-b0 text-sm text-t3 font-mono">
              {{ app.language }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { CubeIcon } from '@heroicons/vue/24/outline';
import type { App } from '../../apps/types';
import { templatesApi } from '../../../api/template';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';

const props = defineProps<{ templateId: number }>();
const router = useRouter();
const apps = ref<App[]>([]);
const loading = ref(false);

onMounted(async () => {
  loading.value = true;
  try {
    apps.value = await templatesApi.getApps(props.templateId);
  } finally {
    loading.value = false;
  }
});
</script>
