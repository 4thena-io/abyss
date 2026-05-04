<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-t1">Dashboard</h1>
      <p class="text-t2 mt-1">Overview of your development platform</p>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatsCard 
        v-for="stat in stats" 
        :key="stat.label" 
        :label="stat.label"
        :value="stat.value"
        :icon="stat.icon"
        :to="stat.to"
        :loading="loading"
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Activity -->
      <ActivityList 
        class="lg:col-span-2"
        title="Recent Activity"
        :items="activityItems"
        :loading="loadingActivity"
      />

      <!-- Quick Actions -->
      <div class="bg-panel rounded-lg border border-b1">
        <div class="px-5 py-4 border-b border-b1">
          <h2 class="text-lg font-medium text-t1">Quick Actions</h2>
        </div>
        <div class="p-4 space-y-2">
          <button
            v-for="action in quickActions"
            :key="action.label"
            @click="action.onClick"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left
                   bg-surface hover:bg-hover transition-colors group"
          >
            <component :is="action.icon" class="w-5 h-5 text-t3 group-hover:text-t1" />
            <span class="text-t2 group-hover:text-t1 text-sm">{{ action.label }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Recent Builds -->
    <div class="bg-panel rounded-lg border border-b1">
      <div class="px-5 py-4 border-b border-b1 flex items-center justify-between">
        <h2 class="text-lg font-medium text-t1">Recent Builds</h2>
        <router-link to="/apps" class="text-sm text-acc hover:opacity-80">View all</router-link>
      </div>

      <div v-if="loadingBuilds" class="px-5 py-8 flex justify-center">
        <Spinner size="md" />
      </div>

      <div v-else-if="recentBuilds.length === 0" class="px-5 py-8 text-center text-t3">
        No recent builds
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="text-left text-t3 text-sm border-b border-b1">
              <th class="px-5 py-3 font-medium">Application</th>
              <th class="px-5 py-3 font-medium">Build</th>
              <th class="px-5 py-3 font-medium">Status</th>
              <th class="px-5 py-3 font-medium">Branch</th>
              <th class="px-5 py-3 font-medium">Duration</th>
              <th class="px-5 py-3 font-medium">Time</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-b1">
            <tr
              v-for="build in recentBuilds"
              :key="`${build.appId}-${build.id}`"
              class="hover:bg-hover transition-colors"
            >
              <td class="px-5 py-4">
                <router-link :to="`/apps/${build.appId}`" class="text-t1 text-sm hover:text-acc">
                  {{ build.appName }}
                </router-link>
              </td>
              <td class="px-5 py-4">
                <a
                  :href="build.link"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-acc hover:opacity-80 text-sm"
                >
                  #{{ build.number }}
                </a>
              </td>
              <td class="px-5 py-4">
                <StatusBadge :status="build.status" />
              </td>
              <td class="px-5 py-4 text-t2 text-sm">{{ build.branch || '-' }}</td>
              <td class="px-5 py-4 text-t2 text-sm">{{ formatDuration(build.duration) }}</td>
              <td class="px-5 py-4 text-t3 text-sm">{{ formatTime(build.startedAt) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Project Modal -->
    <Modal :open="showProjectModal" title="New Project" @close="showProjectModal = false">
      <form @submit.prevent="createProject" class="space-y-4">
        <FormField 
          v-model="projectForm.name" 
          label="Name" 
          required 
          placeholder="my-project" 
        />
        <FormField 
          v-model="projectForm.description" 
          label="Description" 
          type="textarea" 
          :rows="3" 
          placeholder="Project description..." 
        />
        <ErrorAlert v-if="projectError" :message="projectError" />
        <FormActions 
          submit-label="Create Project"
          submitting-label="Creating..."
          :disabled="!projectForm.name"
          :saving="savingProject"
          @cancel="showProjectModal = false"
        />
      </form>
    </Modal>

    <!-- Create App Modal -->
    <Modal :open="showAppModal" title="New Application" size="lg" @close="showAppModal = false">
      <form @submit.prevent="createApp" class="space-y-4">
        <FormField 
          v-model="appForm.name" 
          label="Name" 
          required 
          placeholder="my-awesome-api" 
        />
        <Dropdown 
          v-model="appForm.projectId" 
          label="Project" 
          variant="form"
          required
          placeholder="Select project"
          :options="projectOptions"
        />
        <Dropdown 
          v-model="appForm.templateId" 
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
        <FormField 
          v-model="appForm.description" 
          label="Description" 
          type="textarea" 
          :rows="2" 
          placeholder="Application description..." 
        />
        <ErrorAlert v-if="appError" :message="appError" />
        <FormActions 
          submit-label="Create Application"
          submitting-label="Creating..."
          :disabled="!canCreateApp"
          :saving="savingApp"
          @cancel="showAppModal = false"
        />
      </form>
    </Modal>

    <!-- Create Template Modal -->
    <Modal :open="showTemplateModal" title="Add Template" size="lg" @close="showTemplateModal = false">
      <form @submit.prevent="createTemplate" class="space-y-4">
        <RepositoryPicker 
          v-model="selectedRepo"
          :repos="repoOptions"
          label="Repository"
          required
          @select="onRepoSelect"
        />
        <FormField 
          v-model="templateForm.name" 
          label="Template Name" 
          required 
          placeholder="Go API Template" 
        />
        <div class="grid grid-cols-2 gap-4">
          <Dropdown 
            v-model="templateForm.kind" 
            label="Kind" 
            variant="form"
            required
            placeholder="Select kind"
            :options="kindOptions"
          />
          <Dropdown 
            v-model="templateForm.language" 
            label="Language" 
            variant="form"
            required
            placeholder="Select language"
            :options="homeLanguageOptions"
          />
        </div>
        <FormField 
          v-model="templateForm.description" 
          label="Description" 
          type="textarea" 
          :rows="2" 
          placeholder="Template description..." 
        />
        <ErrorAlert v-if="templateError" :message="templateError" />
        <FormActions 
          submit-label="Add Template"
          submitting-label="Adding..."
          :disabled="!canCreateTemplate"
          :saving="savingTemplate"
          @cancel="showTemplateModal = false"
        />
      </form>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import {
  FolderIcon,
  CubeIcon,
  WrenchScrewdriverIcon,
  DocumentDuplicateIcon,
  PlusIcon,
  RocketLaunchIcon,
  Cog6ToothIcon,
} from '@heroicons/vue/24/outline';
import { appsApi, projectsApi, templatesApi, repoApi, type CreateAppRequest, type CreateProjectRequest, type CreateTemplateRequest } from '../api';
import type { App } from '../features/apps/types';
import type { Build } from '../features/apps/types';
import type { Project } from '../features/projects/types';
import type { Template } from '../features/templates/types';
import type { Repo } from '../types/Repo';
import Modal from '../components/ui/Modal.vue';
import Spinner from '../components/ui/Spinner.vue';
import StatusBadge from '../components/ui/StatusBadge.vue';
import ErrorAlert from '../components/ui/ErrorAlert.vue';
import FormField from '../components/ui/FormField.vue';
import Dropdown from '../components/ui/Dropdown.vue';
import FormActions from '../components/ui/FormActions.vue';
import StatsCard from '../components/ui/StatsCard.vue';
import ActivityList from '../components/ui/ActivityList.vue';
import RepositoryPicker, { type Repository } from '../components/ui/RepositoryPicker.vue';
import { useFormatters } from '../composables/useFormatters';

const router = useRouter();
const { formatDuration, formatRelativeTime: formatTime } = useFormatters();

// Loading states
const loading = ref(true);
const loadingActivity = ref(true);
const loadingBuilds = ref(true);

// Data
const apps = ref<App[]>([]);
const projects = ref<Project[]>([]);
const templates = ref<Template[]>([]);
const repos = ref<Repo[]>([]);

// Stats
const stats = computed(() => [
  { label: 'Projects', value: projects.value.length, icon: FolderIcon, to: '/projects' },
  { label: 'Applications', value: apps.value.length, icon: CubeIcon, to: '/apps' },
  { label: 'Templates', value: templates.value.length, icon: DocumentDuplicateIcon, to: '/templates' },
  { label: 'Total Builds', value: totalBuilds.value, icon: WrenchScrewdriverIcon, to: '/apps' },
]);

// Builds
const allBuilds = ref<(Build & { appId: number; appName: string })[]>([]);
const totalBuilds = computed(() => allBuilds.value.length);

const recentBuilds = computed(() => 
  allBuilds.value
    .sort((a, b) => new Date(b.startedAt).getTime() - new Date(a.startedAt).getTime())
    .slice(0, 5)
);

// Activity (derived from builds) - formatted for ActivityList component
const activityItems = computed(() => 
  allBuilds.value
    .sort((a, b) => new Date(b.startedAt).getTime() - new Date(a.startedAt).getTime())
    .slice(0, 5)
    .map(build => ({
      id: `${build.appId}-${build.id}`,
      title: build.appName,
      message: `build #${build.number} ${build.status}`,
      time: formatTime(build.startedAt),
      status: build.status,
      to: `/apps/${build.appId}`,
    }))
);

// Quick actions
const quickActions = computed(() => [
  { label: 'Create new project', icon: PlusIcon, onClick: () => showProjectModal.value = true },
  { label: 'Create new application', icon: RocketLaunchIcon, onClick: openAppModal },
  { label: 'Add template', icon: DocumentDuplicateIcon, onClick: openTemplateModal },
  { label: 'Settings', icon: Cog6ToothIcon, onClick: () => router.push('/settings') },
]);

// Project modal
const showProjectModal = ref(false);
const savingProject = ref(false);
const projectError = ref('');
const projectForm = reactive<CreateProjectRequest>({
  name: '',
  description: '',
});

// App modal
const showAppModal = ref(false);
const savingApp = ref(false);
const appError = ref('');
const appForm = reactive({
  name: '',
  description: '',
  projectId: '' as number | '',
  templateId: '' as number | '',
});

const selectedTemplate = computed(() => 
  templates.value.find(t => t.id === appForm.templateId)
);

const canCreateApp = computed(() => 
  appForm.name && appForm.projectId && appForm.templateId
);

// Template modal
const showTemplateModal = ref(false);
const savingTemplate = ref(false);
const templateError = ref('');
const selectedRepo = ref<Repository | null>(null);
const templateForm = reactive({
  name: '',
  description: '',
  kind: '',
  language: '',
});

// Repos formatted for RepositoryPicker
const repoOptions = computed(() => 
  repos.value.map(r => ({ id: r.id, fullName: r.fullName, url: r.url }))
);

const canCreateTemplate = computed(() => 
  templateForm.name && templateForm.kind && templateForm.language && selectedRepo.value
);

// Options for form selects
const projectOptions = computed(() =>
  projects.value.map(p => ({ value: p.id, label: p.name }))
);

const templateOptions = computed(() =>
  templates.value.map(t => ({ value: t.id, label: `${t.name} (${t.language})` }))
);

const kindOptions = [
  { value: 'api', label: 'API' },
  { value: 'web', label: 'Web' },
  { value: 'worker', label: 'Worker' },
  { value: 'cli', label: 'CLI' },
  { value: 'library', label: 'Library' },
];

const homeLanguageOptions = [
  { value: 'go', label: 'Go' },
  { value: 'typescript', label: 'TypeScript' },
  { value: 'python', label: 'Python' },
  { value: 'java', label: 'Java' },
  { value: 'rust', label: 'Rust' },
];

// Methods
const fetchData = async () => {
  try {
    loading.value = true;
    const [appsData, projectsData, templatesData] = await Promise.all([
      appsApi.getAll(),
      projectsApi.getAll(),
      templatesApi.getAll(),
    ]);
    apps.value = appsData;
    projects.value = projectsData;
    templates.value = templatesData;
  } catch (e) {
    console.error('Failed to fetch data:', e);
  } finally {
    loading.value = false;
  }
};

const fetchBuilds = async () => {
  loadingActivity.value = true;
  loadingBuilds.value = true;
  
  try {
    const buildsPromises = apps.value.map(async (app) => {
      try {
        const builds = await appsApi.getBuilds(app.id);
        return builds.map(build => ({
          ...build,
          appId: app.id,
          appName: app.name,
        }));
      } catch {
        return [];
      }
    });

    const results = await Promise.all(buildsPromises);
    allBuilds.value = results.flat();
  } catch (e) {
    console.error('Failed to fetch builds:', e);
  } finally {
    loadingActivity.value = false;
    loadingBuilds.value = false;
  }
};

const openAppModal = async () => {
  showAppModal.value = true;
  // Data already fetched
};

const openTemplateModal = async () => {
  showTemplateModal.value = true;
  if (repos.value.length === 0) {
    try {
      repos.value = await repoApi.getRepos();
    } catch (e) {
      console.error('Failed to fetch repos:', e);
    }
  }
};

const onRepoSelect = (repo: Repository) => {
  if (!templateForm.name) {
    templateForm.name = repo.fullName.split('/').pop() || '';
  }
};

const createProject = async () => {
  projectError.value = '';
  savingProject.value = true;

  try {
    const project = await projectsApi.create(projectForm);
    projects.value.push(project);
    showProjectModal.value = false;
    projectForm.name = '';
    projectForm.description = '';
    router.push(`/projects/${project.id}`);
  } catch (e: any) {
    projectError.value = e.message;
  } finally {
    savingProject.value = false;
  }
};

const createApp = async () => {
  if (!selectedTemplate.value) return;

  appError.value = '';
  savingApp.value = true;

  try {
    const request: CreateAppRequest = {
      name: appForm.name,
      description: appForm.description,
      projectId: appForm.projectId as number,
      templateId: appForm.templateId as number,
      kind: selectedTemplate.value.kind,
      language: selectedTemplate.value.language,
    };

    const app = await appsApi.create(request);
    apps.value.push(app);
    showAppModal.value = false;
    appForm.name = '';
    appForm.description = '';
    appForm.projectId = '';
    appForm.templateId = '';
    router.push(`/apps/${app.id}`);
  } catch (e: any) {
    appError.value = e.message;
  } finally {
    savingApp.value = false;
  }
};

const createTemplate = async () => {
  if (!selectedRepo.value) return;

  templateError.value = '';
  savingTemplate.value = true;

  try {
    const request: CreateTemplateRequest = {
      name: templateForm.name,
      description: templateForm.description,
      kind: templateForm.kind,
      language: templateForm.language,
      repoUrl: selectedRepo.value.url,
    };

    const template = await templatesApi.create(request);
    templates.value.push(template);
    showTemplateModal.value = false;
    templateForm.name = '';
    templateForm.description = '';
    templateForm.kind = '';
    templateForm.language = '';
    selectedRepo.value = null;
  } catch (e: any) {
    templateError.value = e.message;
  } finally {
    savingTemplate.value = false;
  }
};

// Lifecycle
onMounted(async () => {
  await fetchData();
  await fetchBuilds();
});
</script>
