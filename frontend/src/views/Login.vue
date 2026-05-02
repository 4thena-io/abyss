<template>
  <div class="min-h-screen bg-gray-950 flex items-center justify-center p-4">
    <div class="w-full max-w-sm">
      <div class="text-center mb-8">
        <img :src="logo" alt="Abyss" class="h-12 mx-auto mb-4" />
        <h1 class="text-2xl font-semibold text-white">Welcome to Abyss</h1>
        <p class="text-gray-400 mt-2 text-sm">Sign in with your forge account to continue</p>
      </div>

      <div v-if="error" class="flex items-start gap-3 bg-red-500/10 border border-red-500/30 rounded-lg px-4 py-3 mb-4">
        <ExclamationTriangleIcon class="w-5 h-5 text-red-400 shrink-0 mt-0.5" />
        <p class="text-red-300 text-sm">{{ errorMessage }}</p>
      </div>

      <div class="bg-gray-800 border border-gray-700 rounded-lg p-6">
        <a
          href="/auth/login"
          class="flex items-center justify-center gap-3 w-full bg-blue-600 hover:bg-blue-500
                 text-white font-medium py-2.5 px-4 rounded-lg transition-colors"
        >
          <component :is="forgeIcon" class="w-5 h-5" />
          Sign in with {{ forgeLabel }}
        </a>
      </div>

      <p class="text-center text-gray-600 text-xs mt-6">
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
