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
      <EntityCard 
        v-for="template in filteredTemplates" 
        :key="template.id"
        :icon="DocumentDuplicateIcon"
        :title="template.name"
        :subtitle="template.kind"
        :description="template.description"
        :hoverable="false"
      >
        <template #header-action>
          <button @click="deleteTemplate(template.id)" class="text-t3 hover:text-fail transition-colors p-1">
            <TrashIcon class="w-4 h-4" />
          </button>
        </template>
        <template #badges>
          <TagGroup :tags="[template.language]" />
        </template>
        <template #footer>
          <a :href="template.repoUrl" target="_blank"
            class="flex items-center gap-2 text-t3 hover:text-t1 text-sm transition-colors">
            <CodeBracketIcon class="w-4 h-4" />
            <span>View repo</span>
          </a>
        </template>
      </EntityCard>
    </div>

    <!-- Add Template Modal -->
    <Modal :open="showModal" title="Add Template" size="lg" @close="closeModal">
      <form @submit.prevent="createTemplate" class="space-y-4">
        <!-- Repo Selection -->
        <div>
          <label class="block text-t3 text-sm mb-2">Repository *</label>
          <div class="relative">
            <input v-model="repoSearch" type="text" placeholder="Search repositories..."
              @focus="showRepoDropdown = true" class="w-full bg-bg border border-b1 rounded-lg px-4 py-2.5 text-t1
                     placeholder:text-t3 focus:outline-none focus:border-b2" />

            <!-- Repo dropdown -->
            <div v-if="showRepoDropdown && filteredRepos.length > 0"
              class="absolute z-10 w-full mt-1 bg-bg border border-b1 rounded-lg max-h-48 overflow-y-auto">
              <button v-for="repo in filteredRepos" :key="repo.id" type="button" @click="selectRepo(repo)"
                class="w-full px-4 py-2.5 text-left hover:bg-panel transition-colors">
                <p class="text-t1 text-sm">{{ repo.fullName }}</p>
              </button>
            </div>

            <!-- No repos found -->
            <div v-if="showRepoDropdown && repoSearch && filteredRepos.length === 0 && !loadingRepos"
              class="absolute z-10 w-full mt-1 bg-bg border border-b1 rounded-lg p-4 text-center">
              <p class="text-t3 text-sm">No repositories found</p>
            </div>

            <!-- Loading repos -->
            <div v-if="loadingRepos"
              class="absolute z-10 w-full mt-1 bg-bg border border-b1 rounded-lg p-4 text-center">
              <Spinner size="sm" class="mx-auto" />
            </div>
          </div>

          <!-- Selected repo chip -->
          <div v-if="selectedRepo"
            class="mt-2 p-3 bg-bg border border-b1 rounded-lg flex items-center justify-between">
            <p class="text-t1 text-sm">{{ selectedRepo.fullName }}</p>
            <button type="button" @click="clearRepo" class="text-t3 hover:text-t1 transition-colors">
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

        <FormActions 
          submit-label="Add Template"
          submitting-label="Adding..."
          :disabled="!canCreate"
          :saving="saving"
          @cancel="closeModal"
        />
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
import type { Template } from '../features/templates/types';
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
import FormActions from '../components/ui/FormActions.vue';
import EntityCard from '../components/ui/EntityCard.vue';
import TagGroup from '../components/ui/TagGroup.vue';

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
