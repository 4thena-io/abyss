<template>
  <div class="space-y-6">
    <PageHeader title="Projects" :subtitle="`${projectList.length} projects`">
      <template #actions>
        <Button @click="showModal = true" variant="secondary">
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
      <router-link v-for="project in filteredProjects" :key="project.id" :to="`/projects/${project.id}`" class="bg-gray-800 border border-gray-700 rounded-lg p-5 
          hover:border-gray-600 transition-all group">
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-gray-700 flex items-center justify-center">
              <FolderIcon class="w-5 h-5 text-gray-400" />
            </div>
            <div>
              <h3 class="text-white font-medium group-hover:text-blue-400 transition-colors">
                {{ project.name }}
              </h3>
            </div>
          </div>
        </div>

        <p v-if="project.description" class="mt-4 text-gray-400 text-sm line-clamp-2">
          {{ project.description }}
        </p>
        <p v-else class="mt-4 text-gray-500 text-sm italic">
          No description
        </p>
      </router-link>
    </div>

    <!-- Create Project Modal -->
    <Modal :open="showModal" title="New Project" @close="closeModal">
      <form @submit.prevent="createProject" class="space-y-4">
        <FormField 
          v-model="form.name" 
          label="Name" 
          required 
          placeholder="my-project" 
        />

        <FormField 
          v-model="form.description" 
          label="Description" 
          type="textarea" 
          :rows="3" 
          placeholder="Project description..." 
        />

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
            {{ saving ? 'Creating...' : 'Create Project' }}
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
  FolderIcon
} from '@heroicons/vue/24/outline';
import type { Project } from '../types/Project';
import { projectsApi, type CreateProjectRequest } from '../api';
import Modal from '../components/ui/Modal.vue';
import Button from '../components/ui/Button.vue';
import Spinner from '../components/ui/Spinner.vue';
import SearchInput from '../components/ui/SearchInput.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import ErrorAlert from '../components/ui/ErrorAlert.vue';
import FormField from '../components/ui/FormField.vue';
import PageHeader from '../components/ui/PageHeader.vue';

// List state
const projectList = ref<Project[]>([]);
const loading = ref(true);
const search = ref('');

// Modal state
const showModal = ref(false);
const saving = ref(false);
const error = ref('');

// Form state
const form = reactive<CreateProjectRequest>({
  name: '',
  description: '',
});

// Computed
const filteredProjects = computed(() => {
  return projectList.value.filter(project => {
    return !search.value ||
      project.name.toLowerCase().includes(search.value.toLowerCase()) ||
      project.description?.toLowerCase().includes(search.value.toLowerCase());
  });
});

const canCreate = computed(() => form.name.trim() !== '');

// Methods
const resetForm = () => {
  form.name = '';
  form.description = '';
  error.value = '';
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
  } catch (e: any) {
    error.value = e.message;
  } finally {
    saving.value = false;
  }
};

// Lifecycle
onMounted(fetchProjects);
</script>
