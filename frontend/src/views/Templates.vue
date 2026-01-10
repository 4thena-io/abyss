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
      <select v-model="kindFilter" class="bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-300 
          focus:outline-none focus:border-gray-600">
        <option value="">All types</option>
        <option value="api">API</option>
        <option value="cli">CLI</option>
        <option value="web">Web</option>
        <option value="mobile">Mobile</option>
      </select>
      <select v-model="languageFilter" class="bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-300 
          focus:outline-none focus:border-gray-600">
        <option value="">All languages</option>
        <option value="go">Go</option>
        <option value="vue">Vue</option>
        <option value="flutter">Flutter</option>
        <option value="react">React</option>
        <option value="python">Python</option>
      </select>
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
            <div :class="[
              'w-10 h-10 rounded-lg flex items-center justify-center',
              getKindColor(template.kind)
            ]">
              <component :is="getKindIcon(template.kind)" class="w-5 h-5" />
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

        <div class="mt-4 pt-4 border-t border-gray-700 flex items-center justify-between">
          <a :href="template.repo_url" target="_blank"
            class="flex items-center gap-2 text-gray-400 hover:text-white text-sm transition-colors">
            <CodeBracketIcon class="w-4 h-4" />
            <span>View repo</span>
          </a>
          <span class="text-gray-500 text-xs">
            Added {{ formatDate(template.created_at) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Add Template Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/60" @click="closeModal" />
      <div class="relative bg-gray-800 border border-gray-700 rounded-lg w-full max-w-md p-6">
        <h2 class="text-xl font-semibold text-white mb-4">Add Template</h2>

        <div class="space-y-4">
          <div>
            <label class="block text-gray-400 text-sm mb-2">Name</label>
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
              <select v-model="form.kind" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                  focus:outline-none focus:border-gray-600">
                <option value="">Select type</option>
                <option value="api">API</option>
                <option value="cli">CLI</option>
                <option value="web">Web</option>
                <option value="mobile">Mobile</option>
              </select>
            </div>
            <div>
              <label class="block text-gray-400 text-sm mb-2">Language</label>
              <select v-model="form.language" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                  focus:outline-none focus:border-gray-600">
                <option value="">Select language</option>
                <option value="go">Go</option>
                <option value="vue">Vue</option>
                <option value="flutter">Flutter</option>
                <option value="react">React</option>
                <option value="python">Python</option>
              </select>
            </div>
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Repository URL</label>
            <input v-model="form.repo_url" type="url" placeholder="https://git.company.io/templates/go-api" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                placeholder-gray-500 focus:outline-none focus:border-gray-600" />
          </div>

          <div>
            <label class="block text-gray-400 text-sm mb-2">Clone URL</label>
            <input v-model="form.clone_url" type="url" placeholder="git@git.company.io:templates/go-api.git" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                placeholder-gray-500 focus:outline-none focus:border-gray-600" />
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
import { ref, computed, onMounted } from 'vue';
import {
  PlusIcon,
  MagnifyingGlassIcon,
  DocumentDuplicateIcon,
  CodeBracketIcon,
  TrashIcon,
  CubeIcon,
  CommandLineIcon,
  ComputerDesktopIcon,
  DevicePhoneMobileIcon
} from '@heroicons/vue/24/outline';
import type { Template } from '../types/Template';

const templateList = ref<Template[]>([]);
const loading = ref(true);
const search = ref('');
const kindFilter = ref('');
const languageFilter = ref('');
const showModal = ref(false);
const saving = ref(false);
const error = ref('');

const form = ref({
  name: '',
  description: '',
  kind: '',
  language: '',
  repo_url: '',
  clone_url: '',
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

const canCreate = computed(() => {
  return form.value.name && form.value.kind && form.value.language && form.value.repo_url && form.value.clone_url;
});

const getKindIcon = (kind: string) => {
  const icons: Record<string, any> = {
    api: CubeIcon,
    cli: CommandLineIcon,
    web: ComputerDesktopIcon,
    mobile: DevicePhoneMobileIcon,
  };
  return icons[kind] || DocumentDuplicateIcon;
};

const getKindColor = (kind: string) => {
  const colors: Record<string, string> = {
    api: 'bg-blue-500/20 text-blue-400',
    cli: 'bg-green-500/20 text-green-400',
    web: 'bg-purple-500/20 text-purple-400',
    mobile: 'bg-orange-500/20 text-orange-400',
  };
  return colors[kind] || 'bg-gray-700 text-gray-400';
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (days === 0) return 'today';
  if (days === 1) return 'yesterday';
  if (days < 7) return `${days} days ago`;
  return date.toLocaleDateString();
};

const resetForm = () => {
  form.value = { name: '', description: '', kind: '', language: '', repo_url: '', clone_url: '' };
  error.value = '';
};

const closeModal = () => {
  showModal.value = false;
  resetForm();
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
  error.value = '';
  saving.value = true;

  try {
    const response = await fetch('/api/templates', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value),
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

onMounted(fetchTemplates);
</script>
