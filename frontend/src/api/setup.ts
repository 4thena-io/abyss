export interface SetupDefaults {
  db_type: string;
  db_path: string;
  db_host: string;
  db_port: string;
  db_user: string;
  db_name: string;
  forge_type: string;
  forge_host: string;
  forge_owner: string;
  ci_type: string;
  ci_host: string;
  client_id: string;
}

export interface SetupStatus {
  configured: boolean;
  defaults?: SetupDefaults;
}

export interface SetupRequest {
  // Database
  db_type: string;
  db_path: string;
  db_host: string;
  db_port: string;
  db_user: string;
  db_password: string;
  db_name: string;

  // Forge
  forge_type: string;
  forge_host: string;
  forge_token: string;
  forge_owner: string;

  // OAuth
  client_id: string;
  client_secret: string;
  callback_url: string;

  // CI
  ci_type: string;
  ci_host: string;
  ci_token: string;
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options);
  const data = await res.json().catch(() => ({ message: 'Request failed' }));
  if (!res.ok) throw new Error((data as { message: string }).message ?? 'Request failed');
  return data as T;
}

export const setupApi = {
  status: () => request<SetupStatus>('/api/setup/status'),

  configure: (data: SetupRequest) =>
    request<{ message: string }>('/api/setup/configure', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    }),
};
