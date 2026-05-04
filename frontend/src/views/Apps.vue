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
      <EntityCard 
        v-for="app in filteredApps" 
        :key="app.id"
        :icon="CubeIcon"
        :title="app.name"
        :subtitle="app.kind"
        :description="app.description"
        :show-no-description="false"
        :to="`/apps/${app.id}`"
      >
        <template #badges>
          <TagGroup :tags="[app.kind, app.language].filter(Boolean) as string[]" />
        </template>
      </EntityCard>
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

        <FormActions 
          submit-label="Create Application"
          submitting-label="Creating..."
          :disabled="!canCreate"
          :saving="saving"
          @cancel="closeModal"
        />
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
import type { App } from '../features/apps/types';
import type { Project } from '../features/projects/types';
import type { Template } from '../features/templates/types';
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
import FormActions from '../components/ui/FormActions.vue';
import EntityCard from '../components/ui/EntityCard.vue';
import TagGroup from '../components/ui/TagGroup.vue';

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
