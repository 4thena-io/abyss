<template>
  <component
    :is="to ? 'router-link' : 'div'"
    :to="to"
    :class="[
      'bg-panel border border-b1 rounded-lg p-5 block',
      hoverable && 'hover:border-b2 transition-all',
      to && 'group',
    ]"
  >
    <div class="flex items-start justify-between">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-lg bg-surface flex items-center justify-center">
          <component :is="icon" class="w-5 h-5 text-t3" />
        </div>
        <div>
          <h3 :class="['text-t1 font-medium', to && 'group-hover:text-acc transition-colors']">
            {{ title }}
          </h3>
          <p v-if="subtitle" class="text-t3 text-xs mt-0.5">{{ subtitle }}</p>
        </div>
      </div>
      <slot name="header-action" />
    </div>

    <p v-if="description" class="mt-4 text-t2 text-sm line-clamp-2">
      {{ description }}
    </p>
    <p v-else-if="showNoDescription" class="mt-4 text-t3 text-sm italic">
      {{ noDescriptionText }}
    </p>

    <div v-if="$slots.badges" class="mt-4 flex items-center gap-2 flex-wrap">
      <slot name="badges" />
    </div>

    <div v-if="$slots.footer" class="mt-4 pt-4 border-t border-b1">
      <slot name="footer" />
    </div>
  </component>
</template>

<script setup lang="ts">
import type { Component } from 'vue';

withDefaults(
  defineProps<{
    icon: Component;
    title: string;
    subtitle?: string;
    description?: string;
    noDescriptionText?: string;
    showNoDescription?: boolean;
    to?: string;
    hoverable?: boolean;
  }>(),
  {
    noDescriptionText: 'No description',
    showNoDescription: true,
    hoverable: true,
  },
);
</script>
