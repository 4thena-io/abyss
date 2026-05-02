<template>
  <div class="min-h-screen bg-gray-950 flex items-center justify-center p-4">
    <div class="w-full max-w-lg">

      <!-- Header -->
      <div class="text-center mb-8">
        <img :src="logo" alt="Abyss" class="h-12 mx-auto mb-4" />
        <h1 class="text-2xl font-semibold text-white">Set up Abyss</h1>
        <p class="text-gray-400 text-sm mt-1">Configure your instance to get started</p>
      </div>

      <!-- Step indicator -->
      <div class="flex items-center justify-center gap-2 mb-8">
        <template v-for="(label, i) in steps" :key="i">
          <div class="flex items-center gap-2">
            <div
              class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-semibold transition-colors"
              :class="stepCircleClass(i)"
            >
              <CheckIcon v-if="i < step" class="w-4 h-4" />
              <span v-else>{{ i + 1 }}</span>
            </div>
            <span class="text-xs hidden sm:block" :class="i === step ? 'text-white' : 'text-gray-500'">
              {{ label }}
            </span>
          </div>
          <div v-if="i < steps.length - 1" class="flex-1 h-px max-w-12" :class="i < step ? 'bg-blue-500' : 'bg-gray-700'" />
        </template>
      </div>

      <!-- Card -->
      <div class="bg-gray-800 border border-gray-700 rounded-lg p-6 space-y-5">

        <!-- Step 0: Database -->
        <template v-if="step === 0">
          <div>
            <h2 class="text-lg font-medium text-white">Database</h2>
            <p class="text-gray-400 text-sm mt-1">Choose where Abyss stores its data.</p>
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Database type</label>
            <div class="grid grid-cols-3 gap-3">
              <button
                v-for="opt in dbOptions"
                :key="opt.value"
                type="button"
                @click="form.db_type = opt.value"
                class="flex items-center justify-center px-4 py-3 rounded-lg border text-sm font-medium transition-colors"
                :class="form.db_type === opt.value
                  ? 'border-blue-500 bg-blue-500/10 text-white'
                  : 'border-gray-700 bg-gray-900 text-gray-300 hover:border-gray-600'"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>

          <!-- SQLite -->
          <template v-if="form.db_type === 'sqlite'">
            <div>
              <label class="block text-gray-400 text-sm mb-2">Data file path</label>
              <input
                v-model="form.db_path"
                type="text"
                placeholder="./abyss.db"
                class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
              />
              <p class="text-xs text-gray-500 mt-1">Relative to the working directory where Abyss runs.</p>
            </div>
          </template>

          <!-- PostgreSQL / MySQL -->
          <template v-if="form.db_type === 'postgres' || form.db_type === 'mysql'">
            <div class="grid grid-cols-3 gap-3">
              <div class="col-span-2">
                <label class="block text-gray-400 text-sm mb-2">Host</label>
                <input
                  v-model="form.db_host"
                  type="text"
                  placeholder="localhost"
                  class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
                />
              </div>
              <div>
                <label class="block text-gray-400 text-sm mb-2">Port</label>
                <input
                  v-model="form.db_port"
                  type="text"
                  :placeholder="form.db_type === 'postgres' ? '5432' : '3306'"
                  class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
                />
              </div>
            </div>
            <div>
              <label class="block text-gray-400 text-sm mb-2">Database name</label>
              <input
                v-model="form.db_name"
                type="text"
                placeholder="abyss"
                class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
              />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-gray-400 text-sm mb-2">User</label>
                <input
                  v-model="form.db_user"
                  type="text"
                  placeholder="abyss"
                  class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
                />
              </div>
              <div>
                <label class="block text-gray-400 text-sm mb-2">Password</label>
                <input
                  v-model="form.db_password"
                  type="password"
                  placeholder="••••••••"
                  class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
                />
              </div>
            </div>
          </template>
        </template>

        <!-- Step 1: Forge -->
        <template v-if="step === 1">
          <div>
            <h2 class="text-lg font-medium text-white">Forge connection</h2>
            <p class="text-gray-400 text-sm mt-1">Connect Abyss to your self-hosted forge instance.</p>
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Forge type</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                v-for="opt in forgeOptions"
                :key="opt.value"
                type="button"
                @click="form.forge_type = opt.value"
                class="flex items-center justify-center px-4 py-3 rounded-lg border text-sm font-medium transition-colors"
                :class="form.forge_type === opt.value
                  ? 'border-blue-500 bg-blue-500/10 text-white'
                  : 'border-gray-700 bg-gray-900 text-gray-300 hover:border-gray-600'"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Forge URL</label>
            <input
              v-model="form.forge_host"
              type="text"
              placeholder="https://git.example.com"
              class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
            />
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Bot API token</label>
            <input
              v-model="form.forge_token"
              type="password"
              placeholder="Your service account token"
              class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
            />
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Organization / owner</label>
            <input
              v-model="form.forge_owner"
              type="text"
              placeholder="my-org"
              class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
            />
            <p class="text-xs text-gray-500 mt-1">The organization or user that owns repositories created by Abyss.</p>
          </div>
        </template>

        <!-- Step 2: OAuth App -->
        <template v-if="step === 2">
          <div>
            <h2 class="text-lg font-medium text-white">OAuth application</h2>
            <p class="text-gray-400 text-sm mt-1">Create an OAuth app in your forge and paste the credentials here.</p>
          </div>

          <div class="bg-gray-900 border border-gray-700 rounded-lg p-4 text-sm space-y-2">
            <p class="text-gray-300 font-medium">How to create the OAuth app:</p>
            <ol class="list-decimal list-inside text-gray-400 space-y-1">
              <li>Go to <span class="text-gray-300">{{ form.forge_host || 'your forge' }}</span> → Settings → Applications</li>
              <li>Click <span class="text-gray-300">"Manage OAuth2 Applications"</span></li>
              <li>Create a new app and set the callback URL below</li>
            </ol>
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">
              Callback URL
              <span class="text-gray-500 text-xs ml-1">(copy this into your forge app)</span>
            </label>
            <div class="flex items-center gap-2">
              <input
                :value="form.callback_url"
                readonly
                class="flex-1 bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-gray-300 text-sm font-mono"
              />
              <button
                type="button"
                @click="copyCallback"
                class="px-3 py-2.5 rounded-lg border border-gray-700 text-gray-400 hover:text-white hover:border-gray-600 transition-colors text-sm"
              >
                {{ copied ? 'Copied!' : 'Copy' }}
              </button>
            </div>
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Client ID</label>
            <input
              v-model="form.client_id"
              type="text"
              placeholder="Paste client ID from your forge"
              class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
            />
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Client secret</label>
            <input
              v-model="form.client_secret"
              type="password"
              placeholder="Paste client secret from your forge"
              class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
            />
          </div>
        </template>

        <!-- Step 3: CI -->
        <template v-if="step === 3">
          <div>
            <h2 class="text-lg font-medium text-white">CI provider</h2>
            <p class="text-gray-400 text-sm mt-1">Choose how pipelines are triggered for your apps.</p>
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Provider</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                v-for="opt in ciOptions"
                :key="opt.value"
                type="button"
                @click="form.ci_type = opt.value"
                class="flex items-center justify-center px-4 py-3 rounded-lg border text-sm font-medium transition-colors"
                :class="form.ci_type === opt.value
                  ? 'border-blue-500 bg-blue-500/10 text-white'
                  : 'border-gray-700 bg-gray-900 text-gray-300 hover:border-gray-600'"
              >
                {{ opt.label }}
              </button>
            </div>
          </div>

          <template v-if="ciNeedsCredentials">
            <div>
              <label class="block text-gray-400 text-sm mb-2">{{ ciProviderLabel }} URL</label>
              <input
                v-model="form.ci_host"
                type="text"
                :placeholder="`https://ci.example.com`"
                class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
              />
            </div>

            <div>
              <label class="block text-gray-400 text-sm mb-2">API token</label>
              <input
                v-model="form.ci_token"
                type="password"
                placeholder="Your CI API token"
                class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-gray-600"
              />
            </div>
          </template>

          <div v-else class="bg-gray-900 border border-gray-700 rounded-lg p-4 text-sm text-gray-400">
            {{ ciProviderLabel }} is built into your forge — no extra credentials needed.
          </div>
        </template>

        <!-- Step 4: Restarting -->
        <template v-if="step === 4">
          <div class="text-center py-6">
            <template v-if="!restartReady">
              <div class="w-12 h-12 border-2 border-blue-400/30 border-t-blue-400 rounded-full animate-spin mx-auto mb-4" />
              <h2 class="text-lg font-medium text-white">Restarting server</h2>
              <p class="text-gray-400 text-sm mt-2">Applying the new configuration…</p>
            </template>
            <template v-else>
              <div class="w-14 h-14 rounded-full bg-green-500/10 border border-green-500/30 flex items-center justify-center mx-auto mb-4">
                <CheckIcon class="w-7 h-7 text-green-400" />
              </div>
              <h2 class="text-lg font-medium text-white">All set!</h2>
              <p class="text-gray-400 text-sm mt-2">Redirecting to login…</p>
            </template>
          </div>
        </template>

        <!-- Error -->
        <div v-if="error" class="flex items-start gap-3 bg-red-500/10 border border-red-500/30 rounded-lg px-4 py-3">
          <ExclamationTriangleIcon class="w-5 h-5 text-red-400 shrink-0 mt-0.5" />
          <p class="text-red-300 text-sm">{{ error }}</p>
        </div>

        <!-- Actions -->
        <div v-if="step < 4" class="flex items-center justify-between pt-2">
          <button
            v-if="step > 0"
            type="button"
            @click="step--"
            class="px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors"
          >
            Back
          </button>
          <div v-else />

          <button
            type="button"
            @click="advance"
            :disabled="!canAdvance || saving"
            class="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed
                   text-white text-sm font-medium rounded-lg transition-colors flex items-center gap-2"
          >
            <span v-if="saving" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
            {{ step === 3 ? 'Save configuration' : 'Next' }}
          </button>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue';
import { CheckIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/solid';
import logo from '../assets/logo-white.svg';
import { setupApi } from '../api/setup';

const steps = ['Database', 'Forge', 'OAuth app', 'CI'];
const step = ref(0);
const saving = ref(false);
const error = ref('');
const copied = ref(false);
const restartReady = ref(false);

const form = reactive({
  // Database
  db_type: 'sqlite',
  db_path: './abyss.db',
  db_host: '',
  db_port: '',
  db_user: '',
  db_password: '',
  db_name: '',

  // Forge
  forge_type: 'gitea',
  forge_host: '',
  forge_token: '',
  forge_owner: '',

  // OAuth
  client_id: '',
  client_secret: '',
  callback_url: `${window.location.origin}/auth/callback`,

  // CI
  ci_type: 'gitea-actions',
  ci_host: '',
  ci_token: '',
});

const dbOptions = [
  { value: 'sqlite',   label: 'SQLite' },
  { value: 'postgres', label: 'PostgreSQL' },
  { value: 'mysql',    label: 'MySQL' },
];

const forgeOptions = [
  { value: 'gitea',   label: 'Gitea' },
  { value: 'forgejo', label: 'Forgejo' },
];

const ciOptions = [
  { value: 'gitea-actions',  label: 'Gitea Actions' },
  { value: 'woodpecker',     label: 'Woodpecker CI' },
  { value: 'drone',          label: 'Drone CI' },
  { value: 'github-actions', label: 'GitHub Actions' },
  { value: 'gitlab-ci',      label: 'GitLab CI' },
];

const ciNeedsCredentials = computed(() =>
  form.ci_type === 'woodpecker' || form.ci_type === 'drone'
);

const ciProviderLabel = computed(() =>
  ciOptions.find(o => o.value === form.ci_type)?.label ?? form.ci_type
);

const canAdvance = computed(() => {
  if (step.value === 0) {
    if (form.db_type === 'postgres' || form.db_type === 'mysql') {
      return !!(form.db_host && form.db_user && form.db_name);
    }
    return true; // SQLite always valid (path has a default)
  }
  if (step.value === 1) return !!(form.forge_type && form.forge_host && form.forge_token);
  if (step.value === 2) return !!(form.client_id && form.client_secret);
  if (step.value === 3) {
    if (ciNeedsCredentials.value) return !!(form.ci_host && form.ci_token);
    return !!form.ci_type;
  }
  return false;
});

function stepCircleClass(i: number) {
  if (i < step.value) return 'bg-blue-600 text-white';
  if (i === step.value) return 'bg-blue-600 text-white ring-2 ring-blue-400/30';
  return 'bg-gray-700 text-gray-400';
}

async function copyCallback() {
  await navigator.clipboard.writeText(form.callback_url);
  copied.value = true;
  setTimeout(() => (copied.value = false), 2000);
}

async function advance() {
  error.value = '';
  if (step.value < 3) {
    step.value++;
    return;
  }
  saving.value = true;
  try {
    await setupApi.configure({ ...form });
    step.value = 4;
    await waitForRestart();
    restartReady.value = true;
    await new Promise(resolve => setTimeout(resolve, 800));
    // Full-page redirect clears any stale JS state from before the restart.
    window.location.href = '/login';
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Setup failed';
    step.value = 3;
  } finally {
    saving.value = false;
  }
}

async function waitForRestart() {
  // Give the server a moment to begin shutting down.
  await new Promise(resolve => setTimeout(resolve, 1500));

  const deadline = Date.now() + 30_000;
  while (Date.now() < deadline) {
    try {
      const res = await fetch('/api/setup/status');
      if (res.ok) {
        const data = await res.json();
        if (data.configured) return;
      }
    } catch {
      // Server is mid-restart — keep polling.
    }
    await new Promise(resolve => setTimeout(resolve, 500));
  }

  throw new Error('Server did not come back up in time. Please restart it manually.');
}
</script>
