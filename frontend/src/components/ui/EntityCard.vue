<template>
  <component 
    :is="to ? 'router-link' : 'div'" 
    :to="to"
    :class="[
      'bg-gray-800 border border-gray-700 rounded-lg p-5 block',
      hoverable && 'hover:border-gray-600 transition-all',
      to && 'group'
    ]"
  >
    <div class="flex items-start justify-between">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-lg bg-gray-700 flex items-center justify-center">
          <component :is="icon" class="w-5 h-5 text-gray-400" />
        </div>
        <div>
          <h3 :class="[
            'text-white font-medium',
            to && 'group-hover:text-blue-400 transition-colors'
          ]">
            {{ title }}
          </h3>
          <p v-if="subtitle" class="text-gray-500 text-xs mt-0.5">{{ subtitle }}</p>
        </div>
      </div>
      <slot name="header-action" />
    </div>

    <p v-if="description" class="mt-4 text-gray-400 text-sm line-clamp-2">
      {{ description }}
    </p>
    <p v-else-if="showNoDescription" class="mt-4 text-gray-500 text-sm italic">
      {{ noDescriptionText }}
    </p>

    <div v-if="$slots.badges" class="mt-4 flex items-center gap-2 flex-wrap">
      <slot name="badges" />
    </div>

    <div v-if="$slots.footer" class="mt-4 pt-4 border-t border-gray-700">
      <slot name="footer" />
    </div>
  </component>
</template>

<script setup lang="ts">
import type { Component } from 'vue';

withDefaults(defineProps<{
  icon: Component;
  title: string;
  subtitle?: string;
  description?: string;
  noDescriptionText?: string;
  showNoDescription?: boolean;
  to?: string;
  hoverable?: boolean;
}>(), {
  noDescriptionText: 'No description',
  showNoDescription: true,
  hoverable: true,
});
</script>
