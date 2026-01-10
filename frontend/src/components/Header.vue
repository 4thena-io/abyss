<template>
  <header class="bg-gray-950 h-12 flex items-center justify-between">
    <div class="flex w-12 items-center justify-center shrink-0">
      <button 
        @click="$emit('toggle-sidebar')"
        class="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-gray-800"
        aria-label="Toggle navigation menu"
      >
        <Bars3Icon class="w-5 h-5" />
      </button>
    </div>
    <div class="flex-1 flex items-center justify-start min-w-0">
      <img :src="logo" alt="logo" class="w-auto h-12 cursor-pointer hover:opacity-80 transition-opacity" />
      <div class="flex-1 flex justify-center px-8 min-w-0">
        <div class="w-full max-w-2xl relative">
          <input type="text" placeholder="Search for projects and applications"
            class="w-full text-gray-400 rounded-lg px-4 h-8 pl-10 border border-gray-400 focus:outline-none" />
          <MagnifyingGlassIcon class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 transform -translate-y-1/2" />
        </div>
      </div>
    </div>
    <div class="flex items-center justify-end w-24 space-x-4 pr-4">
      <button class="text-gray-400 hover:text-white transition-colors relative p-2 rounded-lg hover:bg-gray-800"
        aria-label="Notifications">
        <BellIcon class="w-5 h-5" />
      </button>
      <button class="text-gray-400 hover:text-white transition-colors p-2 rounded-lg hover:bg-gray-800"
        aria-label="Settings">
        <Cog6ToothIcon class="w-5 h-5" />
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import logo from '../assets/logo-white.svg'
import { Bars3Icon, BellIcon, Cog6ToothIcon, MagnifyingGlassIcon } from '@heroicons/vue/24/outline';

defineEmits<{
  'toggle-sidebar': []
}>();
</script>
