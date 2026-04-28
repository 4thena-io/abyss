<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <!-- Not Found -->
    <EmptyState
      v-else-if="!team"
      :icon="UsersIcon"
      title="Team not found"
    >
      <router-link to="/teams" class="text-blue-400 hover:text-blue-300">
        Back to teams
      </router-link>
    </EmptyState>

    <template v-else>
      <!-- Header -->
      <DetailHeader
        :icon="UsersIcon"
        :title="team.name"
        :subtitle="team.description"
      />

      <!-- Tabs -->
      <Tabs v-model="activeTab" :tabs="tabs" />

      <!-- Tab Content -->
      <div class="mt-6">
        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'">
          <DefinitionList
            title="Details"
            :items="[
              { label: 'Name', value: team.name },
              { label: 'Members', value: String(team.memberCount) },
              { label: 'Projects', value: String(team.projectCount) },
              { label: 'Applications', value: String(team.appCount) },
            ]"
          />
        </div>

        <!-- Projects Tab -->
        <div v-if="activeTab === 'projects'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-medium text-white">Projects</h2>
            <RefreshButton :loading="loadingProjects" @click="fetchProjects" />
          </div>

          <div v-if="loadingProjects" class="flex items-center justify-center py-8">
            <Spinner size="md" />
          </div>

          <EmptyState
            v-else-if="projects.length === 0"
            :icon="FolderIcon"
            title="No projects"
            message="This team has no projects yet"
          />

          <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            <EntityCard
              v-for="project in projects"
              :key="project.id"
              :icon="FolderIcon"
              :title="project.name"
              :description="project.description"
              :to="`/projects/${project.id}`"
            />
          </div>
        </div>

        <!-- Members Tab -->
        <div v-if="activeTab === 'members'" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-medium text-white">Members</h2>
            <RefreshButton :loading="loadingMembers" @click="fetchMembers" />
          </div>

          <div v-if="loadingMembers" class="flex items-center justify-center py-8">
            <Spinner size="md" />
          </div>

          <EmptyState
            v-else-if="members.length === 0"
            :icon="UsersIcon"
            title="No members"
            message="This team has no members yet"
          />

          <DataTable
            v-else
            :columns="memberColumns"
            :rows="members"
            row-key="username"
          >
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
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { UsersIcon, FolderIcon } from '@heroicons/vue/24/outline';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import RefreshButton from '../components/ui/RefreshButton.vue';
import DetailHeader from '../components/ui/DetailHeader.vue';
import DefinitionList from '../components/ui/DefinitionList.vue';
import DataTable from '../components/ui/DataTable.vue';
import EntityCard from '../components/ui/EntityCard.vue';
import { useFormatters } from '../composables/useFormatters';
import type { Team, TeamMember } from '../types/Team';
import type { Project } from '../types/Project';

const route = useRoute();
const { formatRelativeTime: formatTime } = useFormatters();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'projects', label: 'Projects' },
  { id: 'members', label: 'Members' },
];

const memberColumns = [
  { key: 'username', label: 'Member' },
  { key: 'role', label: 'Role' },
  { key: 'joinedAt', label: 'Joined' },
];

const team = ref<Team | null>(null);
const projects = ref<Project[]>([]);
const members = ref<TeamMember[]>([]);
const loading = ref(true);
const loadingProjects = ref(false);
const loadingMembers = ref(false);
const activeTab = ref('overview');

const projectsFetched = ref(false);
const membersFetched = ref(false);

const fetchTeam = async () => {
  try {
    loading.value = true;
    const id = Number(route.params.id);
    const response = await fetch(`/api/teams/${id}`);
    if (!response.ok) throw new Error('Team not found');
    team.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch team:', error);
    team.value = null;
  } finally {
    loading.value = false;
  }
};

const fetchProjects = async () => {
  if (!team.value) return;
  try {
    loadingProjects.value = true;
    const response = await fetch(`/api/teams/${team.value.id}/projects`);
    if (!response.ok) throw new Error('Failed to fetch projects');
    projects.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch projects:', error);
    projects.value = [];
  } finally {
    loadingProjects.value = false;
  }
};

const fetchMembers = async () => {
  if (!team.value) return;
  try {
    loadingMembers.value = true;
    const response = await fetch(`/api/teams/${team.value.id}/members`);
    if (!response.ok) throw new Error('Failed to fetch members');
    members.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch members:', error);
    members.value = [];
  } finally {
    loadingMembers.value = false;
  }
};

watch(activeTab, (tab) => {
  if (tab === 'projects' && !projectsFetched.value) {
    projectsFetched.value = true;
    fetchProjects();
  }
  if (tab === 'members' && !membersFetched.value) {
    membersFetched.value = true;
    fetchMembers();
  }
});

onMounted(fetchTeam);
</script>
