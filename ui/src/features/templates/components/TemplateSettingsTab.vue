<template>
  <div class="space-y-6">
    <!-- General -->
    <div class="bg-panel border border-b1 rounded-xl p-5 space-y-4">
      <h3 class="text-sm font-semibold text-t1">General</h3>

      <FormField v-model="form.name" label="Name" required />
      <FormField v-model="form.description" label="Description" type="textarea" :rows="2" />

      <div class="grid grid-cols-2 gap-4">
        <Dropdown
          v-model="form.kind"
          label="Kind"
          variant="form"
          required
          placeholder="Select kind"
          :options="kindOptions"
        />
        <Dropdown
          v-model="form.language"
          label="Language"
          variant="form"
          required
          placeholder="Select language"
          :options="languageOptions"
        />
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
          <p class="text-sm font-medium text-t1">Delete this template</p>
          <p class="text-xs text-t3 mt-0.5">
            Permanently removes the template. Apps created from it are not affected.
          </p>
        </div>
        <button
          @click="showDelete = true"
          class="ml-4 shrink-0 px-4 py-2 text-sm font-medium text-fail border border-fail/40 hover:bg-fail-s rounded-lg transition-colors"
        >
          Delete template
        </button>
      </div>
    </div>
  </div>

  <ConfirmModal
    :open="showDelete"
    title="Delete template"
    :message="`This will permanently delete ${template.name}. Apps created from it are not affected. This cannot be undone.`"
    confirm-label="Delete template"
    :confirm-name="template.name"
    @confirm="deleteTemplate"
    @cancel="showDelete = false"
  />
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline';
import type { Template } from '../types';
import { templatesApi } from '../../../api/template';
import FormField from '../../../components/ui/FormField.vue';
import Dropdown from '../../../components/ui/Dropdown.vue';
import ConfirmModal from '../../../components/ui/ConfirmModal.vue';

const props = defineProps<{ template: Template }>();
const emit = defineEmits<{ updated: [template: Template] }>();

const router = useRouter();

const form = reactive({
  name: props.template.name,
  description: props.template.description ?? '',
  kind: props.template.kind,
  language: props.template.language,
});

const saving = ref(false);
const saved = ref(false);
const saveError = ref('');
const showDelete = ref(false);

const kindOptions = [
  { value: 'api', label: 'API' },
  { value: 'web', label: 'Web' },
  { value: 'worker', label: 'Worker' },
  { value: 'cli', label: 'CLI' },
  { value: 'library', label: 'Library' },
  { value: 'job', label: 'Job' },
];

const languageOptions = [
  { value: 'go', label: 'Go' },
  { value: 'typescript', label: 'TypeScript' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'python', label: 'Python' },
  { value: 'java', label: 'Java' },
  { value: 'rust', label: 'Rust' },
  { value: 'csharp', label: 'C#' },
];

const dirty = computed(
  () =>
    form.name !== props.template.name ||
    form.description !== (props.template.description ?? '') ||
    form.kind !== props.template.kind ||
    form.language !== props.template.language,
);

watch(
  () => props.template,
  (t) => {
    form.name = t.name;
    form.description = t.description ?? '';
    form.kind = t.kind;
    form.language = t.language;
  },
  { deep: true },
);

function reset() {
  form.name = props.template.name;
  form.description = props.template.description ?? '';
  form.kind = props.template.kind;
  form.language = props.template.language;
  saveError.value = '';
}

async function save() {
  saving.value = true;
  saved.value = false;
  saveError.value = '';
  try {
    const updated = await templatesApi.update(props.template.id, {
      name: form.name,
      description: form.description,
      kind: form.kind,
      language: form.language,
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

async function deleteTemplate() {
  await templatesApi.delete(props.template.id);
  router.push('/templates');
}
</script>
