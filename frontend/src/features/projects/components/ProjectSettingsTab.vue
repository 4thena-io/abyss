<template>
  <div class="space-y-6">
    <!-- General -->
    <div class="bg-panel border border-b1 rounded-xl p-5 space-y-4">
      <h3 class="text-sm font-semibold text-t1">General</h3>

      <FormField v-model="form.name" label="Name" required />
      <FormField v-model="form.description" label="Description" type="textarea" :rows="2" />

      <div class="space-y-1.5">
        <label class="block text-sm font-medium text-t2">Team</label>
        <Dropdown
          v-model="form.teamId"
          variant="form"
          placeholder="No team (standalone)"
          :options="teamOptions"
        />
        <p class="text-xs text-t3">Assign or reassign this project to a team.</p>
      </div>

      <div v-if="saveError" class="flex items-start gap-2 bg-fail-s border border-fail/30 rounded-lg px-3 py-2.5">
        <ExclamationTriangleIcon class="w-4 h-4 text-fail shrink-0 mt-0.5" />
        <p class="text-fail text-xs">{{ saveError }}</p>
      </div>

      <div class="flex items-center gap-2">
        <button
          @click="save"
          :disabled="saving || !dirty"
          class="px-4 py-2 text-sm font-medium text-white bg-acc hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition-opacity flex items-center gap-2"
        >
          <span v-if="saving" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
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
          <p class="text-sm font-medium text-t1">Delete this project</p>
          <p class="text-xs text-t3 mt-0.5">Permanently removes the project. Applications in this project will become standalone.</p>
        </div>
        <button
          @click="showDelete = true"
          class="ml-4 shrink-0 px-4 py-2 text-sm font-medium text-fail border border-fail/40 hover:bg-fail-s rounded-lg transition-colors"
        >
          Delete project
        </button>
      </div>
    </div>
  </div>

  <ConfirmModal
    :open="showDelete"
    title="Delete project"
    :message="`This will permanently delete ${project.name}. Applications in this project will become standalone. This action cannot be undone.`"
    confirm-label="Delete project"
    :confirm-name="project.name"
    @confirm="deleteProject"
    @cancel="showDelete = false"
  />
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline';
import type { Project } from '../types';
import type { Team } from '../../teams/types';
import { projectsApi } from '../../../api/project';
import { teamsApi } from '../../../api/team';
import FormField from '../../../components/ui/FormField.vue';
import Dropdown from '../../../components/ui/Dropdown.vue';
import ConfirmModal from '../../../components/ui/ConfirmModal.vue';

const props = defineProps<{ project: Project }>();
const emit = defineEmits<{ updated: [project: Project] }>();

const router = useRouter();
const teams = ref<Team[]>([]);

const form = reactive({
  name: props.project.name,
  description: props.project.description ?? '',
  teamId: (props.project.teamId ?? '') as number | '',
});

const saving = ref(false);
const saved = ref(false);
const saveError = ref('');
const showDelete = ref(false);

const teamOptions = computed(() =>
  teams.value.map(t => ({ value: t.id, label: t.name }))
);

const dirty = computed(() =>
  form.name !== props.project.name ||
  form.description !== (props.project.description ?? '') ||
  form.teamId !== (props.project.teamId ?? '')
);

watch(() => props.project, (p) => {
  form.name = p.name;
  form.description = p.description ?? '';
  form.teamId = p.teamId ?? '';
}, { deep: true });

function reset() {
  form.name = props.project.name;
  form.description = props.project.description ?? '';
  form.teamId = props.project.teamId ?? '';
  saveError.value = '';
}

async function save() {
  saving.value = true;
  saved.value = false;
  saveError.value = '';
  try {
    const updated = await projectsApi.update(props.project.id, {
      name: form.name,
      description: form.description,
      teamId: form.teamId as number || null,
    });
    emit('updated', updated);
    saved.value = true;
    setTimeout(() => (saved.value = false), 2000);
  } catch (e: any) {
    saveError.value = e.message ?? 'Failed to save';
  } finally {
    saving.value = false;
  }
}

async function deleteProject() {
  await projectsApi.delete(props.project.id);
  router.push('/projects');
}

onMounted(async () => {
  teams.value = await teamsApi.getAll();
});
</script>
