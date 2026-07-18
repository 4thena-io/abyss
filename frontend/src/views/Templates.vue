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
      <Dropdown v-model="kindFilter" :options="kindFilterOptions" all-label="All kinds" />
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
      :message="
        search || kindFilter || languageFilter
          ? 'Try adjusting your filters'
          : 'Add your first template to get started'
      "
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
        :to="`/templates/${template.id}`"
      >
        <template #badges>
          <TagGroup :tags="[template.language]" />
        </template>
      </EntityCard>
    </div>

    <!-- Add Template Modal -->
    <Modal :open="showModal" title="Add Template" size="lg" @close="closeModal">
      <form @submit.prevent="createTemplate" class="space-y-4">
        <!-- Source toggle -->
        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-t2">Source</label>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              @click="source = 'repo'"
              class="flex items-center gap-2 px-3 py-2.5 rounded-lg border text-sm font-medium transition-colors"
              :class="
                source === 'repo'
                  ? 'border-acc bg-acc-s text-t1'
                  : 'border-b1 bg-bg text-t2 hover:border-b2'
              "
            >
              <CodeBracketIcon class="w-4 h-4 shrink-0" />
              <div class="text-left">
                <p class="font-medium">From existing repo</p>
                <p class="text-xs font-normal opacity-70">Register a repo as a template</p>
              </div>
            </button>
            <button
              type="button"
              @click="source = 'blank'"
              class="flex items-center gap-2 px-3 py-2.5 rounded-lg border text-sm font-medium transition-colors"
              :class="
                source === 'blank'
                  ? 'border-acc bg-acc-s text-t1'
                  : 'border-b1 bg-bg text-t2 hover:border-b2'
              "
            >
              <DocumentPlusIcon class="w-4 h-4 shrink-0" />
              <div class="text-left">
                <p class="font-medium">Blank template</p>
                <p class="text-xs font-normal opacity-70">Create a new empty repo</p>
              </div>
            </button>
          </div>
        </div>

        <FormField v-model="form.name" label="Name" required placeholder="go-api-template" />
        <FormField
          v-model="form.description"
          label="Description"
          type="textarea"
          :rows="2"
          placeholder="Standard Go REST API..."
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

        <!-- Repo picker — only for "repo" source -->
        <template v-if="source === 'repo'">
          <div>
            <label class="block text-sm font-medium text-t2 mb-1.5">Repository *</label>
            <div class="relative">
              <input
                v-model="repoSearch"
                type="text"
                placeholder="Search repositories..."
                @focus="showRepoDropdown = true"
                class="w-full bg-bg border border-b1 rounded-lg px-4 py-2.5 text-t1 placeholder:text-t3 focus:outline-none focus:border-b2"
              />
              <div
                v-if="showRepoDropdown && filteredRepos.length > 0"
                class="absolute z-10 w-full mt-1 bg-bg border border-b1 rounded-lg max-h-48 overflow-y-auto shadow-lg"
              >
                <button
                  v-for="repo in filteredRepos"
                  :key="repo.id"
                  type="button"
                  @click="selectRepo(repo)"
                  class="w-full px-4 py-2.5 text-left hover:bg-panel transition-colors text-t1 text-sm"
                >
                  {{ repo.fullName }}
                </button>
              </div>
              <div
                v-if="showRepoDropdown && repoSearch && filteredRepos.length === 0 && !loadingRepos"
                class="absolute z-10 w-full mt-1 bg-bg border border-b1 rounded-lg p-4 text-center"
              >
                <p class="text-t3 text-sm">No repositories found</p>
              </div>
              <div
                v-if="loadingRepos"
                class="absolute z-10 w-full mt-1 bg-bg border border-b1 rounded-lg p-4 text-center"
              >
                <Spinner size="sm" class="mx-auto" />
              </div>
            </div>
            <div
              v-if="selectedRepo"
              class="mt-2 p-3 bg-bg border border-b1 rounded-lg flex items-center justify-between"
            >
              <p class="text-t1 text-sm font-mono">{{ selectedRepo.fullName }}</p>
              <button
                type="button"
                @click="clearRepo"
                class="text-t3 hover:text-t1 transition-colors"
              >
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </template>

        <!-- Blank info -->
        <div
          v-else
          class="flex items-start gap-2 bg-acc-s border border-acc-b rounded-lg px-3 py-2.5"
        >
          <InformationCircleIcon class="w-4 h-4 text-acc shrink-0 mt-0.5" />
          <p class="text-acc text-xs">
            A new empty repository named <strong>{{ form.name || '…' }}</strong> will be created in
            the forge and registered as a template.
          </p>
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
  DocumentPlusIcon,
  XMarkIcon,
  InformationCircleIcon,
} from '@heroicons/vue/24/outline';
import type { Template } from '../features/templates/types';
import type { Repo } from '../types/Repo';
import { templatesApi, type CreateTemplateRequest } from '../api/template';
import { repoApi } from '../api/repo';
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

