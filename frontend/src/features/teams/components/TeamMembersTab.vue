<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium text-t1">Members</h2>
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

    <div v-else class="bg-panel border border-b1 rounded-xl overflow-hidden">
      <table class="w-full border-collapse">
        <thead>
          <tr>
            <th class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1">Member</th>
            <th class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1">Role</th>
            <th class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1">Joined</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="member in members"
            :key="member.id"
            class="cursor-pointer hover:bg-surface transition-colors"
            @click="router.push(`/users/${member.userId}`)"
          >
            <td class="px-4 py-3 border-b border-b0">
              <div class="flex items-center gap-3">
                <img v-if="member.avatarUrl" :src="member.avatarUrl" :alt="member.username"
                  class="w-8 h-8 rounded-full object-cover shrink-0" />
                <div v-else
                  class="w-8 h-8 rounded-full bg-acc-s border border-acc-b flex items-center justify-center shrink-0">
                  <span class="text-xs font-semibold text-acc font-mono">{{ member.username.slice(0, 2) }}</span>
                </div>
                <div>
                  <p class="text-sm font-medium text-t1">{{ member.username }}</p>
                  <p v-if="member.email" class="text-xs text-t3 font-mono">{{ member.email }}</p>
                </div>
              </div>
            </td>
            <td class="px-4 py-3 border-b border-b0">
              <span class="px-2 py-0.5 rounded-full text-xs font-medium"
                :class="member.role === 'owner' ? 'bg-ok-s text-ok' : 'bg-surface border border-b1 text-t2'">
                {{ member.role }}
              </span>
            </td>
            <td class="px-4 py-3 border-b border-b0 text-sm text-t3">
              {{ formatTime(member.joinedAt) }}
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
import { UsersIcon } from '@heroicons/vue/24/outline';
import type { TeamMember } from '../types';
import { useFormatters } from '../../../composables/useFormatters';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';
import RefreshButton from '../../../components/ui/RefreshButton.vue';

const props = defineProps<{ teamId: number }>();
const router = useRouter();
const { formatRelativeTime: formatTime } = useFormatters();

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
