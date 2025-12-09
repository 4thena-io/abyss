<template>
  <div>
    <h1>Applications</h1>
    <ul>
      <li v-for="app in appList" :key="app.id">
        <strong>{{ app.name }}</strong> ({{ app.language }}) - Project: {{ app.project }}
        <p>Repo: <a :href="app.repo_url">{{ app.repo_url }}</a></p>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { App } from '../types/App';

const appList = ref<App[]>([]);

const fetchApps = async () => {
  try {
    const response = await fetch('/api/apps');

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    appList.value = data as App[];

  } catch (error) {
    console.error("Failed to fetch apps:", error);
  }
};

onMounted(fetchApps);
</script>
