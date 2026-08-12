<template>
  <div class="space-y-6">
    <!-- General -->
    <div class="bg-panel border border-b1 rounded-xl p-5 space-y-4">
      <h3 class="text-sm font-semibold text-t1">General</h3>

      <FormField v-model="form.name" label="Name" required />
      <FormField v-model="form.description" label="Description" type="textarea" :rows="2" />

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
          <p class="text-sm font-medium text-t1">Delete this team</p>
          <p class="text-xs text-t3 mt-0.5">
            Permanently removes the team and all memberships. Projects remain but lose their team
            association.
          </p>
        </div>
        <button
          @click="showDelete = true"
          class="ml-4 shrink-0 px-4 py-2 text-sm font-medium text-fail border border-fail/40 hover:bg-fail-s rounded-lg transition-colors"
        >
          Delete team
        </button>
      </div>
    </div>
  </div>

  <ConfirmModal
    :open="showDelete"
    title="Delete team"
    :message="`This will permanently delete ${team.name} and remove all members. Projects will remain but lose their team association. This cannot be undone.`"
    confirm-label="Delete team"
    :confirm-name="team.name"
    @confirm="deleteTeam"
    @cancel="showDelete = false"
  />
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline';
import type { Team } from '../types';
import { teamsApi } from '../../../api/team';
import FormField from '../../../components/ui/FormField.vue';
import ConfirmModal from '../../../components/ui/ConfirmModal.vue';

const props = defineProps<{ team: Team }>();
const emit = defineEmits<{ updated: [team: Team] }>();

const router = useRouter();

const form = reactive({
  name: props.team.name,
  description: props.team.description ?? '',
});

const saving = ref(false);
const saved = ref(false);
const saveError = ref('');
const showDelete = ref(false);

const dirty = computed(
  () => form.name !== props.team.name || form.description !== (props.team.description ?? ''),
);

watch(
  () => props.team,
  (t) => {
    form.name = t.name;
    form.description = t.description ?? '';
  },
  { deep: true },
);

function reset() {
  form.name = props.team.name;
  form.description = props.team.description ?? '';
  saveError.value = '';
}

async function save() {
  saving.value = true;
  saved.value = false;
  saveError.value = '';
  try {
    const updated = await teamsApi.update(props.team.id, {
      name: form.name,
      description: form.description,
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

async function deleteTeam() {
  await teamsApi.delete(props.team.id);
  router.push('/teams');
}
</script>
