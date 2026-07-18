<template>
  <div class="space-y-4">
    <!-- Hero -->
    <div
      class="bg-panel border border-b1 rounded-xl p-6 grid grid-cols-[1fr_auto] gap-5 items-start"
    >
      <div>
        <p class="text-sm text-t2 leading-relaxed mb-4">
          {{ app.description || 'No description.' }}
        </p>
        <div class="flex flex-wrap gap-2">
          <RouterLink
            v-if="app.projectId"
            :to="`/projects/${app.projectId}`"
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-surface border border-b1 text-xs text-t2 hover:text-acc transition-colors"
          >
            <FolderIcon class="w-3 h-3 text-t3" />
            {{ app.projectName || `Project #${app.projectId}` }}
          </RouterLink>
          <span
            v-if="app.kind || app.language"
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-surface border border-b1 text-xs text-t2"
          >
            <CodeBracketIcon class="w-3 h-3 text-t3" />
            {{ [app.language, app.kind].filter(Boolean).join(' · ') }}
          </span>
          <RouterLink
            v-if="app.templateId"
            :to="`/templates/${app.templateId}`"
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-surface border border-b1 text-xs text-t2 hover:text-acc transition-colors"
          >
            <DocumentDuplicateIcon class="w-3 h-3 text-t3" />
            {{ app.templateName || `Template #${app.templateId}` }}
          </RouterLink>
          <span
            v-if="latestBuild"
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-surface border border-b1 text-xs text-t2"
          >
            <ClockIcon class="w-3 h-3 text-t3" />
            Build #{{ latestBuild.number }} · {{ formatRelativeTime(latestBuild.startedAt) }}
          </span>
          <span
            v-if="app.creatorUsername"
            class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-surface border border-b1 text-xs text-t2"
          >
            <UserIcon class="w-3 h-3 text-t3" />
            {{ app.creatorUsername }}
          </span>
        </div>
      </div>
      <!-- Build rate widget -->
      <div
        v-if="builds.length"
        class="bg-ok-s border border-ok/20 rounded-xl px-5 py-4 text-center min-w-[120px]"
      >
        <div class="text-3xl font-bold text-ok tracking-tight">{{ buildRate }}%</div>
        <div class="text-[10px] font-semibold uppercase tracking-wider text-ok opacity-70 mt-1">
          Build rate
        </div>
        <div class="text-[11px] text-ok opacity-60 mt-1.5">{{ builds.length }} total</div>
      </div>
      <div v-else-if="buildsLoading" class="min-w-[120px] flex items-center justify-center py-6">
        <Spinner size="sm" />
      </div>
    </div>

    <!-- Environment cards -->
    <div v-if="deployments.length || deploymentsLoading" class="grid grid-cols-3 gap-3">
      <template v-if="deploymentsLoading">
        <div
          v-for="i in 3"
          :key="i"
          class="bg-panel border border-b1 rounded-xl p-4 animate-pulse h-24"
        />
      </template>
      <template v-else>
        <div
          v-for="env in envCards"
          :key="env.name"
          class="bg-panel border border-b1 rounded-xl p-4"
        >
          <div class="text-[10px] font-semibold uppercase tracking-wider text-t3 mb-2">
            {{ env.name }}
          </div>
          <div class="flex items-center justify-between mb-2">
            <span class="font-mono text-xs text-t1">{{ env.commit || '—' }}</span>
            <span
              :class="statusBadge(env.status).class"
              class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium"
            >
              <span :class="statusBadge(env.status).dot" class="w-1.5 h-1.5 rounded-full" />
              {{ env.status }}
            </span>
          </div>
          <div class="text-[11px] text-t3">{{ env.time }}</div>
        </div>
      </template>
    </div>

    <!-- Recent activity -->
    <div class="bg-panel border border-b1 rounded-xl overflow-hidden">
      <div class="px-4 py-3 border-b border-b1 bg-surface">
        <span class="text-sm font-semibold text-t1">Recent activity</span>
      </div>
      <div v-if="buildsLoading" class="flex items-center justify-center py-8">
        <Spinner size="sm" />
      </div>
      <div v-else-if="activity.length === 0" class="px-4 py-6 text-sm text-t3 text-center">
        No activity yet
      </div>
      <div v-else>
        <div
          v-for="(item, i) in activity"
          :key="i"
          class="flex items-start gap-3 px-4 py-3 border-b border-b0 last:border-0"
        >
          <span :class="item.dotColor" class="w-2 h-2 rounded-full mt-1.5 shrink-0" />
          <div class="flex-1 min-w-0">
            <p class="text-sm text-t1" v-html="item.text" />
            <p class="text-[11px] text-t3 mt-0.5">{{ item.time }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import {
  FolderIcon,
  CodeBracketIcon,
  DocumentDuplicateIcon,
  ClockIcon,
  UserIcon,
} from '@heroicons/vue/24/outline';
import type { App, Build, Deployment } from '../types';
import { appsApi } from '../../../api/app';
import { useFormatters } from '../../../composables/useFormatters';
import Spinner from '../../../components/ui/Spinner.vue';

const props = defineProps<{ app: App }>();
const { formatRelativeTime, formatDuration } = useFormatters();

const builds = ref<Build[]>([]);
const deployments = ref<Deployment[]>([]);
const buildsLoading = ref(false);
const deploymentsLoading = ref(false);

const latestBuild = computed(() => builds.value[0] ?? null);

const buildRate = computed(() => {
  if (!builds.value.length) return 0;
  const ok = builds.value.filter((b) => b.status === 'success').length;
  return Math.round((ok / builds.value.length) * 100);
});

const ENV_ORDER = ['production', 'staging', 'dev'];

const envCards = computed(() => {
  const byEnv = new Map<string, Deployment>();
  for (const d of deployments.value) {
    const key = d.environment;
    if (!byEnv.has(key) || new Date(d.deployedAt) > new Date(byEnv.get(key)!.deployedAt)) {
      byEnv.set(key, d);
    }
  }
  const envs = ENV_ORDER.filter((e) => byEnv.has(e));
  if (!envs.length) return [];
  return envs.map((name) => {
    const d = byEnv.get(name)!;
    return {
      name,
      commit: d.commit?.slice(0, 7) ?? '—',
      status: d.status,
      time: `Deployed ${formatRelativeTime(d.deployedAt)} · ${formatDuration(d.duration)}`,
    };
  });
});

function statusBadge(status: string) {
  if (status === 'success') return { class: 'bg-ok-s text-ok', dot: 'bg-ok' };
  if (status === 'failure') return { class: 'bg-fail-s text-fail', dot: 'bg-fail' };
  if (status === 'running') return { class: 'bg-run-s text-run', dot: 'bg-run animate-pulse' };
  return { class: 'bg-surface text-t2', dot: 'bg-t3' };
}

type ActivityItem = { dotColor: string; text: string; time: string };

const activity = computed((): ActivityItem[] => {
  const items: ActivityItem[] = [];
  for (const b of builds.value.slice(0, 5)) {
    const ok = b.status === 'success';
    const fail = b.status === 'failure';
    const run = b.status === 'running';
    items.push({
      dotColor: ok ? 'bg-ok' : fail ? 'bg-fail' : run ? 'bg-run' : 'bg-t3',
      text: `<strong>Build #${b.number}</strong> ${b.status} on <strong>${b.branch}</strong>${b.duration ? ` — ${formatDuration(b.duration)}` : ''}`,
      time: formatRelativeTime(b.startedAt),
    });
  }
  return items;
});

onMounted(async () => {
  buildsLoading.value = true;
  deploymentsLoading.value = true;
  try {
    const [b, d] = await Promise.all([
      appsApi.getBuilds(props.app.id),
      appsApi.getDeployments(props.app.id),
    ]);
    builds.value = b;
    deployments.value = d;
  } finally {
    buildsLoading.value = false;
    deploymentsLoading.value = false;
  }
});
</script>