const templateList = ref<Template[]>([]);
const loading = ref(true);
const search = ref('');
const kindFilter = ref('');
const languageFilter = ref('');

const showModal = ref(false);
const saving = ref(false);
const error = ref('');
const source = ref<'repo' | 'blank'>('repo');

const repos = ref<Repo[]>([]);
const repoSearch = ref('');
const showRepoDropdown = ref(false);
const loadingRepos = ref(false);
const selectedRepo = ref<Repo | null>(null);

const form = reactive({ name: '', description: '', kind: '', language: '' });

const filteredTemplates = computed(() =>
  templateList.value.filter((t) => {
    const matchesSearch =
      !search.value ||
      t.name.toLowerCase().includes(search.value.toLowerCase()) ||
      t.description?.toLowerCase().includes(search.value.toLowerCase());
    return (
      matchesSearch &&
      (!kindFilter.value || t.kind === kindFilter.value) &&
      (!languageFilter.value || t.language === languageFilter.value)
    );
  }),
);

const uniqueKinds = computed(() =>
  [...new Set(templateList.value.map((t) => t.kind).filter(Boolean))].sort(),
);
const uniqueLanguages = computed(() =>
  [...new Set(templateList.value.map((t) => t.language).filter(Boolean))].sort(),
);
const kindFilterOptions = computed(() => uniqueKinds.value.map((k) => ({ value: k, label: k })));
const languageFilterOptions = computed(() =>
  uniqueLanguages.value.map((l) => ({ value: l, label: l })),
);

const canCreate = computed(() => {
  if (!form.name || !form.kind || !form.language) return false;
  if (source.value === 'repo') return !!selectedRepo.value;
  return true;
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
  return repos.value
    .filter((r) => r.fullName.toLowerCase().includes(repoSearch.value.toLowerCase()))
    .slice(0, 10);
});

function resetForm() {
  form.name = '';
  form.description = '';
  form.kind = '';
  form.language = '';
  source.value = 'repo';
  selectedRepo.value = null;
  repoSearch.value = '';
  showRepoDropdown.value = false;
  error.value = '';
}

function closeModal() {
  showModal.value = false;
  resetForm();
}

async function openModal() {
  showModal.value = true;
  loadingRepos.value = true;
  try {
    repos.value = await repoApi.getRepos();
  } catch (e) {
    console.error('Failed to fetch repos:', e);
  } finally {
    loadingRepos.value = false;
  }
}

function selectRepo(repo: Repo) {
  selectedRepo.value = repo;
  repoSearch.value = '';
  showRepoDropdown.value = false;
  if (!form.name) form.name = repo.fullName.split('/').pop() || '';
}

function clearRepo() {
  selectedRepo.value = null;
}

async function fetchTemplates() {
  try {
    loading.value = true;
    templateList.value = await templatesApi.getAll();
  } catch (e) {
    console.error('Failed to fetch templates:', e);
  } finally {
    loading.value = false;
  }
}

async function createTemplate() {
  error.value = '';
  saving.value = true;
  try {
    const req: CreateTemplateRequest = {
      source: source.value,
      name: form.name,
      description: form.description,
      kind: form.kind,
      language: form.language,
      ...(source.value === 'repo' && selectedRepo.value ? { repoUrl: selectedRepo.value.url } : {}),
    };
    const template = await templatesApi.create(req);
    templateList.value.push(template);
    closeModal();
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'An error occurred';
  } finally {
    saving.value = false;
  }
}

onMounted(fetchTemplates);
</script>
