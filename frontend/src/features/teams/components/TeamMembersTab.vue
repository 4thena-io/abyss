<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium text-t1">Members</h2>
      <div class="flex items-center gap-2">
        <RefreshButton :loading="loading" @click="fetchMembers" />
        <button
          @click="showAddModal = true"
          class="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium bg-acc hover:opacity-90 text-white rounded-lg transition-opacity"
        >
          <PlusIcon class="w-4 h-4" />
          Add member
        </button>
      </div>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <Spinner size="md" />
    </div>

    <EmptyState
      v-else-if="members.length === 0"
      :icon="UsersIcon"
      title="No members"
      message="This team has no members yet"
    />

    <div v-else class="bg-panel border border-b1 rounded-xl overflow-hidden">
      <table class="w-full border-collapse">
        <thead>
          <tr>
            <th class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1">Member</th>
            <th class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1">Role</th>
            <th class="text-left text-xs font-semibold uppercase tracking-wide text-t3 px-4 py-2.5 border-b border-b1">Joined</th>
            <th class="w-24 px-4 py-2.5 border-b border-b1" />
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="member in members"
            :key="member.id"
            class="group hover:bg-surface transition-colors"
          >
            <td class="px-4 py-3 border-b border-b0 cursor-pointer" @click="router.push(`/users/${member.userId}`)">
              <div class="flex items-center gap-3">
                <img v-if="member.avatarUrl" :src="member.avatarUrl" :alt="member.username"
                  class="w-8 h-8 rounded-full object-cover shrink-0" />
                <div v-else
                  class="w-8 h-8 rounded-full bg-acc-s border border-acc-b flex items-center justify-center shrink-0">
                  <span class="text-xs font-semibold text-acc font-mono">{{ member.username.slice(0, 2) }}</span>
                </div>
                <div>
                  <p class="text-sm font-medium text-t1">{{ member.username }}</p>
                  <p v-if="member.email" class="text-xs text-t3 font-mono">{{ member.email }}</p>
                </div>
              </div>
            </td>
            <td class="px-4 py-3 border-b border-b0">
              <span class="px-2 py-0.5 rounded-full text-xs font-medium"
                :class="member.role === 'owner' ? 'bg-ok-s text-ok' : 'bg-surface border border-b1 text-t2'">
                {{ member.role }}
              </span>
            </td>
            <td class="px-4 py-3 border-b border-b0 text-sm text-t3">
              {{ formatTime(member.joinedAt) }}
            </td>
            <td class="px-4 py-3 border-b border-b0 text-right">
              <template v-if="confirmRemoveId === member.id">
                <div class="flex items-center justify-end gap-1">
                  <button
                    @click="removeMember(member.id)"
                    :disabled="removing === member.id"
                    class="px-2 py-1 text-xs font-medium text-white bg-fail hover:opacity-90 rounded transition-opacity disabled:opacity-50"
                  >
                    {{ removing === member.id ? '…' : 'Confirm' }}
                  </button>
                  <button
                    @click="confirmRemoveId = null"
                    class="px-2 py-1 text-xs text-t2 hover:text-t1 border border-b1 rounded transition-colors"
                  >
                    Cancel
                  </button>
                </div>
              </template>
              <button
                v-else
                @click.stop="confirmRemoveId = member.id"
                class="opacity-0 group-hover:opacity-100 px-2 py-1 text-xs text-fail hover:bg-fail-s border border-transparent hover:border-fail/30 rounded transition-all"
              >
                Remove
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add member modal -->
    <Modal :open="showAddModal" title="Add member" @close="closeAddModal">
      <form @submit.prevent="addMember" class="space-y-4">
        <FormField
          v-model="addForm.username"
          label="Username"
          required
          placeholder="username"
          hint="Must match their forge username"
        />

        <div class="space-y-1.5">
          <label class="block text-sm font-medium text-t2">Role</label>
          <div class="flex gap-2">
            <button
              v-for="opt in roleOptions"
              :key="opt.value"
              type="button"
              @click="addForm.role = opt.value"
              class="flex-1 py-2 px-3 rounded-lg border text-sm font-medium transition-colors"
              :class="addForm.role === opt.value
                ? 'border-acc bg-acc-s text-t1'
                : 'border-b1 bg-bg text-t2 hover:border-b2'"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>

        <div v-if="addError" class="flex items-start gap-2 bg-fail-s border border-fail/30 rounded-lg px-3 py-2.5">
          <ExclamationTriangleIcon class="w-4 h-4 text-fail shrink-0 mt-0.5" />
          <p class="text-fail text-xs">{{ addError }}</p>
        </div>

        <FormActions
          submit-label="Add member"
          submitting-label="Adding..."
          :disabled="!addForm.username.trim()"
          :saving="adding"
          @cancel="closeAddModal"
        />
      </form>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { UsersIcon, PlusIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline';
import type { TeamMember } from '../types';
import { teamsApi } from '../../../api/team';
import { useFormatters } from '../../../composables/useFormatters';
import Spinner from '../../../components/ui/Spinner.vue';
import EmptyState from '../../../components/ui/EmptyState.vue';
import RefreshButton from '../../../components/ui/RefreshButton.vue';
import Modal from '../../../components/ui/Modal.vue';
import FormField from '../../../components/ui/FormField.vue';
import FormActions from '../../../components/ui/FormActions.vue';

const props = defineProps<{ teamId: number }>();
const router = useRouter();
const { formatRelativeTime: formatTime } = useFormatters();

const members = ref<TeamMember[]>([]);
const loading = ref(false);

// Remove state
const confirmRemoveId = ref<number | null>(null);
const removing = ref<number | null>(null);

// Add modal state
const showAddModal = ref(false);
const adding = ref(false);
const addError = ref('');
const addForm = reactive({ username: '', role: 'member' });

const roleOptions = [
  { value: 'member', label: 'Member' },
  { value: 'owner', label: 'Owner' },
];

const fetchMembers = async () => {
  try {
    loading.value = true;
    members.value = await teamsApi.getMembers(props.teamId);
  } catch (error) {
    console.error('Failed to fetch members:', error);
    members.value = [];
  } finally {
    loading.value = false;
  }
};

const closeAddModal = () => {
  showAddModal.value = false;
  addForm.username = '';
  addForm.role = 'member';
  addError.value = '';
};

const addMember = async () => {
  adding.value = true;
  addError.value = '';
  try {
    const member = await teamsApi.addMember(props.teamId, addForm);
    members.value.push(member);
    closeAddModal();
  } catch (e: any) {
    addError.value = e.message ?? 'Failed to add member';
  } finally {
    adding.value = false;
  }
};

const removeMember = async (memberId: number) => {
  removing.value = memberId;
  try {
    await teamsApi.removeMember(props.teamId, memberId);
    members.value = members.value.filter(m => m.id !== memberId);
    confirmRemoveId.value = null;
  } catch (e) {
    console.error('Failed to remove member:', e);
  } finally {
    removing.value = null;
  }
};

onMounted(fetchMembers);
</script>
