<template>
  <div class="space-y-4">
    <!-- Stat cards -->
    <div class="grid grid-cols-3 gap-3">
      <div class="bg-panel border border-b1 rounded-xl p-5">
        <div class="text-[11px] font-semibold uppercase tracking-wider text-t3 mb-2">Members</div>
        <div class="text-3xl font-bold text-t1 tracking-tight">{{ team.memberCount }}</div>
        <div class="text-xs text-t3 mt-1">active</div>
      </div>
      <div class="bg-panel border border-b1 rounded-xl p-5">
        <div class="text-[11px] font-semibold uppercase tracking-wider text-t3 mb-2">Projects</div>
        <div class="text-3xl font-bold text-t1 tracking-tight">{{ team.projectCount }}</div>
        <div class="text-xs text-t3 mt-1">in this team</div>
      </div>
      <div class="bg-panel border border-b1 rounded-xl p-5">
        <div class="text-[11px] font-semibold uppercase tracking-wider text-t3 mb-2">
          Applications
        </div>
        <div class="text-3xl font-bold text-t1 tracking-tight">{{ team.appCount }}</div>
        <div class="text-xs text-t3 mt-1">across all projects</div>
      </div>
    </div>

    <!-- Members + Projects -->
    <div class="grid grid-cols-2 gap-4">
      <!-- Members snapshot -->
      <div class="bg-panel border border-b1 rounded-xl overflow-hidden">
        <div class="flex items-center justify-between px-4 py-3 border-b border-b1 bg-surface">
          <span class="text-sm font-semibold text-t1">Members</span>
          <RouterLink
            :to="`/teams/${team.id}?tab=members`"
            class="text-xs text-acc hover:underline"
          >
            View all →
          </RouterLink>
        </div>
        <div v-if="membersLoading" class="flex items-center justify-center py-8">
          <Spinner size="sm" />
        </div>
        <div v-else-if="members.length === 0" class="px-4 py-6 text-sm text-t3 text-center">
          No members
        </div>
        <div v-else class="p-3 space-y-1">
          <div
            v-for="m in members.slice(0, 5)"
            :key="m.id"
            class="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-surface transition-colors"
          >
            <div
              class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 font-mono text-[11px] font-semibold"
              :class="
                m.role === 'owner'
                  ? 'bg-acc-s border border-acc-b text-acc'
                  : 'bg-surface border border-b1 text-t2'
              "
            >
              {{ initials(m.username) }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-t1 truncate">{{ m.username }}</p>
              <p class="text-[11px] text-t3 font-mono">@{{ m.username }}</p>
            </div>
            <span
              class="text-[11px] font-medium px-2 py-0.5 rounded-full"
              :class="
                m.role === 'owner' ? 'bg-ok-s text-ok' : 'bg-surface text-t2 border border-b1'
              "
              >{{ m.role }}</span
            >
          </div>
          <div v-if="members.length > 5" class="px-2 pt-1 text-xs text-t3 text-center">
            +{{ members.length - 5 }} more
          </div>
        </div>
      </div>

      <!-- Projects snapshot -->
      <div class="bg-panel border border-b1 rounded-xl overflow-hidden">
        <div class="flex items-center justify-between px-4 py-3 border-b border-b1 bg-surface">
          <span class="text-sm font-semibold text-t1">Projects</span>
          <RouterLink
            :to="`/teams/${team.id}?tab=projects`"
            class="text-xs text-acc hover:underline"
          >
            View all →
          </RouterLink>
        </div>
        <div v-if="projectsLoading" class="flex items-center justify-center py-8">
          <Spinner size="sm" />
        </div>
        <div v-else-if="projects.length === 0" class="px-4 py-6 text-sm text-t3 text-center">
          No projects
        </div>
        <div v-else class="p-2 space-y-0.5">
          <RouterLink
            v-for="p in projects.slice(0, 5)"
            :key="p.id"
            :to="`/projects/${p.id}`"
            class="flex items-center gap-3 p-2.5 rounded-lg hover:bg-surface transition-colors"
          >
            <div
              class="w-7 h-7 rounded-lg bg-surface border border-b1 flex items-center justify-center shrink-0"
            >
              <FolderIcon class="w-3.5 h-3.5 text-t3" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-t1 truncate">{{ p.name }}</p>
              <p class="text-xs text-t3 truncate">{{ p.description || 'No description' }}</p>
            </div>
          </RouterLink>
          <div v-if="projects.length > 5" class="px-3 pt-1 text-xs text-t3 text-center">
            +{{ projects.length - 5 }} more
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import { FolderIcon } from '@heroicons/vue/24/outline';
import type { Team, TeamMember } from '../types';
import type { Project } from '../../projects/types';
import { teamsApi } from '../../../api/team';
import Spinner from '../../../components/ui/Spinner.vue';

const props = defineProps<{ team: Team }>();

const members = ref<TeamMember[]>([]);
const projects = ref<Project[]>([]);
const membersLoading = ref(false);
const projectsLoading = ref(false);

function initials(username: string) {
  return username.slice(0, 2).toLowerCase();
}

onMounted(async () => {
  membersLoading.value = true;
  projectsLoading.value = true;
  try {
    const [m, p] = await Promise.all([
      teamsApi.getMembers(props.team.id),
      teamsApi.getProjects(props.team.id),
    ]);
    members.value = m;
    projects.value = p;
  } finally {
    membersLoading.value = false;
    projectsLoading.value = false;
  }
});
</script>
