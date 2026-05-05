<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-2xl font-semibold text-t1">Settings</h1>
      <p class="text-t2 mt-1">Manage your preferences and account</p>
    </div>

    <Tabs v-model="tab" :tabs="tabs" />

    <!-- Appearance -->
    <div v-if="tab === 'appearance'" class="space-y-4">
      <div class="bg-panel border border-b1 rounded-lg p-5 space-y-4">
        <div>
          <p class="text-sm font-medium text-t1">Theme</p>
          <p class="text-xs text-t3 mt-0.5">Choose how Abyss looks for you</p>
        </div>
        <div class="grid grid-cols-3 gap-3">
          <button
            v-for="opt in themeOptions"
            :key="opt.value"
            @click="setTheme(opt.value)"
            class="flex flex-col items-center gap-2 p-3 rounded-lg border text-sm font-medium transition-colors"
            :class="theme === opt.value
              ? 'border-acc bg-acc-s text-t1'
              : 'border-b1 bg-bg text-t2 hover:border-b2'"
          >
            <component :is="opt.icon" class="w-5 h-5" />
            {{ opt.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Language -->
    <div v-if="tab === 'language'" class="space-y-4">
      <div class="bg-panel border border-b1 rounded-lg p-5">
        <p class="text-sm font-medium text-t1">Language</p>
        <p class="text-xs text-t3 mt-0.5 mb-4">Choose the display language</p>
        <div class="flex items-center gap-3 p-3 rounded-lg bg-bg border border-b1 text-sm text-t3">
          <GlobeAltIcon class="w-5 h-5 shrink-0" />
          Language support is coming soon.
        </div>
      </div>
    </div>

    <!-- API Token -->
    <div v-if="tab === 'token'" class="space-y-4">
      <div class="bg-panel border border-b1 rounded-lg p-5 space-y-4">
        <div>
          <p class="text-sm font-medium text-t1">Personal API Token</p>
          <p class="text-xs text-t3 mt-0.5">
            Use this token to authenticate API requests with
            <code class="font-mono bg-surface px-1 rounded">Authorization: Bearer &lt;token&gt;</code>
            or <code class="font-mono bg-surface px-1 rounded">?access_token=&lt;token&gt;</code>.
          </p>
        </div>

        <!-- Token revealed after generate -->
        <div v-if="newToken" class="space-y-2">
          <p class="text-xs text-ok font-medium">Token generated — copy it now, it won't be shown again.</p>
          <div class="flex items-center gap-2">
            <input
              :value="newToken"
              readonly
              class="flex-1 font-mono text-xs bg-bg border border-b1 rounded-lg px-3 py-2 text-t1 focus:outline-none"
            />
            <button
              @click="copyToken"
              class="px-3 py-2 rounded-lg border border-b1 text-xs text-t2 hover:text-t1 hover:border-b2 transition-colors whitespace-nowrap"
            >
              {{ copied ? 'Copied!' : 'Copy' }}
            </button>
          </div>
        </div>

        <!-- Status when no new token shown -->
        <div v-else class="flex items-center gap-2 text-sm">
          <span
            class="inline-flex items-center gap-1.5 px-2 py-1 rounded-full text-xs font-medium"
            :class="hasToken ? 'bg-ok-s text-ok' : 'bg-surface text-t3'"
          >
            <span class="w-1.5 h-1.5 rounded-full" :class="hasToken ? 'bg-ok' : 'bg-t3'" />
            {{ hasToken ? 'Token active' : 'No token' }}
          </span>
        </div>

        <div class="flex items-center gap-2 pt-1">
          <button
            @click="generate"
            :disabled="generating"
            class="px-4 py-2 bg-acc hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm font-medium rounded-lg transition-opacity flex items-center gap-2"
          >
            <span v-if="generating" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
            {{ hasToken ? 'Regenerate token' : 'Generate token' }}
          </button>
          <button
            v-if="hasToken"
            @click="revoke"
            :disabled="revoking"
            class="px-4 py-2 border border-b1 text-sm text-fail hover:bg-fail-s disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-colors flex items-center gap-2"
          >
            <span v-if="revoking" class="w-4 h-4 border-2 border-fail/30 border-t-fail rounded-full animate-spin" />
            Revoke
          </button>
        </div>

        <div v-if="tokenError" class="flex items-start gap-2 bg-fail-s border border-fail/30 rounded-lg px-3 py-2.5">
          <ExclamationTriangleIcon class="w-4 h-4 text-fail shrink-0 mt-0.5" />
          <p class="text-fail text-xs">{{ tokenError }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { SunIcon, MoonIcon, GlobeAltIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline';
import Tabs from '../components/ui/Tabs.vue';
import { useTheme } from '../composables/useTheme';
import { userApi } from '../api/user';

const route = useRoute();
const { theme, toggle } = useTheme();
const validTabs = ['appearance', 'language', 'token'];
const tab = ref(validTabs.includes(route.query.tab as string) ? (route.query.tab as string) : 'appearance');
const tabs = [
  { id: 'appearance', label: 'Appearance' },
  { id: 'language',   label: 'Language' },
  { id: 'token',      label: 'API Token' },
];

const themeOptions = [
  { value: 'light', label: 'Light', icon: SunIcon },
  { value: 'dark',  label: 'Dark',  icon: MoonIcon },
];

function setTheme(value: string) {
  if (theme.value !== value) toggle();
}

// Token state
const hasToken  = ref(false);
const newToken  = ref('');
const copied    = ref(false);
const generating = ref(false);
const revoking  = ref(false);
const tokenError = ref('');

onMounted(async () => {
  try {
    const status = await userApi.tokenStatus();
    hasToken.value = status.has_token;
  } catch {
    // non-fatal
  }
});

async function generate() {
  generating.value = true;
  tokenError.value = '';
  newToken.value = '';
  try {
    const res = await userApi.generateToken();
    newToken.value = res.token;
    hasToken.value = true;
  } catch (e) {
    tokenError.value = e instanceof Error ? e.message : 'Failed to generate token';
  } finally {
    generating.value = false;
  }
}

async function revoke() {
  revoking.value = true;
  tokenError.value = '';
  try {
    await userApi.revokeToken();
    hasToken.value = false;
    newToken.value = '';
  } catch (e) {
    tokenError.value = e instanceof Error ? e.message : 'Failed to revoke token';
  } finally {
    revoking.value = false;
  }
}

async function copyToken() {
  await navigator.clipboard.writeText(newToken.value);
  copied.value = true;
  setTimeout(() => (copied.value = false), 2000);
}
</script>
