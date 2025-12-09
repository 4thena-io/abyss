<template>
  <aside class="w-12 bg-gray-950 flex flex-col items-center py-4 space-y-8">
    <nav class="flex-1 flex-col flex space-y-6 items-center">
      <router-link to="/" class="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-gray-800">
        <HomeIcon class="w-5 h-5" />
      </router-link>
      <router-link to="/projects" class="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-gray-800">
        <FolderIcon class="w-5 h-5" />
      </router-link>
      <router-link to="/apps" class="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-gray-800">
        <CubeIcon class="w-5 h-5" />
      </router-link>
      <router-link to="/templates" class="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-gray-800">
        <DocumentDuplicateIcon class="w-5 h-5" />
      </router-link>
      <router-link to="/learn" class="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-gray-800">
        <AcademicCapIcon class="w-5 h-5" />
      </router-link>

    </nav>
  </aside>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';
import { AcademicCapIcon, CubeIcon, DocumentDuplicateIcon, FolderIcon, HomeIcon } from '@heroicons/vue/24/outline';
</script>
