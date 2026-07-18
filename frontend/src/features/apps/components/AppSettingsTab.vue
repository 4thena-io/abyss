<template>
  <div class="space-y-6">
    <!-- General -->
    <div class="bg-panel border border-b1 rounded-xl p-5 space-y-4">
      <h3 class="text-sm font-semibold text-t1">General</h3>

      <FormField v-model="form.name" label="Name" required />
      <FormField v-model="form.description" label="Description" type="textarea" :rows="2" />

      <div class="space-y-1.5">
        <label class="block text-sm font-medium text-t2">Project</label>
        <Dropdown
          v-model="form.projectId"
          variant="form"
          placeholder="No project (standalone)"
          :options="projectOptions"
        />
        <p class="text-xs text-t3">
          Reassign this app to a different project, or remove it from any project.
        </p>
      </div>

      <div
        v-if="saveError"
        class="flex items-start gap-2 bg-fail-s border border-fail/30 rounded-lg px-3 py-2.5"
      >
        <ExclamationTriangleIcon class="w-4 h-4 text-fail shrink-0 mt-0.5" />
        <p class="text-fail text-xs">{{ saveError }}</p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="save"
          :disabled="saving || !dirty"
          class="px-4 py-2 text-sm font-medium text-white bg-acc hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-opacity flex items-center gap-2"
        >
          <span
            v-if="saving"
            class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
          />
          {{ saved ? 'Saved!' : 'Save changes' }}
        </button>
        <button
          v-if="dirty"
          @click="reset"
          class="px-4 py-2 text-sm text-t2 border border-b1 hover:text-t1 hover:border-b2 rounded-lg transition-colors"
        >
          Discard
        </button>
      </div>
    </div>

    <!-- Danger zone -->
    <div class="bg-panel border border-fail/40 rounded-xl p-5 space-y-3">
      <h3 class="text-sm font-semibold text-fail">Danger zone</h3>
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-medium text-t1">Delete this application</p>
          <p class="text-xs text-t3 mt-0.5">
            Permanently removes the app, its CI pipeline, and its repository. This cannot be undone.
          </p>
        </div>
        <button
          @click="showDelete = true"
          class="ml-4 shrink-0 px-4 py-2 text-sm font-medium text-fail border border-fail/40 hover:bg-fail-s rounded-lg transition-colors"
        >
          Delete app
        </button>
      </div>
    </div>
  </div>

  <ConfirmModal
    :open="showDelete"
    title="Delete application"
    :message="`This will permanently delete ${app.name}, its CI pipeline, and its repository. This action cannot be undone.`"
    confirm-label="Delete application"
    :confirm-name="app.name"
    @confirm="deleteApp"
    @cancel="showDelete = false"
  />
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline';
import type { App } from '../types';
import type { Project } from '../../projects/types';
import { appsApi } from '../../../api/app';
import { projectsApi } from '../../../api/project';
import FormField from '../../../components/ui/FormField.vue';
import Dropdown from '../../../components/ui/Dropdown.vue';
import ConfirmModal from '../../../components/ui/ConfirmModal.vue';

const props = defineProps<{ app: App }>();
const emit = defineEmits<{ updated: [app: App] }>();

const router = useRouter();
const projects = ref<Project[]>([]);

const form = reactive({
  name: props.app.name,
  description: props.app.description ?? '',
  projectId: props.app.projectId as number | '',
});

const saving = ref(false);
const saved = ref(false);
const saveError = ref('');
const showDelete = ref(false);

const projectOptions = computed(() => projects.value.map((p) => ({ value: p.id, label: p.name })));

const dirty = computed(
  () =>
    form.name !== props.app.name ||
    form.description !== (props.app.description ?? '') ||
    form.projectId !== (props.app.projectId || ''),
);

watch(
  () => props.app,
  (a) => {
    form.name = a.name;
    form.description = a.description ?? '';
    form.projectId = a.projectId || '';
  },
  { deep: true },
);

function reset() {
  form.name = props.app.name;
  form.description = props.app.description ?? '';
  form.projectId = props.app.projectId || '';
  saveError.value = '';
}

async function save() {
  saving.value = true;
  saved.value = false;
  saveError.value = '';
  try {
    const updated = await appsApi.update(props.app.id, {
      name: form.name,
      description: form.description,
      projectId: (form.projectId as number) || null,
    });
    emit('updated', updated);
    saved.value = true;
    setTimeout(() => (saved.value = false), 2000);
  } catch (e) {
    saveError.value = e instanceof Error ? e.message : 'Failed to save';
  } finally {
    saving.value = false;
  }
}

async function deleteApp() {
  await appsApi.delete(props.app.id);
  router.push('/apps');
}

onMounted(async () => {
  projects.value = await projectsApi.getAll();
});
</script>
