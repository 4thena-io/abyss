<template>
  <div class="space-y-6">
    <PageHeader title="Templates" :subtitle="`${templateList.length} templates available`">
      <template #actions>
        <Button @click="openModal" variant="secondary">
          <PlusIcon class="w-4 h-4" />
          Add Template
        </Button>
      </template>
    </PageHeader>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <SearchInput v-model="search" placeholder="Filter templates..." />
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
      v-else-if="filteredTemplates.length === 0"
      :icon="DocumentDuplicateIcon"
      title="No templates found"
      :message="search || kindFilter || languageFilter ? 'Try adjusting your filters' : 'Add your first template to get started'"
    />

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
    <Modal :open="showModal" title="Add Template" size="lg" @close="closeModal">
      <form @submit.prevent="createTemplate" class="space-y-4">
        <!-- Repo Selection -->
        <div>
          <label class="block text-gray-400 text-sm mb-2">Repository *</label>
          <div class="relative">
            <input v-model="repoSearch" type="text" placeholder="Search repositories..."
              @focus="showRepoDropdown = true" class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                     placeholder-gray-500 focus:outline-none focus:border-gray-600" />

            <!-- Repo dropdown -->
            <div v-if="showRepoDropdown && filteredRepos.length > 0"
              class="absolute z-10 w-full mt-1 bg-gray-900 border border-gray-700 rounded-lg max-h-48 overflow-y-auto">
              <button v-for="repo in filteredRepos" :key="repo.id" type="button" @click="selectRepo(repo)"
                class="w-full px-4 py-2.5 text-left hover:bg-gray-800 transition-colors">
                <p class="text-white text-sm">{{ repo.fullName }}</p>
              </button>
            </div>

            <!-- No repos found -->
            <div v-if="showRepoDropdown && repoSearch && filteredRepos.length === 0 && !loadingRepos"
              class="absolute z-10 w-full mt-1 bg-gray-900 border border-gray-700 rounded-lg p-4 text-center">
              <p class="text-gray-500 text-sm">No repositories found</p>
            </div>

            <!-- Loading repos -->
            <div v-if="loadingRepos"
              class="absolute z-10 w-full mt-1 bg-gray-900 border border-gray-700 rounded-lg p-4 text-center">
              <Spinner size="sm" class="mx-auto" />
            </div>
          </div>

          <!-- Selected repo chip -->
          <div v-if="selectedRepo"
            class="mt-2 p-3 bg-gray-900 border border-gray-700 rounded-lg flex items-center justify-between">
            <p class="text-white text-sm">{{ selectedRepo.fullName }}</p>
            <button type="button" @click="clearRepo" class="text-gray-500 hover:text-white">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </div>

        <FormField 
          v-model="form.name" 
          label="Template Name" 
          required 
          placeholder="Go API Template" 
        />

        <FormField 
          v-model="form.description" 
          label="Description" 
          type="textarea" 
          :rows="2" 
          placeholder="Standard Go API with Echo framework..." 
        />

        <div class="grid grid-cols-2 gap-4">
          <Dropdown 
            v-model="form.kind" 
            label="Kind" 
            variant="form"
            required
            placeholder="Select kind"
            :options="kindOptions"
          />
          <Dropdown 
            v-model="form.language" 
            label="Language" 
            variant="form"
            required
            placeholder="Select language"
            :options="languageOptions"
          />
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
            {{ saving ? 'Adding...' : 'Add Template' }}
          </button>
        </div>
      </form>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue';
import {
  PlusIcon,
  DocumentDuplicateIcon,
  CodeBracketIcon,
  TrashIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline';
import type { Template } from '../types/Template';
import type { Repo } from '../types/Repo';
import { templatesApi, repoApi, type CreateTemplateRequest } from '../api';
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
const templateList = ref<Template[]>([]);
const loading = ref(true);

// Filter state
const search = ref('');
const kindFilter = ref('');
const languageFilter = ref('');

// Modal state
const showModal = ref(false);
const saving = ref(false);
const error = ref('');

// Repo picker state
const repos = ref<Repo[]>([]);
const repoSearch = ref('');
const showRepoDropdown = ref(false);
const loadingRepos = ref(false);
const selectedRepo = ref<Repo | null>(null);

// Form state
const form = reactive({
  name: '',
  description: '',
  kind: '',
  language: '',
});

// Computed
const filteredTemplates = computed(() => {
  return templateList.value.filter(t => {
    const matchesSearch = !search.value ||
      t.name.toLowerCase().includes(search.value.toLowerCase()) ||
      t.description?.toLowerCase().includes(search.value.toLowerCase());
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

const kindFilterOptions = computed(() => 
  uniqueKinds.value.map(k => ({ value: k, label: k }))
);

const languageFilterOptions = computed(() => 
  uniqueLanguages.value.map(l => ({ value: l, label: l }))
);

const canCreate = computed(() => {
  return form.name && form.kind && form.language && selectedRepo.value;
});

const kindOptions = [
  { value: 'api', label: 'API' },
  { value: 'web', label: 'Web' },
  { value: 'worker', label: 'Worker' },
  { value: 'cli', label: 'CLI' },
  { value: 'library', label: 'Library' },
  { value: 'job', label: 'Job' },
];

const languageOptions = [
  { value: 'go', label: 'Go' },
  { value: 'typescript', label: 'TypeScript' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'python', label: 'Python' },
  { value: 'java', label: 'Java' },
  { value: 'rust', label: 'Rust' },
  { value: 'csharp', label: 'C#' },
];

const filteredRepos = computed(() => {
  if (!repoSearch.value) return repos.value.slice(0, 10);
  return repos.value.filter(r =>
    r.fullName.toLowerCase().includes(repoSearch.value.toLowerCase())
  ).slice(0, 10);
});

// Methods
const resetForm = () => {
  form.name = '';
  form.description = '';
  form.kind = '';
  form.language = '';
  selectedRepo.value = null;
  repoSearch.value = '';
  showRepoDropdown.value = false;
  error.value = '';
};

const closeModal = () => {
  showModal.value = false;
  resetForm();
};

const openModal = async () => {
  showModal.value = true;
  await fetchRepos();
};

const selectRepo = (repo: Repo) => {
  selectedRepo.value = repo;
  repoSearch.value = '';
  showRepoDropdown.value = false;
  if (!form.name) {
    form.name = repo.fullName.split('/').pop() || '';
  }
};

const clearRepo = () => {
  selectedRepo.value = null;
};

const fetchRepos = async () => {
  loadingRepos.value = true;
  try {
    repos.value = await repoApi.getRepos();
  } catch (e) {
    console.error('Failed to fetch repos:', e);
  } finally {
    loadingRepos.value = false;
  }
};

const fetchTemplates = async () => {
  try {
    loading.value = true;
    templateList.value = await templatesApi.getAll();
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
    const request: CreateTemplateRequest = {
      name: form.name,
      description: form.description,
      kind: form.kind,
      language: form.language,
      repoUrl: selectedRepo.value.url,
    };

    const template = await templatesApi.create(request);
    templateList.value.push(template);
    closeModal();
  } catch (e: any) {
    error.value = e.message;
  } finally {
    saving.value = false;
  }
};

const deleteTemplate = async (id: number) => {
  if (!confirm('Are you sure you want to delete this template?')) return;

  try {
    await templatesApi.delete(id);
    templateList.value = templateList.value.filter(t => t.id !== id);
  } catch (e: any) {
    console.error('Failed to delete template:', e);
  }
};

// Lifecycle
onMounted(() => {
  fetchTemplates();
});
</script>
