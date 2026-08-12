<template>
  <div v-if="loading" class="flex items-center justify-center py-16">
    <Spinner size="lg" />
  </div>

  <div v-else-if="error" class="flex items-center justify-center py-16">
    <p class="text-t3 text-sm">{{ error }}</p>
  </div>

  <div v-else-if="profile" class="space-y-4">
    <!-- Page header -->
    <div class="flex items-start justify-between mb-2">
      <div>
        <p class="text-xs font-semibold uppercase tracking-wide text-t3 mb-1">Account</p>
        <h1 class="text-2xl font-bold text-t1 tracking-tight">Profile</h1>
      </div>
      <a
        v-if="user?.forge_host"
        :href="`${user.forge_host}/${profile.username}`"
        target="_blank"
        rel="noopener"
        class="flex items-center gap-1.5 px-3 py-1.5 text-sm text-t2 border border-b2 rounded-lg hover:bg-surface hover:text-t1 transition-colors"
      >
        <ArrowTopRightOnSquareIcon class="w-3.5 h-3.5" />
        View on {{ capitalize(user.forge_type) }}
      </a>
    </div>

    <!-- Identity hero -->
    <div class="bg-panel border border-b1 rounded-xl p-7">
      <div class="flex items-center gap-5 mb-5">
        <!-- Avatar -->
        <div class="shrink-0">
          <img
            v-if="profile.avatarUrl"
            :src="profile.avatarUrl"
            :alt="profile.username"
            class="w-18 h-18 rounded-full object-cover ring-2 ring-b1"
            style="width: 72px; height: 72px"
          />
          <div
            v-else
            class="w-18 h-18 rounded-full bg-acc-s border border-acc-b flex items-center justify-center"
            style="width: 72px; height: 72px"
          >
            <span class="text-2xl font-semibold text-acc font-mono">{{
              profile.username.slice(0, 2)
            }}</span>
          </div>
        </div>

        <!-- Identity -->
        <div class="flex-1 min-w-0">
          <p class="text-xl font-bold text-t1 tracking-tight truncate">{{ profile.username }}</p>
          <p class="text-sm text-t3 font-mono mt-0.5 truncate">
            @{{ profile.username }}
            <span v-if="profile.email"> · {{ profile.email }}</span>
          </p>
          <div class="flex items-center gap-2 mt-2.5 flex-wrap">
            <span
              v-if="profile.isAdmin"
              class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium bg-ok-s text-ok"
            >
              admin
            </span>
            <span
              v-if="user?.forge_type"
              class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium bg-surface border border-b1 text-t2"
            >
              <ServerIcon class="w-3 h-3" />
              via {{ capitalize(user.forge_type) }} OAuth
            </span>
          </div>
        </div>

        <!-- Stats block -->
        <div class="flex border border-b1 rounded-lg overflow-hidden shrink-0">
          <div class="px-5 py-3 text-center border-r border-b1">
            <p class="text-xl font-bold text-t1">{{ profile.teamCount }}</p>
            <p class="text-xs font-semibold uppercase tracking-wide text-t3 mt-0.5">Teams</p>
          </div>
          <div class="px-5 py-3 text-center border-r border-b1">
            <p class="text-xl font-bold text-t1">{{ profile.projectCount }}</p>
            <p class="text-xs font-semibold uppercase tracking-wide text-t3 mt-0.5">Projects</p>
          </div>
          <div class="px-5 py-3 text-center">
            <p class="text-xl font-bold text-t1">{{ profile.appCount }}</p>
            <p class="text-xs font-semibold uppercase tracking-wide text-t3 mt-0.5">Apps</p>
          </div>
        </div>
      </div>

      <!-- Notice -->
      <div class="flex items-center gap-2.5 px-3.5 py-2.5 bg-surface border border-b1 rounded-lg">
        <InformationCircleIcon class="w-4 h-4 text-t3 shrink-0" />
        <p class="text-xs text-t3">
          Your identity is managed by
          <strong class="text-t2">{{ capitalize(user?.forge_type ?? 'your forge') }}</strong>
          via OAuth2. To change your name, email or avatar, update them directly there and sign back
          in.
        </p>
      </div>
    </div>

    <!-- Two-column layout -->
    <div class="grid grid-cols-2 gap-4">
      <!-- Left: Teams -->
      <div class="bg-panel border border-b1 rounded-xl overflow-hidden">
        <div class="flex items-center justify-between px-4 py-3 border-b border-b1 bg-surface">
          <span class="text-sm font-semibold text-t1">Teams</span>
        </div>
        <div v-if="profile.teams.length === 0" class="px-4 py-6 text-center text-sm text-t3">
          Not a member of any team
        </div>
        <div v-else class="p-2 space-y-0.5">
          <RouterLink
            v-for="team in profile.teams"
            :key="team.id"
            :to="`/teams/${team.id}`"
            class="flex items-center gap-3 p-2.5 rounded-lg hover:bg-surface transition-colors cursor-pointer"
          >
            <div
              class="w-8 h-8 rounded-lg bg-acc-s border border-acc-b flex items-center justify-center shrink-0"
            >
              <UsersIcon class="w-4 h-4 text-acc" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-t1 truncate">{{ team.name }}</p>
              <p class="text-xs text-t3 mt-0.5">
                {{ team.memberCount }} members · {{ team.projectCount }} projects
              </p>
            </div>
            <span
              class="px-2 py-0.5 rounded-full text-xs font-medium"
              :class="
                team.role === 'owner' ? 'bg-ok-s text-ok' : 'bg-surface border border-b1 text-t2'
              "
            >
              {{ team.role }}
            </span>
          </RouterLink>
        </div>
      </div>

      <!-- Right: Forge identity -->
      <div class="bg-panel border border-b1 rounded-xl overflow-hidden">
        <div class="flex items-center justify-between px-4 py-3 border-b border-b1 bg-surface">
          <span class="text-sm font-semibold text-t1">Forge identity</span>
        </div>
        <dl>
          <div class="flex px-4 py-2.5 border-b border-b0">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Provider</dt>
            <dd class="flex items-center gap-2 text-xs font-mono text-t1">
              <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-ok-s text-ok">
                {{ capitalize(user?.forge_type ?? '—') }}
              </span>
              <span class="text-t3">{{ user?.forge_host }}</span>
            </dd>
          </div>
          <div class="flex px-4 py-2.5 border-b border-b0">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Username</dt>
            <dd class="text-xs font-mono text-t1">{{ profile.username }}</dd>
          </div>
          <div class="flex px-4 py-2.5 border-b border-b0">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Email</dt>
            <dd class="text-xs font-mono text-t1">{{ profile.email || '—' }}</dd>
          </div>
          <div class="flex px-4 py-2.5">
            <dt class="w-36 shrink-0 text-xs font-medium text-t3 flex items-center">Auth method</dt>
            <dd class="text-xs text-t1">OAuth2</dd>
          </div>
        </dl>

        <!-- API Token shortcut (only for own profile) -->
        <template v-if="isOwnProfile">
          <div
            class="flex items-center justify-between px-4 py-3 border-t border-b1 bg-surface mt-0"
          >
            <span class="text-sm font-semibold text-t1">API Token</span>
            <RouterLink to="/settings?tab=token" class="text-xs text-acc hover:underline">
              Manage in settings →
            </RouterLink>
          </div>
          <div class="px-4 py-3">
            <div class="flex items-center gap-2 px-3 py-2 bg-surface rounded-lg">
              <span
                class="w-2 h-2 rounded-full shrink-0"
                :class="tokenActive ? 'bg-ok' : 'bg-t3'"
              />
              <span class="text-xs" :class="tokenActive ? 'text-ok' : 'text-t3'">
                {{ tokenActive ? 'Token active' : 'No token generated' }}
              </span>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, RouterLink } from 'vue-router';
import {
  ArrowTopRightOnSquareIcon,
  InformationCircleIcon,
  ServerIcon,
} from '@heroicons/vue/24/outline';
import { UsersIcon } from '@heroicons/vue/24/outline';
import type { UserProfile } from '../features/users/types';
import { userApi } from '../api/user';
import { useCurrentUser } from '../features/auth/composables/useCurrentUser';
import Spinner from '../components/ui/Spinner.vue';

const route = useRoute();
const { user } = useCurrentUser();

const profile = ref<UserProfile | null>(null);
const loading = ref(true);
const error = ref('');
const tokenActive = ref(false);

const isOwnProfile = computed(() => user.value?.id === profile.value?.id);

function capitalize(s: string) {
  return s ? s.charAt(0).toUpperCase() + s.slice(1) : s;
}

onMounted(async () => {
  const id = Number(route.params.id);
  try {
    profile.value = await userApi.getProfile(id);
    if (user.value?.id === id) {
      const status = await userApi.tokenStatus();
      tokenActive.value = status.has_token;
    }
  } catch {
    error.value = 'Failed to load profile';
  } finally {
    loading.value = false;
  }
});
</script>
