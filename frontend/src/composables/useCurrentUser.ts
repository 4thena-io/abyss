import { ref } from 'vue';

interface CurrentUser {
  id: number;
  username: string;
  is_admin: boolean;
  avatar_url: string;
}

const user = ref<CurrentUser | null>(null);
let fetched = false;

export function useCurrentUser() {
  async function fetch() {
    if (fetched) return;
    fetched = true;
    try {
      const res = await window.fetch('/api/me');
      if (res.ok) user.value = await res.json();
    } catch {
      // ignore
    }
  }

  function clear() {
    user.value = null;
    fetched = false;
  }

  return { user, fetch, clear };
}
