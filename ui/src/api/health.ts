export interface ServiceHealth {
  ok: boolean;
  error?: string;
}

export interface HealthStatus {
  forge: ServiceHealth;
  ci: ServiceHealth;
}

export async function fetchHealth(): Promise<HealthStatus> {
  const res = await fetch('/api/health');
  if (!res.ok) throw new Error('health check failed');
  return res.json();
}
