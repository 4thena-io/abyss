<template>
  <div class="space-y-6">
    <PageHeader title="Teams" :subtitle="`${teamList.length} teams`">
      <template #actions>
        <Button @click="showModal = true" variant="secondary">
          <PlusIcon class="w-4 h-4" />
          New Team
        </Button>
      </template>
    </PageHeader>

    <!-- Filters -->
    <div class="flex items-center gap-4">
      <SearchInput v-model="search" placeholder="Filter teams..." />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <Spinner size="lg" />
    </div>

    <!-- Empty State -->
    <EmptyState
      v-else-if="filteredTeams.length === 0"
      :icon="UsersIcon"
      title="No teams found"
      :message="search ? 'Try adjusting your search' : 'Create your first team to get started'"
    />

    <!-- Team Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <EntityCard
        v-for="team in filteredTeams"
        :key="team.id"
        :icon="UsersIcon"
        :title="team.name"
        :subtitle="`${team.memberCount} members`"
        :description="team.description"
        :to="`/teams/${team.id}`"
      >
        <template #badges>
          <TagGroup :tags="[`${team.projectCount} projects`, `${team.appCount} apps`]" />
        </template>
      </EntityCard>
    </div>

    <!-- Create Team Modal -->
    <Modal :open="showModal" title="New Team" @close="closeModal">
      <form @submit.prevent="createTeam" class="space-y-4">
        <FormField
          v-model="form.name"
          label="Name"
          required
          placeholder="my-team"
        />

        <FormField
          v-model="form.description"
          label="Description"
          type="textarea"
          :rows="3"
          placeholder="Team description..."
        />

        <ErrorAlert v-if="error" :message="error" />

        <FormActions
          submit-label="Create Team"
          submitting-label="Creating..."
          :disabled="!canCreate"
          :saving="saving"
          @cancel="closeModal"
        />
      </form>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue';
import { PlusIcon, UsersIcon } from '@heroicons/vue/24/outline';
import type { Team } from '../types/Team';
import { teamsApi, type CreateTeamRequest } from '../api';
import Modal from '../components/ui/Modal.vue';
import Button from '../components/ui/Button.vue';
import Spinner from '../components/ui/Spinner.vue';
import SearchInput from '../components/ui/SearchInput.vue';
import EmptyState from '../components/ui/EmptyState.vue';
import ErrorAlert from '../components/ui/ErrorAlert.vue';
import FormField from '../components/ui/FormField.vue';
import PageHeader from '../components/ui/PageHeader.vue';
import FormActions from '../components/ui/FormActions.vue';
import EntityCard from '../components/ui/EntityCard.vue';
import TagGroup from '../components/ui/TagGroup.vue';

const teamList = ref<Team[]>([]);
const loading = ref(true);
const search = ref('');

const showModal = ref(false);
const saving = ref(false);
const error = ref('');

const form = reactive<CreateTeamRequest>({
  name: '',
  description: '',
});

const filteredTeams = computed(() => {
  return teamList.value.filter(team => {
    return !search.value ||
      team.name.toLowerCase().includes(search.value.toLowerCase()) ||
      team.description?.toLowerCase().includes(search.value.toLowerCase());
  });
});

const canCreate = computed(() => form.name.trim() !== '');

const resetForm = () => {
  form.name = '';
  form.description = '';
  error.value = '';
};

const closeModal = () => {
  showModal.value = false;
  resetForm();
};

const fetchTeams = async () => {
  try {
    loading.value = true;
    teamList.value = await teamsApi.getAll();
  } catch (e) {
    console.error('Failed to fetch teams:', e);
  } finally {
    loading.value = false;
  }
};

const createTeam = async () => {
  error.value = '';
  saving.value = true;

  try {
    const team = await teamsApi.create(form);
    teamList.value.push(team);
    closeModal();
  } catch (e: any) {
    error.value = e.message;
  } finally {
    saving.value = false;
  }
};

onMounted(fetchTeams);
</script>
