<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-white">Templates</h1>
        <p class="text-gray-400 mt-1">{{ templateList.length }} templates available</p>
      </div>
      <button @click="showModal = true" class="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 
          text-white text-sm font-medium rounded-lg transition-colors">
        <PlusIcon class="w-4 h-4" />
        Add Template
      </button>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="relative flex-1 max-w-xs">
        <MagnifyingGlassIcon class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input v-model="search" type="text" placeholder="Filter templates..." class="w-full bg-gray-800 border border-gray-700 rounded-lg pl-10 pr-4 py-2 
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
    <div v-else-if="filteredTemplates.length === 0"
      class="bg-gray-800 border border-gray-700 rounded-lg py-16 text-center">
      <DocumentDuplicateIcon class="w-12 h-12 text-gray-600 mx-auto" />
      <h3 class="mt-4 text-lg font-medium text-white">No templates found</h3>
      <p class="mt-2 text-gray-400 text-sm">
        {{ search || kindFilter || languageFilter ? 'Try adjusting your filters' : 'Add your first template to get started' }}
      </p>
    </div>

    <!-- Template Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div v-for="template in filteredTemplates" :key="template.id"
        class="bg-gray-800 border border-gray-700 rounded-lg p-5 hover:border-gray-600 transition-all">
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-gray-700 flex items-center justify-center">
              <DocumentDuplicateIcon class="w-5 h-5 text-gray-400" />
            </div>
            <div>
              <h3 class="text-white font-medium">{{ template.name }}</h3>
              <p class="text-gray-500 text-xs mt-0.5">{{ template.kind }}</p>
            </div>
          </div>
          <button @click="deleteTemplate(template.id)" class="text-gray-500 hover:text-red-400 transition-colors p-1">
            <TrashIcon class="w-4 h-4" />
          </button>
        </div>

        <p v-if="template.description" class="mt-4 text-gray-400 text-sm line-clamp-2">
          {{ template.description }}
        </p>
        <p v-else class="mt-4 text-gray-500 text-sm italic">No description</p>

        <div class="mt-4 flex items-center gap-2">
          <span class="px-2 py-1 bg-gray-700 rounded text-xs text-gray-300">
            {{ template.language }}
          </span>
        </div>

        <div class="mt-4 pt-4 border-t border-gray-700">
          <a :href="template.repoUrl" target="_blank"
            class="flex items-center gap-2 text-gray-400 hover:text-white text-sm transition-colors">
            <CodeBracketIcon class="w-4 h-4" />
            <span>View repo</span>
          </a>
        </div>
      </div>
    </div>

    <!-- Add Template Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/60" @click="closeModal" />
      <div class="relative bg-gray-800 border border-gray-700 rounded-lg w-full max-w-lg p-6">
        <h2 class="text-xl font-semibold text-white mb-4">Add Template</h2>

        <div class="space-y-4">
          <!-- Repo Selection -->
          <div>
            <label class="block text-gray-400 text-sm mb-2">Repository</label>
            <div class="relative">
              <input v-model="repoSearch" type="text" placeholder="Search repositories..."
                @focus="showRepoDropdown = true" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                  placeholder-gray-500 focus:outline-none focus:border-gray-600" />
              <div v-if="showRepoDropdown && filteredRepos.length > 0" class="absolute z-10 w-full mt-1 bg-gray-900 border border-gray-700 rounded-lg 
                  max-h-48 overflow-y-auto">
                <button v-for="repo in filteredRepos" :key="repo.id" @click="selectRepo(repo)"
                  class="w-full px-4 py-2.5 text-left hover:bg-gray-800 transition-colors">
                  <p class="text-white text-sm">{{ repo.full_name }}</p>
                  <p v-if="repo.description" class="text-gray-500 text-xs truncate">{{ repo.description }}</p>
                </button>
              </div>
              <div v-if="showRepoDropdown && repoSearch && filteredRepos.length === 0 && !loadingRepos"
                class="absolute z-10 w-full mt-1 bg-gray-900 border border-gray-700 rounded-lg p-4 text-center">
                <p class="text-gray-500 text-sm">No repositories found</p>
              </div>
              <div v-if="loadingRepos"
                class="absolute z-10 w-full mt-1 bg-gray-900 border border-gray-700 rounded-lg p-4 text-center">
                <div class="w-5 h-5 border-2 border-gray-600 border-t-blue-500 rounded-full animate-spin mx-auto" />
              </div>
            </div>
            <div v-if="selectedRepo"
              class="mt-2 p-3 bg-gray-900 border border-gray-700 rounded-lg flex items-center justify-between">
              <div>
                <p class="text-white text-sm">{{ selectedRepo.full_name }}</p>
                <p class="text-gray-500 text-xs">{{ selectedRepo.clone_url }}</p>
              </div>
              <button @click="clearRepo" class="text-gray-500 hover:text-white">
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Template Name</label>
            <input v-model="form.name" type="text" placeholder="Go API Template" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                placeholder-gray-500 focus:outline-none focus:border-gray-600" />
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Description</label>
            <textarea v-model="form.description" rows="2" placeholder="Standard Go API with Echo framework..." class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                placeholder-gray-500 focus:outline-none focus:border-gray-600 resize-none" />
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-gray-400 text-sm mb-2">Kind</label>
              <input v-model="form.kind" type="text" placeholder="api, cli, web..." class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                  placeholder-gray-500 focus:outline-none focus:border-gray-600" />
            </div>
            <div>
              <label class="block text-gray-400 text-sm mb-2">Language</label>
              <input v-model="form.language" type="text" placeholder="go, vue, python..." class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                  placeholder-gray-500 focus:outline-none focus:border-gray-600" />
            </div>
          </div>
        </div>

        <div v-if="error" class="mt-4 p-3 bg-red-500/10 border border-red-500/50 rounded-lg">
          <p class="text-red-400 text-sm">{{ error }}</p>
        </div>

        <div class="flex items-center justify-end gap-3 mt-6">
          <button @click="closeModal" class="px-4 py-2 text-gray-400 hover:text-white transition-colors">
            Cancel
          </button>
          <button @click="createTemplate" :disabled="!canCreate || saving" :class="[
            'px-4 py-2 rounded-lg font-medium transition-colors',
            canCreate && !saving
              ? 'bg-blue-600 hover:bg-blue-500 text-white'
              : 'bg-gray-700 text-gray-500 cursor-not-allowed'
          ]">
            {{ saving ? 'Adding...' : 'Add Template' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import {
  PlusIcon,
  MagnifyingGlassIcon,
  DocumentDuplicateIcon,
  CodeBracketIcon,
  TrashIcon,
  XMarkIcon,
  ChevronDownIcon
} from '@heroicons/vue/24/outline';
import type { Template } from '../types/Template';

interface Repo {
  id: number;
  full_name: string;
  description: string;
  html_url: string;
  clone_url: string;
}

const templateList = ref<Template[]>([]);
const loading = ref(true);
const search = ref('');
const kindFilter = ref('');
const languageFilter = ref('');
const showKindDropdown = ref(false);
const showLanguageDropdown = ref(false);
const showModal = ref(false);
const saving = ref(false);
const error = ref('');

// Repo picker state
const repos = ref<Repo[]>([]);
const repoSearch = ref('');
const showRepoDropdown = ref(false);
const loadingRepos = ref(false);
const selectedRepo = ref<Repo | null>(null);

const form = ref({
  name: '',
  description: '',
  kind: '',
  language: '',
});

const filteredTemplates = computed(() => {
  return templateList.value.filter(t => {
    const matchesSearch = !search.value ||
      t.name.toLowerCase().includes(search.value.toLowerCase()) ||
      t.description.toLowerCase().includes(search.value.toLowerCase());
    const matchesKind = !kindFilter.value || t.kind === kindFilter.value;
    const matchesLanguage = !languageFilter.value || t.language === languageFilter.value;
    return matchesSearch && matchesKind && matchesLanguage;
  });
});

const uniqueKinds = computed(() => {
  return [...new Set(templateList.value.map(t => t.kind).filter(Boolean))].sort();
});

const uniqueLanguages = computed(() => {
  return [...new Set(templateList.value.map(t => t.language).filter(Boolean))].sort();
});

const canCreate = computed(() => {
  return form.value.name && form.value.kind && form.value.language && selectedRepo.value;
});

const filteredRepos = computed(() => {
  if (!repoSearch.value) return repos.value.slice(0, 10);
  return repos.value.filter(r =>
    r.full_name.toLowerCase().includes(repoSearch.value.toLowerCase())
  ).slice(0, 10);
});

const resetForm = () => {
  form.value = { name: '', description: '', kind: '', language: '' };
  selectedRepo.value = null;
  repoSearch.value = '';
  showRepoDropdown.value = false;
  error.value = '';
};

const closeModal = () => {
  showModal.value = false;
  resetForm();
};

const selectRepo = (repo: Repo) => {
  selectedRepo.value = repo;
  repoSearch.value = '';
  showRepoDropdown.value = false;
  // Auto-fill name from repo if empty
  if (!form.value.name) {
    form.value.name = repo.full_name.split('/').pop() || '';
  }
};

const clearRepo = () => {
  selectedRepo.value = null;
};

const fetchRepos = async () => {
  loadingRepos.value = true;
  try {
    const response = await fetch('/api/forge/repos');
    if (!response.ok) throw new Error('Failed to fetch repos');
    repos.value = await response.json();
  } catch (e) {
    console.error('Failed to fetch repos:', e);
  } finally {
    loadingRepos.value = false;
  }
};

const fetchTemplates = async () => {
  try {
    loading.value = true;
    const response = await fetch('/api/templates');
    if (!response.ok) throw new Error(`HTTP error: ${response.status}`);
    templateList.value = await response.json();
  } catch (e) {
    console.error('Failed to fetch templates:', e);
  } finally {
    loading.value = false;
  }
};

const createTemplate = async () => {
  if (!selectedRepo.value) return;

  error.value = '';
  saving.value = true;

  try {
    const response = await fetch('/api/templates', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: form.value.name,
        description: form.value.description,
        kind: form.value.kind,
        language: form.value.language,
        repoUrl: selectedRepo.value.html_url,
      }),
    });

    if (!response.ok) {
      const data = await response.json();
      throw new Error(data.error || 'Failed to create template');
    }

    const template = await response.json();
    templateList.value.push(template);
    closeModal();
  } catch (e: any) {
    error.value = e.message;
  } finally {
    saving.value = false;
  }
};

const deleteTemplate = async (id: string) => {
  if (!confirm('Are you sure you want to delete this template?')) return;

  try {
    const response = await fetch(`/api/templates/${id}`, { method: 'DELETE' });
    if (!response.ok) throw new Error('Failed to delete');
    templateList.value = templateList.value.filter(t => t.id !== id);
  } catch (e) {
    console.error('Failed to delete template:', e);
  }
};

onMounted(() => {
  fetchTemplates();
  fetchRepos();
  document.addEventListener('click', closeDropdowns);
});

onUnmounted(() => {
  document.removeEventListener('click', closeDropdowns);
});

const closeDropdowns = (e: Event) => {
  const target = e.target as HTMLElement;
  if (!target.closest('.relative')) {
    showKindDropdown.value = false;
    showLanguageDropdown.value = false;
  }
};
</script>
