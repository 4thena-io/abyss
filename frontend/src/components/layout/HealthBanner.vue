<template>
  <div v-if="alerts.length" class="flex flex-col gap-px">
    <div
      v-for="alert in alerts"
      :key="alert.key"
      class="flex items-center justify-between gap-4 px-4 py-2 bg-warn-s border-b border-warn/30 text-warn text-sm"
    >
      <div class="flex items-center gap-2">
        <ExclamationTriangleIcon class="w-4 h-4 shrink-0" />
        <span>{{ alert.message }}</span>
      </div>
      <button
        @click="dismiss(alert.key)"
        class="shrink-0 text-warn/60 hover:text-warn transition-colors"
        aria-label="Dismiss"
      >
        <XMarkIcon class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { ExclamationTriangleIcon, XMarkIcon } from '@heroicons/vue/24/outline';
import { useHealth } from '../../composables/useHealth';

const { status } = useHealth();
const dismissed = ref<Set<string>>(new Set());

const alerts = computed(() => {
  if (!status.value) return [];
  const list: { key: string; message: string }[] = [];
  if (!status.value.forge.ok)
    list.push({
      key: 'forge',
      message: `Forge is unreachable — app creation and webhooks are unavailable.${status.value.forge.error ? ' (' + status.value.forge.error + ')' : ''}`,
    });
  if (!status.value.ci.ok)
    list.push({
      key: 'ci',
      message: `CI is unreachable — build and deployment data may be unavailable.${status.value.ci.error ? ' (' + status.value.ci.error + ')' : ''}`,
    });
  return list.filter((a) => !dismissed.value.has(a.key));
});

// Auto-restore dismissed alerts when the service comes back up, so the next
// outage shows the banner again.
watch(status, (next) => {
  if (next?.forge.ok) dismissed.value.delete('forge');
  if (next?.ci.ok) dismissed.value.delete('ci');
});

function dismiss(key: string) {
  dismissed.value = new Set([...dismissed.value, key]);
}
</script>
