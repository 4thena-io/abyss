<template>
  <div class="bg-gray-800 border border-gray-700 rounded-lg p-6">
    <h2 v-if="title" class="text-lg font-medium text-white mb-4">{{ title }}</h2>
    <dl class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="item in items" :key="item.label">
        <dt class="text-gray-400 text-sm">{{ item.label }}</dt>
        <dd class="text-white mt-1">
          <slot :name="item.key || item.label" :item="item">
            <component 
              v-if="item.to" 
              :is="item.external ? 'a' : 'router-link'" 
              :to="!item.external ? item.to : undefined"
              :href="item.external ? item.to : undefined"
              :target="item.external ? '_blank' : undefined"
              :rel="item.external ? 'noopener noreferrer' : undefined"
              class="text-blue-400 hover:text-blue-300 break-all"
            >
              {{ item.value || 'Not configured' }}
            </component>
            <span v-else :class="item.value ? '' : 'text-gray-500'">
              {{ item.value || item.fallback || 'Not configured' }}
            </span>
          </slot>
        </dd>
      </div>
    </dl>
  </div>
</template>

<script setup lang="ts">
export interface DefinitionItem {
  key?: string;
  label: string;
  value?: string | number | null;
  fallback?: string;
  to?: string;
  external?: boolean;
}

defineProps<{
  title?: string;
  items: DefinitionItem[];
}>();
</script>
