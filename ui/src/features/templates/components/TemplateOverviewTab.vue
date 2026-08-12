<template>
  <div class="space-y-4">
    <!-- Stats -->
    <div class="grid grid-cols-3 gap-4">
      <div class="bg-panel border border-b1 rounded-xl p-5 text-center">
        <p class="text-2xl font-bold text-t1">{{ appCount }}</p>
        <p class="text-xs font-semibold uppercase tracking-wide text-t3 mt-1">Apps created</p>
      </div>
      <div class="bg-panel border border-b1 rounded-xl p-5 text-center">
        <p class="text-2xl font-bold text-t1">{{ template.kind }}</p>
        <p class="text-xs font-semibold uppercase tracking-wide text-t3 mt-1">Kind</p>
      </div>
      <div class="bg-panel border border-b1 rounded-xl p-5 text-center">
        <p class="text-2xl font-bold text-t1">{{ template.language }}</p>
        <p class="text-xs font-semibold uppercase tracking-wide text-t3 mt-1">Language</p>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <!-- Details -->
      <div class="bg-panel border border-b1 rounded-xl overflow-hidden">
        <div class="px-4 py-3 border-b border-b1 bg-surface">
          <span class="text-sm font-semibold text-t1">Details</span>
        </div>
        <dl>
          <div class="flex px-4 py-2.5 border-b border-b0">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Repository</dt>
            <dd class="text-xs font-mono text-t1 truncate">
              <a
                v-if="template.repoUrl"
                :href="template.repoUrl"
                target="_blank"
                rel="noopener"
                class="text-acc hover:underline flex items-center gap-1"
              >
                {{ repoLabel }}
                <ArrowTopRightOnSquareIcon class="w-3 h-3 shrink-0" />
              </a>
              <span v-else class="text-t3">—</span>
            </dd>
          </div>
          <div class="flex px-4 py-2.5 border-b border-b0">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Kind</dt>
            <dd class="text-xs text-t1">{{ template.kind || '—' }}</dd>
          </div>
          <div class="flex px-4 py-2.5 border-b border-b0">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Language</dt>
            <dd class="text-xs font-mono text-t1">{{ template.language || '—' }}</dd>
          </div>
          <div class="flex px-4 py-2.5 border-b border-b0">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Registered</dt>
            <dd class="text-xs text-t3">{{ formatTime(template.createdAt) }}</dd>
          </div>
          <div class="flex px-4 py-2.5">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Created by</dt>
            <dd class="text-xs text-t1">{{ template.creatorUsername || '—' }}</dd>
          </div>
        </dl>
      </div>

      <!-- Used by -->
      <div class="bg-panel border border-b1 rounded-xl overflow-hidden">
        <div class="flex items-center justify-between px-4 py-3 border-b border-b1 bg-surface">
          <span class="text-sm font-semibold text-t1">Used by</span>
          <button
            v-if="apps.length > 3"
            @click="$emit('show-apps')"
            class="text-xs text-acc hover:underline"
          >
            View all →
          </button>
        </div>
        <div v-if="loading" class="flex items-center justify-center py-6">
          <Spinner size="sm" />
        </div>
        <div v-else-if="apps.length === 0" class="px-4 py-6 text-center text-sm text-t3">
          No applications yet
        </div>
        <div v-else class="p-2 space-y-0.5">
          <RouterLink
            v-for="app in apps.slice(0, 3)"
            :key="app.id"
            :to="`/apps/${app.id}`"
            class="flex items-center gap-3 p-2.5 rounded-lg hover:bg-surface transition-colors"
          >
            <div
              class="w-7 h-7 rounded-lg bg-surface border border-b1 flex items-center justify-center shrink-0"
            >
              <CubeIcon class="w-3.5 h-3.5 text-t3" />
            </div>
            <p class="text-sm font-medium text-t1 truncate flex-1">{{ app.name }}</p>
            <TagGroup v-if="app.kind" :tags="[app.kind]" />
          </RouterLink>
          <div v-if="apps.length > 3" class="px-3 pt-1 pb-1 border-t border-b0 text-center">
            <button
              @click="$emit('show-apps')"
              class="text-xs text-t3 hover:text-t1 transition-colors"
            >
              +{{ apps.length - 3 }} more
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { CubeIcon, ArrowTopRightOnSquareIcon } from '@heroicons/vue/24/outline';
import type { Template } from '../types';
import type { App } from '../../apps/types';
import { templatesApi } from '../../../api/template';
import { useFormatters } from '../../../composables/useFormatters';
import Spinner from '../../../components/ui/Spinner.vue';
import TagGroup from '../../../components/ui/TagGroup.vue';

const props = defineProps<{ template: Template }>();
defineEmits<{ 'show-apps': [] }>();

const { formatRelativeTime: formatTime } = useFormatters();
const apps = ref<App[]>([]);
const loading = ref(false);

const appCount = computed(() => apps.value.length);

const repoLabel = computed(() => {
  if (!props.template.repoUrl) return '';
  try {
    const url = new URL(props.template.repoUrl);
    return url.hostname + url.pathname;
  } catch {
    return props.template.repoUrl;
  }
});

onMounted(async () => {
  loading.value = true;
  try {
    apps.value = await templatesApi.getApps(props.template.id);
  } finally {
    loading.value = false;
  }
});
</script>
