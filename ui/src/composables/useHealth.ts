import { ref, onMounted, onUnmounted } from 'vue';
import { fetchHealth, type HealthStatus } from '../api/health';

const POLL_MS = 30_000;

export function useHealth() {
  const status = ref<HealthStatus | null>(null);

  async function check() {
    try {
      status.value = await fetchHealth();
    } catch {
      // If the health endpoint itself is unreachable, assume everything is down.
      status.value = {
        forge: { ok: false, error: 'unreachable' },
        ci: { ok: false, error: 'unreachable' },
      };
    }
  }

  let timer: ReturnType<typeof setInterval>;

  onMounted(() => {
    check();
    timer = setInterval(check, POLL_MS);
  });

  onUnmounted(() => clearInterval(timer));

  return { status, check };
}
