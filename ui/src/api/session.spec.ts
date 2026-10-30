import { afterEach, describe, expect, it, vi } from 'vitest';
import { checkAuthenticated, checkConfigured } from './session';

afterEach(() => vi.unstubAllGlobals());

describe('session checks during outages', () => {
  it('only treats 401 as logged out and retries after recovery', async () => {
    const fetch = vi
      .fn()
      .mockResolvedValueOnce(new Response('', { status: 503 }))
      .mockResolvedValueOnce(new Response('{}', { status: 200 }))
      .mockResolvedValueOnce(new Response('', { status: 401 }));
    vi.stubGlobal('fetch', fetch);
    await expect(checkAuthenticated()).rejects.toThrow('temporarily unavailable');
    await expect(checkAuthenticated()).resolves.toBe(true);
    await expect(checkAuthenticated()).resolves.toBe(false);
  });

  it('never treats a setup outage as an unconfigured installation', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('', { status: 503 })));
    await expect(checkConfigured()).rejects.toThrow('temporarily unavailable');
  });

  it('does not mistake a network failure for logout', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('offline')));
    await expect(checkAuthenticated()).rejects.toThrow('offline');
  });
});
