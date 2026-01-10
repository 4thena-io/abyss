<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl font-semibold text-white">Dashboard</h1>
      <p class="text-gray-400 mt-1">Overview of your development platform</p>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="stat in stats" :key="stat.label" class="bg-gray-800 rounded-lg p-5 border border-gray-700">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-400 text-sm">{{ stat.label }}</p>
            <p class="text-2xl font-semibold text-white mt-1">{{ stat.value }}</p>
          </div>
          <component :is="stat.icon" class="w-8 h-8 text-gray-600" />
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Activity -->
      <div class="lg:col-span-2 bg-gray-800 rounded-lg border border-gray-700">
        <div class="px-5 py-4 border-b border-gray-700">
          <h2 class="text-lg font-medium text-white">Recent Activity</h2>
        </div>
        <div class="divide-y divide-gray-700">
          <div v-for="activity in recentActivity" :key="activity.id"
            class="px-5 py-4 flex items-start gap-4 hover:bg-gray-750 transition-colors">
            <div :class="[
              'w-2 h-2 rounded-full mt-2 shrink-0',
              activity.type === 'deploy' ? 'bg-green-500' :
                activity.type === 'build' ? 'bg-blue-500' :
                  activity.type === 'error' ? 'bg-red-500' : 'bg-gray-500'
            ]" />
            <div class="flex-1 min-w-0">
              <p class="text-white text-sm">
                <span class="font-medium">{{ activity.app }}</span>
                <span class="text-gray-400"> {{ activity.message }}</span>
              </p>
              <p class="text-gray-500 text-xs mt-1">{{ activity.time }}</p>
            </div>
          </div>
          <div v-if="recentActivity.length === 0" class="px-5 py-8 text-center text-gray-500">
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
          <button v-for="action in quickActions" :key="action.label" class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left
              bg-gray-700/50 hover:bg-gray-700 transition-colors group">
            <component :is="action.icon" class="w-5 h-5 text-gray-400 group-hover:text-white" />
            <span class="text-gray-300 group-hover:text-white text-sm">{{ action.label }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Build Status -->
    <div class="bg-gray-800 rounded-lg border border-gray-700">
      <div class="px-5 py-4 border-b border-gray-700 flex items-center justify-between">
        <h2 class="text-lg font-medium text-white">Recent Builds</h2>
        <router-link to="/apps" class="text-sm text-blue-400 hover:text-blue-300">View all</router-link>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="text-left text-gray-400 text-sm border-b border-gray-700">
              <th class="px-5 py-3 font-medium">Application</th>
              <th class="px-5 py-3 font-medium">Status</th>
              <th class="px-5 py-3 font-medium">Branch</th>
              <th class="px-5 py-3 font-medium">Duration</th>
              <th class="px-5 py-3 font-medium">Time</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-700">
            <tr v-for="build in recentBuilds" :key="build.id" class="hover:bg-gray-750 transition-colors">
              <td class="px-5 py-4">
                <span class="text-white text-sm">{{ build.app }}</span>
              </td>
              <td class="px-5 py-4">
                <span :class="[
                  'inline-flex items-center px-2 py-1 rounded text-xs font-medium',
                  build.status === 'success' ? 'bg-green-500/20 text-green-400' :
                    build.status === 'running' ? 'bg-blue-500/20 text-blue-400' :
                      build.status === 'failed' ? 'bg-red-500/20 text-red-400' :
                        'bg-gray-500/20 text-gray-400'
                ]">
                  {{ build.status }}
                </span>
              </td>
              <td class="px-5 py-4 text-gray-400 text-sm">{{ build.branch }}</td>
              <td class="px-5 py-4 text-gray-400 text-sm">{{ build.duration }}</td>
              <td class="px-5 py-4 text-gray-500 text-sm">{{ build.time }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import {
  FolderIcon,
  CubeIcon,
  WrenchScrewdriverIcon,
  CheckCircleIcon,
  PlusIcon,
  CloudArrowUpIcon,
  DocumentPlusIcon,
  Cog6ToothIcon
} from '@heroicons/vue/24/outline';

const stats = ref([
  { label: 'Projects', value: '0', icon: FolderIcon },
  { label: 'Applications', value: '0', icon: CubeIcon },
  { label: 'Builds Today', value: '0', icon: WrenchScrewdriverIcon },
  { label: 'Deployments', value: '0', icon: CheckCircleIcon },
]);

const recentActivity = ref<{ id: number; app: string; message: string; type: string; time: string }[]>([]);

const quickActions = ref([
  { label: 'Create new project', icon: PlusIcon },
  { label: 'Import from forge', icon: CloudArrowUpIcon },
  { label: 'Create from template', icon: DocumentPlusIcon },
  { label: 'Configure integrations', icon: Cog6ToothIcon },
]);

const recentBuilds = ref<{ id: number; app: string; status: string; branch: string; duration: string; time: string }[]>([]);
</script>
