<template>
  <div class="space-y-4">
    <!-- Stat cards -->
    <div class="grid grid-cols-3 gap-3">
      <div class="bg-panel border border-b1 rounded-xl p-5">
        <div class="text-[11px] font-semibold uppercase tracking-wider text-t3 mb-2">
          Applications
        </div>
        <div class="text-3xl font-bold text-t1 tracking-tight">{{ apps.length }}</div>
        <div class="text-xs text-t3 mt-1">in this project</div>
      </div>
      <div class="bg-panel border border-b1 rounded-xl p-5">
        <div class="text-[11px] font-semibold uppercase tracking-wider text-t3 mb-2">
          Created by
        </div>
        <div class="text-lg font-semibold text-t1 mt-1 truncate">
          {{ project.creatorUsername || '—' }}
        </div>
      </div>
      <div v-if="project.teamId" class="bg-panel border border-b1 rounded-xl p-5">
        <div class="text-[11px] font-semibold uppercase tracking-wider text-t3 mb-2">Team</div>
        <RouterLink
          :to="`/teams/${project.teamId}`"
          class="text-lg font-semibold text-acc hover:opacity-80 mt-1 block truncate"
        >
          {{ project.teamName || `Team #${project.teamId}` }}
        </RouterLink>
      </div>
      <div v-else class="bg-panel border border-b1 rounded-xl p-5">
        <div class="text-[11px] font-semibold uppercase tracking-wider text-t3 mb-2">Team</div>
        <div class="text-sm text-t3 mt-1">Standalone</div>
      </div>
    </div>

    <!-- App health cards -->
    <div v-if="loading" class="grid grid-cols-2 gap-3">
      <div
        v-for="i in 2"
        :key="i"
        class="bg-panel border border-b1 rounded-xl p-4 h-28 animate-pulse"
      />
    </div>

    <div v-else-if="apps.length" class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <RouterLink
        v-for="app in apps"
        :key="app.id"
        :to="`/apps/${app.id}`"
        class="bg-panel border border-b1 rounded-xl overflow-hidden hover:border-b2 transition-colors cursor-pointer"
      >
        <div class="flex items-center justify-between px-4 py-3 border-b border-b1 bg-surface">
          <div class="flex items-center gap-2">
            <span class="text-sm font-semibold text-t1">{{ app.name }}</span>
            <span
              v-if="app.kind || app.language"
              class="text-[11px] px-1.5 py-0.5 rounded bg-surface border border-b1 text-t2 font-mono"
            >
              {{ [app.language, app.kind].filter(Boolean).join(' · ') }}
            </span>
          </div>
          <span
            class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium bg-surface text-t2 border border-b1"
          >
            <CubeIcon class="w-3 h-3" />
            app
          </span>
        </div>
        <div class="px-4 py-3">
          <p class="text-xs text-t3 truncate">{{ app.description || 'No description' }}</p>
          <div class="flex items-center gap-2 mt-2">
            <a
              v-if="app.ciUrl"
              :href="app.ciUrl"
              target="_blank"
              rel="noopener"
              class="text-xs text-acc hover:underline"
              @click.stop
              >CI pipeline →</a
            >
            <span v-if="app.repoUrl" class="text-xs text-t3 font-mono truncate">{{
              repoLabel(app.repoUrl)
            }}</span>
          </div>
        </div>
      </RouterLink>
    </div>

    <div v-else class="bg-panel border border-b1 rounded-xl px-4 py-10 text-center text-sm text-t3">
      No applications yet
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { CubeIcon } from '@heroicons/vue/24/outline';
import type { Project } from '../types';
import type { App } from '../../apps/types';
import { projectsApi } from '../../../api/project';

const props = defineProps<{ project: Project }>();

const apps = ref<App[]>([]);
const loading = ref(false);

function repoLabel(url: string) {
  try {
    return new URL(url).pathname.replace(/^\//, '');
  } catch {
    return url;
  }
}

onMounted(async () => {
  loading.value = true;
  try {
    apps.value = await projectsApi.getApps(props.project.id);
  } finally {
    loading.value = false;
  }
});
</script>
