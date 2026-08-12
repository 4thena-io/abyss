import { describe, it, expect, vi, afterEach } from 'vitest';
import { mount } from '@vue/test-utils';
import AppDocsTab from './AppDocsTab.vue';
import type { RenderedDocs } from '../types';

afterEach(() => {
  vi.unstubAllGlobals();
});

function mockFetchOnce(body: unknown, status = 200) {
  vi.stubGlobal(
    'fetch',
    vi.fn().mockResolvedValue({
      status,
      ok: status >= 200 && status < 300,
      json: () => Promise.resolve(body),
    }),
  );
}

describe('AppDocsTab', () => {
  it('renders the nav tree and selects index.md by default', async () => {
    const docs: RenderedDocs = {
      pages: [
        { path: 'index.md', html: '<h1>Home</h1>' },
        { path: 'guide.md', html: '<h1>Guide</h1>' },
      ],
      nav: [
        { title: 'Overview', path: 'index.md', children: null },
        { title: 'Guide', path: 'guide.md', children: null },
      ],
    };
    mockFetchOnce(docs);

    const wrapper = mount(AppDocsTab, { props: { appId: 1 } });
    await new Promise((resolve) => setTimeout(resolve, 0));
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain('Overview');
    expect(wrapper.text()).toContain('Guide');
    expect(wrapper.html()).toContain('Home');
  });

  it('renders without throwing when a cached docs.json predates the nav field', async () => {
    mockFetchOnce({ pages: [{ path: 'index.md', html: '<h1>Home</h1>' }] });

    const wrapper = mount(AppDocsTab, { props: { appId: 1 } });
    await new Promise((resolve) => setTimeout(resolve, 0));
    await wrapper.vm.$nextTick();

    expect(wrapper.html()).toContain('Home');
  });
});
