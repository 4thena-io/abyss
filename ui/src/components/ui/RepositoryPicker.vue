<template>
  <div>
    <label v-if="label" class="block text-t3 text-sm mb-2">
      {{ label }}{{ required ? ' *' : '' }}
    </label>
    <div class="relative">
      <input
        v-model="search"
        type="text"
        :placeholder="placeholder"
        @focus="showDropdown = true"
        class="w-full bg-bg border border-b1 rounded-lg px-4 py-2.5 text-t1 placeholder:text-t3 focus:outline-none focus:border-b2"
      />
      <div
        v-if="showDropdown && filteredRepos.length > 0"
        class="absolute z-10 w-full mt-1 bg-bg border border-b1 rounded-lg max-h-48 overflow-y-auto"
      >
        <button
          v-for="repo in filteredRepos"
          :key="repo.id"
          type="button"
          @click="selectRepo(repo)"
          class="w-full px-4 py-2.5 text-left hover:bg-panel transition-colors"
        >
          <p class="text-t1 text-sm">{{ repo.fullName }}</p>
        </button>
      </div>
    </div>
    <div
      v-if="modelValue"
      class="mt-2 p-3 bg-bg border border-b1 rounded-lg flex items-center justify-between"
    >
      <p class="text-t1 text-sm">{{ modelValue.fullName }}</p>
      <button type="button" @click="clearSelection" class="text-t3 hover:text-t1 transition-colors">
        <XMarkIcon class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { XMarkIcon } from '@heroicons/vue/24/outline';

export interface Repository {
  id: number | string;
  fullName: string;
  url: string;
}

const props = withDefaults(
  defineProps<{
    modelValue: Repository | null;
    repos: Repository[];
    label?: string;
    placeholder?: string;
    required?: boolean;
    maxResults?: number;
  }>(),
  {
    placeholder: 'Search repositories...',
    maxResults: 10,
  },
);

const emit = defineEmits<{
  'update:modelValue': [repo: Repository | null];
  select: [repo: Repository];
}>();

const search = ref('');
const showDropdown = ref(false);

const filteredRepos = computed(() => {
  if (!search.value) return props.repos.slice(0, props.maxResults);
  return props.repos
    .filter((r) => r.fullName.toLowerCase().includes(search.value.toLowerCase()))
    .slice(0, props.maxResults);
});

const selectRepo = (repo: Repository) => {
  emit('update:modelValue', repo);
  emit('select', repo);
  search.value = '';
  showDropdown.value = false;
};

const clearSelection = () => {
  emit('update:modelValue', null);
};
</script>
