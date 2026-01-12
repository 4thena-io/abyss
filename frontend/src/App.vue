<template>
  <div class="app flex flex-col h-screen overflow-hidden">
    <Header class="shrink-0" @toggle-sidebar="sidebarExpanded = !sidebarExpanded" />
    <div class="flex flex-1 min-h-0">
      <Sidebar class="shrink-0" :expanded="sidebarExpanded" />
      <main class="flex-1 overflow-y-auto p-8">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { RouterView } from 'vue-router';
import Header from './components/layout/Header.vue';
import Sidebar from './components/layout/Sidebar.vue';

const sidebarExpanded = ref(false);
</script>
