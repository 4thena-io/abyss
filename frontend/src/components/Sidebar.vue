<template>
  <aside 
    :class="[
      'bg-gray-950 flex flex-col items-center py-4 transition-all duration-300 ease-in-out',
      expanded ? 'w-48' : 'w-12'
    ]"
  >
    <nav class="flex-1 flex-col flex space-y-2 w-full px-2">
      <router-link 
        v-for="item in navItems" 
        :key="item.to"
        :to="item.to" 
        v-slot="{ isExactActive }"
        custom
      >
        <button
          @click="$router.push(item.to)"
          :class="[
            'flex items-center gap-3 w-full transition-colors p-2 rounded-lg',
            expanded ? 'justify-start' : 'justify-center',
            isExactActive 
              ? 'text-white bg-gray-800' 
              : 'text-gray-400 hover:text-white hover:bg-gray-800'
          ]"
        >
          <component :is="item.icon" class="w-5 h-5 shrink-0" />
          <span 
            v-if="expanded" 
            class="text-sm font-medium whitespace-nowrap"
          >
            {{ item.label }}
          </span>
        </button>
      </router-link>
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { 
  AcademicCapIcon, 
  CubeIcon, 
  DocumentDuplicateIcon, 
  FolderIcon, 
  HomeIcon 
} from '@heroicons/vue/24/outline';

defineProps<{
  expanded: boolean
}>();

const navItems = [
  { to: '/', label: 'Home', icon: HomeIcon },
  { to: '/projects', label: 'Projects', icon: FolderIcon },
  { to: '/apps', label: 'Applications', icon: CubeIcon },
  { to: '/templates', label: 'Templates', icon: DocumentDuplicateIcon },
  { to: '/learn', label: 'Learn', icon: AcademicCapIcon },
];
</script>
