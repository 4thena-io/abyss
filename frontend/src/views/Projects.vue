<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-white">Projects</h1>
        <p class="text-gray-400 mt-1">{{ projectList.length }} projects</p>
      </div>
      <button class="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-500 
        text-white text-sm font-medium rounded-lg transition-colors">
        <PlusIcon class="w-4 h-4" />
        New Project
      </button>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <div class="relative flex-1 max-w-xs">
        <MagnifyingGlassIcon class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input v-model="search" type="text" placeholder="Filter projects..." class="w-full bg-gray-800 border border-gray-700 rounded-lg pl-10 pr-4 py-2 
            text-sm text-white placeholder-gray-500 focus:outline-none focus:border-gray-600" />
      </div>
      <select v-model="statusFilter" class="bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-300 
          focus:outline-none focus:border-gray-600">
        <option value="">All statuses</option>
        <option value="active">Active</option>
        <option value="inactive">Inactive</option>
        <option value="archived">Archived</option>
      </select>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="w-8 h-8 border-2 border-gray-600 border-t-blue-500 rounded-full animate-spin" />
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredProjects.length === 0"
      class="bg-gray-800 border border-gray-700 rounded-lg py-16 text-center">
      <FolderIcon class="w-12 h-12 text-gray-600 mx-auto" />
      <h3 class="mt-4 text-lg font-medium text-white">No projects found</h3>
      <p class="mt-2 text-gray-400 text-sm">
        {{ search || statusFilter ? 'Try adjusting your filters' : 'Create your first project to get started' }}
      </p>
    </div>

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
              <p class="text-gray-500 text-xs mt-0.5">
                Created {{ formatDate(project.created_at) }}
              </p>
            </div>
          </div>
          <span :class="[
            'px-2 py-1 rounded text-xs font-medium',
            project.status === 'active' ? 'bg-green-500/20 text-green-400' :
              project.status === 'inactive' ? 'bg-gray-500/20 text-gray-400' :
                project.status === 'archived' ? 'bg-yellow-500/20 text-yellow-400' :
                  'bg-gray-500/20 text-gray-400'
          ]">
            {{ project.status }}
          </span>
        </div>

        <p v-if="project.description" class="mt-4 text-gray-400 text-sm line-clamp-2">
          {{ project.description }}
        </p>
        <p v-else class="mt-4 text-gray-500 text-sm italic">
          No description
        </p>

        <div class="mt-4 pt-4 border-t border-gray-700 flex items-center justify-between">
          <div class="flex items-center gap-2 text-gray-500 text-sm">
            <CubeIcon class="w-4 h-4" />
            <span>{{ getAppCount(project.id) }} apps</span>
          </div>
          <span class="text-gray-500 text-xs">
            Updated {{ formatDate(project.updated_at) }}
          </span>
        </div>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import {
  PlusIcon,
  MagnifyingGlassIcon,
  FolderIcon,
  CubeIcon
} from '@heroicons/vue/24/outline';
import type { Project } from '../types/Project';

const projectList = ref<Project[]>([]);
const loading = ref(true);
const search = ref('');
const statusFilter = ref('');

// TODO: Replace with actual app counts from API
const appCounts = ref<Record<string, number>>({});

const filteredProjects = computed(() => {
  return projectList.value.filter(project => {
    const matchesSearch = !search.value ||
      project.name.toLowerCase().includes(search.value.toLowerCase()) ||
      project.description.toLowerCase().includes(search.value.toLowerCase());
    const matchesStatus = !statusFilter.value || project.status === statusFilter.value;
    return matchesSearch && matchesStatus;
  });
});

const getAppCount = (projectId: string): number => {
  return appCounts.value[projectId] || 0;
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (days === 0) return 'today';
  if (days === 1) return 'yesterday';
  if (days < 7) return `${days} days ago`;
  return date.toLocaleDateString();
};

const fetchProjects = async () => {
  try {
    loading.value = true;
    const response = await fetch('/api/projects');
    if (!response.ok) throw new Error(`HTTP error: ${response.status}`);
    projectList.value = await response.json();
  } catch (error) {
    console.error('Failed to fetch projects:', error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchProjects);
</script>
