<template>
  <header class="bg-gray-950 h-12 flex items-center justify-between">
    <div class="flex w-12 items-center justify-center shrink-0">
      <button @click="$emit('toggle-sidebar')"
        class="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-gray-800"
        aria-label="Toggle navigation menu">
        <Bars3Icon class="w-5 h-5" />
      </button>
    </div>
    <div class="flex-1 flex items-center justify-start min-w-0">
      <img :src="logo" alt="logo" class="w-auto h-12 cursor-pointer hover:opacity-80 transition-opacity" />
      <div class="flex-1 flex justify-center px-8 min-w-0">
        <SearchInput v-model="search" placeholder="Search for projects and applications..." max-width="max-w-2xl" />
      </div>
    </div>
    <div class="flex items-center justify-end w-24 space-x-2 pr-4">
      <button class="text-gray-400 hover:text-white transition-colors relative p-2 rounded-lg hover:bg-gray-800"
        aria-label="Notifications">
        <BellIcon class="w-5 h-5" />
      </button>

      <!-- User menu -->
      <div class="relative" ref="menuRef">
        <button @click="menuOpen = !menuOpen"
          class="flex items-center justify-center w-8 h-8 rounded-full overflow-hidden ring-2 ring-transparent hover:ring-gray-600 transition-all"
          aria-label="User menu">
          <img v-if="user?.avatar_url" :src="user.avatar_url" :alt="user.username" class="w-full h-full object-cover" />
          <div v-else class="w-full h-full bg-gray-700 flex items-center justify-center">
            <UserCircleIcon class="w-5 h-5 text-gray-400" />
          </div>
        </button>

        <div v-if="menuOpen"
          class="absolute right-0 mt-2 w-48 bg-gray-800 border border-gray-700 rounded-lg shadow-lg overflow-hidden z-50">
          <div class="px-4 py-3 border-b border-gray-700">
            <p class="text-sm font-medium text-white truncate">{{ user?.username }}</p>
          </div>
          <div class="border-b border-gray-700">
            <button
              class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-gray-300 hover:bg-gray-700 hover:text-white transition-colors">
              <UserIcon class="w-4 h4" />
              Profile
            </button>
            <button
              class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-gray-300 hover:bg-gray-700 hover:text-white transition-colors">
              <Cog6ToothIcon class="w-4 h4" />
              Settings
            </button>
          </div>
          <button @click="logout"
            class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-gray-300 hover:bg-gray-700 hover:text-white transition-colors">
            <ArrowRightStartOnRectangleIcon class="w-4 h-4" />
            Sign out
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Bars3Icon, BellIcon } from '@heroicons/vue/24/outline';
import SearchInput from '../ui/SearchInput.vue';
import { UserCircleIcon, ArrowRightStartOnRectangleIcon, UserIcon, Cog6ToothIcon } from '@heroicons/vue/24/outline';
import logo from '../../assets/logo-white.svg';
import { useClickOutside } from '../../composables/useClickOutside';
import { useCurrentUser } from '../../features/auth/composables/useCurrentUser';

defineEmits<{
  'toggle-sidebar': []
}>();

const { user, fetch, clear } = useCurrentUser();
const search = ref('');
const menuOpen = ref(false);
const menuRef = ref<HTMLElement | null>(null);

useClickOutside(menuRef, () => { menuOpen.value = false; });

onMounted(fetch);

async function logout() {
  clear();
  await window.fetch('/auth/logout', { method: 'POST' });
  window.location.href = '/login';
}
</script>
