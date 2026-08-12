import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import DocsNavTree from './DocsNavTree.vue';
import type { NavNode } from '../types';

describe('DocsNavTree', () => {
  it('renders nested sections at full depth instead of collapsing to one level', () => {
    const nodes: NavNode[] = [
      {
        title: 'A',
        children: [
          {
            title: 'B',
            children: [{ title: 'C', path: 'a/b/c.md', children: null }],
          },
        ],
      },
    ];

    const wrapper = mount(DocsNavTree, { props: { nodes, activePath: '' } });

    expect(wrapper.text()).toContain('A');
    expect(wrapper.text()).toContain('B');
    const leaf = wrapper.findAll('button').find((b) => b.text() === 'C');
    expect(leaf).toBeTruthy();
  });

  it('renders a section with an index page as a clickable header and emits select', async () => {
    const nodes: NavNode[] = [
      {
        title: 'Getting Started',
        path: 'getting-started/index.md',
        children: [
          { title: 'Installation', path: 'getting-started/installation.md', children: null },
        ],
      },
    ];

    const wrapper = mount(DocsNavTree, { props: { nodes, activePath: '' } });

    const header = wrapper.findAll('button').find((b) => b.text().includes('Getting Started'));
    expect(header).toBeTruthy();
    await header!.trigger('click');
    expect(wrapper.emitted('select')?.[0]).toEqual(['getting-started/index.md']);
  });

  it('renders a section without an index page as a static, non-clickable label', () => {
    const nodes: NavNode[] = [
      {
        title: 'API Reference',
        children: [{ title: 'Users', path: 'api/users.md', children: null }],
      },
    ];

    const wrapper = mount(DocsNavTree, { props: { nodes, activePath: '' } });

    const header = wrapper.findAll('button').find((b) => b.text().includes('API Reference'));
    expect(header).toBeUndefined();
    expect(wrapper.find('div').text()).toContain('API Reference');
  });

  it('gives a folder-with-only-an-index-page the same header treatment as one with siblings', () => {
    // Regression test: a folder's clickable-header styling used to depend on
    // whether it happened to have sibling files (children.length > 0),
    // rather than on whether it's a directory at all. An index-only folder
    // (children: []) must render identically to one with siblings, not fall
    // back to the plain-leaf style.
    const nodes: NavNode[] = [
      { title: 'Architecture', path: 'architecture/index.md', children: [] },
    ];

    const wrapper = mount(DocsNavTree, { props: { nodes, activePath: '' } });

    const header = wrapper.findAll('button').find((b) => b.text().includes('Architecture'));
    expect(header).toBeTruthy();
    expect(header!.classes()).toContain('text-t1');
  });
});
