import { ref } from 'vue';

interface CurrentUser {
  id: number;
  username: string;
  email: string;
  is_admin: boolean;
  avatar_url: string;
  forge_type: string;
  forge_host: string;
}

const user = ref<CurrentUser | null>(null);
let fetched = false;

export function useCurrentUser() {
  async function fetch() {
    if (fetched) return;
    try {
      const res = await window.fetch('/api/me');
      if (res.ok) {
        user.value = await res.json();
        fetched = true;
      } else if (res.status === 401) {
        user.value = null;
        fetched = true;
      }
    } catch {
      // Leave the result uncached so recovery can retry without logging out.
    }
  }

  function clear() {
    user.value = null;
    fetched = false;
  }

  return { user, fetch, clear };
}
