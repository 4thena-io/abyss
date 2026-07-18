<template>
  <div class="space-y-6">
    <PageHeader title="Projects" :subtitle="`${projectList.length} projects`">
      <template #actions>
        <Button @click="openModal" variant="secondary">
          <PlusIcon class="w-4 h-4" />
          New Project
        </Button>
      </template>
    </PageHeader>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <SearchInput v-model="search" placeholder="Filter projects..." />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <!-- Empty State -->
    <EmptyState
      v-else-if="filteredProjects.length === 0"
      :icon="FolderIcon"
      title="No projects found"
      :message="search ? 'Try adjusting your search' : 'Create your first project to get started'"
    />

    <!-- Project Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <EntityCard
        v-for="project in filteredProjects"
        :key="project.id"
        :icon="FolderIcon"
        :title="project.name"
        :description="project.description"
        :to="`/projects/${project.id}`"
      />
    </div>

    <!-- Create Project Modal -->
    <Modal :open="showModal" title="New Project" @close="closeModal">
      <form @submit.prevent="createProject" class="space-y-4">
        <FormField v-model="form.name" label="Name" required placeholder="my-project" />

        <FormField
          v-model="form.description"
          label="Description"
          type="textarea"
          :rows="3"
          placeholder="Project description..."
        />

        <Dropdown
          v-model="form.teamId"
          label="Team"
          variant="form"
          placeholder="No team (standalone)"
          :options="teamOptions"
        />

        <ErrorAlert v-if="error" :message="error" />

        <FormActions
          submit-label="Create Project"
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
import { PlusIcon, FolderIcon } from '@heroicons/vue/24/outline';
import type { Project } from '../features/projects/types';
import type { Team } from '../features/teams/types';
import { projectsApi, teamsApi, type CreateProjectRequest } from '../api';
import Modal from '../components/ui/Modal.vue';
import Button from '../components/ui/Button.vue';
import Spinner from '../components/ui/Spinner.vue';
import SearchInput from '../components/ui/SearchInput.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import ErrorAlert from '../components/ui/ErrorAlert.vue';
import FormField from '../components/ui/FormField.vue';
import Dropdown from '../components/ui/Dropdown.vue';
import PageHeader from '../components/ui/PageHeader.vue';
import FormActions from '../components/ui/FormActions.vue';
import EntityCard from '../components/ui/EntityCard.vue';

// List state
const projectList = ref<Project[]>([]);
const loading = ref(true);
const search = ref('');

// Modal state
const showModal = ref(false);
const saving = ref(false);
const error = ref('');
const teams = ref<Team[]>([]);

// Form state
const form = reactive<CreateProjectRequest>({
  name: '',
  description: '',
  teamId: null,
});

// Computed
const filteredProjects = computed(() => {
  return projectList.value.filter((project) => {
    return (
      !search.value ||
      project.name.toLowerCase().includes(search.value.toLowerCase()) ||
      project.description?.toLowerCase().includes(search.value.toLowerCase())
    );
  });
});

const canCreate = computed(() => form.name.trim() !== '');

// Methods
const teamOptions = computed(() => teams.value.map((t) => ({ value: t.id, label: t.name })));

const resetForm = () => {
  form.name = '';
  form.description = '';
  form.teamId = null;
  error.value = '';
};

const openModal = async () => {
  showModal.value = true;
  try {
    teams.value = await teamsApi.getAll();
  } catch (e) {
    console.error('Failed to fetch teams:', e);
  }
};

const closeModal = () => {
  showModal.value = false;
  resetForm();
};

const fetchProjects = async () => {
  try {
    loading.value = true;
    projectList.value = await projectsApi.getAll();
  } catch (e) {
    console.error('Failed to fetch projects:', e);
  } finally {
    loading.value = false;
  }
};

const createProject = async () => {
  error.value = '';
  saving.value = true;

  try {
    const project = await projectsApi.create(form);
    projectList.value.push(project);
    closeModal();
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'An error occurred';
  } finally {
    saving.value = false;
  }
};

// Lifecycle
onMounted(fetchProjects);
</script>
