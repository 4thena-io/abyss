<template>
  <div v-if="loading" class="flex items-center justify-center py-12">
    <Spinner size="lg" />
  </div>

  <EmptyState v-else-if="!template" :icon="DocumentDuplicateIcon" title="Template not found">
    <router-link to="/templates" class="text-acc hover:opacity-80">Back to templates</router-link>
  </EmptyState>

  <div v-else class="space-y-6">
    <DetailHeader
      :icon="DocumentDuplicateIcon"
      :title="template.name"
      :subtitle="
        [template.description, template.kind, template.language].filter(Boolean).join(' · ')
      "
    >
    </DetailHeader>

    <Tabs v-model="activeTab" :tabs="tabs" />

    <div class="mt-6">
      <TemplateOverviewTab
        v-if="activeTab === 'overview'"
        :template="template"
        @show-apps="activeTab = 'apps'"
      />
      <TemplateAppsTab v-else-if="activeTab === 'apps'" :template-id="template.id" />
      <TemplateSettingsTab
        v-else-if="activeTab === 'settings'"
        :template="template"
        @updated="template = $event"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { DocumentDuplicateIcon } from '@heroicons/vue/24/outline';
import type { Template } from '../features/templates/types';
import { templatesApi } from '../api/template';
import Spinner from '../components/ui/Spinner.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import Tabs from '../components/ui/Tabs.vue';
import DetailHeader from '../components/ui/DetailHeader.vue';
import TemplateOverviewTab from '../features/templates/components/TemplateOverviewTab.vue';
import TemplateAppsTab from '../features/templates/components/TemplateAppsTab.vue';
import TemplateSettingsTab from '../features/templates/components/TemplateSettingsTab.vue';

const route = useRoute();
const router = useRouter();

const tabs = [
  { id: 'overview', label: 'Overview' },
  { id: 'apps', label: 'Applications' },
  { id: 'settings', label: 'Settings' },
];

const validTabs = tabs.map((t) => t.id);
const activeTab = ref(
  validTabs.includes(route.query.tab as string) ? (route.query.tab as string) : 'overview',
);

watch(activeTab, (tab) => router.replace({ query: { tab } }));

const template = ref<Template | null>(null);
const loading = ref(true);

onMounted(async () => {
  const id = Number(route.params.id);
  try {
    template.value = await templatesApi.getById(id);
  } catch {
    template.value = null;
  } finally {
    loading.value = false;
  }
});
</script>
