<template>
  <header class="bg-panel h-13 flex items-center justify-between border-b border-b1">
    <div class="flex w-12 items-center justify-center shrink-0">
      <button @click="$emit('toggle-sidebar')"
        class="text-t3 hover:text-t1 transition-colors p-2 rounded-lg hover:bg-surface"
        aria-label="Toggle navigation menu">
        <Bars3Icon class="w-5 h-5" />
      </button>
    </div>
    <div class="flex-1 flex items-center justify-start min-w-0">
      <img :src="logo" alt="logo" class="w-auto h-8 cursor-pointer hover:opacity-80 transition-opacity ml-1" />
      <div class="flex-1 flex justify-center px-8 min-w-0">
        <SearchInput v-model="search" placeholder="Search for projects and applications..." max-width="max-w-2xl" />
      </div>
    </div>
    <div class="flex items-center justify-end gap-1 pr-4">
      <button class="text-t3 hover:text-t1 transition-colors relative p-2 rounded-lg hover:bg-surface"
        aria-label="Notifications">
        <BellIcon class="w-5 h-5" />
      </button>

      <!-- User menu -->
      <div class="relative" ref="menuRef">
        <button @click="menuOpen = !menuOpen"
          class="flex items-center justify-center w-8 h-8 rounded-full overflow-hidden ring-2 ring-transparent hover:ring-b2 transition-all"
          aria-label="User menu">
          <img v-if="user?.avatar_url" :src="user.avatar_url" :alt="user.username" class="w-full h-full object-cover" />
          <div v-else class="w-full h-full bg-surface flex items-center justify-center">
            <UserCircleIcon class="w-5 h-5 text-t3" />
          </div>
        </button>

        <div v-if="menuOpen"
          class="absolute right-0 mt-2 w-48 bg-panel border border-b2 rounded-xl shadow-lg overflow-hidden z-50">
          <div class="px-4 py-3 border-b border-b1">
            <p class="text-sm font-medium text-t1 truncate">{{ user?.username }}</p>
          </div>
          <div class="border-b border-b1">
            <button @click="goToProfile"
              class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-t2 hover:bg-surface hover:text-t1 transition-colors">
              <UserIcon class="w-4 h-4" />
              Profile
            </button>
            <button @click="goToSettings"
              class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-t2 hover:bg-surface hover:text-t1 transition-colors">
              <Cog6ToothIcon class="w-4 h-4" />
              Settings
            </button>
          </div>
          <button @click="logout"
            class="w-full flex items-center gap-2 px-4 py-2.5 text-sm text-fail hover:bg-fail-s transition-colors">
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
import { useRouter } from 'vue-router';
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
const router = useRouter();
const search = ref('');
const menuOpen = ref(false);
const menuRef = ref<HTMLElement | null>(null);

useClickOutside(menuRef, () => { menuOpen.value = false; });

function goToProfile() {
  menuOpen.value = false;
  if (user.value) router.push(`/users/${user.value.id}`);
}

function goToSettings() {
  menuOpen.value = false;
  router.push('/settings');
}

onMounted(fetch);

async function logout() {
  clear();
  await window.fetch('/auth/logout', { method: 'POST' });
  window.location.href = '/login';
}
</script>
