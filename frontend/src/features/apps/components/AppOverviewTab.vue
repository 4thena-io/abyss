<template>
  <div class="space-y-6">
    <DefinitionList title="Details" :items="detailItems" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { App } from '../types';
import { useFormatters } from '../../../composables/useFormatters';
import DefinitionList from '../../../components/ui/DefinitionList.vue';

const props = defineProps<{ app: App }>();
const { formatUrl } = useFormatters();

const detailItems = computed(() => [
  { label: 'Description', value: props.app.description, fallback: 'No description' },
  { label: 'Project', value: `Project #${props.app.projectId}`, to: `/projects/${props.app.projectId}` },
  { label: 'Repository URL', value: props.app.repoUrl, to: formatUrl(props.app.repoUrl || ''), external: true },
  { label: 'CI URL', value: props.app.ciUrl, to: props.app.ciUrl ? formatUrl(props.app.ciUrl) : undefined, external: true },
  { label: 'Template', value: `Template #${props.app.templateId}`, to: `/templates/${props.app.templateId}` },
]);
</script>
