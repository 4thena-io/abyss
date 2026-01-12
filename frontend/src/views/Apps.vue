<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-white">Applications</h1>
        <p class="text-gray-400 mt-1">{{ appList.length }} applications across all projects</p>
      </div>
      <Button @click="openModal" variant="secondary">
        <PlusIcon class="w-4 h-4" />
        New Application
      </Button>
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

    <!-- Create App Modal -->
    <Modal :open="showModal" title="New Application" size="lg" @close="closeModal">
      <form @submit.prevent="createApp" class="space-y-4">
        <div>
          <label class="block text-gray-400 text-sm mb-2">Name *</label>
          <input v-model="form.name" type="text" placeholder="my-awesome-api" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                   placeholder-gray-500 focus:outline-none focus:border-gray-600" />
          <p class="text-xs text-gray-500 mt-1">Will be used as repository name</p>
        </div>

        <div>
          <label class="block text-gray-400 text-sm mb-2">Project *</label>
          <select v-model="form.projectId" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                   focus:outline-none focus:border-gray-600">
            <option value="">Select project</option>
            <option v-for="project in projects" :key="project.id" :value="project.id">
              {{ project.name }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-gray-400 text-sm mb-2">Template *</label>
          <select v-model="form.templateId" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                   focus:outline-none focus:border-gray-600">
            <option value="">Select template</option>
            <option v-for="template in templates" :key="template.id" :value="template.id">
              {{ template.name }} ({{ template.language }})
            </option>
          </select>
        </div>

        <!-- Auto-filled from template -->
        <div v-if="selectedTemplate" class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-gray-400 text-sm mb-2">Kind</label>
            <input :value="selectedTemplate.kind" disabled
              class="w-full bg-gray-900/50 border border-gray-700 rounded-lg px-4 py-2.5 text-gray-400" />
          </div>
          <div>
            <label class="block text-gray-400 text-sm mb-2">Language</label>
            <input :value="selectedTemplate.language" disabled
              class="w-full bg-gray-900/50 border border-gray-700 rounded-lg px-4 py-2.5 text-gray-400" />
          </div>
        </div>

        <div>
          <label class="block text-gray-400 text-sm mb-2">Description</label>
          <textarea v-model="form.description" rows="2" placeholder="Application description..." class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                   placeholder-gray-500 focus:outline-none focus:border-gray-600 resize-none" />
        </div>

        <!-- Info box -->
        <div v-if="selectedTemplate && selectedProject" class="bg-blue-500/10 border border-blue-500/30 rounded-lg p-3">
          <p class="text-sm text-blue-400">
            This will create a new repository from the <strong>{{ selectedTemplate.name }}</strong> template,
            set up CI, and add it to the <strong>{{ selectedProject.name }}</strong> project.
          </p>
        </div>

        <div v-if="error" class="p-3 bg-red-500/10 border border-red-500/50 rounded-lg">
          <p class="text-red-400 text-sm">{{ error }}</p>
        </div>

        <div class="flex items-center justify-end gap-3 pt-2">
          <button type="button" @click="closeModal" class="px-4 py-2 text-gray-400 hover:text-white transition-colors">
            Cancel
          </button>
          <button type="submit" :disabled="!canCreate || saving" :class="[
            'px-4 py-2 rounded-lg font-medium transition-colors',
            canCreate && !saving
              ? 'bg-blue-600 hover:bg-blue-500 text-white'
              : 'bg-gray-700 text-gray-500 cursor-not-allowed'
          ]">
            {{ saving ? 'Creating...' : 'Create Application' }}
          </button>
        </div>
      </form>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted, onUnmounted } from 'vue';
import {
  PlusIcon,
  MagnifyingGlassIcon,
  CubeIcon,
  ChevronDownIcon
} from '@heroicons/vue/24/outline';
import type { App } from '../types/App';
import type { Project } from '../types/Project';
import type { Template } from '../types/Template';
import { appsApi, projectsApi, templatesApi, type CreateAppRequest } from '../api';
import Modal from '../components/ui/Modal.vue';
import Button from '../components/ui/Button.vue';

// List state
const appList = ref<App[]>([]);
const loading = ref(true);

// Filter state
const search = ref('');
const kindFilter = ref('');
const languageFilter = ref('');
const showKindDropdown = ref(false);
const showLanguageDropdown = ref(false);

// Modal state
const showModal = ref(false);
const saving = ref(false);
const error = ref('');

// Related data for form
const projects = ref<Project[]>([]);
const templates = ref<Template[]>([]);

// Form state
const form = reactive({
  name: '',
  description: '',
  projectId: '' as number | '',
  templateId: '' as number | '',
});

// Computed
const filteredApps = computed(() => {
  return appList.value.filter(app => {
    const matchesSearch = !search.value ||
      app.name.toLowerCase().includes(search.value.toLowerCase()) ||
      app.description?.toLowerCase().includes(search.value.toLowerCase());
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

const selectedProject = computed(() =>
  projects.value.find(p => p.id === form.projectId)
);

const selectedTemplate = computed(() =>
  templates.value.find(t => t.id === form.templateId)
);

const canCreate = computed(() =>
  form.name.trim() !== '' && form.projectId !== '' && form.templateId !== ''
);

// Methods
const resetForm = () => {
  form.name = '';
  form.description = '';
  form.projectId = '';
  form.templateId = '';
  error.value = '';
};

const closeModal = () => {
  showModal.value = false;
  resetForm();
};

const openModal = async () => {
  showModal.value = true;
  // Fetch projects and templates for dropdowns
  try {
    const [projectsData, templatesData] = await Promise.all([
      projectsApi.getAll(),
      templatesApi.getAll(),
    ]);
    projects.value = projectsData;
    templates.value = templatesData;
  } catch (e) {
    console.error('Failed to fetch data:', e);
  }
};

const fetchApps = async () => {
  try {
    loading.value = true;
    appList.value = await appsApi.getAll();
  } catch (e) {
    console.error('Failed to fetch apps:', e);
  } finally {
    loading.value = false;
  }
};

const createApp = async () => {
  if (!selectedTemplate.value) return;

  error.value = '';
  saving.value = true;

  try {
    const request: CreateAppRequest = {
      name: form.name,
      description: form.description,
      projectId: form.projectId as number,
      templateId: form.templateId as number,
      kind: selectedTemplate.value.kind,
      language: selectedTemplate.value.language,
    };

    const app = await appsApi.create(request);
    appList.value.push(app);
    closeModal();
  } catch (e: any) {
    error.value = e.message;
  } finally {
    saving.value = false;
  }
};

const closeDropdowns = (e: Event) => {
  const target = e.target as HTMLElement;
  if (!target.closest('.relative')) {
    showKindDropdown.value = false;
    showLanguageDropdown.value = false;
  }
};

// Lifecycle
onMounted(() => {
  fetchApps();
  document.addEventListener('click', closeDropdowns);
});

onUnmounted(() => {
  document.removeEventListener('click', closeDropdowns);
});
</script>
