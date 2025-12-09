<template>
  <div class="app flex flex-col h-screen overflow-hidden">
    <Header class="shrink-0" />
    <div class="flex flex-1 min-h-0">
      <Sidebar class="shrink-0" />
      <main class="flex-1 overflow-y-auto p-8">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { RouterView } from 'vue-router';
import Header from './components/Header.vue';
import Sidebar from './components/Sidebar.vue';
</script>
