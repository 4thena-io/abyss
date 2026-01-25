<template>
  <div class="space-y-6">
    <PageHeader title="Applications" :subtitle="`${appList.length} applications across all projects`">
      <template #actions>
        <Button @click="openModal" variant="secondary">
          <PlusIcon class="w-4 h-4" />
          New Application
        </Button>
      </template>
    </PageHeader>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <SearchInput v-model="search" placeholder="Filter applications..." />
      <Dropdown 
        v-model="kindFilter" 
        :options="kindFilterOptions" 
        all-label="All kinds" 
      />
      <Dropdown 
        v-model="languageFilter" 
        :options="languageFilterOptions" 
        all-label="All languages" 
      />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <!-- Empty State -->
    <EmptyState 
      v-else-if="filteredApps.length === 0"
      :icon="CubeIcon"
      title="No applications found"
      :message="search || kindFilter || languageFilter ? 'Try adjusting your filters' : 'Create your first application to get started'"
    />

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
        <FormField 
          v-model="form.name" 
          label="Name" 
          required 
          placeholder="my-awesome-api"
          hint="Will be used as repository name"
        />

        <Dropdown 
          v-model="form.projectId" 
          label="Project" 
          variant="form"
          required
          placeholder="Select project"
          :options="projectOptions"
        />

        <Dropdown 
          v-model="form.templateId" 
          label="Template" 
          variant="form"
          required
          placeholder="Select template"
          :options="templateOptions"
        />

        <!-- Auto-filled from template -->
        <div v-if="selectedTemplate" class="grid grid-cols-2 gap-4">
          <FormField 
            :model-value="selectedTemplate.kind" 
            label="Kind" 
            disabled
          />
          <FormField 
            :model-value="selectedTemplate.language" 
            label="Language" 
            disabled
          />
        </div>

        <FormField 
          v-model="form.description" 
          label="Description" 
          type="textarea" 
          :rows="2" 
          placeholder="Application description..."
        />

        <!-- Info box -->
        <div v-if="selectedTemplate && selectedProject" class="bg-blue-500/10 border border-blue-500/30 rounded-lg p-3">
          <p class="text-sm text-blue-400">
            This will create a new repository from the <strong>{{ selectedTemplate.name }}</strong> template,
            set up CI, and add it to the <strong>{{ selectedProject.name }}</strong> project.
          </p>
        </div>

        <ErrorAlert v-if="error" :message="error" />

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
import { ref, computed, reactive, onMounted } from 'vue';
import {
  PlusIcon,
  CubeIcon
} from '@heroicons/vue/24/outline';
import type { App } from '../types/App';
import type { Project } from '../types/Project';
import type { Template } from '../types/Template';
import { appsApi, projectsApi, templatesApi, type CreateAppRequest } from '../api';
import Modal from '../components/ui/Modal.vue';
import Button from '../components/ui/Button.vue';
import Spinner from '../components/ui/Spinner.vue';
import SearchInput from '../components/ui/SearchInput.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import ErrorAlert from '../components/ui/ErrorAlert.vue';
import FormField from '../components/ui/FormField.vue';
import PageHeader from '../components/ui/PageHeader.vue';
import Dropdown from '../components/ui/Dropdown.vue';

// List state
const appList = ref<App[]>([]);
const loading = ref(true);

// Filter state
const search = ref('');
const kindFilter = ref('');
const languageFilter = ref('');

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

const kindFilterOptions = computed(() => 
  uniqueKinds.value.map(k => ({ value: k, label: k }))
);

const languageFilterOptions = computed(() => 
  uniqueLanguages.value.map(l => ({ value: l, label: l }))
);

const selectedProject = computed(() =>
  projects.value.find(p => p.id === form.projectId)
);

const selectedTemplate = computed(() =>
  templates.value.find(t => t.id === form.templateId)
);

const canCreate = computed(() =>
  form.name.trim() !== '' && form.projectId !== '' && form.templateId !== ''
);

const projectOptions = computed(() =>
  projects.value.map(p => ({ value: p.id, label: p.name }))
);

const templateOptions = computed(() =>
  templates.value.map(t => ({ value: t.id, label: `${t.name} (${t.language})` }))
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

// Lifecycle
onMounted(() => {
  fetchApps();
});
</script>
