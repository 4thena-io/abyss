import { describe, it, expect, beforeEach, vi } from 'vitest';
import { api } from './client';

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

describe('api client', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn());
  });

  it('returns parsed JSON on success', async () => {
    vi.mocked(fetch).mockResolvedValue(jsonResponse({ id: 1, name: 'app' }));

    const result = await api.get<{ id: number; name: string }>('/apps/1');

    expect(result).toEqual({ id: 1, name: 'app' });
    expect(fetch).toHaveBeenCalledWith(
      '/api/apps/1',
      expect.objectContaining({
        headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
      }),
    );
  });

  it('returns undefined for a 204 No Content response without parsing a body', async () => {
    // A real 204 has no body; simulate that rather than an empty JSON payload,
    // since response.json() on a truly empty body would throw.
    const response = new Response(null, { status: 204 });
    vi.mocked(fetch).mockResolvedValue(response);

    const result = await api.delete<void>('/apps/1');

    expect(result).toBeUndefined();
  });

  it('throws the server-provided message on a non-2xx response', async () => {
    vi.mocked(fetch).mockResolvedValue(jsonResponse({ message: 'app not found' }, 404));

    await expect(api.get('/apps/999')).rejects.toThrow('app not found');
  });

  it('falls back to a generic message when the error body is not valid JSON', async () => {
    const response = new Response('not json', { status: 500 });
    vi.mocked(fetch).mockResolvedValue(response);

    await expect(api.get('/apps/1')).rejects.toThrow('Request failed');
  });

  it('sends a JSON body for post', async () => {
    vi.mocked(fetch).mockResolvedValue(jsonResponse({ ok: true }));

    await api.post('/apps', { name: 'new-app' });

    expect(fetch).toHaveBeenCalledWith(
      '/api/apps',
      expect.objectContaining({ method: 'POST', body: JSON.stringify({ name: 'new-app' }) }),
    );
  });
});
