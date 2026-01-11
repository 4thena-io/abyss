<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-white">Applications</h1>
        <p class="text-gray-400 mt-1">{{ appList.length }} applications across all projects</p>
      </div>
      <button class="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 
        text-white text-sm font-medium rounded-lg transition-colors">
        <PlusIcon class="w-4 h-4" />
        New Application
      </button>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="relative flex-1 max-w-xs">
        <MagnifyingGlassIcon class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input v-model="search" type="text" placeholder="Filter applications..." class="w-full bg-gray-800 border border-gray-700 rounded-lg pl-10 pr-4 py-2 
            text-sm text-white placeholder-gray-500 focus:outline-none focus:border-gray-600" />
      </div>

      <!-- Kind Dropdown -->
      <div class="relative">
        <button @click="showKindDropdown = !showKindDropdown" class="flex items-center gap-2 bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 
            text-sm text-gray-300 hover:border-gray-600 transition-colors min-w-32">
          <span>{{ kindFilter || 'All kinds' }}</span>
          <ChevronDownIcon class="w-4 h-4 ml-auto" />
        </button>
        <div v-if="showKindDropdown"
          class="absolute z-10 w-full mt-1 bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
          <button @click="kindFilter = ''; showKindDropdown = false" :class="['w-full px-3 py-2 text-left text-sm transition-colors',
            !kindFilter ? 'bg-gray-700 text-white' : 'text-gray-300 hover:bg-gray-700']">
            All kinds
          </button>
          <button v-for="kind in uniqueKinds" :key="kind" @click="kindFilter = kind; showKindDropdown = false" :class="['w-full px-3 py-2 text-left text-sm transition-colors',
            kindFilter === kind ? 'bg-gray-700 text-white' : 'text-gray-300 hover:bg-gray-700']">
            {{ kind }}
          </button>
        </div>
      </div>

      <!-- Language Dropdown -->
      <div class="relative">
        <button @click="showLanguageDropdown = !showLanguageDropdown" class="flex items-center gap-2 bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 
            text-sm text-gray-300 hover:border-gray-600 transition-colors min-w-32">
          <span>{{ languageFilter || 'All languages' }}</span>
          <ChevronDownIcon class="w-4 h-4 ml-auto" />
        </button>
        <div v-if="showLanguageDropdown"
          class="absolute z-10 w-full mt-1 bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
          <button @click="languageFilter = ''; showLanguageDropdown = false" :class="['w-full px-3 py-2 text-left text-sm transition-colors',
            !languageFilter ? 'bg-gray-700 text-white' : 'text-gray-300 hover:bg-gray-700']">
            All languages
          </button>
          <button v-for="lang in uniqueLanguages" :key="lang"
            @click="languageFilter = lang; showLanguageDropdown = false" :class="['w-full px-3 py-2 text-left text-sm transition-colors',
              languageFilter === lang ? 'bg-gray-700 text-white' : 'text-gray-300 hover:bg-gray-700']">
            {{ lang }}
          </button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="w-8 h-8 border-2 border-gray-600 border-t-blue-500 rounded-full animate-spin" />
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredApps.length === 0" class="bg-gray-800 border border-gray-700 rounded-lg py-16 text-center">
      <CubeIcon class="w-12 h-12 text-gray-600 mx-auto" />
      <h3 class="mt-4 text-lg font-medium text-white">No applications found</h3>
      <p class="mt-2 text-gray-400 text-sm">
        {{ search || kindFilter || languageFilter ? 'Try adjusting your filters' : 'Create your first application to get started' }}
      </p>
    </div>

    <!-- App Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <router-link v-for="app in filteredApps" :key="app.id" :to="`/apps/${app.id}`" class="bg-gray-800 border border-gray-700 rounded-lg p-5 
          hover:border-gray-600 hover:bg-gray-750 transition-all group">
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-gray-700 flex items-center justify-center">
              <CubeIcon class="w-5 h-5 text-gray-400" />
            </div>
            <div>
              <h3 class="text-white font-medium group-hover:text-blue-400 transition-colors">
                {{ app.name }}
              </h3>
              <p class="text-gray-500 text-sm">{{ app.kind }}</p>
            </div>
          </div>
        </div>

        <p v-if="app.description" class="mt-3 text-gray-400 text-sm line-clamp-2">
          {{ app.description }}
        </p>

        <div class="mt-4 flex items-center gap-2 flex-wrap">
          <span class="px-2 py-1 bg-gray-700 rounded text-xs text-gray-300">
            {{ app.kind }}
          </span>
          <span v-if="app.language" class="px-2 py-1 bg-gray-700 rounded text-xs text-gray-300">
            {{ app.language }}
          </span>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import {
  PlusIcon,
  MagnifyingGlassIcon,
  CubeIcon,
  ChevronDownIcon
} from '@heroicons/vue/24/outline';
import type { App } from '../types/App';

const appList = ref<App[]>([]);
const loading = ref(true);
const search = ref('');
const kindFilter = ref('');
const languageFilter = ref('');
const showKindDropdown = ref(false);
const showLanguageDropdown = ref(false);

const filteredApps = computed(() => {
  return appList.value.filter(app => {
    const matchesSearch = !search.value ||
      app.name.toLowerCase().includes(search.value.toLowerCase()) ||
      app.description.toLowerCase().includes(search.value.toLowerCase());
    const matchesKind = !kindFilter.value || app.kind === kindFilter.value;
    const matchesLanguage = !languageFilter.value || app.language === languageFilter.value;
    return matchesSearch && matchesKind && matchesLanguage;
  });
});

const uniqueKinds = computed(() => {
  return [...new Set(appList.value.map(a => a.kind).filter(Boolean))].sort();
});

const uniqueLanguages = computed(() => {
  return [...new Set(appList.value.map(a => a.language).filter(Boolean))].sort();
});

const fetchApps = async () => {
  try {
    loading.value = true;
    const response = await fetch('/api/apps');
    if (!response.ok) throw new Error(`HTTP error: ${response.status}`);
    appList.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch apps:', error);
  } finally {
    loading.value = false;
  }
};

const closeDropdowns = (e: Event) => {
  const target = e.target as HTMLElement;
  if (!target.closest('.relative')) {
    showKindDropdown.value = false;
    showLanguageDropdown.value = false;
  }
};

onMounted(() => {
  fetchApps();
  document.addEventListener('click', closeDropdowns);
});

onUnmounted(() => {
  document.removeEventListener('click', closeDropdowns);
});
</script>
