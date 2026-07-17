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
        <!-- Source toggle -->
        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-t2">Source</label>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              @click="source = 'template'"
              class="flex items-center gap-2 px-3 py-2.5 rounded-lg border text-sm font-medium transition-colors"
              :class="source === 'template'
                ? 'border-acc bg-acc-s text-t1'
                : 'border-b1 bg-bg text-t2 hover:border-b2'"
            >
              <DocumentDuplicateIcon class="w-4 h-4 shrink-0" />
              <div class="text-left">
                <p class="font-medium">From template</p>
                <p class="text-xs font-normal opacity-70">Create a new repo from a template</p>
              </div>
            </button>
            <button
              type="button"
              @click="source = 'repo'"
              class="flex items-center gap-2 px-3 py-2.5 rounded-lg border text-sm font-medium transition-colors"
              :class="source === 'repo'
                ? 'border-acc bg-acc-s text-t1'
                : 'border-b1 bg-bg text-t2 hover:border-b2'"
            >
              <CodeBracketIcon class="w-4 h-4 shrink-0" />
              <div class="text-left">
                <p class="font-medium">From existing repo</p>
                <p class="text-xs font-normal opacity-70">Link an already existing repository</p>
              </div>
            </button>
          </div>
        </div>

        <FormField
          v-model="form.name"
          label="Name"
          required
          placeholder="my-awesome-api"
          hint="Will be used as repository name"
        />

        <!-- Team filter — optional, narrows project list -->
        <Dropdown
          v-model="teamFilter"
          label="Team"
          variant="form"
          placeholder="All teams"
          :options="teamOptions"
        />

        <!-- Project — optional -->
        <Dropdown
          v-model="form.projectId"
          label="Project"
          variant="form"
          placeholder="No project (standalone)"
          :options="projectOptions"
        />

        <!-- Template source -->
        <template v-if="source === 'template'">
          <Dropdown
            v-model="form.templateId"
            label="Template"
            variant="form"
            required
            placeholder="Select template"
            :options="templateOptions"
          />

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
        </template>

        <!-- Repo source -->
        <template v-else>
          <Dropdown
            v-model="form.repoId"
            label="Repository"
            variant="form"
            required
            placeholder="Select repository"
            :options="repoOptions"
          />
          <div class="grid grid-cols-2 gap-4">
            <FormField v-model="form.kind" label="Kind" placeholder="service, library…" />
            <FormField v-model="form.language" label="Language" placeholder="go, python…" />
          </div>
        </template>

        <FormField
          v-model="form.description"
          label="Description"
          type="textarea"
          :rows="2"
          placeholder="Application description..."
        />

        <!-- Summary info box -->
        <div v-if="source === 'template' && selectedTemplate" class="bg-acc-s border border-acc-b rounded-lg p-3">
          <p class="text-sm text-acc">
            Creates a new repository from <strong>{{ selectedTemplate.name }}</strong>
            <template v-if="selectedProject">
              and adds it to <strong>{{ selectedProject.name }}</strong>
              <template v-if="selectedTeam"> under <strong>{{ selectedTeam.name }}</strong></template>
            </template>.
          </p>
        </div>
        <div v-else-if="source === 'repo' && selectedRepo" class="bg-acc-s border border-acc-b rounded-lg p-3">
          <p class="text-sm text-acc">
            Links <strong>{{ selectedRepo.fullName }}</strong>
            <template v-if="selectedProject">
              to <strong>{{ selectedProject.name }}</strong>
              <template v-if="selectedTeam"> under <strong>{{ selectedTeam.name }}</strong></template>
            </template>.
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
import { ref, computed, reactive, watch, onMounted } from 'vue';
import {
  PlusIcon,
  CubeIcon,
  DocumentDuplicateIcon,
  CodeBracketIcon,
} from '@heroicons/vue/24/outline';
import type { App } from '../features/apps/types';
import type { Project } from '../features/projects/types';
import type { Template } from '../features/templates/types';
import type { Team } from '../features/teams/types';
import type { Repo } from '../types/Repo';
import { appsApi, projectsApi, templatesApi, teamsApi, type CreateAppRequest } from '../api';
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
const source = ref<'template' | 'repo'>('template');

// Related data for form
const projects = ref<Project[]>([]);
const teams = ref<Team[]>([]);
const templates = ref<Template[]>([]);
const repos = ref<Repo[]>([]);

// Form state
const teamFilter = ref<number | ''>('');
watch(teamFilter, () => {
  const project = projects.value.find(p => p.id === form.projectId);
  if (project && teamFilter.value && project.teamId !== teamFilter.value) {
    form.projectId = '';
  }
});
const form = reactive({
  name: '',
  description: '',
  projectId: '' as number | '',
  templateId: '' as number | '',
  repoId: '' as number | '',
  kind: '',
  language: '',
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

const selectedTeam = computed(() =>
  teams.value.find(t => t.id === teamFilter.value)
);

const selectedProject = computed(() =>
  projects.value.find(p => p.id === form.projectId)
);

const selectedTemplate = computed(() =>
  templates.value.find(t => t.id === form.templateId)
);

const selectedRepo = computed(() =>
  repos.value.find(r => r.id === form.repoId)
);

const canCreate = computed(() => {
  if (!form.name.trim()) return false;
  if (source.value === 'template') return form.templateId !== '';
  return form.repoId !== '';
});

const teamOptions = computed(() =>
  teams.value.map(t => ({ value: t.id, label: t.name }))
);

const projectOptions = computed(() => {
  const filtered = teamFilter.value
    ? projects.value.filter(p => p.teamId === teamFilter.value)
    : projects.value;
  return filtered.map(p => ({ value: p.id, label: p.name }));
});

const templateOptions = computed(() =>
  templates.value.map(t => ({ value: t.id, label: `${t.name} (${t.language})` }))
);

const repoOptions = computed(() =>
  repos.value.map(r => ({ value: r.id, label: r.fullName }))
);

// Methods
const resetForm = () => {
  form.name = '';
  form.description = '';
  form.projectId = '';
  form.templateId = '';
  form.repoId = '';
  form.kind = '';
  form.language = '';
  source.value = 'template';
  teamFilter.value = '';
  error.value = '';
};

const closeModal = () => {
  showModal.value = false;
  resetForm();
};

const openModal = async () => {
  showModal.value = true;
  try {
    const [projectsData, teamsData, templatesData, reposData] = await Promise.all([
      projectsApi.getAll(),
      teamsApi.getAll(),
      templatesApi.getAll(),
      repoApi.getRepos(),
    ]);
    projects.value = projectsData;
    teams.value = teamsData;
    templates.value = templatesData;
    repos.value = reposData;
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
  error.value = '';
  saving.value = true;

  try {
    let request: CreateAppRequest;

    if (source.value === 'template') {
      if (!selectedTemplate.value) return;
      request = {
        name: form.name,
        description: form.description,
        projectId: form.projectId as number || undefined,
        templateId: form.templateId as number,
        kind: selectedTemplate.value.kind,
        language: selectedTemplate.value.language,
      };
    } else {
      request = {
        name: form.name,
        description: form.description,
        projectId: form.projectId as number || undefined,
        repoId: form.repoId as number,
        kind: form.kind,
        language: form.language,
      };
    }

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
