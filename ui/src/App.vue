<template>
  <!-- Standalone pages (setup, login) render without the app shell -->
  <RouterView v-if="isStandalone" />

  <div v-else class="app flex flex-col h-screen overflow-hidden">
    <Header class="shrink-0" @toggle-sidebar="sidebarExpanded = !sidebarExpanded" />
    <HealthBanner class="shrink-0" />
    <div class="flex flex-1 min-h-0">
      <Sidebar class="shrink-0" :expanded="sidebarExpanded" />
      <main class="flex-1 overflow-y-auto p-8">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useRoute, RouterView } from 'vue-router';
import Header from './components/layout/Header.vue';
import Sidebar from './components/layout/Sidebar.vue';
import HealthBanner from './components/layout/HealthBanner.vue';
import { useTheme } from './composables/useTheme';

useTheme();

const route = useRoute();
const sidebarExpanded = ref(false);

const standaloneRoutes = ['/setup', '/login'];
const isStandalone = computed(() => standaloneRoutes.includes(route.path));
</script>
