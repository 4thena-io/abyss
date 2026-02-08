<template>
  <div class="bg-gray-800 rounded-lg border border-gray-700">
    <div class="px-5 py-4 border-b border-gray-700">
      <h2 class="text-lg font-medium text-white">{{ title }}</h2>
    </div>
    <div class="divide-y divide-gray-700">
      <div v-if="loading" class="px-5 py-8 flex justify-center">
        <Spinner size="md" />
      </div>
      <template v-else-if="items.length > 0">
        <div 
          v-for="item in items" 
          :key="item.id"
          class="px-5 py-4 flex items-start gap-4 hover:bg-gray-750 transition-colors"
        >
          <div :class="[
            'w-2 h-2 rounded-full mt-2 shrink-0',
            statusColors[item.status] || 'bg-gray-500'
          ]" />
          <div class="flex-1 min-w-0">
            <p class="text-white text-sm">
              <router-link 
                v-if="item.to" 
                :to="item.to" 
                class="font-medium hover:text-blue-400"
              >
                {{ item.title }}
              </router-link>
              <span v-else class="font-medium">{{ item.title }}</span>
              <span class="text-gray-400"> {{ item.message }}</span>
            </p>
            <p class="text-gray-500 text-xs mt-1">{{ item.time }}</p>
          </div>
        </div>
      </template>
      <div v-else class="px-5 py-8 text-center text-gray-500">
        {{ emptyMessage }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Spinner from './Spinner.vue';

export interface ActivityItem {
  id: string | number;
  title: string;
  message: string;
  time: string;
  status: string;
  to?: string;
}

withDefaults(defineProps<{
  title: string;
  items: ActivityItem[];
  loading?: boolean;
  emptyMessage?: string;
}>(), {
  emptyMessage: 'No recent activity'
});

const statusColors: Record<string, string> = {
  success: 'bg-green-500',
  running: 'bg-blue-500',
  failure: 'bg-red-500',
  pending: 'bg-gray-500',
};
</script>
