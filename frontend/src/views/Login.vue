<template>
  <div class="min-h-screen bg-bg flex items-center justify-center p-4">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <img :src="logo" alt="Abyss" class="h-10 mx-auto mb-4" />
        <h1 class="text-2xl font-semibold text-t1">Welcome to Abyss</h1>
        <p class="text-t3 mt-2 text-sm">Sign in with your forge account to continue</p>
      </div>

      <div v-if="error" class="flex items-start gap-3 bg-fail-s border border-fail/30 rounded-lg px-4 py-3 mb-4">
        <ExclamationTriangleIcon class="w-5 h-5 text-fail shrink-0 mt-0.5" />
        <p class="text-fail text-sm">{{ errorMessage }}</p>
      </div>

      <div class="bg-panel border border-b1 rounded-lg p-6">
        <a
          href="/auth/login"
          class="flex items-center justify-center gap-3 w-full bg-acc hover:opacity-90
                 text-white font-medium py-2.5 px-4 rounded-lg transition-opacity"
        >
          <component :is="forgeIcon" class="w-5 h-5" />
          Sign in with {{ forgeLabel }}
        </a>
      </div>

      <p class="text-center text-t4 text-xs mt-6">
        Abyss — self-hosted developer platform
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { ServerIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline';
import logo from '../assets/logo-white.svg';

const forgeLabel = computed(() => {
  const host = window.location.hostname;
  if (host.includes('github')) return 'GitHub';
  if (host.includes('forgejo')) return 'Forgejo';
  return 'Gitea';
});

const forgeIcon = ServerIcon;

const error = new URLSearchParams(window.location.search).get('error');
const errorMessage = computed(() => {
  if (error === 'access_denied') return 'You are not a member of the configured organization.';
  return 'An error occurred. Please try again.';
});
</script>
