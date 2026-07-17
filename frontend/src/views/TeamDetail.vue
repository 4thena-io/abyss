<template>
  <div class="space-y-6">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <EmptyState
      v-else-if="!team"
      :icon="UsersIcon"
      title="Team not found"
    >
      <router-link to="/teams" class="text-acc hover:opacity-80">
        Back to teams
      </router-link>
    </EmptyState>

    <template v-else>
      <DetailHeader :icon="UsersIcon" :title="team.name" :subtitle="team.description" />

      <Tabs v-model="activeTab" :tabs="tabs" />

      <div class="mt-6">
        <TeamOverviewTab v-if="activeTab === 'overview'" :team="team" />
        <TeamProjectsTab v-else-if="activeTab === 'projects'" :team-id="team.id" />
        <TeamMembersTab v-else-if="activeTab === 'members'" :team-id="team.id" />
        <TeamSettingsTab v-else-if="activeTab === 'settings'" :team="team" @updated="team = $event" />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { UsersIcon } from '@heroicons/vue/24/outline';
import type { Team } from '../features/teams/types';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import DetailHeader from '../components/ui/DetailHeader.vue';
import TeamOverviewTab from '../features/teams/components/TeamOverviewTab.vue';
import TeamProjectsTab from '../features/teams/components/TeamProjectsTab.vue';
import TeamMembersTab from '../features/teams/components/TeamMembersTab.vue';
import TeamSettingsTab from '../features/teams/components/TeamSettingsTab.vue';

const route = useRoute();
const router = useRouter();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'projects', label: 'Projects' },
  { id: 'members', label: 'Members' },
  { id: 'settings', label: 'Settings' },
];

const validTabs = tabs.map(t => t.id);
const team = ref<Team | null>(null);
const loading = ref(true);
const activeTab = ref(validTabs.includes(route.query.tab as string) ? route.query.tab as string : 'overview');

watch(activeTab, tab => router.replace({ query: { tab } }));

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

onMounted(fetchTeam);
</script>
