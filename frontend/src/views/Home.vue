<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-white">Dashboard</h1>
      <p class="text-gray-400 mt-1">Overview of your development platform</p>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <router-link 
        v-for="stat in stats" 
        :key="stat.label" 
        :to="stat.to"
        class="bg-gray-800 rounded-lg p-5 border border-gray-700 hover:border-gray-600 transition-colors"
      >
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-400 text-sm">{{ stat.label }}</p>
            <p class="text-2xl font-semibold text-white mt-1">
              <span v-if="loading" class="inline-block w-8 h-6 bg-gray-700 rounded animate-pulse"></span>
              <span v-else>{{ stat.value }}</span>
            </p>
          </div>
          <component :is="stat.icon" class="w-8 h-8 text-gray-600" />
        </div>
      </router-link>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Activity -->
      <div class="lg:col-span-2 bg-gray-800 rounded-lg border border-gray-700">
        <div class="px-5 py-4 border-b border-gray-700">
          <h2 class="text-lg font-medium text-white">Recent Activity</h2>
        </div>
        <div class="divide-y divide-gray-700">
          <div v-if="loadingActivity" class="px-5 py-8 flex justify-center">
            <Spinner size="md" />
          </div>
          <template v-else-if="recentActivity.length > 0">
            <div 
              v-for="activity in recentActivity" 
              :key="activity.id"
              class="px-5 py-4 flex items-start gap-4 hover:bg-gray-750 transition-colors"
            >
              <div :class="[
                'w-2 h-2 rounded-full mt-2 shrink-0',
                activity.status === 'success' ? 'bg-green-500' :
                activity.status === 'running' ? 'bg-blue-500' :
                activity.status === 'failure' ? 'bg-red-500' : 'bg-gray-500'
              ]" />
              <div class="flex-1 min-w-0">
                <p class="text-white text-sm">
                  <router-link :to="`/apps/${activity.appId}`" class="font-medium hover:text-blue-400">
                    {{ activity.appName }}
                  </router-link>
                  <span class="text-gray-400"> {{ activity.message }}</span>
                </p>
                <p class="text-gray-500 text-xs mt-1">{{ activity.time }}</p>
              </div>
            </div>
          </template>
          <div v-else class="px-5 py-8 text-center text-gray-500">
            No recent activity
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-gray-800 rounded-lg border border-gray-700">
        <div class="px-5 py-4 border-b border-gray-700">
          <h2 class="text-lg font-medium text-white">Quick Actions</h2>
        </div>
        <div class="p-4 space-y-2">
          <button 
            v-for="action in quickActions" 
            :key="action.label" 
            @click="action.onClick"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left
                   bg-gray-700/50 hover:bg-gray-700 transition-colors group"
          >
            <component :is="action.icon" class="w-5 h-5 text-gray-400 group-hover:text-white" />
            <span class="text-gray-300 group-hover:text-white text-sm">{{ action.label }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Recent Builds -->
    <div class="bg-gray-800 rounded-lg border border-gray-700">
      <div class="px-5 py-4 border-b border-gray-700 flex items-center justify-between">
        <h2 class="text-lg font-medium text-white">Recent Builds</h2>
        <router-link to="/apps" class="text-sm text-blue-400 hover:text-blue-300">View all</router-link>
      </div>
      
      <div v-if="loadingBuilds" class="px-5 py-8 flex justify-center">
        <Spinner size="md" />
      </div>
      
      <div v-else-if="recentBuilds.length === 0" class="px-5 py-8 text-center text-gray-500">
        No recent builds
      </div>
      
      <div v-else class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="text-left text-gray-400 text-sm border-b border-gray-700">
              <th class="px-5 py-3 font-medium">Application</th>
              <th class="px-5 py-3 font-medium">Build</th>
              <th class="px-5 py-3 font-medium">Status</th>
              <th class="px-5 py-3 font-medium">Branch</th>
              <th class="px-5 py-3 font-medium">Duration</th>
              <th class="px-5 py-3 font-medium">Time</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-700">
            <tr 
              v-for="build in recentBuilds" 
              :key="`${build.appId}-${build.id}`" 
              class="hover:bg-gray-750 transition-colors"
            >
              <td class="px-5 py-4">
                <router-link :to="`/apps/${build.appId}`" class="text-white text-sm hover:text-blue-400">
                  {{ build.appName }}
                </router-link>
              </td>
              <td class="px-5 py-4">
                <a 
                  :href="build.link" 
                  target="_blank" 
                  rel="noopener noreferrer"
                  class="text-blue-400 hover:text-blue-300 text-sm"
                >
                  #{{ build.number }}
                </a>
              </td>
              <td class="px-5 py-4">
                <StatusBadge :status="build.status" />
              </td>
              <td class="px-5 py-4 text-gray-400 text-sm">{{ build.branch || '-' }}</td>
              <td class="px-5 py-4 text-gray-400 text-sm">{{ formatDuration(build.duration) }}</td>
              <td class="px-5 py-4 text-gray-500 text-sm">{{ formatTime(build.startedAt) }}</td>
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
        <div class="flex items-center justify-end gap-3 pt-2">
          <button type="button" @click="showProjectModal = false" class="px-4 py-2 text-gray-400 hover:text-white">
            Cancel
          </button>
          <button 
            type="submit"
            :disabled="!projectForm.name || savingProject"
            :class="[
              'px-4 py-2 rounded-lg font-medium transition-colors',
              projectForm.name && !savingProject
                ? 'bg-blue-600 hover:bg-blue-500 text-white'
                : 'bg-gray-700 text-gray-500 cursor-not-allowed'
            ]"
          >
            {{ savingProject ? 'Creating...' : 'Create Project' }}
          </button>
        </div>
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
        <div class="flex items-center justify-end gap-3 pt-2">
          <button type="button" @click="showAppModal = false" class="px-4 py-2 text-gray-400 hover:text-white">
            Cancel
          </button>
          <button 
            type="submit"
            :disabled="!canCreateApp || savingApp"
            :class="[
              'px-4 py-2 rounded-lg font-medium transition-colors',
              canCreateApp && !savingApp
                ? 'bg-blue-600 hover:bg-blue-500 text-white'
                : 'bg-gray-700 text-gray-500 cursor-not-allowed'
            ]"
          >
            {{ savingApp ? 'Creating...' : 'Create Application' }}
          </button>
        </div>
      </form>
    </Modal>

    <!-- Create Template Modal -->
    <Modal :open="showTemplateModal" title="Add Template" size="lg" @close="showTemplateModal = false">
      <form @submit.prevent="createTemplate" class="space-y-4">
        <div>
          <label class="block text-gray-400 text-sm mb-2">Repository *</label>
          <div class="relative">
            <input 
              v-model="repoSearch" 
              type="text" 
              placeholder="Search repositories..."
              @focus="showRepoDropdown = true"
              class="w-full bg-gray-900 border border-gray-700 rounded-lg px-4 py-2.5 text-white 
                     placeholder-gray-500 focus:outline-none focus:border-gray-600" 
            />
            <div v-if="showRepoDropdown && filteredRepos.length > 0" 
              class="absolute z-10 w-full mt-1 bg-gray-900 border border-gray-700 rounded-lg max-h-48 overflow-y-auto">
              <button 
                v-for="repo in filteredRepos" 
                :key="repo.id" 
                type="button"
                @click="selectRepo(repo)"
                class="w-full px-4 py-2.5 text-left hover:bg-gray-800 transition-colors"
              >
                <p class="text-white text-sm">{{ repo.fullName }}</p>
              </button>
            </div>
          </div>
          <div v-if="selectedRepo" class="mt-2 p-3 bg-gray-900 border border-gray-700 rounded-lg flex items-center justify-between">
            <p class="text-white text-sm">{{ selectedRepo.fullName }}</p>
            <button type="button" @click="selectedRepo = null" class="text-gray-500 hover:text-white">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
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
        <div class="flex items-center justify-end gap-3 pt-2">
          <button type="button" @click="showTemplateModal = false" class="px-4 py-2 text-gray-400 hover:text-white">
            Cancel
          </button>
          <button 
            type="submit"
            :disabled="!canCreateTemplate || savingTemplate"
            :class="[
              'px-4 py-2 rounded-lg font-medium transition-colors',
              canCreateTemplate && !savingTemplate
                ? 'bg-blue-600 hover:bg-blue-500 text-white'
                : 'bg-gray-700 text-gray-500 cursor-not-allowed'
            ]"
          >
            {{ savingTemplate ? 'Adding...' : 'Add Template' }}
          </button>
        </div>
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
  XMarkIcon
} from '@heroicons/vue/24/outline';
import { appsApi, projectsApi, templatesApi, repoApi, type CreateAppRequest, type CreateProjectRequest, type CreateTemplateRequest } from '../api';
import type { App } from '../types/App';
import type { Project } from '../types/Project';
import type { Template } from '../types/Template';
import type { Build } from '../types/Build';
import type { Repo } from '../types/Repo';
import Modal from '../components/ui/Modal.vue';
import Spinner from '../components/ui/Spinner.vue';
import StatusBadge from '../components/ui/StatusBadge.vue';
import ErrorAlert from '../components/ui/ErrorAlert.vue';
import FormField from '../components/ui/FormField.vue';
import Dropdown from '../components/ui/Dropdown.vue';
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

// Activity (derived from builds)
const recentActivity = computed(() => 
  allBuilds.value
    .sort((a, b) => new Date(b.startedAt).getTime() - new Date(a.startedAt).getTime())
    .slice(0, 5)
    .map(build => ({
      id: `${build.appId}-${build.id}`,
      appId: build.appId,
      appName: build.appName,
      status: build.status,
      message: `build #${build.number} ${build.status}`,
      time: formatTime(build.startedAt),
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
const repoSearch = ref('');
const showRepoDropdown = ref(false);
const selectedRepo = ref<Repo | null>(null);
const templateForm = reactive({
  name: '',
  description: '',
  kind: '',
  language: '',
});

const filteredRepos = computed(() => {
  if (!repoSearch.value) return repos.value.slice(0, 10);
  return repos.value.filter(r =>
    r.fullName.toLowerCase().includes(repoSearch.value.toLowerCase())
  ).slice(0, 10);
});

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

const selectRepo = (repo: Repo) => {
  selectedRepo.value = repo;
  repoSearch.value = '';
  showRepoDropdown.value = false;
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
