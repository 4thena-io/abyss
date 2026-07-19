import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import AppOverviewTab from './AppOverviewTab.vue';
import type { App, Build } from '../types';

vi.mock('../../../api/app', () => ({
  appsApi: {
    getBuilds: vi.fn(),
    getDeployments: vi.fn(),
  },
}));

// Regression test: activity text used to be built as a raw HTML string and
// rendered via v-html, letting a malicious branch name (fully attacker
// controlled by anyone with push access) inject markup into every viewer's
// browser. It's now plain Vue interpolation, which escapes by default.
describe('AppOverviewTab', () => {
  it('renders a malicious branch name as literal text, not markup', async () => {
    const { appsApi } = await import('../../../api/app');
    const maliciousBuild: Build = {
      id: 1,
      number: 42,
      status: 'success',
      branch: '<img src=x onerror="window.__pwned=true">',
      commit: 'abc123',
      duration: 30,
      startedAt: new Date().toISOString(),
      link: 'https://ci.example/1',
    };
    vi.mocked(appsApi.getBuilds).mockResolvedValue([maliciousBuild]);
    vi.mocked(appsApi.getDeployments).mockResolvedValue([]);

    const app: App = {
      id: 1,
      name: 'my-app',
      description: '',
      kind: '',
      language: '',
      repoUrl: '',
      ciUrl: '',
      projectId: 0,
      creatorId: 1,
      creatorUsername: 'alice',
    };

    const wrapper = mount(AppOverviewTab, {
      props: { app },
      global: { stubs: { RouterLink: true } },
    });

    await new Promise((resolve) => setTimeout(resolve, 0));
    await wrapper.vm.$nextTick();

    expect(wrapper.find('img').exists()).toBe(false);
    expect(wrapper.text()).toContain('<img src=x onerror="window.__pwned=true">');
    expect((window as unknown as { __pwned?: boolean }).__pwned).toBeUndefined();
  });
});
