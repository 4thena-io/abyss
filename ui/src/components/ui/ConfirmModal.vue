<template>
  <Modal :open="open" :title="title" @close="$emit('cancel')">
    <div class="space-y-4">
      <p class="text-sm text-t2">{{ message }}</p>

      <div v-if="confirmName" class="space-y-1.5">
        <label class="block text-xs text-t3">
          Type <strong class="text-t1 font-mono">{{ confirmName }}</strong> to confirm
        </label>
        <input
          v-model="typed"
          type="text"
          :placeholder="confirmName"
          class="w-full px-3 py-2 text-sm bg-bg border border-b1 rounded-lg text-t1 placeholder:text-t3 focus:outline-none focus:border-fail/60 transition-colors"
          @keydown.enter.prevent="tryConfirm"
        />
      </div>

      <div class="flex justify-end gap-2 pt-1">
        <button
          type="button"
          @click="$emit('cancel')"
          class="px-4 py-2 text-sm text-t2 border border-b1 rounded-lg hover:text-t1 hover:border-b2 transition-colors"
        >
          Cancel
        </button>
        <button
          type="button"
          @click="tryConfirm"
          :disabled="!!confirmName && typed !== confirmName"
          class="px-4 py-2 text-sm font-medium text-white bg-fail hover:opacity-90 disabled:opacity-40 disabled:cursor-not-allowed rounded-lg transition-opacity"
        >
          {{ confirmLabel }}
        </button>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import Modal from './Modal.vue';

const props = withDefaults(
  defineProps<{
    open: boolean;
    title: string;
    message: string;
    confirmLabel?: string;
    confirmName?: string;
  }>(),
  {
    confirmLabel: 'Delete',
  },
);

const emit = defineEmits<{
  confirm: [];
  cancel: [];
}>();

const typed = ref('');

watch(
  () => props.open,
  (v) => {
    if (!v) typed.value = '';
  },
);

function tryConfirm() {
  if (props.confirmName && typed.value !== props.confirmName) return;
  emit('confirm');
}
</script>
