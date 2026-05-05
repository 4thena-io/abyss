import type { UserProfile } from '../features/users/types';

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options);
  if (res.status === 204) return undefined as T;
  const data = await res.json().catch(() => ({ message: 'Request failed' }));
  if (!res.ok) throw new Error((data as { message: string }).message ?? 'Request failed');
  return data as T;
}

export const userApi = {
  tokenStatus: () => request<{ has_token: boolean }>('/api/user/token'),
  generateToken: () => request<{ token: string }>('/api/user/token', { method: 'POST' }),
  revokeToken: () => request<void>('/api/user/token', { method: 'DELETE' }),
  getProfile: (id: number) => request<UserProfile>(`/api/users/${id}`),
};
