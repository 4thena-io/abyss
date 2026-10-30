// An unavailable backend is neither an unconfigured instance nor a logout.
export async function checkConfigured(): Promise<boolean> {
  const response = await fetch('/api/setup/status');
  if (!response.ok) throw new Error('Service temporarily unavailable');
  const data = await response.json();
  if (typeof data.configured !== 'boolean') throw new Error('Invalid setup status');
  return data.configured;
}

export async function checkAuthenticated(): Promise<boolean> {
  const response = await fetch('/api/me');
  if (response.status === 401) return false;
  if (!response.ok) throw new Error('Service temporarily unavailable');
  return true;
}
